package hype

import (
	"context"
	"fmt"
	"strings"
)

type PostExecuter interface {
	PostExecute(ctx context.Context, d *Document, err error) error
}

func (list Nodes) PostExecute(ctx context.Context, d *Document, err error) error {
	var err2 error

	for _, n := range list {
		if nodes, ok := n.(Nodes); ok {
			err2 = nodes.PostExecute(ctx, d, err)
			if err2 != nil {
				return err2
			}
			continue
		}

		pe, ok := n.(PostExecuter)
		if ok {
			err2 = pe.PostExecute(ctx, d, err)
			if err2 != nil {
				return PostExecuteError{
					OrigErr:      err,
					Err:          err2,
					PostExecuter: pe,
				}
			}
		}

		err2 = n.Children().PostExecute(ctx, d, err)
		if err2 != nil {
			// the error should already be wrapped
			return err2
		}
	}

	return err
}

type PostExecuteFn func(ctx context.Context, d *Document, err error) error

func (fn PostExecuteFn) PostExecute(ctx context.Context, d *Document, err error) error {
	return fn(ctx, d, err)
}

type PostExecuteError struct {
	Err          error
	OrigErr      error
	PostExecuter PostExecuter
}

func (e PostExecuteError) Error() string {
	var errs []string

	if e.Err != nil {
		errs = append(errs, e.Err.Error())
	}

	if e.OrigErr != nil {
		errs = append(errs, e.OrigErr.Error())
	}

	return fmt.Sprintf("post execute error: [%T]: %v", e.PostExecuter, strings.Join(errs, "; "))
}
