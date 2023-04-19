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

type MDNode interface {
	Node
	MD() string
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

func (list Nodes) MD() string {
	var s string
	for _, n := range list {
		if st, ok := n.(MDNode); ok {
			s += st.MD()
			continue
		}
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

func (list Nodes) updateFileName(dir string) {
	type namer interface {
		updateFileName(dir string)
	}

	for _, n := range list {
		if n, ok := n.(namer); ok {
			n.updateFileName(dir)
		}
		n.Children().updateFileName(dir)
	}
}

func ToNodes[T Node](list []T) Nodes {
	nodes := make(Nodes, len(list))
	for i, n := range list {
		nodes[i] = n
	}

	return nodes
}

func (list Nodes) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		for _, n := range list {
			fmt.Fprintf(f, "%v", n)
			n.Children().Format(f, verb)
		}
	default:
		fmt.Fprintf(f, "%s", list.String())
	}
}
