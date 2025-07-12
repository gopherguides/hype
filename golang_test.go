package hype

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Golang_Cmds(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	p := NewParser(nil)

	in := strings.NewReader(`<go doc="context.Context"></go>`)
	nodes, err := p.ParseFragment(in)
	r.NoError(err)

	ctx := context.Background()
	doc := &Document{
		Nodes:  nodes,
		Parser: p,
	}

	err = doc.Execute(ctx)

	r.NoError(err)

	act := doc.String()
	act = strings.TrimSpace(act)

	exp := `type Context interface {`
	r.Contains(act, exp)
}

func Test_Golang_Multiple(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/auto/commands"
	cab := os.DirFS(root)

	p := NewParser(cab)
	p.Root = root

	in := strings.NewReader(`<go doc="-short context,-short errors" run="." src="greet/src"></go>`)
	doc, err := p.Parse(in)
	r.NoError(err)

	ctx := context.Background()

	err = doc.Execute(ctx)

	r.NoError(err)

	act := doc.String()
	act = strings.TrimSpace(act)

	exp := `$ go doc -short context`
	r.Contains(act, exp)

	exp = `type Context interface{`
	r.Contains(act, exp)

	exp = `$ go doc -short errors`
	r.Contains(act, exp)

	exp = `func New(text string) error`
	r.Contains(act, exp)

	r.Contains(act, "Hello, World!")

}

func Test_Golang_Sym(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/golang"
	cab := os.DirFS(root)

	p := NewParser(cab)
	p.Root = root

	in := strings.NewReader(`<go sym="Foo" src="sym" figure-id="xyz"></go>`)

	doc, err := p.Parse(in)
	r.NoError(err)

	ctx := context.Background()

	err = doc.Execute(ctx)

	r.NoError(err)

	act := doc.String()
	act = strings.TrimSpace(act)

	exp := `<html><head></head><body><page>
<cmd exec="go doc -cmd -u -src -short Foo" figure-id="xyz" hide-cmd="" language="go" src="sym" sym="Foo"><pre><code class="language-go" language="go">// Foo is a foo.
func Foo() string {
	return "foo"
}</code></pre></cmd>
</page>
</body></html>`

	// fmt.Println(act)
	r.Equal(exp, act)

}

func Test_Golang_Sym_main(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/golang"
	cab := os.DirFS(root)

	p := NewParser(cab)
	p.Root = root

	in := strings.NewReader(`<go sym="main" src="sym/cmd"></go>`)

	doc, err := p.Parse(in)
	r.NoError(err)

	ctx := context.Background()

	err = doc.Execute(ctx)

	r.NoError(err)

	act := doc.String()
	act = strings.TrimSpace(act)

	exp := `<html><head></head><body><page>
<cmd exec="go doc -cmd -u -src -short main" hide-cmd="" language="go" src="sym/cmd" sym="main"><pre><code class="language-go" language="go">func main() {
	Greet()
}</code></pre></cmd>
</page>
</body></html>`

	// fmt.Println(act)
	r.Equal(exp, act)

}
