package hype

import "strings"

type Link struct {
	*Element
}

func (l *Link) Href() (string, bool) {
	if l == nil {
		return "", false
	}

	return l.Get("href")
}

func NewLink(el *Element) (*Link, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	l := &Link{
		Element: el,
	}

	h, ok := l.Href()
	if !ok {
		return nil, ErrAttrNotFound("href")
	}

	if len(h) == 0 {
		return nil, ErrAttrEmpty("href")
	}

	if !strings.HasPrefix(h, "http") {
		return l, nil
	}

	if _, ok := l.Get("target"); !ok {
		l.Set("target", "_blank")
	}

	return l, nil
}

func NewLinkNodes(p *Parser, el *Element) (Nodes, error) {
	l, err := NewLink(el)
	if err != nil {
		return nil, err
	}

	return Nodes{l}, nil
}
