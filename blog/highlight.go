package blog

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

type Highlighter struct {
	style       string
	lineNumbers bool
	formatter   *html.Formatter
}

func NewHighlighter(style string, lineNumbers bool) *Highlighter {
	if style == "" {
		style = "monokai"
	}

	opts := []html.Option{
		html.WithClasses(true),
		html.TabWidth(4),
	}
	if lineNumbers {
		opts = append(opts, html.WithLineNumbers(true))
	}

	return &Highlighter{
		style:       style,
		lineNumbers: lineNumbers,
		formatter:   html.New(opts...),
	}
}

func (h *Highlighter) Highlight(code, language string) (string, error) {
	lexer := getLexer(language)
	lexer = chroma.Coalesce(lexer)

	style := styles.Get(h.style)
	if style == nil {
		style = styles.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return escapeHTML(code), nil
	}

	var buf bytes.Buffer
	if err := h.formatter.Format(&buf, style, iterator); err != nil {
		return escapeHTML(code), nil
	}

	return buf.String(), nil
}

func (h *Highlighter) CSS() string {
	style := styles.Get(h.style)
	if style == nil {
		style = styles.Fallback
	}

	var buf bytes.Buffer
	if err := h.formatter.WriteCSS(&buf, style); err != nil {
		return ""
	}
	return buf.String()
}

func getLexer(language string) chroma.Lexer {
	lang := strings.ToLower(strings.TrimPrefix(language, "."))

	switch lang {
	case "templ":
		lang = "html"
	case "sh":
		lang = "bash"
	case "dockerfile":
		lang = "docker"
	case "mod":
		lang = "gomod"
	case "yml":
		lang = "yaml"
	case "rb":
		lang = "ruby"
	case "py":
		lang = "python"
	case "js":
		lang = "javascript"
	case "ts":
		lang = "typescript"
	}

	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	return lexer
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
