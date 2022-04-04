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
// 		// multiple sources
// 		`<code src="foo.go,bar.go"></code>`
// 		// sources with hash tag snippets
// 		`<code src="foo.go#example"></code>`
// 	snippet: the name of a snippet to use.
// 		`<code src="foo.go" snippet="example"></code>`
// 		// multiple sources with a snippet attribute
// 		// will apply the snippet to all sources
// 		`<code src="foo.go,bar.go" snippet="example"></code>`
// 	lang: the language of the code. Defaults to the file extension.
type SourceCode struct {
	*Node
	Snippets Snippets // map of snippets in the file
	Body     string   // Full source of file
	lang     string   // language of the file
	multiSrc bool
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
	if c.Node == nil {
		return ""
	}

	if c.multiSrc {
		return ""
	}

	t := c.Node.StartTag()

	return fmt.Sprintf("<p><pre>%s", t)
}

func (c *SourceCode) EndTag() string {
	if c.Node == nil {
		return ""
	}

	if c.multiSrc {
		return ""
	}

	t := c.Node.EndTag()

	return fmt.Sprintf("%s</pre></p>", t)
}

func (c SourceCode) Markdown() string {
	bb := &bytes.Buffer{}
	fmt.Fprintf(bb, "```%s\n", c.Lang())
	text := c.GetChildren().Markdown()
	text = html.UnescapeString(text)
	fmt.Fprintln(bb, text)
	fmt.Fprintln(bb, "```")
	return bb.String()
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
	sc := &SourceCode{
		Node: node,
	}

	if err := sc.Validate(nil); err != nil {
		return nil, err
	}

	return sc, nil
}

func (sc *SourceCode) Finalize(p *Parser) error {
	if sc == nil {
		return fmt.Errorf("source code is nil")
	}

	if p == nil {
		return fmt.Errorf("parser is nil")
	}

	cab := p.FS
	rules := p.snippetRules

	src, err := sc.Get("src")
	if err != nil {
		return err
	}

	srcs := strings.Split(src, ",")
	if len(srcs) > 1 {
		// handle multiple sources
		if err := sc.handleSources(p, srcs); err != nil {
			return err
		}
		return nil
	}

	if err := sc.handleLang(); err != nil {
		return err
	}

	fn := strings.Split(src, "#")[0]

	b, err := fs.ReadFile(cab, fn)
	if err != nil {
		return err
	}

	sc.Body = string(bytes.TrimSpace(b))

	if err := sc.handleSnippets(src, b, rules); err != nil {
		return err
	}

	return nil
}

func (sc *SourceCode) handleSnippets(src string, b []byte, rules map[string]string) error {
	snips, err := ParseSnippets(src, b, rules)
	if err != nil {
		return err
	}
	sc.Snippets = snips

	var name string

	split := strings.Split(src, "#")
	if len(split) > 1 {
		name = split[1]
	}

	if n, ok := sc.attrs["snippet"]; ok {
		if len(name) > 0 {
			return fmt.Errorf("snippet and snippet name cannot be used together %s", n)
		}
		name = n
	}

	if len(name) > 0 { // no snippet name
		snip, ok := sc.Snippets[name]
		if !ok {
			return fmt.Errorf("could not find snippet %q in %q", name, src)
		}

		b = []byte(snip.String())
	}

	esc := html.EscapeString(string(b))

	sc.Children = Tags{QuickText(esc)}

	return nil
}

func (sc *SourceCode) handleSources(p *Parser, srcs []string) error {
	if sc == nil {
		return fmt.Errorf("nil source code")
	}

	node := sc.Node
	if node == nil {
		return fmt.Errorf("nil node")
	}

	if p == nil {
		return fmt.Errorf("nil parser")
	}

	sc.multiSrc = true

	for _, src := range srcs {
		kn := node.Clone()
		kn.attrs["src"] = src

		kid, err := NewSourceCode(p.FS, kn, p.snippetRules)
		if err != nil {
			return err
		}

		if err := kid.Finalize(p); err != nil {
			return err
		}

		sc.Children = append(sc.Children, kid)
	}

	return nil
}

func (sc *SourceCode) handleLang() error {
	if sc == nil {
		return fmt.Errorf("nil source code")
	}

	if lang, ok := sc.attrs["lang"]; ok {
		sc.lang = lang
	}

	lang := sc.Lang()
	sc.Set("language", lang)
	sc.Set("class", fmt.Sprintf("language-%s", lang))

	return nil
}
