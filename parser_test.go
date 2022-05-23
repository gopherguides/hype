package hype

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

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

import "fmt"

func main() {
	fmt.Println("Hello second!")
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
	compareOutput(t, act, exp)
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

func Test_Full_Module(t *testing.T) {
	t.Parallel()
	t.Skip()
	r := require.New(t)

	cab := os.DirFS("testdata/full")

	p := &Parser{
		FS: cab,
	}

	doc, err := p.ParseExecuteFile(context.Background(), "module.md")
	r.NoError(err)

	act := doc.String()
	act = strings.TrimSpace(act)

	b, err := fs.ReadFile(cab, "output.txt")
	r.NoError(err)

	exp := string(b)
	exp = strings.TrimSpace(exp)

	r.Equal(exp, act)

}
