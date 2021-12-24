package hype

import (
	"encoding/json"
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

func Test_Parser_NewComment(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	node := htmx.CommentNode("ssh")

	c, err := NewComment(node)
	r.NoError(err)

	r.Equal(`<!-- ssh -->`, c.String())

	_, err = NewComment(nil)
	r.Error(err)
	_, err = NewComment(htmx.TextNode("hello"))
	r.Error(err)
}

func Test_Comment_JSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	c := &Comment{
		Node: NewNode(htmx.CommentNode("ssh")),
	}

	b, err := json.Marshal(c)
	r.NoError(err)

	exp := `{"data":"ssh","type":"comment"}`
	act := string(b)

	r.Equal(exp, act)

}
