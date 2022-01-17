package hype

import (
	"io"
	"strings"
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_Parser_NewDocument(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	// snippet: errors
	_, err := p.NewDocument(nil)
	r.Error(err)

	_, err = p.NewDocument(htmx.TextNode(""))
	r.Error(err)
	// snippet: errors

	f, err := testdata.Open("html5.html")
	r.NoError(err)
	defer f.Close()

	n, err := html.Parse(f)
	r.NoError(err)

	doc, err := p.NewDocument(n)
	r.NoError(err)
	r.NotNil(doc)

	r.Len(doc.Children, 2)

	dt, ok := doc.Children[0].(*DocType)
	r.True(ok)

	r.True(IsAtom(dt, "html5"))

	html, ok := doc.Children[1].(*Element)
	r.True(ok)

	r.True(IsAtom(html, "html"))

	r.Len(html.Children, 3)

	head := html.Children[0]
	r.NotNil(head)
	r.True(IsAtom(head, "head"))

	r.Len(head.GetChildren(), 29)

}

func Test_Document_Body(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	var d *Document
	_, err := d.Body()
	r.Error(err)

	doc, err := p.ParseFile("html5.html")
	r.NoError(err)
	r.NotNil(doc)

	body, err := doc.Body()
	r.NoError(err)
	r.NotNil(body)

	r.Len(body.Children, 13)
}

func Test_Document_Meta(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	doc, err := p.ParseFile("big.html")
	r.NoError(err)
	r.NotNil(doc)

	data := doc.Meta()
	r.Len(data, 19)
}

func Test_Document_Overview(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	in := `<overview>hi</overview>`

	p := testParser(t, testdata)
	doc, err := p.ParseReader(io.NopCloser(strings.NewReader(in)))
	r.NoError(err)
	r.NotNil(doc)

	ov := doc.Overview()
	r.Equal("hi", ov)
}

func Test_Document_JSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	doc, err := p.ParseFile("html5.html")
	r.NoError(err)

	b, err := doc.MarshalJSON()
	r.NoError(err)

	act := string(b)

	r.Contains(act, `"document":{"children":[`)
}

func Test_Document_Pages(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		exp  int
	}{
		{name: "pages.md", exp: 4},
		{name: "html5.html", exp: 1},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			doc := ParseFile(t, testdata, tt.name)
			r.NotNil(doc)

			r.Len(doc.Pages(), tt.exp)
		})
	}

}
