package hype

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gopherguides/hype/internal/lone"
)

type SourceCode struct {
	*Element
	Lang    string
	Src     string
	Snippet Snippet
}

func (code *SourceCode) MarshalJSON() ([]byte, error) {
	if code == nil {
		return nil, ErrIsNil("code")
	}

	code.RLock()
	defer code.RUnlock()

	m, err := code.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = toType(code)

	if len(code.Lang) > 0 {
		m["lang"] = code.Lang
	}

	if len(code.Src) > 0 {
		m["filepath"] = code.Src
	}

	return json.MarshalIndent(m, "", "  ")
}

func (code *SourceCode) String() string {
	if code == nil {
		return "<code></code>"
	}

	bb := &bytes.Buffer{}
	// Build attribute string in consistent order
	attrs := []string{}
	if code.Lang != "" {
		attrs = append(attrs, fmt.Sprintf("class=\"language-%s\"", code.Lang))
	}
	if code.Lang != "" {
		attrs = append(attrs, fmt.Sprintf("language=\"%s\"", code.Lang))
	}
	if code.Src != "" {
		attrs = append(attrs, fmt.Sprintf("src=\"%s\"", code.Src))
	}
	if code.Snippet.Name != "" {
		attrs = append(attrs, fmt.Sprintf("snippet=\"%s\"", code.Snippet.Name))
	}
	if rng, ok := code.Get("range"); ok && rng != "" {
		attrs = append(attrs, fmt.Sprintf("range=\"%s\"", rng))
	}
	attrStr := ""
	if len(attrs) > 0 {
		attrStr = " " + strings.Join(attrs, " ")
	}

	fmt.Fprintf(bb, "<code%s>", attrStr)

	// Use the content of all children as the code body
	for _, n := range code.Nodes {
		if st, ok := n.(fmt.Stringer); ok {
			bb.WriteString(st.String())
		} else {
			bb.WriteString(n.Children().String())
		}
	}

	bb.WriteString("</code>")
	return bb.String()
}

func (code *SourceCode) MD() string {
	if code == nil {
		return ""
	}

	bb := &bytes.Buffer{}
	fmt.Fprintf(bb, "```%s\n", code.Lang)

	body := code.Children().MD()
	body = html.UnescapeString(body)

	body = strings.TrimSpace(body)

	fmt.Fprintln(bb, body)

	fmt.Fprintln(bb, "```")

	if len(code.Snippet.Name) > 0 {
		fmt.Fprintf(bb, "> *source: %s:%s*\n", code.Src, code.Snippet.Name)
	} else {
		fmt.Fprintf(bb, "> *source: %s*\n", code.Src)
	}

	return bb.String()
}

func (code *SourceCode) Execute(ctx context.Context, d *Document) error {
	if d == nil {
		return ErrIsNil("document")
	}

	if code == nil {
		return ErrIsNil("code")
	}

	if code.Element == nil {
		return ErrIsNil("element")
	}

	code.Lock()

	src, ok := code.Get("src")
	if !ok {
		code.Unlock()
		return nil
	}

	code.Src = src

	if len(code.Lang) == 0 {
		ext := filepath.Ext(src)
		ext = strings.TrimPrefix(ext, ".")
		if ext == "mod" {
			ext = "go"
		}
		code.Lang = Language(code.Attrs(), ext)
	}

	if err := code.Set("language", code.Lang); err != nil {
		return err
	}

	if err := code.Set("class", "language-"+code.Lang); err != nil {
		return err
	}

	bits := strings.Split(src, "#")
	if len(bits) >= 2 {
		code.Unlock()
		return code.parseSnippets(d, bits[0], bits[1])
	}

	if x, ok := code.Get("snippet"); ok {
		code.Unlock()
		return code.parseSnippets(d, src, x)
	}

	if x, ok := code.Get("range"); ok {
		code.Unlock()
		return code.parseRange(d, src, x)
	}

	defer code.Unlock()

	b, err := fs.ReadFile(d.FS, src)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", src, err)
	}

	b, err = d.Snippets.TrimComments(src, b)
	if err != nil {
		return fmt.Errorf("failed to trim comments from file %q: %w", src, err)
	}

	s := html.EscapeString(string(b))

	code.Nodes = Nodes{
		Text(s),
	}

	return nil
}

