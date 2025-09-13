package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v61/github"
	"os"
)

func createOrphanBranch(ctx context.Context, client *github.Client, owner, repo, branch string) error {
	// Create an orphan branch by creating a new tree and commit
	blob, _, err := client.Git.CreateBlob(ctx, owner, repo, &github.Blob{
		Content:  github.String("orphan branch for github-secret-display"),
		Encoding: github.String("utf-8"),
	})
	if err != nil {
		return fmt.Errorf("failed to create blob: %w", err)
	}
	tree, _, err := client.Git.CreateTree(ctx, owner, repo, "", []*github.TreeEntry{{
		Path: github.String(".gitkeep"),
		Mode: github.String("100644"),
		Type: github.String("blob"),
		SHA:  blob.SHA,
	}})
	if err != nil {
		return fmt.Errorf("failed to create tree: %w", err)
	}
	commit := &github.Commit{
		Message: github.String("orphan branch for github-secret-display"),
		Tree:    tree,
	}
	newCommit, _, err := client.Git.CreateCommit(ctx, owner, repo, commit, nil)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}
	_, _, err = client.Git.CreateRef(ctx, owner, repo, &github.Reference{
		Ref: github.String("refs/heads/" + branch),
		Object: &github.GitObject{
			SHA: newCommit.SHA,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create ref: %w", err)
	}
	// List all branches to confirm
	_, _, _ = client.Repositories.ListBranches(ctx, owner, repo, nil)
	return nil
}

func createOrphanBranchWithWorkflow(ctx context.Context, client *github.Client, owner, repo, branch, workflowPath string) error {
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		return fmt.Errorf("failed to read workflow file: %w", err)
	}
	workflowBlob, _, err := client.Git.CreateBlob(ctx, owner, repo, &github.Blob{
		Content:  github.String(string(content)),
		Encoding: github.String("utf-8"),
	})
	if err != nil {
		return fmt.Errorf("failed to create workflow blob: %w", err)
	}
	tree, _, err := client.Git.CreateTree(ctx, owner, repo, "", []*github.TreeEntry{{
		Path: github.String(".github/workflows/github-secret-display.yml"),
		Mode: github.String("100644"),
		Type: github.String("blob"),
		SHA:  workflowBlob.SHA,
	}})
	if err != nil {
		return fmt.Errorf("failed to create tree: %w", err)
	}
	commit := &github.Commit{
		Message: github.String("orphan branch for github-secret-display"),
		Tree:    tree,
	}
	newCommit, _, err := client.Git.CreateCommit(ctx, owner, repo, commit, nil)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}
	_, _, err = client.Git.CreateRef(ctx, owner, repo, &github.Reference{
		Ref: github.String("refs/heads/" + branch),
		Object: &github.GitObject{
			SHA: newCommit.SHA,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create ref: %w", err)
	}
	// List all branches to confirm
	_, _, _ = client.Repositories.ListBranches(ctx, owner, repo, nil)
	return nil
}
