package hype

type Figcaption struct {
	*Element
}

func NewFigcaption(el *Element) (*Figcaption, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	f := &Figcaption{
		Element: el,
	}

	return f, nil
}

func NewFigcaptionNodes(p *Parser, el *Element) (Nodes, error) {
	f, err := NewFigcaption(el)
	if err != nil {
		return nil, err
	}

	return Nodes{f}, nil
}
