package hype

import (
	"fmt"

	"golang.org/x/net/html"
)

type Node interface {
	Children() Nodes
}

type Nodes []Node

type HTMLNode interface {
	Node
	HTML() *html.Node
}

func (list Nodes) String() string {
	var s string
	for _, n := range list {
		if st, ok := n.(fmt.Stringer); ok {
			s += st.String()
			continue
		}
		s += n.Children().String()
	}

	return s
}

func (list Nodes) Children() Nodes {
	return list
}

func (list Nodes) Delete(node Node) Nodes {
	if len(list) == 0 {
		return list
	}

	nodes := make(Nodes, 0, len(list)-1)
	for _, n := range list {
		if n == node {
			continue
		}
		nodes = append(nodes, n)
	}

	return nodes
}

func ToNodes[T Node](list []T) Nodes {
	nodes := make(Nodes, len(list))
	for i, n := range list {
		nodes[i] = n
	}

	return nodes
}
