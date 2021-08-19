package hype

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_NewCode(t *testing.T) {
	t.Parallel()

	table := []struct {
		err  bool
		lang string
		name string
		node *html.Node
	}{
		{name: "nil", err: true},
		{name: "not code node", node: ElementNode(t, "p"), err: true},
		{name: "no lang", node: ElementNode(t, "code"), lang: "plain"},
		{
			name: "no lang, with src", node: AttrNode(t, "code", map[string]string{
				"src": "src/main.go",
			}),
			lang: "go",
		},
		{
			name: "missing src file", node: AttrNode(t, "code", map[string]string{
				"src": "404.go",
			}),
			err: true,
		},
		{
			name: "langauge attr", node: AttrNode(t, "code", map[string]string{
				"language": "bash",
			}),
			lang: "bash",
		},
		{
			name: "class lang-foo attr", node: AttrNode(t, "code", map[string]string{
				"class": "language-bash",
			}),
			lang: "bash",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			p := testParser(t, testdata)
			c, err := p.NewCode(NewNode(tt.node))

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

func Test_Code_String(t *testing.T) {
	t.Parallel()

	const data = `package main`

	table := []struct {
		err  bool
		exp  string
		name string
		node *html.Node
	}{
		{name: "nil", err: true},
		{name: "not code node", node: ElementNode(t, "p"), err: true},
		{
			name: "md node",
			exp:  `<pre><code class="language-go" language="go">package main</code></pre>`,
			node: AttrNode(t, "code", Attributes{
				"class": "language-go",
			}),
		},
		{
			name: "src node",
			exp:  "<pre><code class=\"language-go\" language=\"go\" src=\"src/main.go\">package main\n\nfunc main() {\n\n}</code></pre>",
			node: AttrNode(t, "code", Attributes{
				"src": "src/main.go",
			}),
		},
		{
			name: "inline node",
			exp:  `<code>package main</code>`,
			node: ElementNode(t, "code"),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			node := NewNode(tt.node)
			if tt.node != nil {
				node.Children = append(node.Children, &Text{
					Node: NewNode(TextNode(t, data)),
				})
			}

			p := testParser(t, testdata)

			c, err := p.NewCode(node)
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(c)

			r.Equal(tt.exp, c.String())
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
<h1>Code Test</h1>

<p>This is <code>inline</code> code.</p>

<p>Fenced code block:</p>

<pre><code class="language-sh" language="sh">$ echo hi</code></pre>

<p>A src file:</p>

<p><pre><code class="language-go" language="go" src="src/main.go">package main

func main() {

}</code></pre></p>

</body>
</html>`

	fmt.Println(doc.String())
	r.Equal(exp, doc.String())

}
