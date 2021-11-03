package hype

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"

	"github.com/markbates/sweets"
	"golang.org/x/net/html"
)

var _ Tag = &SourceCode{}

type SourceCode struct {
	*Node
	Snippets Snippets
	Body     string // Full source of file
	lang     string
}

func (c *SourceCode) Source() (Source, bool) {
	c.RLock()
	defer c.RUnlock()
	return SrcAttr(c.attrs)
}

func (c *SourceCode) SetSource(s string) {
	c.Lock()
	defer c.Unlock()
	c.attrs["src"] = s
}

func (c *SourceCode) Lang() string {
	if len(c.lang) > 0 {
		return c.lang
	}

	source, _ := c.Source()
	lang := source.Lang()

	c.Lock()
	c.lang = lang
	c.Unlock()

	return lang
}

func (c *SourceCode) StartTag() string {
	t := c.Node.StartTag()

	return fmt.Sprintf("<p><pre class=\"code-block\">%s", t)
}

func (c *SourceCode) EndTag() string {
	t := c.Node.EndTag()

	return fmt.Sprintf("%s</pre></p>", t)
}

// String returns a properly formatted <code> tag.
// If a snippet is defined on the original <code snippet="foo"> tag, then that snippet's content is used, otherwise the the Source code is used.
func (c *SourceCode) String() string {
	if c.Node == nil {
		return "<code />"
	}

	sb := &strings.Builder{}

	text := c.Children.String()
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		sl := strings.TrimSpace(line)
		if strings.HasPrefix(sl, "// snippet:") {
			continue
		}
		lines = append(lines, line)
	}
	text = strings.Join(lines, "\n")
	text = sweets.TrimLeftSpace(text)
	text = strings.TrimPrefix(text, "\n")

	fmt.Fprint(sb, c.StartTag())
	fmt.Fprint(sb, text)
	fmt.Fprint(sb, c.EndTag())
	return sb.String()
}

func (sc SourceCode) Validate(checks ...ValidatorFn) error {
	fn := func(n *Node) error {

		if _, ok := sc.Source(); !ok {
			return fmt.Errorf("missing source: %v", sc)
		}

		if n, ok := sc.attrs["section"]; ok {
			return fmt.Errorf("section is no longer supported, use snippet instead %s", n)
		}

		return nil
	}

	checks = append(checks, DataValidator("code"), fn)

	return sc.Node.Validate(html.ElementNode, checks...)
}

func (sc SourceCode) ValidateFS(cab fs.FS, checks ...ValidatorFn) error {
	checks = append(checks, SourceValidator(cab, &sc))
	return sc.Validate(checks...)
}

func NewSourceCode(cab fs.ReadFileFS, node *Node, rules map[string]string) (*SourceCode, error) {
	c := &SourceCode{
		Node: node,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	src, err := c.Get("src")
	if err != nil {
		return nil, err
	}

	if lang, ok := c.attrs["lang"]; ok {
		c.lang = lang
	}

	lang := c.Lang()
	c.Set("language", lang)
	c.Set("class", fmt.Sprintf("language-%s", lang))

	b, err := cab.ReadFile(src)
	if err != nil {
		return nil, err
	}
	c.Body = string(bytes.TrimSpace(b))

	snips, err := ParseSnippets(src, b, rules)
	if err != nil {
		return nil, err
	}
	c.Snippets = snips

	if n, ok := c.attrs["snippet"]; ok {
		snip, ok := c.Snippets[n]
		if !ok {
			return nil, fmt.Errorf("could not find snippet %q in %q", n, src)
		}
		b = []byte(snip.String())
	}

	esc := html.EscapeString(string(b))

	tn := &html.Node{
		Data: esc,
		Type: html.TextNode,
	}

	text, err := NewText(tn)
	if err != nil {
		return nil, err
	}

	c.Children = Tags{text}

	return c, c.ValidateFS(cab)
}

func (p *Parser) NewSourceCode(node *Node) (*SourceCode, error) {
	return NewSourceCode(p.FS, node, p.snippetRules)
}
