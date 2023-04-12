package hype

type Paragraph struct {
	*Element
}

func (p *Paragraph) IsEmptyNode() bool {
	if p == nil {
		return true
	}

	kids := p.Children()
	if len(kids) == 0 {
		return true
	}

	return IsEmptyNode(kids)
}

func (p *Paragraph) MD() string {
	if p == nil {
		return ""
	}

	return p.Children().MD()
}

func NewParagraphNodes(p *Parser, el *Element) (Nodes, error) {
	var nodes Nodes

	if el == nil {
		return nil, ErrIsNil("el")
	}

	if IsEmptyNode(el) {
		return nodes, nil
	}

	nodes = append(nodes, &Paragraph{
		Element: el,
	})

	return nodes, nil
}
