package hype

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Figcaption struct {
	*Element
}

func (fc *Figcaption) MarshalJSON() ([]byte, error) {
	if fc == nil {
		return nil, ErrIsNil("figcaption")
	}

	fc.RLock()
	defer fc.RUnlock()

	m, err := fc.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = fmt.Sprintf("%T", fc)

	return json.MarshalIndent(m, "", "  ")
}

func (fc *Figcaption) MD() string {
	if fc == nil {
		return ""
	}

	bb := &strings.Builder{}
	bb.WriteString(fc.StartTag())
	bb.WriteString(fc.Nodes.MD())
	bb.WriteString(fc.EndTag())

	return bb.String()
}

func NewFigcaption(el *Element) (*Figcaption, error) {
	if el == nil {
		return nil, el.WrapErr(ErrIsNil("element"))
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
