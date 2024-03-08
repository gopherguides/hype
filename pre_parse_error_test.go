package hype

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PreParseError(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	ppe := PreParseError{
		Err: io.EOF,
	}

	r.True(ppe.As(&PreParseError{}), ppe)
	r.True(ppe.Is(ppe), ppe)
	r.True(ppe.Unwrap() == io.EOF, ppe)

	var pe PreParseError
	r.True(ppe.As(&pe), ppe)

	pe = PreParseError{}
	r.True(ppe.Is(pe), ppe)

	err := errors.Unwrap(ppe)
	r.Equal(io.EOF, err)
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

	ctx := context.Background()

	tcs := []struct {
		name string
		in   inFn
	}{
		{
			name: "ParseFile",
			in: func() error {
				_, err := tp().ParseFile("hype.md")
				return err
			},
		},
		{
			name: "ParseExecuteFile",
			in: func() error {
				_, err := tp().ParseExecuteFile(ctx, "hype.md")
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

func Test_PreParseError_MarshalJSON(t *testing.T) {
	t.Parallel()

	ppe := PreParseError{
		Err:      io.EOF,
		Filename: "hype.md",
		Root:     "root",
	}

	testJSON(t, "pre_parse_error", ppe)
}

func Test_PreParseError_Error(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	ppe := PreParseError{
		Err:      io.EOF,
		Filename: "hype.md",
		Root:     "root",
	}

	act := ppe.Error()
	act = strings.TrimSpace(act)

	// fmt.Println(act)
	exp := "filepath: root/hype.md\nerror: EOF"

	r.Equal(exp, act)
}
