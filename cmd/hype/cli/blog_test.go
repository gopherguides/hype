package cli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Blog_VersionInfo(t *testing.T) {
	r := require.New(t)

	info := VersionInfo{Version: "v1.0.0", Commit: "abc1234", Date: "2025-01-31"}
	cmd := &Blog{Info: info}

	r.Equal("v1.0.0", cmd.Info.Version)
	r.Equal("abc1234", cmd.Info.Commit)
	r.Equal("2025-01-31", cmd.Info.Date)
	r.Equal("hype version v1.0.0 (commit: abc1234, built: 2025-01-31)", cmd.Info.String())
}

func Test_Blog_VersionInfo_FromNew(t *testing.T) {
	r := require.New(t)

	info := VersionInfo{Version: "v2.0.0", Commit: "def5678", Date: "2025-02-01"}
	app := New(".", info)

	bl, ok := app.Commands["blog"].(*Blog)
	r.True(ok, "blog command should exist")
	r.Equal("v2.0.0", bl.Info.Version)
	r.Equal("def5678", bl.Info.Commit)
	r.Equal("2025-02-01", bl.Info.Date)
}

func Test_Blog_VersionInfo_String(t *testing.T) {
	info := VersionInfo{Version: "v1.0.0", Commit: "abc", Date: "now"}
	cmd := &Blog{Info: info}

	if cmd.Info.Version != "v1.0.0" {
		t.Errorf("got %q, want %q", cmd.Info.Version, "v1.0.0")
	}
	if !strings.Contains(cmd.Info.String(), "hype version v1.0.0") {
		t.Errorf("String() should contain version info")
	}
}
