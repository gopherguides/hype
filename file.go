package hype

import (
	"bytes"
	"fmt"
	"io/fs"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type File struct {
	*Node
}

func (File) Atom() atom.Atom {
	return atomx.File
}

func (c *File) Source() (Source, bool) {
	c.RLock()
	defer c.RUnlock()
	return SrcAttr(c.attrs)
}

func (c *File) SetSource(src string) {
	c.Lock()
	defer c.Unlock()

	c.attrs["src"] = src
}

func (f *File) String() string {
	if f.Node == nil {
		return "<file />"
	}

	bb := &bytes.Buffer{}

	fmt.Fprint(bb, f.StartTag())

	body := f.Children.String()

	fmt.Fprint(bb, body)

	fmt.Fprint(bb, f.EndTag())
	return bb.String()
}

func (f File) Validate(checks ...ValidatorFn) error {
	checks = append(checks, DataValidator("file"))

	return f.Node.Validate(html.ElementNode, checks...)
}

func (f File) ValidateFS(cab fs.FS, checks ...ValidatorFn) error {
	checks = append(checks, SourceValidator(cab, &f))

	return f.Validate(checks...)
}

func NewFile(cab fs.FS, node *Node) (*File, error) {
	fg := &File{
		Node: node,
	}

	err := fg.ValidateFS(cab)

	if err != nil {
		return nil, err
	}

	fg.Node.DataAtom = fg.Atom()

	return fg, fg.ValidateFS(cab)
}

func (p *Parser) NewFile(node *Node) (*File, error) {
	return NewFile(p.FS, node)
}
