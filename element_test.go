package hype

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

func Test_Element_String(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	bd := htmx.AttrNode("div", Attributes{
		"id": "main",
	})

	el := &Element{
		Node: NewNode(bd),
	}
	el.Children = append(el.Children, &Text{
		Node: NewNode(htmx.TextNode("hi")),
	})

	exp := `<div id="main">hi</div>`
	act := el.String()
	r.Equal(exp, act)

}

func Test_Element_JSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	bd := htmx.AttrNode("div", Attributes{
		"id": "main",
	})
	bd.FirstChild = htmx.TextNode("hi")

	el := &Element{
		Node: NewNode(bd),
	}

	b, err := json.Marshal(el)
	r.NoError(err)

	exp := `{"atom":"div","attributes":{"id":"main"},"children":[{"data":"hi","type":"text"}],"data":"div","type":"element"}`

	act := string(b)

	r.Equal(exp, act)
}

var _ Tag = &customTag{}

type customTag struct {
	*Node
}

func (c customTag) String() string {
	return fmt.Sprintf("before\n%s\nafter", c.Children.String())
}

func Test_Parser_ElementNode_Custom(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)
	p.SetCustomTag("foo", func(node *Node) (Tag, error) {
		return customTag{Node: node}, nil
	})

	doc, err := p.ParseFile("assignment.md")
	r.NoError(err)

	act := doc.String()
	// fmt.Println(act)
	exp := "<html><head><meta charset=\"utf-8\" /></head><body>\n<page>\n\n<h1>Assignment 42</h1>\n\n<assignment number=\"42\">\n\n<p>Instructions!</p>\n\n</assignment>\n\n</page><!--BREAK-->\n\n\n</body>\n</html>"

	r.Equal(exp, act)
}
