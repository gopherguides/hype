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

func Test_ExecuteError(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	oce := ExecuteError{
		Err: io.EOF,
	}

	wrapped := fmt.Errorf("error: %w", oce)

	r.True(oce.As(&ExecuteError{}), oce)
	r.True(oce.Is(oce), oce)
	r.True(oce.Unwrap() == io.EOF, oce)

	var ce ExecuteError
	r.True(errors.As(wrapped, &ce), wrapped)

	ce = ExecuteError{}
	r.True(errors.Is(wrapped, ce), wrapped)

	err := errors.Unwrap(oce)
	r.Equal(io.EOF, err)
}

func Test_Execute_Errors(t *testing.T) {
	t.Parallel()

	tp := func() *Parser {
		p := testParser(t, "testdata/parser/errors/execute")

		p.NodeParsers["foo"] = func(p *Parser, el *Element) (Nodes, error) {
			n := newExecuteNode(t, func(ctx context.Context, d *Document) error {
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
			name: "ParseExecute",
			in: func() error {
				p := tp()
				_, err := p.ParseExecute(ctx, strings.NewReader("<foo></foo>"))
				return err
			},
		},
		{
			name: "ParseExecuteFile",
			in: func() error {
				p := tp()
				_, err := p.ParseExecuteFile(ctx, "hype.md")
				return err
			},
		},
		{
			name: "Document.Execute",
			in: func() error {
				p := tp()
				d, err := p.ParseFile("hype.md")
				if err != nil {
					return err
				}

				return d.Execute(ctx)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			err := tc.in()
			r.Error(err)

			var ce ExecuteError
			r.True(errors.As(err, &ce), err)
			r.NotNil(ce.Document, err)

			ce = ExecuteError{}
			r.True(errors.Is(err, ce), err)
		})
	}
}

func Test_ExecuteError_MarshalJSON(t *testing.T) {
	t.Parallel()

	ee := ExecuteError{
		Err:      io.EOF,
		Filename: "hype.md",
		Root:     "testdata/parser/errors/execute",
		Document: &Document{
			Title: "My Title",
		},
		Contents: []byte("foo"),
	}

	testJSON(t, "execute_error", ee)
}

func Test_ExecuteError_String(t *testing.T) {
	t.Parallel()

	detailed := ExecuteError{
		Err:      io.EOF,
		Filename: "hype.md",
		Root:     "testdata/parser/errors/execute",
		Document: &Document{
			Title: "My Title",
		},
		Contents: []byte("foo"),
	}

	tcs := []struct {
		name string
		in   ExecuteError
		exp  string
	}{
		{
			name: "detailed",
			in:   detailed,
			exp:  "filepath: testdata/parser/errors/execute/hype.md\ndocument: My Title\nexecute error: EOF",
		},
		{
			name: "basic",
			in:   ExecuteError{Err: io.EOF},
			exp:  "execute error: EOF",
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			act := tc.in.Error()

			r.Equal(tc.exp, act)
		})
	}
}
