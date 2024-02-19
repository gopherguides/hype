package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_NewHeading(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	h, err := NewHeading(&Element{
		HTMLNode: &html.Node{
			Type: html.ElementNode,
			Data: "h1",
		},
	})

	r.NoError(err)
	r.Equal("<h1></h1>", h.String())
	r.Equal(1, h.Level())
}

func Test_Heading_MarshalJSON(t *testing.T) {
	t.Parallel()

	h := &Heading{
		Element: NewEl("h1", nil),
		level:   1,
	}
	h.Nodes = append(h.Nodes, Text("This is a heading"))

	testJSON(t, "heading", h)
}
