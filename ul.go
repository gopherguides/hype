package hype

type UL struct {
	*Element
}

func (ol *UL) MD() string {
	if ol == nil || ol.Element == nil {
		return ""
	}

	return ol.Children().MD()
}

func NewULNodes(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, nil
	}

	ol := &UL{
		Element: el,
	}

	return Nodes{ol}, nil
}
