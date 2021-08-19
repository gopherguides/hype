package hype

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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

	text := strings.TrimSpace(c.Children.String())
	if _, ok := c.Attrs()["language"]; !ok {
		// This is `inline` code.
		fmt.Fprint(sb, c.StartTag())
		fmt.Fprint(sb, strings.TrimSpace(text))
		fmt.Fprint(sb, c.EndTag())
		return sb.String()
	}

	var inPre bool

	if c.Node.Parent != nil {
		inPre = c.Node.Parent.DataAtom == atom.Pre
	}

	if !inPre {
		fmt.Fprint(sb, "<pre>")
	}

	fmt.Fprint(sb, c.StartTag())
	fmt.Fprint(sb, text)
	fmt.Fprint(sb, c.EndTag())

	if !inPre {
		fmt.Fprint(sb, "</pre>")
	}
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

	for _, v := range ats {
		if !strings.HasPrefix(v, "language-") {
			continue
		}
		ats["language"] = strings.TrimPrefix(v, "language-")
	}

	lang, ok := ats["language"]
	if !ok {
		lang = filepath.Ext(ats["src"])
		lang = strings.TrimPrefix(lang, ".")
	}

	if len(lang) > 0 {
		c.Set("language", lang)
		c.Set("class", fmt.Sprintf("language-%s", lang))
	}

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

	c.Children = Tags{text}

	return c, nil
}
