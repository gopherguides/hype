package cli

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Export_SubdirectoryFile(t *testing.T) {
	r := require.New(t)

	pwd, err := filepath.Abs("testdata/export/subdir")
	r.NoError(err)

	t.Setenv("MARKED_PATH", filepath.Join(pwd, "dummy.md"))

	outFile := filepath.Join(t.TempDir(), "output.md")

	cmd := &Export{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = cmd.Main(ctx, pwd, []string{"-f", ".hype/module.md", "-format", "markdown", "-o", outFile})

	r.NoError(err, "should be able to resolve includes when file is in subdirectory")

	act, err := os.ReadFile(outFile)
	r.NoError(err)
	r.Contains(string(act), "Main Module")
	r.Contains(string(act), "Included Content")
}
