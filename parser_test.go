package hype

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_Parser_ParseFolder(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/parser/folder"
	cab := os.DirFS(root)

	p := NewParser(cab)
	p.Root = root

	docs, err := p.ParseExecuteFolder(context.Background(), root)
	r.NoError(err)

	r.Len(docs, 3)

	exp := `var Canceled = errors.New`

	titles := []string{"ONE", "TWO", "THREE"}

	for i, doc := range docs {
		r.Equal(titles[i], doc.Title)
		act := doc.String()
		r.Contains(act, exp)
	}

}

func Test_Parser_ParseHTMLNode_Error(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		node *html.Node
	}{
		{
			name: "error node",
			node: &html.Node{
				Type: html.ErrorNode,
				Data: "boom",
			},
		},
		{
			name: "unknown node",
			node: &html.Node{
				Type: 42,
				Data: "boom",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			p := testParser(t, "testdata/auto/parser/hello")
			_, err := p.ParseHTMLNode(tc.node, nil)
			r.Error(err)

			var pe ParseError
			r.True(errors.As(err, &pe), err)

			pe = ParseError{}
			r.True(errors.Is(err, pe), err)
		})
	}

}

func Test_Parser_ParseFragment_Error(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/auto/parser/hello")

	_, err := p.ParseFragment(strings.NewReader(`<include`))
	r.Error(err)

	var pe ParseError
	r.True(errors.As(err, &pe), err)

	pe = ParseError{}
	r.True(errors.Is(err, pe), err)
}

func Test_Parser_ParseExecuteFragment_Error(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/auto/parser/hello")
	p.NodeParsers["foo"] = func(p *Parser, el *Element) (Nodes, error) {
		n := newExecuteNode(t, func(ctx context.Context, d *Document) error {
			return fmt.Errorf("boom")
		})
		return Nodes{n}, nil
	}

	ctx := context.Background()
	_, err := p.ParseExecuteFragment(ctx, strings.NewReader(`<include`))
	r.Error(err)

	var pe ParseError
	r.True(errors.As(err, &pe), err)

	pe = ParseError{}
	r.True(errors.Is(err, pe), err)

	_, err = p.ParseExecuteFragment(ctx, strings.NewReader(`<foo></foo>`))
	r.Error(err)

	var ee ExecuteError
	r.True(errors.As(err, &ee), err)

	ee = ExecuteError{}
	r.True(errors.Is(err, ee), err)
}

func Test_Parser_MarshalJSON(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	p := testParser(t, "testdata/auto/snippets/simple")
	p.DisablePages = true
	p.Section = 42

	err := p.Vars.Set("foo", "bar")
	r.NoError(err)

	_, err = p.ParseFile("module.md")
	r.NoError(err)

	testJSON(t, "parser", p)
}
