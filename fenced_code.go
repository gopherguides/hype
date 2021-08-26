package hype

import (
	"fmt"
	"strings"
)

type FencedCode struct {
	*Node
}

func (c FencedCode) String() string {
	sb := &strings.Builder{}

	text := c.Children.String()
	text = strings.TrimSpace(text)
	fmt.Fprint(sb, c.StartTag())
	fmt.Fprint(sb, text)
	fmt.Fprint(sb, c.EndTag())
	return sb.String()
}

func (c *FencedCode) Lang() string {
	ats := c.Attrs()
	if l, ok := ats["language"]; ok {
		return l
	}

	for _, v := range ats {
		if !strings.HasPrefix(v, "language-") {
			continue
		}

		lang := strings.TrimPrefix(v, "language-")
		c.Set("language", lang)
		return lang
	}

	return "plain"
}

func (p *Parser) NewFencedCode(node *Node) (*FencedCode, error) {
	return NewFencedCode(node)
}

func NewFencedCode(node *Node) (*FencedCode, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("fenced code node can not be nil")
	}

	if node.Data != "code" {
		return nil, fmt.Errorf("node is not code %v", node.Data)
	}

	c := &FencedCode{
		Node: node,
	}

	lang := c.Lang()
	c.Set("language", lang)
	c.Set("class", fmt.Sprintf("language-%s", lang))

	return c, nil
}
