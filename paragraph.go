package hype

func NewParagraphNodes(p *Parser, el *Element) (Nodes, error) {
	var nodes Nodes

	if el == nil {
		return nil, ErrIsNil("el")
	}

	if IsEmptyNode(el) {
		return nodes, nil
	}

	nodes = append(nodes, el)

	return nodes, nil
}
