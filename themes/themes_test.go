package themes

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListThemes(t *testing.T) {
	r := require.New(t)

	themes := ListThemes()
	r.Len(themes, 7)
	r.Contains(themes, "github")
	r.Contains(themes, "github-dark")
	r.Contains(themes, "solarized-light")
	r.Contains(themes, "solarized-dark")
	r.Contains(themes, "swiss")
	r.Contains(themes, "air")
	r.Contains(themes, "retro")
}

func TestListThemes_ReturnsCopy(t *testing.T) {
	r := require.New(t)

	themes1 := ListThemes()
	themes1[0] = "modified"

	themes2 := ListThemes()
	r.NotEqual("modified", themes2[0])
}

func TestIsBuiltinTheme(t *testing.T) {
	tests := []struct {
		name     string
		theme    string
		expected bool
	}{
		{"github is builtin", "github", true},
		{"github-dark is builtin", "github-dark", true},
		{"solarized-light is builtin", "solarized-light", true},
		{"solarized-dark is builtin", "solarized-dark", true},
		{"swiss is builtin", "swiss", true},
		{"air is builtin", "air", true},
		{"retro is builtin", "retro", true},
		{"unknown is not builtin", "unknown", false},
		{"empty is not builtin", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			r.Equal(tt.expected, IsBuiltinTheme(tt.theme))
		})
	}
}

func TestGetCSS(t *testing.T) {
	r := require.New(t)

	css, err := GetCSS("github")
	r.NoError(err)
	r.NotEmpty(css)
	r.Contains(css, ".markdown-body")
}

func TestGetCSS_AllThemes(t *testing.T) {
	for _, theme := range ListThemes() {
		t.Run(theme, func(t *testing.T) {
			r := require.New(t)

			css, err := GetCSS(theme)
			r.NoError(err)
			r.NotEmpty(css)
			r.Contains(css, ".markdown-body")
		})
	}
}

func TestGetCSS_UnknownTheme(t *testing.T) {
	r := require.New(t)

	_, err := GetCSS("unknown-theme")
	r.Error(err)
	r.Contains(err.Error(), "unknown theme")
}

func TestLoadCustomCSS(t *testing.T) {
	r := require.New(t)

	tmpDir := t.TempDir()
	cssFile := filepath.Join(tmpDir, "custom.css")
	customCSS := ".custom { color: red; }"
	err := os.WriteFile(cssFile, []byte(customCSS), 0644)
	r.NoError(err)

	css, err := LoadCustomCSS(cssFile)
	r.NoError(err)
	r.Equal(customCSS, css)
}

func TestLoadCustomCSS_FileNotFound(t *testing.T) {
	r := require.New(t)

	_, err := LoadCustomCSS("/nonexistent/path/custom.css")
	r.Error(err)
	r.Contains(err.Error(), "failed to read custom CSS file")
}

func TestRender(t *testing.T) {
	r := require.New(t)

	var buf bytes.Buffer
	data := RenderData{
		Title: "Test Document",
		CSS:   template.CSS(".test { color: blue; }"),
		Body:  template.HTML("<h1>Hello World</h1>"),
	}

	err := Render(&buf, data)
	r.NoError(err)

	output := buf.String()
	r.Contains(output, "<!DOCTYPE html>")
	r.Contains(output, "<title>Test Document</title>")
	r.Contains(output, ".test { color: blue; }")
	r.Contains(output, "<h1>Hello World</h1>")
	r.Contains(output, "<article class=\"markdown-body\">")
}

func TestRender_EscapesTitle(t *testing.T) {
	r := require.New(t)

	var buf bytes.Buffer
	data := RenderData{
		Title: "<script>alert('xss')</script>",
		CSS:   template.CSS(""),
		Body:  template.HTML("<p>Content</p>"),
	}

	err := Render(&buf, data)
	r.NoError(err)

	output := buf.String()
	r.NotContains(output, "<script>alert('xss')</script>")
	r.Contains(output, "&lt;script&gt;")
}

func TestDefaultTheme(t *testing.T) {
	r := require.New(t)
	r.Equal("github", DefaultTheme)
	r.True(IsBuiltinTheme(DefaultTheme))
}

func TestDocumentTemplate_IsValid(t *testing.T) {
	r := require.New(t)
	r.NotEmpty(documentTemplate)
	r.True(strings.Contains(documentTemplate, "{{.Title}}"))
	r.True(strings.Contains(documentTemplate, "{{.CSS}}"))
	r.True(strings.Contains(documentTemplate, "{{.Body}}"))
}
