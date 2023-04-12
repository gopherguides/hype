package hype

type THead struct {
	*Element
}

func (th *THead) String() string {
	if th == nil || th.Element == nil {
		return ""
	}

	if th.IsEmptyNode() {
		return ""
	}

	return th.Element.String()
}

func (th *THead) IsEmptyNode() bool {
	if th == nil {
		return true
	}

	kids := th.Children()

	return len(kids) == 0
}

func NewTHead(el *Element) (*THead, error) {
	if el == nil {
		return nil, ErrIsNil("thead")
	}

	th := &THead{
		Element: el,
	}

	return th, nil
}

func NewTHeadNodes(p *Parser, el *Element) (Nodes, error) {
	th, err := NewTHead(el)
	if err != nil {
		return nil, err
	}

	return Nodes{th}, nil
}
