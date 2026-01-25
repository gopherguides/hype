package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gopherguides/hype/themes"
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

func Test_Export_HTML_DefaultTheme(t *testing.T) {
	r := require.New(t)

	pwd, err := filepath.Abs("testdata/export/html")
	r.NoError(err)

	t.Setenv("MARKED_PATH", filepath.Join(pwd, "dummy.md"))

	outFile := filepath.Join(t.TempDir(), "output.html")

	cmd := &Export{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = cmd.Main(ctx, pwd, []string{"-f", "test.md", "-format", "html", "-o", outFile})

	r.NoError(err)

	act, err := os.ReadFile(outFile)
	r.NoError(err)

	output := string(act)
	r.Contains(output, "<!DOCTYPE html>")
	r.Contains(output, "<title>Test Document</title>")
	r.Contains(output, ".markdown-body")
	r.Contains(output, "<h1")
	r.Contains(output, "Test Document")
}

func Test_Export_HTML_SpecificTheme(t *testing.T) {
	r := require.New(t)

	pwd, err := filepath.Abs("testdata/export/html")
	r.NoError(err)

	t.Setenv("MARKED_PATH", filepath.Join(pwd, "dummy.md"))

	outFile := filepath.Join(t.TempDir(), "output.html")

	cmd := &Export{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = cmd.Main(ctx, pwd, []string{"-f", "test.md", "-format", "html", "-theme", "solarized-dark", "-o", outFile})

	r.NoError(err)

	act, err := os.ReadFile(outFile)
	r.NoError(err)

	output := string(act)
	r.Contains(output, "<!DOCTYPE html>")
	r.Contains(output, "Solarized Dark")
	r.Contains(output, ".markdown-body")
}

func Test_Export_HTML_NoCSS(t *testing.T) {
	r := require.New(t)

	pwd, err := filepath.Abs("testdata/export/html")
	r.NoError(err)

	t.Setenv("MARKED_PATH", filepath.Join(pwd, "dummy.md"))

	outFile := filepath.Join(t.TempDir(), "output.html")

	cmd := &Export{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = cmd.Main(ctx, pwd, []string{"-f", "test.md", "-format", "html", "-no-css", "-o", outFile})

	r.NoError(err)

	act, err := os.ReadFile(outFile)
	r.NoError(err)

	output := string(act)
	r.NotContains(output, "<!DOCTYPE html>")
	r.NotContains(output, "<style>")
	r.Contains(output, "<h1")
}

func Test_Export_HTML_CustomCSS(t *testing.T) {
	r := require.New(t)

	pwd, err := filepath.Abs("testdata/export/html")
	r.NoError(err)

	t.Setenv("MARKED_PATH", filepath.Join(pwd, "dummy.md"))

	tmpDir := t.TempDir()
	cssFile := filepath.Join(tmpDir, "custom.css")
	customCSS := ".custom-class { color: purple; }"
	err = os.WriteFile(cssFile, []byte(customCSS), 0644)
	r.NoError(err)

	outFile := filepath.Join(tmpDir, "output.html")

	cmd := &Export{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = cmd.Main(ctx, pwd, []string{"-f", "test.md", "-format", "html", "-css", cssFile, "-o", outFile})

	r.NoError(err)

	act, err := os.ReadFile(outFile)
	r.NoError(err)

	output := string(act)
	r.Contains(output, "<!DOCTYPE html>")
	r.Contains(output, ".custom-class { color: purple; }")
}

func Test_Export_ListThemes(t *testing.T) {
	r := require.New(t)

	pwd, err := filepath.Abs("testdata/export/html")
	r.NoError(err)

	t.Setenv("MARKED_PATH", filepath.Join(pwd, "dummy.md"))

	cmd := &Export{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	old := os.Stdout
	rf, wf, _ := os.Pipe()
	os.Stdout = wf

	err = cmd.Main(ctx, pwd, []string{"-themes"})

	wf.Close()
	os.Stdout = old

	r.NoError(err)

	buf := make([]byte, 1024)
	n, _ := rf.Read(buf)
	output := string(buf[:n])

	r.Contains(output, "Available themes:")
	for _, theme := range themes.ListThemes() {
		r.True(strings.Contains(output, theme), "should list theme: %s", theme)
	}
	r.Contains(output, "(default)")
}

func Test_Export_HTML_InvalidTheme(t *testing.T) {
	r := require.New(t)

	pwd, err := filepath.Abs("testdata/export/html")
	r.NoError(err)

	t.Setenv("MARKED_PATH", filepath.Join(pwd, "dummy.md"))

	outFile := filepath.Join(t.TempDir(), "output.html")

	cmd := &Export{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = cmd.Main(ctx, pwd, []string{"-f", "test.md", "-format", "html", "-theme", "nonexistent-theme", "-o", outFile})

	r.Error(err)
	r.Contains(err.Error(), "unknown theme")
}

func Test_Export_HTML_InvalidCustomCSS(t *testing.T) {
	r := require.New(t)

	pwd, err := filepath.Abs("testdata/export/html")
	r.NoError(err)

	t.Setenv("MARKED_PATH", filepath.Join(pwd, "dummy.md"))

	outFile := filepath.Join(t.TempDir(), "output.html")

	cmd := &Export{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = cmd.Main(ctx, pwd, []string{"-f", "test.md", "-format", "html", "-css", "/nonexistent/path.css", "-o", outFile})

	r.Error(err)
	r.Contains(err.Error(), "failed to read custom CSS file")
}
