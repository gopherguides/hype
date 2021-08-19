package hype

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/net/html/atom"
)

type Tag interface {
	Attrs() Attributes
	DaNode() *Node
	GetChildren() Tags
	fmt.Stringer
}

type Tags []Tag

func (tags Tags) String() string {
	s := make([]string, 0, len(tags))
	for _, t := range tags {
		s = append(s, t.String())
	}
	return strings.Join(s, "")
}

func (tags Tags) AllAtom(want atom.Atom) Tags {
	var res Tags
	for _, t := range tags {
		if IsAtom(t, want) {
			res = append(res, t)
		}
		res = append(res, t.GetChildren().AllAtom(want)...)
	}
	return res
}

func (tags Tags) AllData(want string) Tags {
	var res Tags
	for _, t := range tags {
		node := t.DaNode()
		if node.Data == want {
			res = append(res, t)
		}
		res = append(res, t.GetChildren().AllData(want)...)
	}
	return res
}

func (tags Tags) AllType(want interface{}) Tags {
	var res Tags

	if want == nil {
		return res
	}

	wt := reflect.TypeOf(want)

	for _, t := range tags {
		tt := reflect.TypeOf(t)
		if tt.AssignableTo(wt) {
			res = append(res, t)
		}

		res = append(res, t.GetChildren().AllType(want)...)
	}

	return res
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
