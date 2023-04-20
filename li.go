package hype

import "bytes"

type LI struct {
	Type string
	*Element
}

func (li *LI) MD() string {
	if li == nil {
		return ""
	}

	bb := &bytes.Buffer{}

	switch li.Type {
	case "ol":
		bb.WriteString("1. ")
	default:
		bb.WriteString("* ")
	}

	bb.WriteString(li.Children().MD())

	return bb.String()
}

func NewLINodes(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, nil
	}

	li := &LI{
		Element: el,
		Type:    "ul",
	}

	if par, ok := el.Parent.(AtomableNode); ok {
		li.Type = par.Atom().String()
	}

	return Nodes{li}, nil
}