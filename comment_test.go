package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_Comment(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := &Parser{}

	node, err := p.ParseHTMLNode(
		&html.Node{
			Type: html.CommentNode,
			Data: "hello",
		},
		nil,
	)

	r.NoError(err)

	comment, ok := node.(Comment)
	r.True(ok)
	r.Equal("<!-- hello -->", comment.String())
}
