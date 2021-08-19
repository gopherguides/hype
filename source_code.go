package hype

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type SourceCode struct {
	*Node
	Snippets Snippets
	lang     string
}

func (c *SourceCode) Src() string {
	return c.attrs["src"]
}

func (c *SourceCode) Lang() string {
	if len(c.lang) > 0 {
		return c.lang
	}

	lang := filepath.Ext(c.Src())
	lang = strings.TrimPrefix(lang, ".")
	c.lang = lang
	return lang
}

func (c *SourceCode) String() string {
	sb := &strings.Builder{}

	text := c.Children.String()
	text = strings.TrimSpace(text)
	fmt.Fprint(sb, "<pre>")
	fmt.Fprint(sb, c.StartTag())
	fmt.Fprint(sb, text)
	fmt.Fprint(sb, c.EndTag())
	fmt.Fprint(sb, "</pre>")
	return sb.String()
}

func (p *Parser) NewSourceCode(node *Node) (*SourceCode, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("source code node can not be nil")
	}

	if node.Data != "code" {
		return nil, fmt.Errorf("node is not code %v", node.Data)
	}

	c := &SourceCode{
		Node: node,
	}

	src, err := c.Get("src")
	if err != nil {
		return nil, err
	}

	lang := c.Lang()
	c.Set("language", lang)
	c.Set("class", fmt.Sprintf("language-%s", lang))

	b, err := p.ReadFile(src)
	if err != nil {
		return nil, err
	}

	tn := &html.Node{
		Data: string(b),
		Type: html.TextNode,
	}

	text, err := p.NewText(tn)
	if err != nil {
		return nil, err
	}

	c.Children = Tags{text}

	return c, nil
}
