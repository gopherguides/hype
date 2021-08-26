package htmx

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

type Attributes map[string]string

func AttributesString(node *html.Node) string {
	return NewAttributes(node).String()
}

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

func AttrNode(name string, ats Attributes) *html.Node {
	node := ElementNode(name)
	node.Attr = ats.Attrs()
	return node
}
