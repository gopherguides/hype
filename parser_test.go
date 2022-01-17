package hype

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"strings"
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

	r.Len(head.GetChildren(), 1)

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
		return el, el.Validate(p)
	})

	p.SetCustomTag("leo:uncle", func(node *Node) (Tag, error) {
		el := &Element{
			Node: node,
		}
		return el, el.Validate(p)
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

func Test_Parser_PreProcessor(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)
	p.PreProcessor = func(r io.Reader) (io.Reader, error) {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		tmpl, err := template.New("").Parse(string(b))
		bb := &bytes.Buffer{}
		err = tmpl.Execute(bb, map[string]any{
			"Name": "World!",
		})

		return bb, nil
	}

	in := `# Hello {{.Name}}`

	doc, err := p.ParseReader(strings.NewReader(in))
	r.NoError(err)

	act := doc.String()
	act = strings.TrimSpace(act)
	// fmt.Println(act)

	exp := `<html><head><meta charset="utf-8" /></head><body>
# Hello World!
</body>
</html>`
	exp = strings.TrimSpace(exp)

	r.Equal(exp, act)
}
