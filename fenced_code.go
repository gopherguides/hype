package hype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
)

type FencedCode struct {
	*Element
}

func (code *FencedCode) MarshalJSON() ([]byte, error) {
	if code == nil {
		return nil, ErrIsNil("fenced code")
	}

	code.RLock()
	defer code.RUnlock()

	m, err := code.JSONMap()
	if err != nil {
		return nil, err
	}

	lang := code.Lang()
	if lang != "" {
		m["lang"] = lang
	}

	m["type"] = toType(code)

	return json.MarshalIndent(m, "", "  ")
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
