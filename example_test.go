package hype

import (
	"fmt"
	"log"
	"os"

	"github.com/gopherguides/hype/atomx"
)

func ExampleParser() {
	p, err := NewParser(os.DirFS("testdata"))
	if err != nil {
		log.Fatal(err)
	}

	doc, err := p.ParseFile("code.md")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(doc.String())

	// output:
	// <html><head></head><body>
	// <page>

	// <h1>Code Test</h1>

	// <p>This is <code>inline</code> code.</p>

	// <p>Fenced code block:</p>

	// <pre><code class="language-sh" language="sh">$ echo hi</code></pre>

	// <p>A src file:</p>

	// <p><pre class="code-block"><code class="language-go" language="go" snippet="main" src="src/main.go">func main() {</code></pre></p>

	// </page><!--BREAK-->

	// </body>
	// </html>
}

func ExampleParser_CustomTag() {
	p, err := NewParser(os.DirFS("testdata"))
	if err != nil {
		log.Fatal(err)
	}

	p.SetCustomTag("custom", func(node *Node) (Tag, error) {
		el := &Element{
			Node: node,
		}

		el.Children = Tags{QuickText("Hello, World!")}
		return el, nil
	})

	in := `
# Custom Tag

<custom>Stuff inside<custom>
`
	doc, err := p.ParseMD([]byte(in))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(doc.String())

	// output:
	// <html><head></head><body>
	// <page>

	// <h1>Custom Tag</h1>

	// <custom>Hello, World!</custom>
	// </page><!--BREAK-->

	// </body>
	// </html>
}

func ExampleByType() {
	p, err := NewParser(os.DirFS("testdata"))
	if err != nil {
		log.Fatal(err)
	}

	doc, err := p.ParseFile("files.md")
	if err != nil {
		log.Fatal(err)
	}

	files := ByType(doc.Children, &File{})
	for _, file := range files {
		fmt.Println(file.String())
	}

	// output:
	// <file src="src/main.go">src/main.go</file>
	// <file src="src/snip.html">src/snip.html</file>
}

func ExampleByAtom() {
	p, err := NewParser(os.DirFS("testdata"))
	if err != nil {
		log.Fatal(err)
	}

	doc, err := p.ParseFile("files.md")
	if err != nil {
		log.Fatal(err)
	}

	files := ByAtom(doc.Children, atomx.File)
	for _, file := range files {
		fmt.Println(file.String())
	}

	// output:
	// <file src="src/main.go">src/main.go</file>
	// <file src="src/snip.html">src/snip.html</file>
}

func ExampleByAttrs() {
	p, err := NewParser(os.DirFS("testdata"))
	if err != nil {
		log.Fatal(err)
	}

	doc, err := p.ParseFile("files.md")
	if err != nil {
		log.Fatal(err)
	}

	results := ByAttrs(doc.Children, Attributes{
		"src": "src/main.go",
	})
	for _, res := range results {
		fmt.Println(res.String())
	}

	// output:
	// <file src="src/main.go">src/main.go</file>
}
