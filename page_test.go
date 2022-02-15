package hype

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Page_ShiftHeadings(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	doc := ParseFile(t, testdata, "pages.md")
	r.NotNil(doc)

	pages := doc.Pages()
	r.Len(pages, 4)

	page := pages[1]

	start := `<page>

<h1>Second H1</h1>

<p>page 2</p>

<h2>H2 under Second H1</h2>

<p>page 2.A</p>

<h3>H3!</h3>

<p>page 2.B</p>

</page><!--BREAK-->`

	start = strings.TrimSpace(start)
	act := page.String()
	act = strings.TrimSpace(act)

	r.Equal(start, act)

	page.ShiftHeadings(1)

	exp := `<page>

<h2>Second H1</h2>

<p>page 2</p>

<h3>H2 under Second H1</h3>

<p>page 2.A</p>

<h4>H3!</h4>

<p>page 2.B</p>

</page><!--BREAK-->`

	act = page.String()
	act = strings.TrimSpace(act)

	r.Equal(exp, act)
}

func Test_Parser_NewPage(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	doc := ParseFile(t, testdata, "pages.md")
	r.NotNil(doc)

	pages := doc.Pages()
	r.Len(pages, 4)

	const exp = `<html><head><meta charset="utf-8" /></head><body>
<page>

<h1>First H1</h1>

<p>page 1</p>

</page><!--BREAK-->

<page>

<h1>Second H1</h1>

<p>page 2</p>

<h2>H2 under Second H1</h2>

<p>page 2.A</p>

<h3>H3!</h3>

<p>page 2.B</p>

</page><!--BREAK-->

<page>

<h1>Code Test</h1>

<p>This is <code>inline</code> code.</p>

<p>Fenced code block:</p>

<pre><code class="language-sh" language="sh">$ echo hi</code></pre>

<p>A src file:</p>

<p><pre><code class="language-go" language="go" snippet="main" src="src/main.go">func main() {</code></pre></p>

</page><!--BREAK-->


<page>

<h1>Last H1</h1>

<p>Last page</p>

</page><!--BREAK-->


</body>
</html>`

	act := doc.String()
	// fmt.Println(act)
	r.Equal(exp, act)
}
