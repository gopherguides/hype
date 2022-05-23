package hype

type InlineCode struct {
	*Element
}

func NewInlineCode(el *Element) (*InlineCode, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	code := &InlineCode{
		Element: el,
	}

	return code, nil
}

func NewInlineCodeNodes(p *Parser, el *Element) (Nodes, error) {
	code, err := NewInlineCode(el)
	if err != nil {
		return nil, err
	}

	return Nodes{code}, nil
}
