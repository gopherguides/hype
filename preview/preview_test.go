package preview

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	r := require.New(t)

	cfg := DefaultConfig()

	r.Equal("hype.md", cfg.File)
	r.Equal(3000, cfg.Port)
	r.Equal([]string{"."}, cfg.WatchDirs)
	r.Equal(300*time.Millisecond, cfg.DebounceDelay)
	r.Equal("github", cfg.Theme)
	r.NotEmpty(cfg.Extensions)
	r.NotEmpty(cfg.ExcludeGlobs)
}

func TestWrapHTML(t *testing.T) {
	r := require.New(t)

	content := "<h1>Test</h1><p>Hello World</p>"

	html := wrapHTML(content, "github")

	r.Contains(html, "<!DOCTYPE html>")
	r.Contains(html, "<title>Hype Preview</title>")
	r.Contains(html, content)
	r.Contains(html, "/_livereload")
	r.Contains(html, "WebSocket")
}

func TestWrapHTML_DarkTheme(t *testing.T) {
	r := require.New(t)

	content := "<h1>Test</h1>"

	html := wrapHTML(content, "github-dark")

	r.Contains(html, "background: #0d1117")
	r.Contains(html, "color: #c9d1d9")
}

func TestGetThemeCSS(t *testing.T) {
	tests := []struct {
		name     string
		theme    string
		contains []string
	}{
		{
			name:     "github theme",
			theme:    "github",
			contains: []string{"background: #fff", "color: #24292f"},
		},
		{
			name:     "github-dark theme",
			theme:    "github-dark",
			contains: []string{"background: #0d1117", "color: #c9d1d9"},
		},
		{
			name:     "unknown theme defaults to github",
			theme:    "unknown",
			contains: []string{"background: #fff", "color: #24292f"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			css := getThemeCSS(tc.theme)
			for _, want := range tc.contains {
				r.Contains(css, want)
			}
		})
	}
}

func TestServer_handlePreview(t *testing.T) {
	r := require.New(t)

	cfg := DefaultConfig()
	srv := New(cfg, nil)
	srv.currentHTML = "<html><body>Test Content</body></html>"

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	srv.handlePreview(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	r.Equal(http.StatusOK, resp.StatusCode)
	r.Equal("text/html; charset=utf-8", resp.Header.Get("Content-Type"))
	r.Contains(string(body), "Test Content")
}

func TestServer_shouldWatch(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		extensions []string
		includes   []string
		excludes   []string
		want       bool
	}{
		{
			name:       "matches extension",
			path:       "/project/file.md",
			extensions: []string{"md", "go"},
			want:       true,
		},
		{
			name:       "no match extension",
			path:       "/project/file.txt",
			extensions: []string{"md", "go"},
			want:       false,
		},
		{
			name:     "excluded by glob",
			path:     "/project/vendor/file.go",
			excludes: []string{"**/vendor/**"},
			want:     false,
		},
		{
			name:     "included by glob",
			path:     "/project/src/main.go",
			includes: []string{"**/src/**"},
			want:     true,
		},
		{
			name:     "not included by glob",
			path:     "/project/other/main.go",
			includes: []string{"**/src/**"},
			want:     false,
		},
		{
			name: "all extensions when none specified and no includes",
			path: "/project/file.xyz",
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			cfg := Config{
				Extensions:   tc.extensions,
				IncludeGlobs: tc.includes,
				ExcludeGlobs: tc.excludes,
			}
			srv := New(cfg, nil)

			got := srv.shouldWatch(tc.path, "/project")
			r.Equal(tc.want, got)
		})
	}
}

func TestServer_Build(t *testing.T) {
	r := require.New(t)

	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "test.md")
	err := os.WriteFile(mdFile, []byte("# Hello World\n\nThis is a test."), 0644)
	r.NoError(err)

	cfg := DefaultConfig()
	cfg.File = "test.md"

	srv := New(cfg, nil)

	ctx := context.Background()
	err = srv.build(ctx, tmpDir)
	r.NoError(err)

	r.Contains(srv.currentHTML, "Hello World")
	r.Contains(srv.currentHTML, "This is a test")
	r.Contains(srv.currentHTML, "<!DOCTYPE html>")
}

func TestServer_SetOutput(t *testing.T) {
	r := require.New(t)

	cfg := DefaultConfig()
	srv := New(cfg, nil)

	var stdoutCalled, stderrCalled bool

	srv.SetOutput(
		func(format string, args ...any) { stdoutCalled = true },
		func(format string, args ...any) { stderrCalled = true },
	)

	srv.stdout("test")
	srv.stderr("test")

	r.True(stdoutCalled)
	r.True(stderrCalled)
}

func TestStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: "",
		},
		{
			name:     "single item",
			input:    []string{"a"},
			expected: "a",
		},
		{
			name:     "multiple items",
			input:    []string{"a", "b", "c"},
			expected: "a,b,c",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			sl := stringSlice(tc.input)
			r.Equal(tc.expected, sl.String())
		})
	}
}

func TestDefaultExtensions(t *testing.T) {
	r := require.New(t)

	r.Contains(defaultExtensions, "md")
	r.Contains(defaultExtensions, "html")
	r.Contains(defaultExtensions, "go")
	r.Contains(defaultExtensions, "css")
	r.Contains(defaultExtensions, "png")
	r.Contains(defaultExtensions, "jpg")
}

func TestDefaultExcludes(t *testing.T) {
	r := require.New(t)

	r.Contains(defaultExcludes, "**/.git/**")
	r.Contains(defaultExcludes, "**/node_modules/**")
	r.Contains(defaultExcludes, "**/vendor/**")
}

type stringSlice []string

func (s stringSlice) String() string {
	return strings.Join(s, ",")
}
