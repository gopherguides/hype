package hype

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gopherguides/hype/htmx"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Nodeable interface {
	DaNode() *Node
}

type Nodes []*Node

type Node struct {
	*html.Node
	*sync.RWMutex
	Children Tags
	attrs    Attributes
}

func (n *Node) Validate(nt html.NodeType, validators ...ValidatorFn) error {
	if n == nil {
		return fmt.Errorf("nil node")
	}

	if n.Node == nil {
		return fmt.Errorf("html node is nil: %v", n)
	}

	if nt == 0 {
		return fmt.Errorf("invalid NodeType provided: %v", nt)
	}

	if n.Type != nt {
		return fmt.Errorf("node type mismatch: %v != %v", n.Type, nt)
	}

	for _, v := range validators {
		if err := v(n); err != nil {
			return err
		}
	}

	return nil
}

func cloneHTMLNode(n *html.Node) *html.Node {
	if n == nil {
		return nil
	}

	return &html.Node{
		Attr:        n.Attr,
		Data:        n.Data,
		DataAtom:    n.DataAtom,
		FirstChild:  cloneHTMLNode(n.FirstChild),
		LastChild:   cloneHTMLNode(n.LastChild),
		Namespace:   n.Namespace,
		NextSibling: cloneHTMLNode(n.NextSibling),
		Parent:      n.Parent,
		PrevSibling: n.PrevSibling,
	}
}

func (n *Node) Clone() *Node {
	node := &Node{
		Children: n.Children,
		Node:     cloneHTMLNode(n.Node),
		RWMutex:  &sync.RWMutex{},
		attrs:    n.Attrs(),
	}
	return node
}

func (n *Node) Atom() atom.Atom {
	if n.Node != nil {
		return n.DataAtom
	}
	return atom.Atom(0)
}

func (n *Node) StartTag() string {
	sb := &strings.Builder{}

	at := n.DataAtom.String()
	if len(at) == 0 {
		at = n.Data
	}

	fmt.Fprintf(sb, "<%s", at)
	ats := n.Attrs().String()
	if len(ats) > 0 {
		fmt.Fprintf(sb, " %s", ats)
	}
	fmt.Fprintf(sb, ">")
	return sb.String()
}

func (n *Node) EndTag() string {
	at := n.DataAtom.String()
	if len(at) == 0 {
		at = n.Data
	}
	return fmt.Sprintf("</%s>", at)
}

func (n *Node) InlineTag() string {
	st := n.StartTag()
	st = strings.TrimSuffix(st, ">")
	st += " />"
	return st
}

func (n *Node) GetChildren() Tags {
	return n.Children
}

func (n *Node) DaNode() *Node {
	return n
}

// Attrs returns a copy of the attributes, not the underlying attributes. Use Set to modify attributes.
func (n *Node) Attrs() Attributes {
	n.RLock()
	defer n.RUnlock()

	if n.attrs == nil {
		return Attributes{}
	}

	ats := Attributes{}
	for k, v := range n.attrs {
		ats[k] = v
	}

	return ats
}

func (n *Node) Set(key string, val string) {
	n.Lock()
	defer n.Unlock()
	if n.attrs == nil {
		n.attrs = Attributes{}
	}
	n.attrs[key] = val
}

// Get a key from the attributes. Will error if the key doesn't exist.
func (n *Node) Get(key string) (string, error) {
	n.Lock()
	if n.attrs == nil {
		n.attrs = Attributes{}
	}
	n.Unlock()

	n.RLock()
	defer n.RUnlock()
	return n.attrs.Get(key)
}

func NewNode(n *html.Node) *Node {
	node := &Node{
		Node:    n,
		RWMutex: &sync.RWMutex{},
		attrs:   NewAttributes(n),
	}

	return node
}

func (g Node) MarshalJSON() ([]byte, error) {
	return htmx.MarshalNode(g.Node)
}

func (p *Parser) ParseNode(node *html.Node) (Tag, error) {
	if node == nil {
		return nil, fmt.Errorf("nil node")
	}

	switch node.Type {
	case html.CommentNode:
		return p.NewComment(node)
	case html.DoctypeNode:
		return p.NewDocType(node)
	case html.DocumentNode:
		return p.NewDocument(node)
	case html.ElementNode:
		return p.ElementNode(node)
	case html.ErrorNode:
		return nil, fmt.Errorf(node.Data)
	case html.TextNode:
		return p.NewText(node)
	}
	return nil, fmt.Errorf("unknown type %v", node)
}
