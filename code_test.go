package hype

import (
	"testing"

	"github.com/gopherguides/hype/htmltest"
	"github.com/stretchr/testify/require"
)

func Test_NewCode(t *testing.T) {
	t.Parallel()

	inline := NewNode(htmltest.ElementNode(t, "code"))
	inline.Children = Tags{
		&Text{
			Node: NewNode(htmltest.TextNode(t, "hello")),
		},
	}

	src := NewNode(htmltest.AttrNode(t, "code", Attributes{
		"src": "src/main.go",
	}))

	fenced := NewNode(htmltest.AttrNode(t, "code", Attributes{
		"class": "language-go",
	}))

	table := []struct {
		err  bool
		lang string
		name string
		node *Node
	}{
		{name: "nil", err: true},

		{name: "nil html node", node: &Node{}, err: true},
		{name: "non code node", node: NewNode(htmltest.ElementNode(t, "p")), err: true},
		{name: "valid inline", node: inline, lang: ""},
		{name: "valid src", lang: "go", node: src},
		{name: "valid fenced", lang: "go", node: fenced},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			p := testParser(t, testdata)
			c, err := p.NewCode(tt.node)

			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(c)
			r.Equal(tt.lang, c.Lang())

		})
	}

}

func Test_Parse_Code(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	p := testParser(t, testdata)

	doc, err := p.ParseFile("code.md")
	r.NoError(err)

	exp := `<html><head></head><body>
<page number="1">

<h1>Code Test</h1>

<p>This is <code>inline</code> code.</p>

<p>Fenced code block:</p>

<pre><code class="language-sh" language="sh">$ echo hi</code></pre>

<p>A src file:</p>

<p><pre><code class="language-go" language="go" src="src/main.go">package main

// snippet: main
func main() {
	// snippet: main
}</code></pre></p>


</page>


</body>
</html>`

	// fmt.Println(doc.String())
	r.Equal(exp, doc.String())

}
