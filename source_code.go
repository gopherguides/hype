package hype

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/markbates/sweets"
	"golang.org/x/net/html"
)

var _ Tag = &SourceCode{}

type SourceCode struct {
	*Node
	Snippets Snippets
	Source   string // Full source of file
	lang     string
}

func (c *SourceCode) Src() string {
	c.RLock()
	defer c.RUnlock()
	return c.attrs["src"]
}

func (c *SourceCode) Lang() string {
	if len(c.lang) > 0 {
		return c.lang
	}

	lang := filepath.Ext(c.Src())
	lang = strings.TrimPrefix(lang, ".")
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

func (p *Parser) NewSourceCode(node *Node) (*SourceCode, error) {
	return NewSourceCode(p.FS, node, p.snippetRules)
}

func NewSourceCode(cab fs.ReadFileFS, node *Node, rules map[string]string) (*SourceCode, error) {
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

	if n, ok := c.attrs["section"]; ok {
		return nil, fmt.Errorf("section is no longer supported, use snippet instead %s", n)
	}

	lang := c.Lang()
	c.Set("language", lang)
	c.Set("class", fmt.Sprintf("language-%s", lang))

	b, err := cab.ReadFile(src)
	if err != nil {
		return nil, err
	}
	c.Source = string(bytes.TrimSpace(b))

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

	return c, nil
}
