package hype

import (
	"bytes"
	"fmt"

	"golang.org/x/net/html/atom"
)

const (
	FileGroup_Atom atom.Atom = 452184562
)

type FileGroup struct {
	*Node
}

func (fg *FileGroup) Name() string {
	fg.RLock()
	defer fg.RUnlock()
	return fg.attrs["name"]
}

func (fg *FileGroup) String() string {
	fg.RLock()
	defer fg.RUnlock()

	bb := &bytes.Buffer{}

	fmt.Fprint(bb, fg.StartTag())
	fmt.Fprint(bb, fg.Children.String())
	fmt.Fprint(bb, fg.EndTag())

	return bb.String()
}

func (p *Parser) NewFileGroup(node *Node) (*FileGroup, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("file node can not be nil")
	}

	if node.Data != "filegroup" {
		return nil, fmt.Errorf("node is not a filegroup %q", node.Data)
	}

	node.DataAtom = FileGroup_Atom

	fg := &FileGroup{
		Node: node,
	}

	if _, err := fg.Get("name"); err != nil {
		return nil, err
	}

	return fg, nil
}
