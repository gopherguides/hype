package hype

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type Code struct {
	*Node
}

func (c *Code) Lang() string {
	if l, ok := c.Attrs()["language"]; ok {
		return l
	}
	return "plain"
}

func (c *Code) String() string {
	sb := &strings.Builder{}
	fmt.Fprintln(sb, "<pre>")
	fmt.Fprintln(sb, c.StartTag())
	fmt.Fprint(sb, c.Children.String())
	fmt.Fprintf(sb, c.EndTag())
	fmt.Fprintln(sb, "</pre>")
	return sb.String()
}

func (p *Parser) NewCode(node *Node) (*Code, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("code node can not be nil")
	}

	if node.Data != "code" {
		return nil, fmt.Errorf("node is not code %v", node.Data)
	}

	c := &Code{
		Node: node,
	}

	ats := c.Attrs()

	lang, ok := ats["language"]
	if !ok {
		lang = filepath.Ext(ats["src"])
		lang = strings.TrimPrefix(lang, ".")
	}

	if len(lang) == 0 {
		lang = "plain"
	}

	c.Set("language", lang)

	src, ok := ats["src"]
	if !ok {
		return c, nil
	}

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

	c.Children = append(c.Children, text)

	return c, nil
}
