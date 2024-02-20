package hype

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PostExecuteError(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	pee := PostExecuteError{
		Err: io.EOF,
	}

	wrapped := fmt.Errorf("error: %w", pee)

	r.True(pee.As(&PostExecuteError{}), pee)
	r.True(pee.Is(pee), pee)
	r.True(pee.Unwrap() == io.EOF, pee)

	var pe PostExecuteError
	r.True(errors.As(wrapped, &pe), wrapped)

	pe = PostExecuteError{}
	r.True(errors.Is(wrapped, pe), wrapped)

	err := errors.Unwrap(pee)
	r.Equal(io.EOF, err)
}

func Test_PostExecute_Errors(t *testing.T) {
	t.Parallel()

	tp := func() *Parser {
		p := testParser(t, "testdata/parser/errors/post_execute")

		p.NodeParsers["foo"] = func(p *Parser, el *Element) (Nodes, error) {
			n := newPostExecuteNode(t, func(ctx context.Context, d *Document, err error) error {
				return fmt.Errorf("boom")
			})
			return Nodes{n}, nil
		}

		return p
	}

	ctx := context.Background()

	type inFn func() error

	tcs := []struct {
		name string
		in   inFn
	}{
		{
			name: "ParseExecuteFile",
			in: func() error {
				p := tp()
				_, err := p.ParseExecuteFile(ctx, "module.md")
				return err
			},
		},
		{
			name: "ParseExecuteFolder",
			in: func() error {
				p := tp()
				_, err := p.ParseExecuteFolder(ctx, "testdata/parser/errors/folder")
				return err
			},
		},
		{
			name: "ParseExecuteFragment",
			in: func() error {
				p := tp()
				_, err := p.ParseExecuteFragment(ctx, strings.NewReader("<foo></foo>"))
				return err
			},
		},
		{
			name: "ParseExecute",
			in: func() error {
				p := tp()
				_, err := p.ParseExecute(ctx, strings.NewReader("<foo></foo>"))
				return err
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			err := tc.in()
			r.Error(err)

			var pe PostExecuteError
			r.True(errors.As(err, &pe), err)

			pe = PostExecuteError{}
			r.True(errors.Is(err, pe), err)

			var ee ExecuteError
			r.True(errors.As(err, &ee), err)

			ee = ExecuteError{}
			r.True(errors.Is(err, ee), err)
		})
	}
}

func Test_PostExecuteError_MarshalJSON(t *testing.T) {
	t.Parallel()

	pee := PostExecuteError{
		Err:      io.EOF,
		OrigErr:  io.ErrClosedPipe,
		Filename: "filename",
		Root:     "root",
		Document: &Document{
			Title: "My Title",
		},
	}

	testJSON(t, "post_execute_error", pee)
}
