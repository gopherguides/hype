package htmx

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

// Attributes is a map of key/value pairs
// respresenting HTML element attributes.
//
// Example:
// 	`<div class="foo" id="bar" />`
// 	map[string]string{"class": "foo", "id": "bar"}
type Attributes map[string]string

// AttributesString returns a string representation of the
// HTMl node's attributes.
func AttributesString(node *html.Node) string {
	return NewAttributes(node).String()
}

// NewAttributes returns a new Attributes from the given node.
func NewAttributes(node *html.Node) Attributes {
	ats := Attributes{}
	if node == nil {
		return ats
	}

	for _, at := range node.Attr {
		ats[at.Key] = at.Val
	}

	return ats
}

// Get returns the value of the attribute with the given key.
// If the attribute does not exist an error is returned.
func (ats Attributes) Get(key string) (string, error) {
	if ats == nil {
		return "", fmt.Errorf("no attributes")
	}

	v, ok := ats[key]
	if !ok {
		return "", ErrAttrNotFound(key)
	}
	return v, nil
}

// HasKey returns true if the attributes contain all of the keys.
func (ats Attributes) HasKeys(keys ...string) bool {
	for _, key := range keys {
		if _, ok := ats[key]; !ok {
			return false
		}
	}
	return true
}

// Matches returns true if the attributes match the given keys/values.
//
// Specials:
//	*: Matches any attribute value.
//	map[string]string{"src": "*"}
func (ats Attributes) Matches(query map[string]string) bool {
	for k, v := range query {
		av, ok := ats[k]
		if !ok {
			return false
		}

		if v == "*" {
			continue
		}

		if av != v {
			return false
		}
	}
	return true
}

// Attrs returns a slice of html.Attribute from the attributes.
func (ats Attributes) Attrs() []html.Attribute {
	if ats == nil {
		return nil
	}

	ha := make([]html.Attribute, 0, len(ats))
	for k, v := range ats {
		ha = append(ha, html.Attribute{
			Key: k,
			Val: v,
		})
	}

	return ha
}

// String returns a string representation of the attributes
// in HTML element attribute format.
//
// Example:
//	`class="foo" id="bar"`
func (ats Attributes) String() string {
	if len(ats) == 0 {
		return ""
	}
	lines := make([]string, 0, len(ats))
	for k, v := range ats {
		lines = append(lines, fmt.Sprintf("%s=%q", k, v))
	}

	sort.Strings(lines)
	return strings.Join(lines, " ")
}

// AttrNode returns a new html.Node with the attributes.
func AttrNode(atom string, ats Attributes) *html.Node {
	node := ElementNode(atom)
	node.Attr = ats.Attrs()
	return node
}
