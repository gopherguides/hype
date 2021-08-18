package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html/atom"
)

func Test_Parser_ParseHTML(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	doc, err := p.ParseFile("html5.html")
	r.NoError(err)
	r.NotNil(doc)

	r.Len(doc.Children, 2)

	dt, ok := doc.Children[0].(*DocType)
	r.True(ok)
	r.Equal("html5", dt.Data)

	html, ok := doc.Children[1].(*Element)
	r.True(ok)
	r.Equal(atom.Html, html.DataAtom)

	r.Len(html.Children, 3)

	head := html.Children[0]
	r.NotNil(head)
	r.Equal(atom.Head, head.DaNode().DataAtom)

	r.Len(head.GetChildren(), 29)

	body, err := doc.Body()
	r.NoError(err)
	r.NotNil(body)

	r.Len(body.Children, 7)
}

func Test_Parser_ParseMD(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, week01)

	doc, err := p.ParseFile("module.md")
	r.NoError(err)
	r.NotNil(doc)

	r.Len(doc.Children, 1)

	html, ok := doc.Children[0].(*Element)
	r.True(ok)
	r.Equal(atom.Html, html.DataAtom)

	r.Len(html.Children, 2)

	head := html.Children[0]
	r.NotNil(head)
	r.Equal(atom.Head, head.DaNode().DataAtom)

	r.Len(head.GetChildren(), 0)

	body, err := doc.Body()
	r.NoError(err)
	r.NotNil(body)

	r.Len(body.Children, 54)

	act := doc.String()
	r.Contains(act, "Basics of Running a Go Program")
}
