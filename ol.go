package hype

type OL struct {
	*Element
}

func (ol *OL) MD() string {
	if ol == nil || ol.Element == nil {
		return ""
	}

	return ol.Children().MD()
}

func NewOLNodes(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, nil
	}

	ol := &OL{
		Element: el,
	}

	return Nodes{ol}, nil
}
