package hype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/gopherguides/hype/atomx"
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
	Filename string // only set when Parser.ParseFile() is used
}

// FileName returns the filename of the element.
// This is only set when Parser.ParseFile() is used.
func (el *Element) FileName() string {
	// Note: this function was added to satisfy an interface.
	return el.Filename
}

func (el *Element) JSONMap() (map[string]any, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	el.RLock()
	defer el.RUnlock()

	m := map[string]any{
		"atom":       el.Atom(),
		"attributes": map[string]string{},
		"filename":   el.Filename,
		"nodes":      Nodes{},
		"tag":        el.StartTag(),
		"type":       fmt.Sprintf("%T", el),
	}

	if len(el.Nodes) > 0 {
		m["nodes"] = el.Nodes
	}

	if el.Attributes.Len() > 0 {
		m["attributes"] = el.Attributes
	}

	hn := el.HTMLNode
	if hn == nil {
		return m, nil
	}

	hnm := map[string]any{
		"data":      hn.Data,
		"data_atom": hn.DataAtom.String(),
		"namespace": hn.Namespace,
		"type":      fmt.Sprintf("%T", hn),
	}

	switch hn.Type {
	case html.ErrorNode:
		hnm["node_type"] = "html.ErrorNode"
	case html.TextNode:
		hnm["node_type"] = "html.TextNode"
	case html.DocumentNode:
		hnm["node_type"] = "html.DocumentNode"
	case html.ElementNode:
		hnm["node_type"] = "html.ElementNode"
	case html.CommentNode:
		hnm["node_type"] = "html.CommentNode"
	case html.DoctypeNode:
		hnm["node_type"] = "html.DoctypeNode"
	case html.RawNode:
		hnm["node_type"] = "html.RawNode"
	}

	if len(hn.Attr) > 0 {
		hnm["attributes"] = hn.Attr
	}

	m["html_node"] = hnm

	return m, nil
}

func (el *Element) MarshalJSON() ([]byte, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	m, err := el.JSONMap()
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(m, "", "  ")
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

		if len(el.Filename) > 0 {
			fmt.Fprintf(f, "file://%s: ", el.Filename)
		}

		fmt.Fprintf(f, "%s", st)

	default:
		fmt.Fprintf(f, "%s", el.String())
	}
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
	if el == nil {
		return ""
	}

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
	if el == nil {
		return ""
	}

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
		fn = e.Filename
	}

	return &Element{
		Attributes: &Attributes{},
		HTMLNode: &html.Node{
			Type:     html.ElementNode,
			Data:     string(at),
			DataAtom: atom.Lookup([]byte(string(at))),
		},
		Parent:   parent,
		Filename: fn,
	}
}

func (el *Element) updateFileName(dir string) {
	if el == nil {
		return
	}

	if strings.HasPrefix(el.Filename, dir) {
		return
	}

	el.Lock()
	defer el.Unlock()
	el.Filename = filepath.Join(dir, el.Filename)
}

func (el *Element) Set(k string, v string) error {
	err := el.Attributes.Set(k, v)

	if err != nil {
		return el.WrapErr(err)
	}

	return nil
}

func (el *Element) MD() string {
	if el == nil {
		return ""
	}

	switch el.Atom() {
	case atomx.Strong, atomx.B:
		return fmt.Sprintf("**%s**", el.Children().MD())
	case atomx.Em, atomx.I:
		return fmt.Sprintf("_%s_", el.Children().MD())
	case atomx.Pre:
		return el.Children().MD()
	case atomx.Hr:
		return "\n\n---\n\n"
	case atomx.Br:
		return "\n\n"
	case atomx.Blockquote:
		b := el.Children().MD()
		b = strings.TrimSpace(b)
		return fmt.Sprintf("> %s", b)
	case atomx.Details:
		return el.String()
	default:
		fmt.Printf("TODO: Element.MD(): %q\n", el.Atom())
	}

	bb := &bytes.Buffer{}
	bb.WriteString(el.StartTag())
	bb.WriteString(el.Children().MD())
	bb.WriteString(el.EndTag())

	return bb.String()
}
