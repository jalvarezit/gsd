package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

// streamAllFilesFromZipToStdout streams the content of all files in the zip (zipData) to stdout, with no header or extra newlines.
func streamAllFilesFromZipToStdout(zipData []byte) error {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	for _, f := range reader.File {
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open %s in zip: %w", f.Name, err)
		}
		defer rc.Close()
		_, err = io.Copy(os.Stdout, rc)
		if err != nil {
			return fmt.Errorf("failed to stream %s: %w", f.Name, err)
		}
	}
	return nil
}
