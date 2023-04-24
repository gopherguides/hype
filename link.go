package hype

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Link struct {
	*Element
}

func (l *Link) MarshalJSON() ([]byte, error) {
	if l == nil {
		return nil, ErrIsNil("link")
	}

	l.RLock()
	defer l.RUnlock()

	m, err := l.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = fmt.Sprintf("%T", l)

	h, err := l.Href()
	if err != nil {
		return nil, err
	}

	m["url"] = h

	return json.Marshal(m)
}

func (l *Link) Href() (string, error) {
	if l == nil {
		return "", ErrIsNil("link")
	}

	return l.ValidAttr("href")
}

func (l *Link) MD() string {
	if l == nil {
		return ""
	}

	h, err := l.Href()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("[%s](%s)", l.Children().MD(), h)
}

func NewLink(el *Element) (*Link, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	l := &Link{
		Element: el,
	}

	h, ok := l.Get("href")
	if !ok {
		return l, nil
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
