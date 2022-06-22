package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WholeFromPath(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	fp := "testdata/whole/simple"

	w, err := WholeFromPath(fp, "book", "chapter")
	r.NoError(err)

	r.Equal("simple", w.Name.String())
	r.Equal("book", w.Ident.String())
	r.Len(w.Parts, 3)

	part, ok := w.Parts["one"]
	r.True(ok)
	r.Equal("chapter", part.Ident.String())
	r.Equal("One", part.Name.String())
	r.Equal(1, part.Number)
	r.Equal(`"Chapter 1: One"`, part.String())

	part, ok = w.Parts["two"]
	r.True(ok)
	r.Equal("chapter", part.Ident.String())
	r.Equal("Two", part.Name.String())
	r.Equal(2, part.Number)
	r.Equal(`"Chapter 2: Two"`, part.String())
}
