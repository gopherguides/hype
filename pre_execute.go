package hype

import (
	"context"
)

type PreExecuter interface {
	PreExecute(ctx context.Context, d *Document) error
}

func (list Nodes) PreExecute(ctx context.Context, d *Document) error {
	if d == nil {
		return ErrIsNil("document")
	}

	var err error

	for _, n := range list {
		if nodes, ok := n.(Nodes); ok {
			err = nodes.PreExecute(ctx, d)
			if err != nil {
				return err
			}
			continue
		}

		pe, ok := n.(PreExecuter)
		if ok {
			err = pe.PreExecute(ctx, d)
			if err != nil {
				return PreExecuteError{
					Err:         err,
					Filename:    d.Filename,
					PreExecuter: pe,
					Root:        d.Root,
				}
			}
		}

		err = n.Children().PreExecute(ctx, d)
		if err != nil {
			return PreExecuteError{
				Err:         err,
				Filename:    d.Filename,
				PreExecuter: pe,
				Root:        d.Root,
			}
		}
	}

	return nil
}

type PreExecuteFn func(ctx context.Context, d *Document) error

func (fn PreExecuteFn) PreExecute(ctx context.Context, d *Document) error {
	return fn(ctx, d)
}
