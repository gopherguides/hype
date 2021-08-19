package hype

import (
	"io/fs"
	"os"
	"testing"

	"github.com/markbates/fsx"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	testdata = os.DirFS("testdata")
	week01   = os.DirFS("testdata/week01")
)

func testParser(t *testing.T, cab fs.FS) *Parser {
	t.Helper()

	r := require.New(t)
	p, err := NewParser(fsx.NewFS(cab))
	r.NoError(err)
	return p
}

func DocNode(t *testing.T) *html.Node {
	t.Helper()
	return &html.Node{
		Type: html.DocumentNode,
	}
}

func DocTypeNode(t *testing.T, value string) *html.Node {
	t.Helper()
	return &html.Node{
		Type: html.DoctypeNode,
		Data: value,
	}
}

func ElementNode(t *testing.T, name string) *html.Node {
	t.Helper()

	return &html.Node{
		Data:     name,
		DataAtom: atom.Lookup([]byte(name)),
		Type:     html.ElementNode,
	}
}

func AttrNode(t *testing.T, name string, ats Attributes) *html.Node {
	t.Helper()
	node := ElementNode(t, name)
	node.Attr = ats.Attrs()
	return node
}

func TextNode(t *testing.T, text string) *html.Node {
	t.Helper()
	return &html.Node{
		Data: text,
		Type: html.TextNode,
	}
}

func CommentNode(t *testing.T, text string) *html.Node {
	t.Helper()
	return &html.Node{
		Type: html.CommentNode,
		Data: text,
	}
}
