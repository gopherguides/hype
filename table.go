package hype

type Table struct {
	*Element
}

func NewTable(el *Element) (*Table, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	h := &Table{
		Element: el,
	}

	return h, nil
}

func NewTableNodes(p *Parser, el *Element) (Nodes, error) {
	panic(el)
	h, err := NewTable(el)
	if err != nil {
		return nil, err
	}

	return Nodes{h}, nil
}
