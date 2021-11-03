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

	files := doc.Children.ByAdam(File_Adam)
	r.Len(files, 2)

	f, ok := files[0].(*File)
	r.True(ok)

	src, ok := f.Source()
	r.True(ok)
	r.Equal("src/main.go", src.String())

	exp := `<file src="src/main.go"></file>`
	act := f.String()

	// fmt.Println(act)
	r.Equal(exp, act)
}
