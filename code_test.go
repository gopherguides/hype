package hype

import (
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

func Test_NewCode(t *testing.T) {
	t.Parallel()

	inline := NewNode(htmx.ElementNode("code"))
	inline.Children = Tags{
		&Text{
			Node: NewNode(htmx.TextNode("hello")),
		},
	}

	src := NewNode(htmx.AttrNode("code", Attributes{
		"src": "src/main.go",
	}))

	fenced := NewNode(htmx.AttrNode("code", Attributes{
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
		{name: "non code node", node: NewNode(htmx.ElementNode("p")), err: true},
		{name: "valid inline", node: inline, lang: ""},
		{name: "valid src", lang: "go", node: src},
		{name: "valid fenced", lang: "go", node: fenced},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			p := testParser(t, testdata)
			c, err := NewCode(p, tt.node)

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

	exp := `<html><head><meta charset="utf-8" /></head><body>
<page>

<h1>Code Test</h1>

<p>This is <code>inline</code> code.</p>

<p>Fenced code block:</p>

<pre><code class="language-sh" language="sh">$ echo hi</code></pre>

<p>A src file:</p>

<p><pre><code class="language-go" language="go" snippet="main" src="src/main.go">func main() {</code></pre></p>

</page><!--BREAK-->


</body>
</html>`

	// fmt.Println(doc.String())
	r.Equal(exp, doc.String())

}

func Test_Code_MultipleSources(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	node := htmx.AttrNode("code", Attributes{
		"src": "src/snippets.go,src/snippets.js",
	})

	sc, err := NewSourceCode(testdata, NewNode(node), nil)
	r.NoError(err)

	kids := sc.GetChildren()
	r.Len(kids, 2)

	const exp = `<p><pre><code src="src/snippets.go,src/snippets.js"><p><pre><code class="language-go" language="go" src="src/snippets.go">package main

import &#34;fmt&#34;

func Hello() {
	fmt.Println(&#34;Hello, World!&#34;)
}


func Goodbye() {
	fmt.Println(&#34;Goodbye, World!&#34;)
}

</code></pre></p><p><pre><code class="language-js" language="js" src="src/snippets.js">function hello() {
    console.log(&#39;Hello, World!&#39;);
}

function goodbye() {
    console.log(&#39;Goodbye, World!&#39;);
}</code></pre></p></code></pre></p>`

	act := sc.String()

	// fmt.Println(act)

	r.Equal(exp, act)

}

func Test_Code_MultipleSources_Errors(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	node := htmx.AttrNode("code", Attributes{
		"src":     "src/snippets.go,src/snippets.js",
		"snippet": "hello",
	})

	_, err := NewCode(p, NewNode(node))
	r.NoError(err)
}
