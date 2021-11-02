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

	groups := doc.Children.ByAtom(FileGroup_Atom)
	r.Len(groups, 1)

	fg, ok := groups[0].(*FileGroup)
	r.True(ok)

	r.Equal("snippets", fg.Name())

	exp := `<filegroup name="snippets">
  <file src="src/snip.html"></file>
  <file src="src/snip.txt"></file>
  <file src="src/snippets.go"></file>
</filegroup>`

	act := fg.String()

	// fmt.Println(act)
	r.Equal(exp, act)
}