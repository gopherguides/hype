package hype

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PreParsers(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var in io.Reader = strings.NewReader(`<html><body><h1>Hello</h1></body></html>`)

	pp := PreParsers{
		PreParseFn(func(p *Parser, r io.Reader) (io.Reader, error) {
			b, err := io.ReadAll(r)
			if err != nil {
				return nil, err
			}
			return bytes.NewReader(bytes.ToUpper(b)), nil
		}),
		PreParseFn(func(p *Parser, r io.Reader) (io.Reader, error) {
			b, err := io.ReadAll(r)
			if err != nil {
				return nil, err
			}

			b = bytes.ReplaceAll(b, []byte("H1"), []byte("h2"))
			return bytes.NewReader(b), nil
		}),
	}

	in, err := pp.PreParse(testParser(t, ""), in)

	r.NoError(err)

	b, err := io.ReadAll(in)
	r.NoError(err)

	act := string(b)
	exp := `<HTML><BODY><h2>HELLO</h2></BODY></HTML>`

	r.Equal(exp, act)
}

func Test_PreParser_Errors(t *testing.T) {
	t.Parallel()

	const root = "testdata/parser/errors"

	tp := func() *Parser {
		p := testParser(t, filepath.Join(root, "pre_parse"))

		fn := PreParseFn(func(p *Parser, r io.Reader) (io.Reader, error) {
			return nil, fmt.Errorf("boom")
		})

		p.PreParsers = append(p.PreParsers, fn)
		return p
	}

	type inFn func() error

	tcs := []struct {
		name string
		in   inFn
	}{
		{
			name: "ParseFile",
			in: func() error {
				_, err := tp().ParseFile("module.md")
				return err
			},
		},
		{
			name: "ParseExecuteFile",
			in: func() error {
				ctx := context.Background()
				_, err := tp().ParseExecuteFile(ctx, "module.md")
				return err
			},
		},
		{
			name: "Parse",
			in: func() error {
				_, err := tp().Parse(strings.NewReader("hello"))
				return err
			},
		},
		{
			name: "ParseExecute",
			in: func() error {
				ctx := context.Background()
				_, err := tp().ParseExecute(ctx, strings.NewReader("hello"))
				return err
			},
		},
		{
			name: "ParseFragment",
			in: func() error {
				_, err := tp().ParseFragment(strings.NewReader("hello"))
				return err
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			err := tc.in()
			r.Error(err)

			var pe ParseError
			r.True(errors.As(err, &pe), err)

			pe = ParseError{}
			r.True(errors.Is(err, pe), err)

			var ppe PreParseError
			r.True(errors.As(err, &ppe), err)

			ppe = PreParseError{}
			r.True(errors.Is(err, ppe), err)
		})
	}

}
