package hype

import (
	"testing"
	"testing/fstest"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_Parser_NewText(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, fstest.MapFS{})

	table := []*html.Node{
		nil,
		htmx.ElementNode("div"),
	}
	for _, node := range table {
		_, err := p.NewText(node)
		r.Error(err)
	}

	node := htmx.TextNode("hello")

	text, err := p.NewText(node)
	r.NoError(err)
	r.NotNil(text)
	r.Equal(`hello`, text.String())
}
