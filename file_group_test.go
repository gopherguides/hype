package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FileGroup(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := testdata

	p := testParser(t, cab)

	doc, err := p.ParseFile("filegroups.md")
	r.NoError(err)
	r.NotNil(doc)

	groups := doc.Children.AllAtom(FileGroup_Atom)
	r.Len(groups, 1)

	fg, ok := groups[0].(*FileGroup)
	r.True(ok)

	r.Equal("snippets", fg.Name())

	exp := `<filegroup name="snippets">
  <file src="src/snip.html"><a href="src/snip.html" target="_blank">src/snip.html</a></file>
  <file src="src/snip.txt"><a href="src/snip.txt" target="_blank">src/snip.txt</a></file>
  <file src="src/snippets.go"><a href="src/snippets.go" target="_blank">src/snippets.go</a></file>
</filegroup>`

	act := fg.String()

	// fmt.Println(act)
	r.Equal(exp, act)
}
