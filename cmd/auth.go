package main

import (
	"context"
	"os"
	"regexp"

	"golang.org/x/oauth2"
	"github.com/google/go-github/v61/github"
)

func getGithubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return github.NewClient(oauth2.NewClient(context.Background(), ts))
}

func getTokenFromGhCli() string {
	// Try to get token from gh CLI config
	configPath := os.ExpandEnv("$HOME/.config/gh/hosts.yml")
	f, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}
	// Find the first oauth_token: <token> (allowing for indented YAML)
	re := regexp.MustCompile(`(?m)^\s*oauth_token:\s*([^\s]+)\s*$`)
	matches := re.FindAllStringSubmatch(string(f), -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		return matches[0][1]
	}
	return ""
}
