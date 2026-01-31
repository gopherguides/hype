package hype

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Mermaid_Execute(t *testing.T) {
	t.Parallel()

	t.Run("basic graph", func(t *testing.T) {
		t.Parallel()
		r := require.New(t)

		p := testParser(t, "testdata/mermaid/graph")

		doc, err := p.ParseFile("hype.md")
		r.NoError(err)

		err = doc.Execute(context.Background())
		r.NoError(err)

		act := doc.String()

		// The rendered output should contain ASCII art
		r.Contains(act, "Start")
		r.Contains(act, "End")
		// Should have box-drawing characters or dashes for arrows
		r.True(strings.Contains(act, "─") || strings.Contains(act, "-"))
	})

	t.Run("sequence diagram", func(t *testing.T) {
		t.Parallel()
		r := require.New(t)

		p := testParser(t, "testdata/mermaid/sequence")

		doc, err := p.ParseFile("hype.md")
		r.NoError(err)

		err = doc.Execute(context.Background())
		r.NoError(err)

		act := doc.String()

		// The rendered output should contain participant names
		r.Contains(act, "Alice")
		r.Contains(act, "Bob")
	})

	t.Run("empty source", func(t *testing.T) {
		t.Parallel()
		r := require.New(t)

		el := NewEl("code", nil)
		el.Nodes = Nodes{}

		_, err := NewMermaid(el)
		r.Error(err)
		r.Contains(err.Error(), "empty")
	})

	t.Run("invalid mermaid syntax", func(t *testing.T) {
		t.Parallel()
		r := require.New(t)

		p := testParser(t, "testdata/mermaid/invalid")

		doc, err := p.ParseFile("hype.md")
		r.NoError(err)

		err = doc.Execute(context.Background())
		r.Error(err)
	})
}

func Test_Mermaid_MD(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/mermaid/graph")

	doc, err := p.ParseFile("hype.md")
	r.NoError(err)

	err = doc.Execute(context.Background())
	r.NoError(err)

	md := doc.MD()

	// Should be wrapped in code fences
	r.Contains(md, "```")
	// Should contain the rendered diagram
	r.Contains(md, "Start")
}

func Test_Mermaid_String(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/mermaid/graph")

	doc, err := p.ParseFile("hype.md")
	r.NoError(err)

	err = doc.Execute(context.Background())
	r.NoError(err)

	html := doc.String()

	// Should be wrapped in pre/code tags
	r.Contains(html, "<pre><code")
	r.Contains(html, "</code></pre>")

	// Should NOT have nested pre tags (bug fix verification)
	r.NotContains(html, "<pre><pre>")
	r.NotContains(html, "</pre></pre>")
}

func Test_Mermaid_MarshalJSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	// Create a Mermaid element directly for testing
	el := NewEl("code", nil)
	el.Nodes = Nodes{Text("graph LR\n    A --> B")}

	m, err := NewMermaid(el)
	r.NoError(err)

	// Manually set rendered output for testing
	m.Rendered = "ASCII diagram output"

	b, err := m.MarshalJSON()
	r.NoError(err)

	json := string(b)
	r.Contains(json, `"type"`)
	r.Contains(json, `"source"`)
	r.Contains(json, `"rendered"`)
}
