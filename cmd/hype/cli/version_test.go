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
		name string
		info VersionInfo
		want string
	}{
		{
			name: "dev defaults",
			info: VersionInfo{Version: "dev", Commit: "none", Date: "unknown"},
			want: "hype version dev (commit: none, built: unknown)",
		},
		{
			name: "release version",
			info: VersionInfo{Version: "v0.5.0", Commit: "abc1234", Date: "2025-01-31T10:00:00Z"},
			want: "hype version v0.5.0 (commit: abc1234, built: 2025-01-31T10:00:00Z)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &Version{Info: tt.info}

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

func Test_Version_Cmd(t *testing.T) {
	app := New(".", VersionInfo{Version: "v1.0.0", Commit: "abc", Date: "now"})

	ver, ok := app.Commands["version"].(*Version)
	if !ok {
		t.Fatal("version command not found")
	}

	if ver.Cmd.Name != "version" {
		t.Errorf("expected name 'version', got %q", ver.Cmd.Name)
	}

	if len(ver.Cmd.Aliases) != 1 || ver.Cmd.Aliases[0] != "v" {
		t.Errorf("expected aliases ['v'], got %v", ver.Cmd.Aliases)
	}

	if ver.Info.Version != "v1.0.0" {
		t.Errorf("expected version 'v1.0.0', got %q", ver.Info.Version)
	}
}
