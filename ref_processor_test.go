package hype

import (
	"context"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RefProcessor_Process(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/auto/refs/figure-styles"
	cab := os.DirFS(root)

	p := NewParser(cab)

	ctx := context.Background()

	doc, err := p.ParseExecuteFile(ctx, "module.md")
	r.NoError(err)

	rp := &RefProcessor{}

	err = rp.Process(doc)
	r.NoError(err)

	b, err := fs.ReadFile(cab, "module.gold")
	r.NoError(err)

	exp := string(b)
	act := doc.String()

	r.Equal(exp, act)

}

func Test_RefProcessor_Current(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	const key = "figure"

	var rp *RefProcessor

	exp := 0
	act := rp.CurIndex(key)

	r.Equal(exp, act)

	rp = &RefProcessor{}

	exp = 0
	act = rp.CurIndex(key)

	r.Equal(exp, act)

	rp.indexes[key] = 1

	exp = 1
	act = rp.CurIndex(key)

	r.Equal(exp, act)
}

func Test_RefProcessor_Next(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	key := "figure"

	rp := &RefProcessor{}

	exp := 0
	act := rp.CurIndex(key)

	r.Equal(exp, act)

	exp = 1
	act = rp.NextIndex(key)

	r.Equal(exp, act)

	key = "table"

	exp = 0
	act = rp.CurIndex(key)

	r.Equal(exp, act)

	exp = 1
	act = rp.NextIndex(key)

	r.Equal(exp, act)
}
