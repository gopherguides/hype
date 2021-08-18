package hype

import (
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
			lang: "go", err: true,
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
	r := require.New(t)

	cn := ElementNode(t, "code")

	p := testParser(t, testdata)

	c, err := p.NewCode(NewNode(cn))
	r.NoError(err)
	r.Equal("plain", c.Lang())

	r.Equal("<pre>\n<code language=\"plain\">\n</code></pre>\n", c.String())

	tn := TextNode(t, "hello")
	c.Children = append(c.Children, &Text{Node: NewNode(tn)})

	r.Equal("<pre>\n<code language=\"plain\">\nhello</code></pre>\n", c.String())
}
