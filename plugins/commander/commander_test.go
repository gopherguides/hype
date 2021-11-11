package commander

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/gopherguides/hype"
	"github.com/stretchr/testify/require"
)

var testdata = os.DirFS("testdata")

func testParser(t testing.TB, cab fs.FS, root string) *hype.Parser {
	t.Helper()

	r := require.New(t)

	p, err := hype.NewParser(cab)
	p.Root = root
	r.NoError(err)

	Register(p)
	return p
}

func assertReaders(t testing.TB, r1 io.Reader, r2 io.Reader) {
	t.Helper()
	r := require.New(t)
	b, err := io.ReadAll(r1)
	r.NoError(err)

	act := string(b)

	b, err = io.ReadAll(r2)
	r.NoError(err)
	exp := string(b)

	r.Equal(exp, act)
}

func assertExp(t testing.TB, name string, act string) {
	t.Helper()

	r := require.New(t)

	b, err := fs.ReadFile(testdata, filepath.Join("exps", name))
	r.NoError(err)

	exp := string(b)
	r.Equal(exp, act)
}
