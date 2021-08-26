package hype

import (
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

func Test_Parser_NewPage(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	doc := ParseFile(t, testdata, "pages.md")
	r.NotNil(doc)

	pages := doc.Pages()
	r.Len(pages, 3)

	const exp = `<html><head></head><body>
<page number="1">

<h1>First H1</h1>

<p>Page 1</p>


</page>

<page number="2">

<h1>Second H1</h1>

<p>Page 2</p>

<h2>H2 under Second H1</h2>

<p>Page 2.A</p>

<h3>H3!</h3>

<p>Page 2.B</p>


</page>

<page number="3">

<h1>Last H1</h1>

<p>Last page</p>


</page>


</body>
</html>`

	r.Equal(exp, doc.String())
}

func Test_Page_Number(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := &Page{
		Node: NewNode(htmx.ElementNode("page")),
	}
	r.Equal(p.Number(), 0)

	p.Set("number", "42")
	r.Equal(p.Number(), 42)
}