func (code *SourceCode) parseRange(d *Document, src string, name string) error {
	if d == nil {
		return ErrIsNil("document")
	}

	if err := code.validate(); err != nil {
		return err
	}

	code.Lock()
	defer code.Unlock()

	b, err := fs.ReadFile(d.FS, src)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", src, err)
	}

	lines := bytes.Split(b, []byte("\n"))
	ll := len(lines)

	ranger := lone.Ranger{
		End: ll,
	}

	if err := ranger.Parse(name); err != nil {
		return code.WrapErr(fmt.Errorf("failed to parse range %q: %w", name, err))
	}

	if err := ranger.Unsigned(); err != nil {
		return code.WrapErr(fmt.Errorf("failed to validate range %q: %w", name, err))
	}

	if err := ranger.Validate(); err != nil {
		return code.WrapErr(fmt.Errorf("failed to validate range %q: %w", name, err))
	}

	if ranger.End > ll {
		return code.WrapErr(fmt.Errorf("range %q extends past end %d of file %q", name, ll, src))
	}

	b = bytes.Join(lines[ranger.Start:ranger.End], []byte("\n"))

	snip := Snippet{
		Lang:    code.Lang,
		Content: string(b),
		Name:    name,
		File:    src,
		Start:   ranger.Start,
		End:     ranger.End,
	}

	code.Nodes = Nodes{snip}

	return nil
}

func (code *SourceCode) parseSnippets(d *Document, src string, name string) error {
	if d == nil {
		return ErrIsNil("document")
	}

	if err := code.validate(); err != nil {
		return err
	}

	code.Lock()
	defer code.Unlock()

	snips, ok := d.Snippets.Get(src)
	if ok {
		snip, ok := snips[name]
		if !ok {
			return fmt.Errorf("snippet %q not found in %q", name, src)
		}

		if err := code.setSnippet(snip); err != nil {
			return err
		}

		return nil
	}

	b, err := fs.ReadFile(d.FS, src)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", src, err)
	}

	snips, err = d.Snippets.Parse(src, b)
	if err != nil {
		return err
	}

	snip, ok := snips[name]
	if !ok {
		return fmt.Errorf("snippet %q not found in %q", name, src)
	}

	if err := code.setSnippet(snip); err != nil {
		return err
	}

	return nil
}

func (code *SourceCode) setSnippet(snippet Snippet) error {
	if err := code.validate(); err != nil {
		return err
	}

	code.Lang = snippet.Lang
	code.Nodes = Nodes{snippet}
	code.Snippet = snippet

	if err := code.Set("language", code.Lang); err != nil {
		return err
	}

	if err := code.Set("class", "language-"+code.Lang); err != nil {
		return err
	}

	return nil
}

func NewSourceCodeNodes(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	var codes Nodes

	code := &SourceCode{
		Element: el,
	}

	if err := code.validate(); err != nil {
		return nil, err
	}

	if sec, ok := el.Get("section"); ok {
		err := fmt.Errorf("`section` is no longer supported, use `snippet` instead: %q", sec)
		return nil, el.WrapErr(err)
	}

	// Wrap in <pre> for API compatibility - tests expect this structure
	pre := NewEl("pre", nil)
	pre.Nodes = append(pre.Nodes, code)
	codes = append(codes, pre)
	return codes, nil
}

func (code *SourceCode) validate() error {
	if code == nil {
		return code.WrapErr(ErrIsNil("code"))
	}

	if code.Element == nil {
		return code.WrapErr(ErrIsNil("element"))
	}

	return nil
}

func (code *SourceCode) updateFileName(dir string) {
	src, _ := code.Get("src")
	if strings.HasPrefix(src, dir) {
		return
	}

	src = filepath.Join(dir, src)
	code.Set("src", src)
}
