package hype

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/net/html/atom"
)

type Tag interface {
	Atomable
	Attrs() Attributes
	GetChildren() Tags
	Nodeable
	fmt.Stringer
}

type StartTagger interface {
	StartTag() string
}

type EndTagger interface {
	EndTag() string
}

type Tagger interface {
	StartTagger
	EndTagger
}

type Tags []Tag

func (tags Tags) String() string {
	s := make([]string, 0, len(tags))
	for _, t := range tags {
		s = append(s, t.String())
	}
	return strings.Join(s, "")
}

func (tags Tags) ByAtom(want atom.Atom) Tags {
	var res Tags
	for _, t := range tags {
		if IsAtom(t, want) {
			res = append(res, t)
		}
		res = append(res, t.GetChildren().ByAtom(want)...)
	}
	return res
}

func (tags Tags) ByData(want string) Tags {
	var res Tags
	for _, t := range tags {
		if t.DaNode().Data == want {
			res = append(res, t)
		}
		res = append(res, t.GetChildren().ByData(want)...)
	}
	return res
}

func (tags Tags) ByType(want interface{}) Tags {
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

		res = append(res, t.GetChildren().ByType(want)...)
	}

	return res
}

func (tags Tags) ByAttrs(query Attributes) Tags {
	var res Tags
	for _, t := range tags {
		ta := t.Attrs()

		if ta.Matches(query) {
			res = append(res, t)
		}

		res = append(res, t.GetChildren().ByAttrs(query)...)
	}
	return res
}
