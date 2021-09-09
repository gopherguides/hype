package hype

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

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

// String returns a properly formatted <code> tag.
// If a snippet is defined on the original <code snippet="foo"> tag, then that snippet's content is used, otherwise the the Source code is used.
func (c *SourceCode) String() string {
	sb := &strings.Builder{}

	text := c.Children.String()
	text = strings.TrimPrefix(text, "\n")
	// text = strings.TrimSuffix(text, "\n")
	fmt.Fprint(sb, "<pre>")
	fmt.Fprint(sb, c.StartTag())
	fmt.Fprint(sb, text)
	fmt.Fprint(sb, c.EndTag())
	fmt.Fprint(sb, "</pre>")
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

	tn := &html.Node{
		Data: string(b),
		Type: html.TextNode,
	}

	text, err := NewText(tn)
	if err != nil {
		return nil, err
	}

	c.Children = Tags{text}

	return c, nil
}
