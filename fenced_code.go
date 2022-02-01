package hype

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

var _ Tag = &FencedCode{}
var _ Validatable = &FencedCode{}

// FencedCode represents a fenced code block.
//
// Example:
// 	```go
// 	fmt.Println("Hello, World!")
// 	```
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

// Lang returns the language of the fenced code block.
// Defaults to plain.
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

func (fc FencedCode) Validate(p *Parser, checks ...ValidatorFn) error {
	checks = append(checks, AtomValidator("code"))
	return fc.Node.Validate(p, html.ElementNode, checks...)
}

// NewFencedCode returns a new FencedCode from the given node.
func NewFencedCode(node *Node) (*FencedCode, error) {
	c := &FencedCode{
		Node: node,
	}

	if err := c.Validate(nil); err != nil {
		return nil, err
	}

	lang := c.Lang()
	c.Set("language", lang)
	c.Set("class", fmt.Sprintf("language-%s", lang))

	s := node.Children.String()
	s = html.EscapeString(s)

	c.Children = Tags{QuickText(s)}

	return c, nil
}
