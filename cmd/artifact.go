package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/go-github/v61/github"
)

const (
	workflowPollInterval = 1 * time.Second
	workflowPollAttempts = 6
	artifactPollInterval = 1 * time.Second
	artifactPollAttempts = 6
)

// pollWorkflowCompletion waits for the workflow run to complete, polling at a fixed interval.
func pollWorkflowCompletion(ctx context.Context, client *github.Client, owner, repo string, runID int64) error {
	for i := 0; i < workflowPollAttempts; i++ {
		run, _, err := client.Actions.GetWorkflowRunByID(ctx, owner, repo, runID)
		if err != nil {
			return fmt.Errorf("failed to get workflow run: %w", err)
		}
		if run.GetStatus() == "completed" {
			return nil
		}
		time.Sleep(workflowPollInterval)
	}
	return fmt.Errorf("workflow run did not complete after polling")
}

// pollArtifactID waits for the artifact to appear and returns its ID.
func pollArtifactID(ctx context.Context, client *github.Client, owner, repo string, runID int64, artifactName string) (int64, error) {
	for i := 0; i < artifactPollAttempts; i++ {
		arts, _, err := client.Actions.ListWorkflowRunArtifacts(ctx, owner, repo, runID, &github.ListOptions{})
		if err != nil {
			return 0, fmt.Errorf("failed to list artifacts: %w", err)
		}
		for _, art := range arts.Artifacts {
			if art.GetName() == artifactName {
				return art.GetID(), nil
			}
		}
		time.Sleep(artifactPollInterval)
	}
	return 0, fmt.Errorf("artifact '%s' not found after polling", artifactName)
}

// waitAndDownloadArtifact waits for workflow completion, polls for the artifact, downloads it, and streams its contents to stdout.
func waitAndDownloadArtifact(ctx context.Context, client *github.Client, owner, repo string, runID int64) (string, error) {
	if err := pollWorkflowCompletion(ctx, client, owner, repo, runID); err != nil {
		return "", err
	}

	artifactID, err := pollArtifactID(ctx, client, owner, repo, runID, "secrets")
	if err != nil {
		return "", err
	}

	url, _, err := client.Actions.DownloadArtifact(ctx, owner, repo, artifactID, 0)
	if err != nil {
		return "", fmt.Errorf("failed to get artifact url: %w", err)
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return "", fmt.Errorf("failed to download artifact: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("artifact download failed with status %d: %s", resp.StatusCode, string(body))
	}

	zipData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := streamAllFilesFromZipToStdout(zipData); err != nil {
		return "", err
	}
	return "", nil
}
