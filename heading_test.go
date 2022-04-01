package hype

import (
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

func Test_Heading_Markdown(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	node := htmx.ElementNode("h3")

	h, err := NewHeading(NewNode(node))
	r.NoError(err)

	h.Children = append(h.Children, QuickText("hello"))

	r.Equal("### hello", h.Markdown())

}
