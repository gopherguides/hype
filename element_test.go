package hype

import (
	"encoding/json"
	"testing"

	"github.com/gopherguides/hype/htmltest"
	"github.com/stretchr/testify/require"
)

func Test_Element_String(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	bd := htmltest.AttrNode(t, "div", Attributes{
		"id": "main",
	})

	el := &Element{
		Node: NewNode(bd),
	}
	el.Children = append(el.Children, &Text{
		Node: NewNode(htmltest.TextNode(t, "hi")),
	})

	exp := `<div id="main">hi</div>`
	act := el.String()
	r.Equal(exp, act)

}

func Test_Element_JSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	bd := htmltest.AttrNode(t, "div", Attributes{
		"id": "main",
	})
	bd.FirstChild = htmltest.TextNode(t, "hi")

	el := &Element{
		Node: NewNode(bd),
	}

	b, err := json.Marshal(el)
	r.NoError(err)

	exp := `{"atom":"div","attributes":{"id":"main"},"children":[{"data":"hi","type":"text"}],"data":"div","type":"element"}`

	act := string(b)

	r.Equal(exp, act)
}
