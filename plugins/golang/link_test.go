package golang

import (
	"os"
	"testing"

	"github.com/gopherguides/hype"
	"github.com/stretchr/testify/require"
)

func Test_Link(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := os.DirFS("testdata")

	p, err := hype.NewParser(cab)
	r.NoError(err)
	p.SetCustomTag(LINK, func(node *hype.Node) (hype.Tag, error) {
		return NewLink(node)
	})

	doc, err := p.ParseFile("link.md")
	r.NoError(err)

	exp := `<html><head></head><body>
<page>

<h1>Links</h1>

<h2>Element links</h2>

<p>This <godoc#a for="context"><a href="https://pkg.go.dev/context" target="_blank" /></godoc#a> link is <godoc#a for="io#Writer"><a href="https://pkg.go.dev/io#Writer" target="_blank" /></godoc#a> SEE ME.</p>

</page><!--BREAK-->


</body>
</html>`
	act := doc.String()

	// fmt.Println(act)
	r.Contains(act, "SEE ME")
	r.Equal(exp, act)

}
