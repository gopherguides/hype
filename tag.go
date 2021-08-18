package hype

import (
	"fmt"
	"strings"

	"golang.org/x/net/html/atom"
)

type Tags []Tag

func (tags Tags) MetaData() MetaData {
	md := MetaData{}
	for _, t := range tags {

		if m, ok := t.(*Meta); ok {
			md[m.Key] = m.Val
		}

		for _, c := range t.GetChildren() {
			for k, v := range c.GetChildren().MetaData() {
				md[k] = v
			}
		}
	}
	return md
}

type Tag interface {
	Attrs() Attributes
	DaNode() *Node
	GetChildren() Tags
	fmt.Stringer
}

func (tags Tags) String() string {
	s := make([]string, 0, len(tags))
	for _, t := range tags {
		s = append(s, t.String())
	}
	return strings.Join(s, "")
}

func IsAtom(tag Tag, want atom.Atom) bool {
	if tag == nil {
		return false
	}
	n := tag.DaNode()
	if n == nil {
		return false
	}

	return n.DataAtom == want
}
