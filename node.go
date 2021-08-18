package hype

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Nodes []*Node

type Node struct {
	*html.Node
	Children Tags
	attrs    Attributes
	moot     *sync.RWMutex
}

func (n *Node) StartTag() string {
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "<%s", n.DataAtom)
	ats := n.Attrs().String()
	if len(ats) > 0 {
		fmt.Fprintf(sb, " %s", ats)
	}
	fmt.Fprintf(sb, ">")
	return sb.String()
}

func (n *Node) EndTag() string {
	return fmt.Sprintf("</%s>", n.DataAtom)
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
	n.moot.Lock()
	defer n.moot.Unlock()
	if n.attrs == nil {
		n.attrs = Attributes{}
	}
	n.attrs[key] = val
}

// Get a key from the attributes. Will error if the key doesn't exist.
func (n *Node) Get(key string) (string, error) {
	n.moot.RLock()
	defer n.moot.RUnlock()

	if v, ok := n.Attrs()[key]; ok {
		return v, nil
	}

	return "", fmt.Errorf("no attribute found %q", key)
}

func NewNode(n *html.Node) *Node {
	node := &Node{
		Node:  n,
		attrs: NewAttributes(n),
		moot:  &sync.RWMutex{},
	}

	return node
}

// func (node Node) Format(state fmt.State, verb rune) {
// 	switch verb {
// 	case 'v':
// 		b, _ := node.MarshalJSON()
// 		state.Write(b)
// 	}
// 	// state.Write(string(node.String()))
// }

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
