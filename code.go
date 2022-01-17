package hype

import (
	"golang.org/x/net/html"
)

// Code represents a type of code block.
type Code interface {
	Tag
	Lang() string
}

// NewCode will return the appropriate code type for the given node.
func NewCode(p *Parser, node *Node) (Code, error) {
	err := node.Validate(p, html.ElementNode, AtomValidator("code"))

	if err != nil {
		return nil, err
	}

	ats := node.Attrs()

	if len(ats) == 0 {
		return NewInlineCode(node)
	}

	if _, ok := ats["src"]; ok {
		return NewSourceCode(p.FS, node, p.snippetRules)
	}

	return NewFencedCode(node)
}
