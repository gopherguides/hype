package hype

import (
	"regexp"
	"strings"

	"github.com/markbates/syncx"
	"golang.org/x/net/html"
)

type AttrNode interface {
	Node
	Attrs() *Attributes
}

type Attributes = syncx.Map[string, string]

// ConvertHTMLAttrs converts a slice of HTML attributes
// to a new Attributes type.
func ConvertHTMLAttrs(attrs []html.Attribute) *Attributes {
	ats := &Attributes{}

	for _, a := range attrs {
		ats.Set(a.Key, a.Val)
	}

	return ats
}

// AttrMatches returns true if the given keys and values in the
// query map are present, and equal, in the given attributes.
// A `*` matches any value.
func AttrMatches(ats *Attributes, query map[string]string) bool {
	if ats == nil {
		return false
	}

	for k, v := range query {
		av, ok := ats.Get(k)
		if !ok {
			return false
		}

		if v == "*" {
			continue
		}

		if len(av) == 0 || len(v) == 0 {
			return false
		}

		rx, err := regexp.Compile(v)
		if err != nil {
			return false
		}

		if !rx.MatchString(av) {
			return false
		}
	}

	return true
}

// Language tries to determine the language of the given
// set of attributes.
// 	- "language" is the first attr tested.
// 	- "language-*" is the second attr tested. (e.g. "language-go", "language-js")
//  - "lang" is the third attr tested.
func Language(ats *Attributes, lang string) string {
	if ats == nil {
		return lang
	}

	if l, ok := ats.Get("language"); ok {
		if len(l) > 0 {
			return l
		}
	}

	var l string
	ats.Range(func(k, v string) bool {
		if !strings.HasPrefix(v, "language-") {
			return true
		}

		l = strings.TrimPrefix(v, "language-")
		return false
	})

	if len(l) > 0 {
		return l
	}

	if l, ok := ats.Get("lang"); ok {
		if len(l) > 0 {
			return l
		}
	}

	return lang
}
