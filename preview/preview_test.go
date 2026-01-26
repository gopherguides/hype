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

	"github.com/gopherguides/hype/themes"
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

func TestServer_wrapHTML(t *testing.T) {
	r := require.New(t)

	cfg := DefaultConfig()
	cfg.Theme = "github"
	srv := New(cfg, nil)

	content := "<h1>Test</h1><p>Hello World</p>"

	html, err := srv.wrapHTML(content)
	r.NoError(err)

	r.Contains(html, "<!DOCTYPE html>")
	r.Contains(html, "<title>Hype Preview</title>")
	r.Contains(html, content)
	r.Contains(html, "/_livereload")
	r.Contains(html, "WebSocket")
}

func TestServer_wrapHTML_DifferentTheme(t *testing.T) {
	r := require.New(t)

	cfg := DefaultConfig()
	cfg.Theme = "github-dark"
	srv := New(cfg, nil)

	content := "<h1>Test</h1>"

	html, err := srv.wrapHTML(content)
	r.NoError(err)

	r.Contains(html, "<!DOCTYPE html>")
	r.Contains(html, content)
}

func TestServer_wrapHTML_CustomCSS(t *testing.T) {
	r := require.New(t)

	tmpDir := t.TempDir()
	cssFile := filepath.Join(tmpDir, "custom.css")
	err := os.WriteFile(cssFile, []byte("body { color: red; }"), 0644)
	r.NoError(err)

	cfg := DefaultConfig()
	cfg.CustomCSS = cssFile
	srv := New(cfg, nil)

	content := "<h1>Test</h1>"

	html, err := srv.wrapHTML(content)
	r.NoError(err)

	r.Contains(html, "body { color: red; }")
	r.Contains(html, content)
}

func TestServer_getCSS(t *testing.T) {
	tests := []struct {
		name      string
		theme     string
		customCSS string
		wantErr   bool
	}{
		{
			name:    "builtin github theme",
			theme:   "github",
			wantErr: false,
		},
		{
			name:    "builtin github-dark theme",
			theme:   "github-dark",
			wantErr: false,
		},
		{
			name:    "unknown theme returns error",
			theme:   "unknown",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			cfg := Config{Theme: tc.theme, CustomCSS: tc.customCSS}
			srv := New(cfg, nil)

			css, err := srv.getCSS()
			if tc.wantErr {
				r.Error(err)
			} else {
				r.NoError(err)
				r.NotEmpty(css)
			}
		})
	}
}

func TestThemesIntegration(t *testing.T) {
	r := require.New(t)

	availableThemes := themes.ListThemes()
	r.NotEmpty(availableThemes)
	r.Contains(availableThemes, "github")
	r.Contains(availableThemes, "github-dark")

	r.True(themes.IsBuiltinTheme("github"))
	r.False(themes.IsBuiltinTheme("nonexistent"))
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

	r.Contains(defaultExcludes, ".git")
	r.Contains(defaultExcludes, "node_modules")
	r.Contains(defaultExcludes, "vendor")
	r.Contains(defaultExcludes, "**/.git/**")
	r.Contains(defaultExcludes, "**/node_modules/**")
	r.Contains(defaultExcludes, "**/vendor/**")
}

func TestServer_handleRequest_StaticFile(t *testing.T) {
	r := require.New(t)

	tmpDir := t.TempDir()
	imgFile := filepath.Join(tmpDir, "test.png")
	err := os.WriteFile(imgFile, []byte("PNG content"), 0644)
	r.NoError(err)

	cfg := DefaultConfig()
	srv := New(cfg, nil)
	srv.pwd = tmpDir
	srv.currentHTML = "<html><body>Preview</body></html>"

	req := httptest.NewRequest("GET", "/test.png", nil)
	w := httptest.NewRecorder()

	srv.handleRequest(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	r.Equal(http.StatusOK, resp.StatusCode)
	r.Equal("PNG content", string(body))
}

func TestServer_handleRequest_FallbackToPreview(t *testing.T) {
	r := require.New(t)

	tmpDir := t.TempDir()

	cfg := DefaultConfig()
	srv := New(cfg, nil)
	srv.pwd = tmpDir
	srv.currentHTML = "<html><body>Preview Content</body></html>"

	req := httptest.NewRequest("GET", "/nonexistent.png", nil)
	w := httptest.NewRecorder()

	srv.handleRequest(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	r.Equal(http.StatusOK, resp.StatusCode)
	r.Contains(string(body), "Preview Content")
}

func TestServer_handleRequest_RootPath(t *testing.T) {
	r := require.New(t)

	tmpDir := t.TempDir()

	cfg := DefaultConfig()
	srv := New(cfg, nil)
	srv.pwd = tmpDir
	srv.currentHTML = "<html><body>Root Preview</body></html>"

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	srv.handleRequest(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	r.Equal(http.StatusOK, resp.StatusCode)
	r.Contains(string(body), "Root Preview")
}

func TestServer_handleRequest_PathTraversal(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "path traversal with ..",
			path:       "/../../../etc/passwd",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "encoded path traversal",
			path:       "/%2e%2e/etc/passwd",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "double dot in middle",
			path:       "/foo/../../../etc/passwd",
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			tmpDir := t.TempDir()

			cfg := DefaultConfig()
			srv := New(cfg, nil)
			srv.pwd = tmpDir
			srv.currentHTML = "<html><body>Preview</body></html>"

			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()

			srv.handleRequest(w, req)

			resp := w.Result()
			r.Equal(tc.wantStatus, resp.StatusCode)
		})
	}
}

func TestServer_Build_WithTimeout(t *testing.T) {
	r := require.New(t)

	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "test.md")
	err := os.WriteFile(mdFile, []byte("# Hello\n\nSimple content."), 0644)
	r.NoError(err)

	cfg := DefaultConfig()
	cfg.File = "test.md"
	cfg.Timeout = 5 * time.Second

	srv := New(cfg, nil)

	ctx := context.Background()
	err = srv.build(ctx, tmpDir)
	r.NoError(err)

	r.Contains(srv.currentHTML, "Hello")
}

type stringSlice []string

func (s stringSlice) String() string {
	return strings.Join(s, ",")
}
