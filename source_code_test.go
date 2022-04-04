package hype

import (
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_Parser_NewSourceCode(t *testing.T) {
	t.Parallel()

	p := testParser(t, testdata)

	table := []struct {
		finalErr bool
		parseErr bool
		exp      string
		lang     string
		name     string
		node     *html.Node
	}{
		{name: "nil", parseErr: true},
		{name: "non code node", node: htmx.ElementNode("p"), parseErr: true},
		{name: "no src attr", node: htmx.ElementNode("code"), parseErr: true},
		{
			name:     "src file missing",
			node:     htmx.AttrNode("code", Attributes{"src": "404.go"}),
			finalErr: true,
		},
		{
			name: "valid Go",
			lang: "go",
			node: htmx.AttrNode("code", Attributes{"src": "src/main.go"}),
			exp:  "<p><pre><code class=\"language-go\" language=\"go\" src=\"src/main.go\">package main\n\nfunc main() {\n}\n</code></pre></p>",
		},
		{
			name: "valid snippet",
			lang: "go",
			node: htmx.AttrNode("code", Attributes{"src": "src/snippets.go", "snippet": "hello"}),
			exp:  "<p><pre><code class=\"language-go\" language=\"go\" snippet=\"hello\" src=\"src/snippets.go\">func Hello() {\n\tfmt.Println(&#34;Hello, World!&#34;)\n}\n</code></pre></p>",
		},
		{
			name: "valid HTML",
			lang: "html",
			node: htmx.AttrNode("code", Attributes{"src": "src/snip.html"}),
			exp:  "<p><pre><code class=\"language-html\" language=\"html\" src=\"src/snip.html\">&lt;!doctype html5&gt;\n&lt;html lang=&#34;en&#34;&gt;\n\n&lt;head&gt;&lt;/head&gt;\n\n&lt;body&gt;\n\n  &lt;!-- your content here... --&gt;\n  &lt;script src=&#34;js/scripts.js&#34;&gt;&lt;/script&gt;\n\n  &lt;div class=&#34;text&#34;&gt;\n    &lt;img src=&#34;assets/foo.png&#34; width=&#34;100%&#34;&gt;\n    &lt;!-- snippet: main --&gt;\n    &lt;p&gt;Hello World&lt;/p&gt;\n    &lt;!-- snippet: main --&gt;\n  &lt;/div&gt;\n\n&lt;/body&gt;\n\n&lt;/html&gt;</code></pre></p>",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			sc, err := NewSourceCode(testdata, NewNode(tt.node), nil)
			if tt.parseErr {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(sc)

			err = sc.Finalize(p)
			if tt.finalErr {
				r.Error(err)
				return
			}

			r.Equal(tt.lang, sc.Lang())
			r.Equal(tt.exp, sc.String())
		})
	}

}

func Test_SourceCode_HashtagSnippets(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	node := htmx.AttrNode("code", Attributes{
		"src": "src/snippets.go#hello,src/snippets.js#goodbye",
	})

	p := testParser(t, testdata)

	sc, err := NewSourceCode(p, NewNode(node), p.snippetRules)
	r.NoError(err)

	r.NoError(sc.Finalize(p))

	kids := sc.GetChildren()
	r.Len(kids, 2)

	const exp = `<p><pre><code class="language-go#hello" language="go#hello" src="src/snippets.go#hello">func Hello() {
	fmt.Println(&#34;Hello, World!&#34;)
}
</code></pre></p><p><pre><code class="language-js#goodbye" language="js#goodbye" src="src/snippets.js#goodbye">function goodbye() {
    console.log(&#39;Goodbye, World!&#39;);
}</code></pre></p>`

	act := sc.String()

	// fmt.Println(act)

	r.Equal(exp, act)
}
