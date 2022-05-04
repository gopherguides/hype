package commander

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cmd_Tag(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	os.RemoveAll("~/.hype")
	p := testParser(t, testdata, "testdata")

	doc, err := p.ParseFile("run.md")
	r.NoError(err)
	r.NotNil(doc)

	act := doc.String()

	// fmt.Println(act)
	assertExp(t, "cmd.exp.tmpl", act)
}

func Test_Cmd_Run_Error(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata, "testdata")
	p.FileName = "run.md"

	_, err := p.ParseMD([]byte(`<cmd exec="go run ." src="cmd-err"></cmd>`))

	r.Error(err)

	act := err.Error()
	act = strings.TrimSpace(act)

	// fmt.Println(act)

	exp := `expected exit code 0, got 2:
<cmd exec="go run ." src="cmd-err">`

	r.Contains(act, exp)

	exp = `./main.go:6:14: undefined: missing`
	r.Contains(act, exp)

}

func Test_Cmd_Code(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	dir := "testdata/set-sources"
	cab := os.DirFS(dir)

	p := testParser(t, cab, dir)

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	act := doc.String()

	fmt.Println(act)

	exp := `<html><head><meta charset="utf-8" /></head><body>
<page>

<h1>Hello World</h1>

<img src="assets/foo.png" />

<p><pre><code class="language-go" language="go" src="src/demo.go">package main

import &#34;fmt&#34;

func main() {
	fmt.Println(&#34;Hello, World!&#34;)
}

</code></pre></p>

</page><!--BREAK-->

<page>

<h1>Foo</h1>

<p><pre><code class="language-js" language="js" src="foo/src/foo.js">function hello() {
    console.log(&#39;Hello, World!&#39;);
}

function goodbye() {
    console.log(&#39;Goodbye, World!&#39;);
}</code></pre></p>

<p>Some text</p>

<div><p><pre><code class="language-go#example" code="main.go#example" exec="go run main.go" hide-duration="" language="go#example" src="foo/src/bar/main.go#example">func main() {
	fmt.Println(&#34;Hello, World!&#34;)
}
</code></pre></p><cmd code="main.go#example" exec="go run main.go" hide-duration="" src="foo/src/bar"><pre class="code-block"><code class="language-plain" language="plain">$ go run main.go

Hello, World!
--------------------------------------------------------------------------------
go: go1.18.1</code></pre></cmd></div>

</page><!--BREAK-->


</body>
</html>`

	r.Equal(exp, act)
}
