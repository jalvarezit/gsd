package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	var repo string
	var secrets multiFlag
	var githubToken string

	flag.StringVar(&repo, "R", "", "GitHub repository in the form org/repo")
	flag.Var(&secrets, "secret", "Secret name to display (can be specified multiple times)")
	flag.StringVar(&githubToken, "github-token", os.Getenv("GITHUB_TOKEN"), "GitHub token (or set GITHUB_TOKEN env)")
	flag.Parse()

	if repo == "" {
		fmt.Fprintln(os.Stderr, "-R org/repo is required")
		os.Exit(1)
	}
	if len(secrets) == 0 {
		fmt.Fprintln(os.Stderr, "--secret is required (at least one)")
		os.Exit(1)
	}

	// Token handling: try flag/env, then gh CLI
	if githubToken == "" {
		githubToken = getTokenFromGhCli()
		if githubToken == "" {
			fmt.Fprintln(os.Stderr, "No GitHub token found via --github-token/GITHUB_TOKEN or gh CLI")
			os.Exit(1)
		}
	}

	client := getGithubClient(githubToken)
	ctx := context.Background()

	// Split repo into owner/repo
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) != 2 {
		fmt.Fprintln(os.Stderr, "-R must be in the form org/repo")
		os.Exit(1)
	}
	owner, repoName := parts[0], parts[1]

	// 1. Render workflow file
	if err := renderWorkflowTemplate(secrets, "github-secret-display.yml"); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to render workflow template:", err)
		os.Exit(1)
	}
	// 2. Create orphan branch with workflow file
	branch := "github-secret-display"
	err := createOrphanBranchWithWorkflow(ctx, client, owner, repoName, branch, "github-secret-display.yml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create orphan branch:", err)
		os.Exit(1)
	}
	// Wait for GitHub to index the workflow (poll for up to 60 seconds)
	workflowFound := false
	for i := 0; i < 12; i++ { // 12 * 5s = 60s
		workflows, _, err := client.Actions.ListWorkflows(ctx, owner, repoName, nil)
		if err == nil {
			for _, wf := range workflows.Workflows {
				if wf.GetName() == "github-secret-display" {
					workflowFound = true
					break
				}
			}
		}
		if workflowFound {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if !workflowFound {
		fmt.Fprintln(os.Stderr, "Failed to dispatch workflow: workflow not found after waiting")
		_ = cleanup(ctx, client, owner, repoName, branch, 0)
		os.Exit(1)
	}
	// 3. Wait for workflow run triggered by push
	runID, err := waitForWorkflowRun(ctx, client, owner, repoName, branch)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to find workflow run:", err)
		_ = cleanup(ctx, client, owner, repoName, branch, 0)
		os.Exit(1)
	}
	// 4. Wait and download artifact
	_, err = waitAndDownloadArtifact(ctx, client, owner, repoName, runID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get artifact:", err)
		_ = cleanup(ctx, client, owner, repoName, branch, runID)
		os.Exit(1)
	}
	// 5. Cleanup
	if err := cleanup(ctx, client, owner, repoName, branch, runID); err != nil {
		fmt.Fprintln(os.Stderr, "Cleanup failed:", err)
		os.Exit(1)
	}
}

type multiFlag []string

func (m *multiFlag) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}
