package hype

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

type executeNode struct {
	*Element
	ExecuteFn
}

func newExecuteNode(t testing.TB, fn ExecuteFn) *executeNode {
	t.Helper()

	return &executeNode{
		ExecuteFn: fn,
		Element:   &Element{},
	}
}

func Test_Nodes_Execute(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	doc := &Document{}

	n1 := newExecuteNode(t, func(ctx context.Context, d *Document) error {
		d.Lock()
		d.Nodes = append(d.Nodes, Text("foo"))
		d.Unlock()
		return nil
	})

	n1.Nodes = append(n1.Nodes, newExecuteNode(t, func(ctx context.Context, d *Document) error {
		d.Lock()
		d.Nodes = append(d.Nodes, Text("bar"))
		d.Unlock()
		return nil
	}))

	nodes := Nodes{n1}

	wg := &errgroup.Group{}
	err := nodes.Execute(wg, nil, doc)
	r.NoError(err)

	err = wg.Wait()
	r.NoError(err)

	r.Len(doc.Nodes, 2)
}

func Test_Nodes_Execute_Errors(t *testing.T) {
	t.Parallel()

	en := newExecuteNode(t, func(ctx context.Context, d *Document) error {
		return fmt.Errorf("boom")
	})

	nen := newExecuteNode(t, func(ctx context.Context, d *Document) error {
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

			wg := &errgroup.Group{}
			err := nodes.Execute(wg, nil, &Document{})
			r.NoError(err)

			err = wg.Wait()
			r.Error(err)

			r.True(errors.Is(err, ExecuteError{}), err)
		})
	}
}
