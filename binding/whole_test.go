package binding

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WholeFromPath(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	fp := "testdata/whole/simple"
	cab := os.DirFS(fp)

	w, err := WholeFromPath(cab, fp, "book", "chapter")
	r.NoError(err)

	r.Equal("simple", w.Name.String())
	r.Equal("book", w.Ident.String())
	r.Len(w.Parts, 3)

	part, ok := w.Parts["one"]
	r.True(ok)
	r.Equal("chapter", part.Ident.String())
	r.Equal("one", part.Name.String())
	r.Equal(1, part.Number)
	r.Equal(`one`, part.Name.String())

	part, ok = w.Parts["two"]
	r.True(ok)
	r.Equal("chapter", part.Ident.String())
	r.Equal("two", part.Name.String())
	r.Equal(2, part.Number)
	r.Equal(`two`, part.Name.String())
}
