package hype

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var _ HTMLNode = &Element{}

type Element struct {
	*Attributes
	sync.RWMutex

	HTMLNode *html.Node
	Nodes    Nodes
	Parent   Node
}

func (el *Element) Clone() (*Element, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	ats, err := el.Attributes.Clone()
	if err != nil {
		return nil, err
	}

	nel := &Element{
		Attributes: ats,
		HTMLNode: &html.Node{
			Attr:      el.HTMLNode.Attr,
			Data:      el.HTMLNode.Data,
			DataAtom:  el.HTMLNode.DataAtom,
			Namespace: el.HTMLNode.Namespace,
			Type:      el.HTMLNode.Type,
		},
		Nodes:  el.Nodes,
		Parent: el.Parent,
	}

	return nel, nil
}

func (el *Element) Atom() Atom {
	if el.HTMLNode == nil {
		return Atom("")
	}

	return Atom(el.HTMLNode.Data)
}

func (el *Element) Children() Nodes {
	el.RLock()
	defer el.RUnlock()
	return el.Nodes
}

func (el *Element) HTML() *html.Node {
	return el.HTMLNode
}

// StartTag returns the start tag for the element.
// For example, for an element with an Atom of "div", the start tag would be "<div>".
func (el *Element) StartTag() string {
	a := el.Atom()
	if len(a) == 0 {
		return ""
	}

	if el.Attributes == nil || el.Attributes.Len() == 0 {
		return fmt.Sprintf("<%s>", a)
	}

	bb := &bytes.Buffer{}

	var lines []string

	el.Attributes.Range(func(k string, v string) bool {
		lines = append(lines, fmt.Sprintf("%s=%q", k, v))
		return true
	})

	sort.Strings(lines)

	fmt.Fprintf(bb, "<%s %s>", a, strings.Join(lines, " "))

	return bb.String()
}

// EndTag returns the end tag for the element.
// For example, for an element with an Atom of "div", the end tag would be "</div>".
func (el *Element) EndTag() string {
	a := el.Atom()
	if len(a) == 0 {
		return ""
	}

	return fmt.Sprintf("</%s>", a)
}

// String returns StartTag() + Children().String() + EndTag()
func (el *Element) String() string {
	s := el.StartTag()
	s += el.Children().String()
	s += el.EndTag()
	return s
}

func (el *Element) Attrs() *Attributes {
	if el == nil {
		return &Attributes{}
	}

	return el.Attributes
}

func NewEl[T ~string](at T, parent Node) *Element {
	return &Element{
		Attributes: &Attributes{},
		HTMLNode: &html.Node{
			Type:     html.ElementNode,
			Data:     string(at),
			DataAtom: atom.Lookup([]byte(string(at))),
		},
		Parent: parent,
	}
}
