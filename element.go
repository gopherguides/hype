package hype

import (
	"bytes"
	"fmt"
	"path/filepath"
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
	FileName string // only set when Parser.ParseFile() is used
}

func (el *Element) Format(f fmt.State, verb rune) {
	if el == nil {
		return
	}

	switch verb {
	case 'v':
		st := el.StartTag()
		if len(st) == 0 {
			return
		}

		if len(el.FileName) > 0 {
			fmt.Fprintf(f, "file://%s: ", el.FileName)
		}

		fmt.Fprintf(f, "%s\n", st)

	default:
		fmt.Fprintf(f, "%s", el.String())
	}
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
		Nodes:    el.Nodes,
		Parent:   el.Parent,
		FileName: el.FileName,
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

func (el *Element) ValidAttr(k string) (string, error) {
	if el == nil {
		return "", ErrIsNil("element")
	}

	v, ok := el.Get(k)
	if !ok {
		return "", el.WrapErr(ErrAttrNotFound(k))
	}

	if len(v) == 0 {
		return "", el.WrapErr(ErrAttrEmpty(k))
	}

	return v, nil
}

func (el *Element) WrapErr(err error) error {
	return WrapNodeErr(el, err)
}

func NewEl[T ~string](at T, parent Node) *Element {
	var fn string

	if e, ok := parent.(*Element); ok {
		fn = e.FileName
	}

	return &Element{
		Attributes: &Attributes{},
		HTMLNode: &html.Node{
			Type:     html.ElementNode,
			Data:     string(at),
			DataAtom: atom.Lookup([]byte(string(at))),
		},
		Parent:   parent,
		FileName: fn,
	}
}

func (el *Element) updateFileName(dir string) {
	if el == nil {
		return
	}

	if strings.HasPrefix(el.FileName, dir) {
		return
	}

	el.Lock()
	defer el.Unlock()
	el.FileName = filepath.Join(dir, el.FileName)
}
