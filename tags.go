package hype

import (
	"bytes"
	"fmt"
	"strings"
)

// T

type Tags []Tag

// Validate the tags and their children.
func (tags Tags) Validate(p *Parser, checks ...ValidatorFn) error {
	for _, t := range tags {
		if v, ok := t.(Validatable); ok {
			if err := v.Validate(p, checks...); err != nil {
				return err
			}
		}

		if err := t.GetChildren().Validate(p, checks...); err != nil {
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

// ByAtom returns all tags that match the given atoms.
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

// ByAtom returns all tags that match the given atoms.
func ByAtom(tags Tags, want ...Atom) Tags {
	return tags.ByAtom(want...)
}

// ByAttrs returns all tags that match the given attributes.
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

// ByAttrs returns all tags that match the given attributes.
func ByAttrs(tags Tags, query Attributes) Tags {
	return tags.ByAttrs(query)
}

// ByType returns all tags that match the given type.
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

func (tags Tags) Delete(atoms ...Atom) Tags {
	var res Tags
	for _, t := range tags {
		if t.Atom().Is(atoms...) {
			continue
		}

		node := t.DaNode()
		node.Children = node.Children.Delete(atoms...)
		res = append(res, t)
	}
	return res
}

type Markdowner interface {
	Tag
	Markdown() string
}

func (tags Tags) Markdown() string {
	bb := &bytes.Buffer{}

	for _, tag := range tags {
		md, ok := tag.(Markdowner)
		if ok {
			fmt.Fprint(bb, md.Markdown())
			continue
		}
		bb.WriteString(tag.GetChildren().Markdown())
	}

	return bb.String()
}
