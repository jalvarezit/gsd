package main

import (
	"bytes"
	"os"
	"strings"
	"text/template"
)

var workflowTmpl = `name: github-secret-display
on:
  push:
    branches:
      - github-secret-display
permissions:
  contents: read
  actions: write
  id-token: write
  issues: none
  checks: none
  deployments: none
  discussions: none
  packages: none
  pull-requests: none
  repository-projects: none
  security-events: none
  statuses: none
jobs:
  display-secrets:
    runs-on: ubuntu-latest
    steps:
      - name: Write secrets to artifact
        run: |
          echo "{{ .EchoLines }}" > secrets.txt
      - name: Upload secrets artifact
        uses: actions/upload-artifact@v4
        with:
          name: secrets
          path: secrets.txt
`

func renderWorkflowTemplate(secrets []string, outPath string) error {
	echoLines := make([]string, len(secrets))
	for i, s := range secrets {
		echoLines[i] = s + `=${{ secrets.` + s + ` }}`
	}
	echoStr := strings.Join(echoLines, `\n`)
	tmpl, err := template.New("workflow").Parse(workflowTmpl)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{"EchoLines": echoStr}); err != nil {
		return err
	}
	return os.WriteFile(outPath, buf.Bytes(), 0644)
}
