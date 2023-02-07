package hype

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Parser_UnknownAtom(t *testing.T) {
	// t.Skip()
	t.Parallel()
	r := require.New(t)

	in := `# Hello

<notes>

<figure id="x">
<go doc="io.EOF"></go>
<figcaption>Docs for <godoc>io#EOF</godoc></figcaption>
</figure>


</notes>`

	p := NewParser(nil)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	doc, err := p.ParseExecute(ctx, strings.NewReader(in))
	r.NoError(err)

	exp := `<a for="io#EOF" href="https://pkg.go.dev/io#EOF" target="_blank"><code>io.EOF</code></a>`

	act := doc.String()
	act = strings.TrimSpace(act)

	r.Contains(act, exp)
}

func Test_Parser(t *testing.T) {
	t.Parallel()

	// t.Skip()
	r := require.New(t)

	cab := helloCab(t, "second")

	modmd := []byte(`# Page 1

This is ` + "`inline`" + ` code.

<include src="second/second.md"></include>

<cmd exec="echo hello"></cmd>

more words`)

	second_md := []byte(`# Second Page

<img src="assets/second.png" />

<code src="src/main.go"></code>

`)

	cab["module.md"] = &fstest.MapFile{
		Data: modmd,
	}

	cab["second/second.md"] = &fstest.MapFile{
		Data: second_md,
	}

	p := NewParser(cab)

	doc, err := p.ParseExecuteFile(context.Background(), "module.md")
	r.NoError(err)
	r.NotNil(doc)

	exp := `<html><head></head><body><page>
<h1>Page 1</h1>

<p>This is <code>inline</code> code.</p>
</page>
<page>
<h1>Second Page</h1>

<img src="second/assets/second.png"></img>

<pre><code class="language-go" language="go" src="second/src/main.go">package main

import &#34;fmt&#34;

func main() {
	fmt.Println(&#34;Hello second!&#34;)
}</code></pre>
</page>

<page>
<cmd exec="echo hello"><pre><code class="language-text" language="text">$ echo hello

hello</code></pre></cmd>

<p>more words</p>
</page>
</body></html>`

	act := doc.String()

	// fmt.Println(act)
	r.Equal(exp, act)
}

func helloCab(t testing.TB, names ...string) fstest.MapFS {
	t.Helper()

	cab := fstest.MapFS{}

	fn := func(name string) {
		cab[fmt.Sprintf("%s/src/go.mod", name)] = &fstest.MapFile{
			Data: []byte(fmt.Sprintf(go_mod, name)),
		}
		cab[fmt.Sprintf("%s/src/main.go", name)] = &fstest.MapFile{
			Data: []byte(fmt.Sprintf(main_go, name)),
		}
		cab[fmt.Sprintf("%[1]s/assets/%[1]s.png", name)] = &fstest.MapFile{
			Data: []byte(fmt.Sprintf("%s.png", name)),
		}
	}

	for _, name := range names {
		fn(name)
	}

	return cab
}

const go_mod = `module %s

go 1.18`

const main_go = `package main

import "fmt"

func main() {
	fmt.Println("Hello %s!")
}`

func Test_Parser_ParseFolder(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/parser/folder"
	cab := os.DirFS(root)

	p := NewParser(cab)
	p.Root = root

	docs, err := p.ParseExecuteFolder(context.Background(), root)
	r.NoError(err)

	r.Len(docs, 3)

	exp := `var Canceled = errors.New`

	titles := []string{"ONE", "TWO", "THREE"}

	for i, doc := range docs {
		r.Equal(titles[i], doc.Title)
		act := doc.String()
		r.Contains(act, exp)
	}

}
