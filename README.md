# gsd: GitHub Secret Display CLI

`gsd` is a Go CLI tool to safely extract and display GitHub Actions secrets from a repository. It works by creating an orphan branch, uploading a workflow that exposes the requested secrets as an artifact, and then downloading and streaming the artifact contents to stdout. The tool cleans up after itself, leaving no trace in the main branch.

## Features
- Authenticates using your GitHub token or `gh` CLI.
- Creates an orphan branch and uploads a workflow using Go templates.
- Triggers the workflow and waits for completion.
- Downloads the workflow artifact and streams its contents to stdout.
- Cleans up the workflow run and orphan branch.
- Supports private repositories.

## Usage

```sh
gsd -R org/repo --secret SECRET_NAME [--secret ANOTHER_SECRET] [--github-token YOUR_TOKEN]
```

- `-R org/repo`: The GitHub repository (required)
- `--secret`: Name of the secret to display (can be specified multiple times)
- `--github-token`: GitHub token (or set `GITHUB_TOKEN` env)

Example:

```sh
gsd -R myorg/myrepo --secret FOO --secret BAR
```

## Requirements
- Go 1.18+
- A GitHub personal access token

## Security
- The tool does not pollute the main branch or leave workflow files behind.
- All secrets are streamed to stdout and not written to disk by default.

## Development
- See `cmd/` for main logic.
- Workflow template is embedded as a Go string.
- Polling intervals and attempts are configurable in `artifact.go`.

## License
CC BY-NC 4.0 (Creative Commons Attribution-NonCommercial 4.0 International)

You may not use the material for commercial purposes. See LICENSE for details.
