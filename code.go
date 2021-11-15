package hype

import (
	"golang.org/x/net/html"
)

type Code interface {
	Tag
	Lang() string
}

func NewCode(node *Node, p *Parser) (Code, error) {
	err := node.Validate(html.ElementNode, AtomValidator("code"))

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
