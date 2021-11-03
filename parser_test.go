package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
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
	r.True(IsAtom(dt, "html5"))

	html, ok := doc.Children[1].(*Element)
	r.True(ok)

	r.True(IsAtom(html, "html"))

	r.Len(html.Children, 3)

	head := html.Children[0]
	r.NotNil(head)

	r.True(IsAtom(head, "head"))

	r.Len(head.GetChildren(), 29)

	body, err := doc.Body()
	r.NoError(err)
	r.NotNil(body)

	r.Len(body.Children, 13)
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
	r.True(IsAtom(html, "html"))

	r.Len(html.Children, 2)

	head := html.Children[0]
	r.NotNil(head)
	r.True(IsAtom(head, "head"))

	r.Len(head.GetChildren(), 0)

	body, err := doc.Body()
	r.NoError(err)
	r.NotNil(body)

	act := doc.String()
	// fmt.Println(act)

	r.Len(body.Children, 20)

	r.Contains(act, "Basics of Running a Go Program")
	r.Contains(act, "// 9 characters (including the space and comma)")
}

func Test_Parser_CustomTag(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)
	p.SetCustomTag("newman", func(node *Node) (Tag, error) {
		el := &Element{
			Node: node,
		}
		return el, el.Validate()
	})

	p.SetCustomTag("leo:uncle", func(node *Node) (Tag, error) {
		el := &Element{
			Node: node,
		}
		return el, el.Validate()
	})

	doc, err := p.ParseFile("custom_tags.md")
	r.NoError(err)
	r.NotNil(doc)

	newmans := doc.Children.ByAtom("newman")
	r.Len(newmans, 1)
	r.Equal(Atom("newman"), newmans[0].Atom())

	leos := doc.Children.ByAtom("leo:uncle")
	r.Len(leos, 1)
	r.Equal(Atom("leo:uncle"), leos[0].Atom())

	leos = doc.Children.ByAtom("leo")
	r.Len(leos, 0)
}
