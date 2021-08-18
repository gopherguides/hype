package hype

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_Parser_NewComment(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, fstest.MapFS{})

	node := CommentNode(t, "ssh")

	c, err := p.NewComment(node)
	r.NoError(err)

	r.Equal(`<!-- ssh -->`, c.String())

	_, err = p.NewComment(nil)
	r.Error(err)
	_, err = p.NewComment(TextNode(t, "hello"))
	r.Error(err)
}
