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

	p := testParser(t, testdata, "testdata")
	p.FileName = "run.md"

	doc, err := p.ParseMD([]byte(`<cmd exec="go run ." src="cmd" code="main.go" snippet="example" hide-duration></cmd>`))
	r.NoError(err)

	act := doc.String()

	fmt.Println(act)

	exp := `<html><head><meta charset="utf-8" /></head><body>
<page>

<div><p><pre><code class="language-go" code="main.go" exec="go run ." hide-duration="" language="go" snippet="example" src="cmd/main.go">func main() {
	fmt.Println(&#34;Hello, world!&#34;)
}
</code></pre></p><cmd code="main.go" exec="go run ." hide-duration="" snippet="example" src="cmd"><pre class="code-block"><code class="language-plain" language="plain">$ go run .

Hello, world!
--------------------------------------------------------------------------------
go: go1.18</code></pre></cmd></div>

</page><!--BREAK-->


</body>
</html>`

	r.Equal(exp, act)
}
