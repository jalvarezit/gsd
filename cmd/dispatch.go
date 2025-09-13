package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v61/github"
)

func waitForWorkflowRun(ctx context.Context, client *github.Client, owner, repo, branch string) (int64, error) {
	workflowFile := "github-secret-display.yml"
	var runID int64
	var latestTime time.Time
	for i := 0; i < 20; i++ {
		runs, _, err := client.Actions.ListWorkflowRunsByFileName(ctx, owner, repo, workflowFile, &github.ListWorkflowRunsOptions{Branch: branch, Event: "push"})
		if err == nil && len(runs.WorkflowRuns) > 0 {
			for _, run := range runs.WorkflowRuns {
				if run.GetHeadBranch() == branch && run.GetCreatedAt().After(latestTime) {
					runID = run.GetID()
					latestTime = run.GetCreatedAt().Time
				}
			}
			if runID != 0 {
				break
			}
		}
		time.Sleep(5 * time.Second)
	}
	if runID == 0 {
		return 0, fmt.Errorf("workflow run not found after push")
	}
	return runID, nil
}
