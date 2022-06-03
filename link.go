package hype

import "strings"

type Link struct {
	*Element
}

func (l *Link) Href() (string, error) {
	if l == nil {
		return "", ErrIsNil("link")
	}

	return l.ValidAttr("href")
}

func NewLink(el *Element) (*Link, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	l := &Link{
		Element: el,
	}

	h, err := l.Href()
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(h, "http") {
		return l, nil
	}

	if _, ok := l.Get("target"); !ok {
		if err := l.Set("target", "_blank"); err != nil {
			return nil, err
		}
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
