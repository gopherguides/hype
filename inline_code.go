package hype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
)

type InlineCode struct {
	*Element
}

func (code *InlineCode) MarshalJSON() ([]byte, error) {
	if code == nil {
		return nil, ErrIsNil("inline code")
	}

	code.RLock()
	defer code.RUnlock()

	m, err := code.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = fmt.Sprintf("%T", code)

	return json.MarshalIndent(m, "", "  ")
}

func (code *InlineCode) StartTag() string {
	if code == nil || code.Element == nil {
		return "<code>"
	}

	return code.Element.StartTag()
}

func (code *InlineCode) EndTag() string {
	if code == nil || code.Element == nil {
		return "</code>"
	}

	return code.Element.EndTag()
}

func (code *InlineCode) String() string {
	if code == nil || code.Element == nil {
		return "<code></code>"
	}

	bb := &bytes.Buffer{}

	fmt.Fprint(bb, code.StartTag())

	body := code.Nodes.String()
	body = html.EscapeString(body)
	fmt.Fprint(bb, body)
	fmt.Fprint(bb, code.EndTag())

	return bb.String()
}

func (code *InlineCode) MD() string {
	if code == nil || code.Element == nil {
		return ""
	}

	return fmt.Sprintf("`%s`", code.Nodes.String())
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
