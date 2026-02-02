package hype

import "strings"

// NewCodeNodes implements the ParseElementFn type
func NewCodeNodes(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	ats := el.Attrs()

	// Check for special code types first (they may be multi-line)
	if _, ok := ats.Get("src"); ok {
		return NewSourceCodeNodes(p, el)
	}

	lang := Language(ats, "")
	if lang == "mermaid" {
		return NewMermaidNodes(p, el)
	}

	// No attributes: check if single-line (inline) or multi-line (fenced)
	// Per CommonMark spec, fenced code blocks are block-level elements
	if ats.Len() == 0 {
		content := el.Nodes.String()
		if strings.Contains(content, "\n") {
			return NewFencedCodeNodes(p, el)
		}
		return NewInlineCodeNodes(p, el)
	}

	// Has language attribute but not mermaid → fenced code
	return NewFencedCodeNodes(p, el)
}
