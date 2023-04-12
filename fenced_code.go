package hype

import (
	"bytes"
	"fmt"
	"html"
)

type FencedCode struct {
	*Element
}

func (code *FencedCode) StartTag() string {
	if code == nil || code.Element == nil {
		return ""
	}

	return code.Element.StartTag()
}

func (code *FencedCode) EndTag() string {
	if code == nil || code.Element == nil {
		return ""
	}

	return "</code>"
}

func (code *FencedCode) String() string {
	return code.StartTag() + code.Children().String() + code.EndTag()
}

func (code *FencedCode) MD() string {
	if code == nil {
		return ""
	}

	bb := &bytes.Buffer{}

	fmt.Fprintf(bb, "```%s\n", code.Lang())

	body := code.Children().MD()
	body = html.UnescapeString(body)

	fmt.Fprintln(bb, body)

	fmt.Fprint(bb, "```")

	return bb.String()
}

func (code *FencedCode) Lang() string {
	lang := "plain"
	if code == nil {
		return lang
	}

	return Language(code.Attrs(), lang)
}

func NewFencedCode(el *Element) (*FencedCode, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	code := &FencedCode{
		Element: el,
	}

	if err := code.Set("language", code.Lang()); err != nil {
		return nil, err
	}

	if err := code.Set("class", "language-"+code.Lang()); err != nil {
		return nil, err
	}

	return code, nil
}

func NewFencedCodeNodes(p *Parser, el *Element) (Nodes, error) {
	code, err := NewFencedCode(el)
	if err != nil {
		return nil, err
	}

	return Nodes{code}, nil
}
