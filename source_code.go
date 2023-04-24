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

	"github.com/gopherguides/hype/atomx"
	"github.com/gopherguides/hype/internal/lone"
)

type SourceCode struct {
	*Element
	Lang string
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

	m["type"] = fmt.Sprintf("%T", code)

	if len(code.Lang) > 0 {
		m["lang"] = code.Lang
	}

	return json.Marshal(m)
}

func (code *SourceCode) StartTag() string {
	if code == nil || code.Element == nil {
		return ""
	}

	return code.Element.StartTag()
}

func (code *SourceCode) EndTag() string {
	if code == nil || code.Element == nil {
		return ""
	}

	return "</code>"
}

func (code *SourceCode) String() string {
	if code == nil {
		return ""
	}

	bb := &bytes.Buffer{}
	bb.WriteString(code.StartTag())

	s := code.Children().String()

	if _, ok := code.Get("esc"); ok {
		s = html.EscapeString(s)
	}

	bb.WriteString(s)
	bb.WriteString(code.EndTag())

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

	fmt.Fprintln(bb, body)

	fmt.Fprint(bb, "```")

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

	if code.Parent != nil {
		var at Atom
		if a, ok := code.Parent.(Atomable); ok {
			at = a.Atom()
		}

		if at == atomx.Pre {
			codes = append(codes, code)
			return codes, nil
		}
	}

	pre := NewEl("pre", el.Parent)
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
