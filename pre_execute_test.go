package hype

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type preExecuteNode struct {
	*Element
	PreExecuteFn
}

func newPreExecuteNode(t testing.TB, fn PreExecuteFn) *preExecuteNode {
	t.Helper()

	return &preExecuteNode{
		PreExecuteFn: fn,
		Element:      &Element{},
	}
}

func Test_Nodes_PreExecute(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	doc := &Document{}

	n1 := newPreExecuteNode(t, func(ctx context.Context, d *Document) error {
		d.Title = "Hello"
		return nil
	})

	n1.Nodes = append(n1.Nodes, newPreExecuteNode(t, func(ctx context.Context, d *Document) error {
		d.Title += " World"
		return nil
	}))

	nodes := Nodes{n1}

	err := nodes.PreExecute(context.Background(), doc)
	r.NoError(err)

	act := doc.Title
	exp := "Hello World"

	r.Equal(exp, act)
}

func Test_Nodes_PreExecute_Errors(t *testing.T) {
	t.Parallel()

	en := newPreExecuteNode(t, func(ctx context.Context, d *Document) error {
		return fmt.Errorf("boom")
	})

	nen := newPreExecuteNode(t, func(ctx context.Context, d *Document) error {
		return nil
	})
	nen.Nodes = append(nen.Nodes, en)

	table := []struct {
		name string
		node Node
	}{
		{"top level", en},
		{"nested", nen},
	}

	for _, tc := range table {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			nodes := Nodes{tc.node}

			err := nodes.PreExecute(context.Background(), nil)
			r.Error(err)

			r.Contains(err.Error(), "boom")
		})
	}
}
