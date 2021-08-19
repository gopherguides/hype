package hype

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Test_Parser_NewDocument(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	// section: errors
	_, err := p.NewDocument(nil)
	r.Error(err)

	_, err = p.NewDocument(TextNode(t, ""))
	r.Error(err)
	// section: errors

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
	r.Equal("html5", dt.Data)

	html, ok := doc.Children[1].(*Element)
	r.True(ok)
	r.Equal("html", html.DataAtom.String())

	r.Len(html.Children, 3)

	head := html.Children[0]
	r.NotNil(head)
	r.Equal(atom.Head, head.Atom())

	r.Len(head.GetChildren(), 29)

}

func Test_Document_Body(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

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

	r.Contains(act, `"fs":{"assets":{`)
	r.Contains(act, `"document":{"children":[`)
}
