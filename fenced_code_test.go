package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parser_NewFencedCode(t *testing.T) {
	t.Parallel()

	valid := NewNode(AttrNode(t, "code", Attributes{"class": "language-go"}))
	valid.Children = Tags{
		&Text{
			Node: NewNode(TextNode(t, "hello")),
		},
	}

	table := []struct {
		err  bool
		exp  string
		name string
		node *Node
	}{
		{name: "nil", err: true},
		{name: "nil html node", node: &Node{}, err: true},
		{name: "non code node", node: NewNode(ElementNode(t, "p")), err: true},
		{name: "valid", node: valid, exp: `<code class="language-go" language="go">hello</code>`},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			p := testParser(t, testdata)

			sc, err := p.NewFencedCode(tt.node)
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(sc)

			r.Equal(tt.exp, sc.String())
		})
	}

}
