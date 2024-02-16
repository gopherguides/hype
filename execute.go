package hype

import (
	"context"
)

type ExecutableNode interface {
	Node
	Execute(ctx context.Context, d *Document) error
}

type ExecuteFn func(ctx context.Context, d *Document) error

func (fn ExecuteFn) Execute(ctx context.Context, d *Document) error {
	return fn(ctx, d)
}

type WaitGrouper interface {
	Go(fn func() error)
}

func (list Nodes) Execute(wg WaitGrouper, ctx context.Context, d *Document) error {
	if d == nil {
		return ErrIsNil("document")
	}

	for _, n := range list {

		if nodes, ok := n.(Nodes); ok {
			err := nodes.Execute(wg, ctx, d)
			if err != nil {
				return err
			}
			continue
		}

		name := d.Filename

		if n, ok := n.(interface{ FileName() string }); ok {
			name = n.FileName()
		}

		cn, ok := n.(ExecutableNode)
		if ok {
			wg.Go(func() error {
				err := cn.Execute(ctx, d)
				if err != nil {
					return ExecuteError{
						Err:      err,
						Filename: name,
						Root:     d.Root,
					}
				}
				return nil
			})
		}

		err := n.Children().Execute(wg, ctx, d)
		if err != nil {
			return ExecuteError{
				Err:      err,
				Root:     d.Root,
				Filename: name,
			}
		}

	}

	return nil

}
