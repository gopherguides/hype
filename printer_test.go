package hype

import (
	"bytes"
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

func Test_Printer_Doc(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	in := `# Hello

<code src="src/main.go"></code>
	`

	doc := ParseMD(t, testdata, []byte(in))

	bb := &bytes.Buffer{}
	p := NewPrinter(bb)

	err := p.Print(doc.Children...)
	r.NoError(err)

	act := bb.String()
	exp := `<html><head></head><body><page number="1">
<h1>Hello</h1>

<p><code class="language-go" language="go" src="src/main.go">package main

// snippet: main
func main() {
	// snippet: main
}
</code></p>

</page>
</body></html>`

	// fmt.Println(act)
	r.Equal(exp, act)

	bb = &bytes.Buffer{}
	p = NewPrinter(bb)

	p.SetTransformer(func(tag Tag) (Tag, error) {
		sc, ok := tag.(*SourceCode)
		if !ok {
			return tag, nil
		}

		hn := htmx.ElementNode("div")
		el := &Element{
			Node: NewNode(hn),
		}

		tn, err := NewText(htmx.TextNode(sc.String()))
		r.NoError(err)

		el.Children = Tags{
			tn,
		}

		return el, nil
	})

	r.NoError(p.Print(doc.Children...))

	act = bb.String()
	exp = `<html><head></head><body><page number="1">
<h1>Hello</h1>

<p><div><pre><code class="language-go" language="go" src="src/main.go">package main

// snippet: main
func main() {
	// snippet: main
}</code></pre></div></p>

</page>
</body></html>`
	// fmt.Println(act)
	r.Equal(exp, act)
}
