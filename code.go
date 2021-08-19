package hype

import (
	"fmt"
)

type Code interface {
	Tag
	Lang() string
}

func (p *Parser) NewCode(node *Node) (Code, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("code node can not be nil")
	}

	if node.Data != "code" {
		return nil, fmt.Errorf("node is not code %v", node.Data)
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
