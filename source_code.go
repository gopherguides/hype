package hype

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gopherguides/hype/atomx"
	"github.com/gopherguides/hype/internal/lone"
)

type SourceCode struct {
	*Element
	Snippet
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
	bb := &bytes.Buffer{}
	bb.WriteString(code.StartTag())

	if !code.Snippet.IsZero() {
		bb.WriteString(code.Snippet.String())
	} else {
		bb.WriteString(code.Children().String())
	}

	bb.WriteString(code.EndTag())

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

	code.Content = string(b)

	code.Nodes = Nodes{}

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

	ranger := lone.Ranger{
		End: len(lines),
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

	if ranger.End > len(lines) {
		return code.WrapErr(fmt.Errorf("range %q extends past end %d of file %q", name, len(lines), src))
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

	code.Snippet = snip

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

	code.Snippet = snippet
	code.Lang = snippet.Lang

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
