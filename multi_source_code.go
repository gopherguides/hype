package hype

import (
	"fmt"
	"io/fs"
	"strings"
)

var _ Tag = &SourceCode{}

type MultiSourceCode struct {
	*Node
}

func (ms MultiSourceCode) Lang() string {
	return "multi"
}

func (ms MultiSourceCode) StartTag() string {
	return ""
}

func (ms MultiSourceCode) EndTag() string {
	return ""
}

func (ms MultiSourceCode) String() string {
	return ms.GetChildren().String()
}

func NewMultiSourceCode(cab fs.FS, node *Node, rules map[string]string) (*MultiSourceCode, error) {
	if node == nil {
		return nil, fmt.Errorf("node cannot be nil")
	}

	ms := &MultiSourceCode{
		Node: node,
	}

	src, err := ms.Get("src")
	if err != nil {
		return nil, err
	}

	srcs := strings.Split(src, ",")

	if _, ok := ms.attrs["snippet"]; ok {
		return nil, fmt.Errorf("snippets can't be combined with multiple sources")
	}

	for _, src := range srcs {
		kn := node.Clone()
		kn.attrs["src"] = src
		kid, err := NewSourceCode(cab, kn, rules)
		if err != nil {
			return nil, err
		}
		ms.Children = append(ms.Children, kid)
	}

	return ms, nil
}
