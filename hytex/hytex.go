package hytex

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os/exec"
	"strings"

	"github.com/gopherguides/hype"
)

// Convert REQUIRES pandoc to be installed in PATH.
func Convert(ctx context.Context, doc *hype.Document) (fs.FS, error) {
	if doc == nil {
		return nil, fmt.Errorf("doc is nil")
	}

	cmd := exec.CommandContext(ctx, "pandoc", "--from", "html", "--to", "latex")

	stdin := strings.NewReader(doc.String())
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(stderr.String())
		return nil, fmt.Errorf("pandoc failed: %w", err)
	}

	return nil, nil
}
