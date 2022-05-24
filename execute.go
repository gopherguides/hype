package hype

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type ExecutableNode interface {
	Node
	Execute(ctx context.Context, d *Document) error
}

type ExecuteFn func(ctx context.Context, d *Document) error

func (fn ExecuteFn) Execute(ctx context.Context, d *Document) error {
	return fn(ctx, d)
}

func (list Nodes) Execute(wg *errgroup.Group, ctx context.Context, d *Document) error {
	for _, n := range list {
		if nodes, ok := n.(Nodes); ok {
			err := nodes.Execute(wg, ctx, d)
			if err != nil {
				return err
			}
			continue
		}

		cn, ok := n.(ExecutableNode)
		if ok {
			wg.Go(func() error {
				return cn.Execute(ctx, d)
			})
		}

		err := n.Children().Execute(wg, ctx, d)
		if err != nil {
			return err
		}
	}

	return nil

}