package hype

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type postExecuterNode struct {
	*Element
	PostExecuteFn
}

func newPostExecuteNode(t testing.TB, fn PostExecuteFn) *postExecuterNode {
	t.Helper()

	return &postExecuterNode{
		PostExecuteFn: fn,
		Element:       &Element{},
	}
}

func Test_Nodes_PostExecute(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	nodes := Nodes{
		newPostExecuteNode(t, func(ctx context.Context, d *Document, err error) error {
			d.Title = "Hello"
			return nil
		}),
		newPostExecuteNode(t, func(ctx context.Context, d *Document, err error) error {
			d.Title += " World"
			return nil
		}),
	}

	doc := &Document{}

	err := nodes.PostExecute(nil, doc, nil)
	r.NoError(err)

	act := doc.Title
	exp := "Hello World"

	r.Equal(exp, act)
}

func Test_Nodes_PostExecute_Errors(t *testing.T) {
	t.Parallel()

	en := newPostExecuteNode(t, func(ctx context.Context, d *Document, err error) error {
		return fmt.Errorf("boom")
	})

	nen := newPostExecuteNode(t, func(ctx context.Context, d *Document, err error) error {
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

			err := nodes.PostExecute(nil, nil, fmt.Errorf("original"))
			r.Error(err)

			r.Contains(err.Error(), "boom")
			r.Contains(err.Error(), "original")
		})
	}
}
