package hype

import (
	"strings"
)

type EmptyableNode interface {
	IsEmptyNode() bool
}

func IsEmptyNode(node Node) bool {
	if node == nil {
		return true
	}

	if n, ok := node.(EmptyableNode); ok {
		return n.IsEmptyNode()
	}

	for _, n := range node.Children() {
		if !IsEmptyNode(n) {
			return false
		}
	}

	s := node.Children().String()
	s = strings.TrimSpace(s)
	return len(s) == 0
}
