package hype

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"

	mermaidcmd "github.com/AlexanderGrooff/mermaid-ascii/cmd"
)

// Mermaid is a tag that renders Mermaid diagrams as ASCII art.
// It processes fenced code blocks with the "mermaid" language identifier
// and converts them to ASCII art using the mermaid-ascii library.
type Mermaid struct {
	*Element

	// Source is the original Mermaid diagram source
	Source string

	// Rendered is the ASCII art output
	Rendered string
}

func (m *Mermaid) MarshalJSON() ([]byte, error) {
	if m == nil {
		return nil, ErrIsNil("mermaid")
	}

	m.RLock()
	defer m.RUnlock()

	mm, err := m.JSONMap()
	if err != nil {
		return nil, err
	}

	mm["type"] = toType(m)

	if len(m.Source) > 0 {
		mm["source"] = m.Source
	}

	if len(m.Rendered) > 0 {
		mm["rendered"] = m.Rendered
	}

	return json.MarshalIndent(mm, "", "  ")
}

// MD returns the markdown representation of the rendered Mermaid diagram.
func (m *Mermaid) MD() string {
	if m == nil {
		return ""
	}

	m.RLock()
	defer m.RUnlock()

	if len(m.Rendered) == 0 {
		return ""
	}

	bb := &bytes.Buffer{}
	fmt.Fprint(bb, "```\n")
	fmt.Fprint(bb, m.Rendered)
	fmt.Fprint(bb, "\n```")

	return bb.String()
}

// String returns the HTML representation of the rendered Mermaid diagram.
func (m *Mermaid) String() string {
	if m == nil {
		return ""
	}

	m.RLock()
	defer m.RUnlock()

	if len(m.Rendered) == 0 {
		return ""
	}

	return fmt.Sprintf("<pre><code class=\"language-plain\">%s</code></pre>",
		html.EscapeString(m.Rendered))
}

// Execute renders the Mermaid diagram to ASCII art.
func (m *Mermaid) Execute(ctx context.Context, doc *Document) error {
	if m == nil {
		return ErrIsNil("mermaid")
	}

	if m.Element == nil {
		return ErrIsNil("element")
	}

	if doc == nil {
		return ErrIsNil("document")
	}

	m.Lock()
	defer m.Unlock()

	// Render the diagram using default config (nil)
	// The mermaid-ascii library uses sensible defaults:
	// - Unicode box-drawing characters
	// - 5px horizontal/vertical padding between nodes
	// - LR (left-to-right) direction
	output, err := mermaidcmd.RenderDiagram(m.Source, nil)
	if err != nil {
		return m.WrapErr(fmt.Errorf("failed to render mermaid diagram: %w", err))
	}

	m.Rendered = output

	return nil
}

// NewMermaid creates a new Mermaid element from the given element.
func NewMermaid(el *Element) (*Mermaid, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	m := &Mermaid{
		Element: el,
	}

	// Extract the mermaid source from the element's children (text content)
	m.Source = html.UnescapeString(el.Children().String())

	if len(m.Source) == 0 {
		return nil, m.WrapErr(fmt.Errorf("mermaid diagram source is empty"))
	}

	return m, nil
}

// NewMermaidNodes creates a new Mermaid node from a fenced code block.
func NewMermaidNodes(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	m, err := NewMermaid(el)
	if err != nil {
		return nil, err
	}

	return Nodes{m}, nil
}
