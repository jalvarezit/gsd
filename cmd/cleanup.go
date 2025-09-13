package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v61/github"
)

func cleanup(ctx context.Context, client *github.Client, owner, repo, branch string, runID int64) error {
	// Delete workflow run
	if runID != 0 {
		_, err := client.Actions.DeleteWorkflowRun(ctx, owner, repo, runID)
		if err != nil {
			return fmt.Errorf("failed to delete workflow run: %w", err)
		}
	}
	// Delete workflow file
	filePath := ".github/workflows/github-secret-display.yml"
	file, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, &github.RepositoryContentGetOptions{Ref: branch})
	if err == nil && file != nil {
		_, _, err = client.Repositories.DeleteFile(ctx, owner, repo, filePath, &github.RepositoryContentFileOptions{
			Message: github.String("delete github-secret-display workflow"),
			SHA:     file.SHA,
			Branch:  github.String(branch),
		})
		if err != nil {
			return fmt.Errorf("failed to delete workflow file: %w", err)
		}
	}
	// Delete branch
	_, err = client.Git.DeleteRef(ctx, owner, repo, "refs/heads/"+branch)
	if err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}
	return nil
}
