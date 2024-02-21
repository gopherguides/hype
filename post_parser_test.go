package hype

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type postParserNode struct {
	*Element
	PostParseFn
}

func newPostParserNode(t testing.TB, fn PostParseFn) *postParserNode {
	t.Helper()

	return &postParserNode{
		PostParseFn: fn,
		Element:     &Element{},
	}
}

func Test_PostParsers_PostParse(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	nodes := Nodes{
		newPostParserNode(t, func(p *Parser, d *Document, err error) error {
			d.Title = "Hello"
			return nil
		}),
	}

	d := &Document{}
	err := nodes.PostParse(nil, d, nil)

	r.NoError(err)

	act := d.Title
	exp := "Hello"

	r.Equal(exp, act)

}

func Test_PostParsers_PostParse_Errors(t *testing.T) {
	t.Parallel()

	efn := newPostParserNode(t, func(p *Parser, d *Document, err error) error {
		d.Title = "Hello"
		return fmt.Errorf("boom")
	})

	nofn := newPostParserNode(t, func(p *Parser, d *Document, err error) error {
		d.Title = "Hello"
		return nil
	})

	table := []struct {
		name string
		fn   Node
		exp  string
	}{
		{name: "no error", fn: nofn, exp: "original"},
		{name: "nodes list", fn: Nodes{nofn}, exp: "original"},
		{name: "extra error", fn: efn, exp: "boom; original"},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			nodes := Nodes{
				tc.fn,
			}

			d := &Document{}

			p := testParser(t, "testdata/whole/simple")
			p.Filename = "module.md"

			err := nodes.PostParse(p, d, fmt.Errorf("original"))
			r.Error(err)

			var pperr PostParseError

			r.ErrorAs(err, &pperr)

			r.Equal("Hello", d.Title)
		})
	}
}
