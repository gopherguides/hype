package hype

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
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

func (fc FencedCode) Validate(checks ...ValidatorFn) error {
	checks = append(checks, AtomValidator("code"))
	return fc.Node.Validate(html.ElementNode, checks...)
}

func (p *Parser) NewFencedCode(node *Node) (*FencedCode, error) {
	return NewFencedCode(node)
}

func NewFencedCode(node *Node) (*FencedCode, error) {
	c := &FencedCode{
		Node: node,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	lang := c.Lang()
	c.Set("language", lang)
	c.Set("class", fmt.Sprintf("language-%s", lang))

	return c, c.Validate()
}
