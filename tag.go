package hype

import (
	"fmt"
	"strings"

	"github.com/gopherguides/hype/atomx"
)

type Tag interface {
	atomx.Atomable
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

func (tags Tags) Validate(checks ...ValidatorFn) error {
	for _, t := range tags {
		if v, ok := t.(Validatable); ok {
			if err := v.Validate(checks...); err != nil {
				return err
			}
		}

		if err := t.GetChildren().Validate(checks...); err != nil {
			return err
		}
	}
	return nil
}

func (tags Tags) String() string {
	s := make([]string, 0, len(tags))
	for _, t := range tags {
		s = append(s, t.String())
	}
	return strings.Join(s, "")
}

func (tags Tags) ByAtom(want ...Atom) Tags {
	var res Tags
	for _, t := range tags {
		for _, w := range want {
			if t.Atom() == w {
				res = append(res, t)
				break
			}
		}
		res = append(res, t.GetChildren().ByAtom(want...)...)
	}
	return res
}

// func (tags Tags) ByType(want any) Tags {
// 	var res Tags

// 	if want == nil {
// 		return res
// 	}

// 	wt := reflect.TypeOf(want)

// 	for _, t := range tags {
// 		tt := reflect.TypeOf(t)
// 		if tt.AssignableTo(wt) {
// 			res = append(res, t)
// 		}

// 		res = append(res, t.GetChildren().ByType(want)...)
// 	}

// 	return res
// }

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

func ByType[T Tag](tags Tags, want T) []T {
	var res []T

	for _, t := range tags {
		if x, ok := t.(T); ok {
			res = append(res, x)
		}

		res = append(res, ByType(t.GetChildren(), want)...)
	}

	return res
}
