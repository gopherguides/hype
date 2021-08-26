package hype

import (
	"encoding/json"
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

func Test_Body_JSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	bd := htmx.AttrNode(t, "body", Attributes{
		"id": "main",
	})
	bd.FirstChild = htmx.TextNode(t, "hi")

	body := &Body{
		Node: NewNode(bd),
	}

	b, err := json.Marshal(body)
	r.NoError(err)

	exp := `{"atom":"body","attributes":{"id":"main"},"children":[{"data":"hi","type":"text"}],"data":"body","type":"element"}`

	act := string(b)

	r.Equal(exp, act)
}

func Test_Body_AsPage(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	doc := ParseFile(t, testdata, "html5.html")
	pages := doc.Pages()
	r.Len(pages, 1)

	body, err := doc.Body()
	r.NoError(err)

	exp := body.Children.String()

	page := body.AsPage()

	// page := pages[0]
	act := page.Children.String()
	r.Equal(exp, act)

}
