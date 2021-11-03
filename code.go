package hype

import (
	"golang.org/x/net/html"
)

type Code interface {
	Tag
	Lang() string
}

func (p *Parser) NewCode(node *Node) (Code, error) {
	err := node.Validate(html.ElementNode, AdamValidator("code"))

	if err != nil {
		return nil, err
	}

	ats := node.Attrs()

	if len(ats) == 0 {
		return p.NewInlineCode(node)
	}

	if _, ok := ats["src"]; ok {
		return p.NewSourceCode(node)
	}

	return p.NewFencedCode(node)
}
