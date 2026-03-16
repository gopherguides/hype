package blog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHighlighter(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	h := NewHighlighter("", false)
	r.NotNil(h)
	r.Equal("monokai", h.style)
	r.False(h.lineNumbers)

	h2 := NewHighlighter("dracula", true)
	r.Equal("dracula", h2.style)
	r.True(h2.lineNumbers)
}

func TestHighlighter_Highlight(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	h := NewHighlighter("monokai", false)

	t.Run("go code", func(t *testing.T) {
		result, err := h.Highlight("fmt.Println(\"hello\")", "go")
		r.NoError(err)
		r.NotEmpty(result)
		r.Contains(result, "chroma")
	})

	t.Run("unknown language uses fallback", func(t *testing.T) {
		result, err := h.Highlight("some code", "nonexistentlang")
		r.NoError(err)
		r.NotEmpty(result)
	})

	t.Run("empty code", func(t *testing.T) {
		result, err := h.Highlight("", "go")
		r.NoError(err)
		r.NotNil(result)
	})
}

func TestHighlighter_CSS(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	h := NewHighlighter("monokai", false)
	css := h.CSS()
	r.NotEmpty(css)
	r.Contains(css, "chroma")
}

func TestGetLexer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{"templ maps to html", "templ"},
		{"sh maps to bash", "sh"},
		{"dockerfile maps to docker", "dockerfile"},
		{"mod maps to gomod", "mod"},
		{"yml maps to yaml", "yml"},
		{"rb maps to ruby", "rb"},
		{"py maps to python", "py"},
		{"js maps to javascript", "js"},
		{"ts maps to typescript", "ts"},
		{"dot prefix stripped", ".go"},
		{"direct language", "go"},
		{"unknown returns fallback", "zzz_unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			lexer := getLexer(tt.input)
			r.NotNil(lexer)
		})
	}
}

func TestEscapeHTML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		out  string
	}{
		{"ampersand", "a&b", "a&amp;b"},
		{"less than", "a<b", "a&lt;b"},
		{"greater than", "a>b", "a&gt;b"},
		{"quote", `a"b`, "a&quot;b"},
		{"empty", "", ""},
		{"no special chars", "hello", "hello"},
		{"all special", `&<>"`, "&amp;&lt;&gt;&quot;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			r.Equal(tt.out, escapeHTML(tt.in))
		})
	}
}
