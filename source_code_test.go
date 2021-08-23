package hype

import (
	"testing"

	"github.com/gopherguides/hype/htmltest"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_Parser_NewSourceCode(t *testing.T) {
	t.Parallel()

	table := []struct {
		err  bool
		exp  string
		lang string
		name string
		node *html.Node
	}{
		{name: "nil", err: true},
		{name: "non code node", node: htmltest.ElementNode(t, "p"), err: true},
		{name: "no src attr", node: htmltest.ElementNode(t, "code"), err: true},
		{
			name: "src file missing",
			node: htmltest.AttrNode(t, "code", Attributes{"src": "404.go"}),
			err:  true,
		},
		{
			name: "valid",
			lang: "go",
			node: htmltest.AttrNode(t, "code", Attributes{"src": "src/main.go"}),
			exp:  "<pre><code class=\"language-go\" language=\"go\" src=\"src/main.go\">package main\n\n// snippet: main\nfunc main() {\n\t// snippet: main\n}</code></pre>",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			p := testParser(t, testdata)

			sc, err := p.NewSourceCode(NewNode(tt.node))
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(sc)

			r.Equal(tt.lang, sc.Lang())
			r.Equal(tt.exp, sc.String())
		})
	}

}
