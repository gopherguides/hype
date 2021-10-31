package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_File(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := testdata

	p := testParser(t, cab)

	doc, err := p.ParseFile("files.md")
	r.NoError(err)
	r.NotNil(doc)

	files := doc.Children.AllAtom(File_Atom)
	r.Len(files, 2)

	f, ok := files[0].(*File)
	r.True(ok)

	r.Equal("src/main.go", f.Src())

	exp := `<file src="src/main.go"><a href="src/main.go" target="_blank">src/main.go</a></file>`
	act := f.String()

	// fmt.Println(act)
	r.Equal(exp, act)
}
