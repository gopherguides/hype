package hype

import (
	"bytes"
	"fmt"

	"golang.org/x/net/html/atom"
)

type FileGroup struct {
	*Node
	name string
}

func (FileGroup) Atom() atom.Atom {
	return FileGroup_Atom
}

func (fg *FileGroup) Name() string {
	fg.Lock()
	name := fg.name
	if len(name) == 0 {
		name = fg.attrs["name"]
		fg.name = name
	}
	fg.Unlock()
	return name
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

	fg := &FileGroup{
		Node: node,
	}
	fg.Node.DataAtom = fg.Atom()

	name, err := fg.Get("name")
	if err != nil {
		return nil, err
	}

	fg.name = name

	return fg, nil
}
