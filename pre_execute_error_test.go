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

func Test_PreExecuteError(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	pee := PreExecuteError{
		Err: io.EOF,
	}

	wrapped := fmt.Errorf("error: %w", pee)

	r.True(pee.As(&PreExecuteError{}), pee)
	r.True(pee.Is(pee), pee)
	r.True(pee.Unwrap() == io.EOF, pee)

	var pe PreExecuteError
	r.True(errors.As(wrapped, &pe), wrapped)

	pe = PreExecuteError{}
	r.True(errors.Is(wrapped, pe), wrapped)

	err := errors.Unwrap(pee)
	r.Equal(io.EOF, err)
}

func Test_PreExecute_Errors(t *testing.T) {
	t.Parallel()

	tp := func() *Parser {
		p := testParser(t, "testdata/parser/errors/pre_execute")

		p.NodeParsers["foo"] = func(p *Parser, el *Element) (Nodes, error) {
			n := newPreExecuteNode(t, func(ctx context.Context, d *Document) error {
				return fmt.Errorf("boom")
			})

			return Nodes{n}, nil
		}

		return p
	}

	type inFn func() error

	ctx := context.Background()

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
			name: "ParseExecute",
			in: func() error {
				p := tp()
				_, err := p.ParseExecute(ctx, strings.NewReader("<foo></foo>"))
				return err
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			err := tc.in()
			r.Error(err)

			var pee PreExecuteError
			r.True(errors.As(err, &pee), err)

			pee = PreExecuteError{}
			r.True(errors.Is(err, pee), err)

			var ee ExecuteError
			r.True(errors.As(err, &ee), err)

			ee = ExecuteError{}
			r.True(errors.Is(err, ee), err)

		})
	}
}

func Test_PreExecuteError_MarshalJSON(t *testing.T) {
	t.Parallel()

	pee := PreExecuteError{
		Err:      io.EOF,
		Filename: "module.md",
		Root:     "root",
	}

	testJSON(t, "pre_execute_error", pee)
}
