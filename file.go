package hype

import (
	"bytes"
	"fmt"
	"io/fs"

	"golang.org/x/net/html/atom"
)

const (
	File_Atom atom.Atom = 1421757657
)

type File struct {
	*Node
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

func (p *Parser) NewFile(node *Node) (*File, error) {
	return NewFile(p.FS, node)
}

func NewFile(cab fs.StatFS, node *Node) (*File, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("file node can not be nil")
	}

	if node.Data != "file" {
		return nil, fmt.Errorf("node is not a file %q", node.Data)
	}

	node.DataAtom = File_Atom

	fg := &File{
		Node: node,
	}

	source, ok := fg.Source()
	if !ok {
		return nil, fmt.Errorf("file node has no src attribute")
	}

	if _, err := source.StatFile(cab); err != nil {
		return nil, err
	}

	return fg, nil
}
