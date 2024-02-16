package hype

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PreParsers_PreParse(t *testing.T) {
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

func Test_PreParsers_PreParse_Error(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	pp := PreParsers{
		PreParseFn(func(p *Parser, r io.Reader) (io.Reader, error) {
			return nil, fmt.Errorf("boom")
		}),
	}

	p := testParser(t, "")

	_, err := pp.PreParse(p, strings.NewReader(""))
	r.Error(err)

	var e2 PreParseError

	r.True(errors.As(err, &e2))

	r.Contains(e2.Error(), "boom")
	r.NotNil(e2.PreParser)
}
