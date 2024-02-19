package hype

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_Parser(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	p := testParser(t, "testdata/parser/hello")

	doc, err := p.ParseExecuteFile(context.Background(), "module.md")
	r.NoError(err)
	r.NotNil(doc)

	exp := `<html><head></head><body><page>
<h1>Page 1</h1>

<p>This is <code>inline</code> code.</p>
</page>
<page>
<h1>Second Page</h1>

<pre><code class="language-go" language="go" src="second/src/main.go">package main

import &#34;fmt&#34;

func main() {
	fmt.Println(&#34;Hello second!&#34;)
}
</code></pre>
</page>

<page>
<cmd exec="echo hello"><pre><code class="language-shell" language="shell">$ echo hello

hello</code></pre></cmd>

<p>more words</p>
</page>
</body></html>`

	act := doc.String()

	// fmt.Println(act)
	r.Equal(exp, act)
}

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

			p := testParser(t, "testdata/parser/hello")
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

	p := testParser(t, "testdata/parser/hello")

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

	p := testParser(t, "testdata/parser/hello")

	ctx := context.Background()
	_, err := p.ParseExecuteFragment(ctx, strings.NewReader(`<include`))
	r.Error(err)

	var pe ParseError
	r.True(errors.As(err, &pe), err)

	pe = ParseError{}
	r.True(errors.Is(err, pe), err)

	_, err = p.ParseExecuteFragment(ctx, strings.NewReader(`<cmd exec="boom"></cmd>`))
	r.Error(err)

	var ee ExecuteError
	r.True(errors.As(err, &ee), err)

	ee = ExecuteError{}
	r.True(errors.Is(err, ee), err)
}
