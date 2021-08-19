package hype

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

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
	n.RLock()
	defer n.RUnlock()

	if v, ok := n.Attrs()[key]; ok {
		return v, nil
	}

	return "", fmt.Errorf("no attribute found %q", key)
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
	return json.Marshal(NewNodeJSON(g.Node))
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
