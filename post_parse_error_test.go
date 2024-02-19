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

func Test_PostParseError(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	pee := PostParseError{
		Err: io.EOF,
	}

	wrapped := fmt.Errorf("error: %w", pee)

	r.True(pee.As(&PostParseError{}), pee)
	r.True(pee.Is(pee), pee)
	r.True(pee.Unwrap() == io.EOF, pee)

	var pe PostParseError
	r.True(errors.As(wrapped, &pe), wrapped)

	pe = PostParseError{}
	r.True(errors.Is(wrapped, pe), wrapped)

	err := errors.Unwrap(pee)
	r.Equal(io.EOF, err)
}

func Test_PostParser_Errors(t *testing.T) {
	t.Parallel()

	tp := func() *Parser {
		p := testParser(t, "testdata/parser/errors/post_parse")

		p.NodeParsers["foo"] = func(p *Parser, el *Element) (Nodes, error) {
			n := newPostParserNode(t, func(p *Parser, d *Document, err error) error {
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
			name: "ParseFile",
			in: func() error {
				p := tp()
				_, err := p.ParseFile("module.md")
				return err
			},
		},
		{
			name: "ParseExecuteFile",
			in: func() error {
				p := tp()
				_, err := p.ParseExecuteFile(ctx, "module.md")
				return err
			},
		},
		{
			name: "ParseFragment",
			in: func() error {
				p := tp()
				_, err := p.ParseFragment(strings.NewReader(`<foo></foo>`))
				return err
			},
		},
		{
			name: "Parse",
			in: func() error {
				p := tp()
				_, err := p.Parse(strings.NewReader("<foo></foo>"))
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

			var ppe PostParseError
			r.True(errors.As(err, &ppe), err)

			ppe = PostParseError{}
			r.True(errors.Is(err, ppe), err)

			var pe ParseError
			r.True(errors.As(err, &pe), err)

			pe = ParseError{}
			r.True(errors.Is(err, pe), err)

		})
	}
}
