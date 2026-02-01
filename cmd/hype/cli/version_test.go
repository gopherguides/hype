package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func Test_Version_Main(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		version string
		commit  string
		date    string
		want    string
	}{
		{
			name:    "dev defaults",
			version: "dev",
			commit:  "none",
			date:    "unknown",
			want:    "hype version dev (commit: none, built: unknown)",
		},
		{
			name:    "release version",
			version: "v0.5.0",
			commit:  "abc1234",
			date:    "2025-01-31T10:00:00Z",
			want:    "hype version v0.5.0 (commit: abc1234, built: 2025-01-31T10:00:00Z)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewVersion(tt.version, tt.commit, tt.date)

			var stdout bytes.Buffer
			cmd.Out = &stdout

			err := cmd.Main(context.Background(), ".", nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := strings.TrimSpace(stdout.String())
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_Version_Alias(t *testing.T) {
	cmd := NewVersion("dev", "none", "unknown")

	if cmd.Name != "version" {
		t.Errorf("expected name 'version', got %q", cmd.Name)
	}

	if len(cmd.Aliases) != 1 || cmd.Aliases[0] != "v" {
		t.Errorf("expected aliases ['v'], got %v", cmd.Aliases)
	}
}
