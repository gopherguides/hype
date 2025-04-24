package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanPath(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "../",
			expected: filepath.Join(cwd, "../"),
		},
		{
			input:    "./../../readme.md",
			expected: filepath.Join(cwd, "./../../readme.md"),
		},
		{
			input:    "/usr/bin/readme.md",
			expected: "/usr/bin/readme.md",
		},
		{
			input:    "/",
			expected: "/",
		},
		{
			input:    "/../",
			expected: "/",
		},
		{
			input:    "../nonexistent/", // <-- filepath sanitizes the last '/' aways unless it's root path
			expected: filepath.Join(cwd, "../", "nonexistent/"),
		},
	}
	for _, tt := range tests {
		act, err := cleanPath(tt.input)
		fmt.Printf("input: %s\ncleaned: %s\n\n", tt.input, act)
		r.NoError(err)
		r.Equal(tt.expected, act)
	}
}

func TestBadOutputFilePaths(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	paths := []string{
		"../",
		cwd,
		"/",
		"./",
		"./nonexistent/readme.md",
		"/nonexistent/readme.md",
	}
	for _, p := range paths {
		fullPath, err := cleanPath(p)
		r.NoError(err)

		err = dirExists(fullPath)
		r.Error(err)
	}
}
