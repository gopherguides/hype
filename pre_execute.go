package hype

import (
	"context"
	"fmt"
)

type PreExecuter interface {
	PreExecute(ctx context.Context, d *Document) error
}

func (list Nodes) PreExecute(ctx context.Context, d *Document) error {
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
					PreExecuter: pe,
				}
			}
		}

		err = n.Children().PreExecute(ctx, d)
		if err != nil {
			return err
		}
	}

	return nil
}

type PreExecuteFn func(ctx context.Context, d *Document) error

func (fn PreExecuteFn) PreExecute(ctx context.Context, d *Document) error {
	return fn(ctx, d)
}

type PreExecuteError struct {
	Err         error
	PreExecuter PreExecuter
}

func (e PreExecuteError) Error() string {
	return fmt.Sprintf("pre execute error: [%T]: %v", e.PreExecuter, e.Err)
}
