package hype

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"

	"github.com/gopherguides/hype/atomx"
	"github.com/markbates/sweets"
	"golang.org/x/net/html"
)

var _ Tag = &SourceCode{}
var _ Validatable = &SourceCode{}
var _ ValidatableFS = &SourceCode{}

// SourceCode represents a code file on disk.
//
// HTML Attributes:
// 	src (required): the path to the source file.
// 		`<code src="foo.go"></code>`
// 	snippet: the name of a snippet to use.
// 		`<code src="foo.go" snippet="example"></code>`
// 	lang: the language of the code. Defaults to the file extension.
type SourceCode struct {
	*Node
	Snippets Snippets // map of snippets in the file
	Body     string   // Full source of file
	lang     string   // language of the file
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

	return fmt.Sprintf("<p><pre>%s", t)
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

func (sc SourceCode) Validate(p *Parser, checks ...ValidatorFn) error {
	fn := func(p *Parser, n *Node) error {

		if _, ok := sc.Source(); !ok {
			return fmt.Errorf("missing source: %v", sc)
		}

		if n, ok := sc.attrs["section"]; ok {
			return fmt.Errorf("section is no longer supported, use snippet instead %s", n)
		}

		return nil
	}

	checks = append(checks, AtomValidator(atomx.Code), fn)

	return sc.Node.Validate(p, html.ElementNode, checks...)
}

func (sc SourceCode) ValidateFS(p *Parser, cab fs.FS, checks ...ValidatorFn) error {
	checks = append(checks, SourceValidator(cab, &sc))
	return sc.Validate(p, checks...)
}

func NewSourceCode(cab fs.FS, node *Node, rules map[string]string) (*SourceCode, error) {
	c := &SourceCode{
		Node: node,
	}

	if err := c.Validate(nil); err != nil {
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

	b, err := fs.ReadFile(cab, src)
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

	c.Children = Tags{QuickText(esc)}

	return c, nil
}
