package hype

import (
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testdata = os.DirFS("testdata")
	week01   = os.DirFS("testdata/week01")
)

func testParser(t *testing.T, cab fs.FS) *Parser {
	t.Helper()

	r := require.New(t)

	p, err := NewParser(cab)
	r.NoError(err)
	return p
}

func ParseFile(t *testing.T, cab fs.FS, name string) *Document {
	t.Helper()

	r := require.New(t)

	p := testParser(t, cab)

	doc, err := p.ParseFile(name)
	r.NoError(err)
	return doc
}

func ParseMD(t *testing.T, cab fs.FS, src []byte) *Document {
	t.Helper()

	r := require.New(t)

	p := testParser(t, cab)

	doc, err := p.ParseMD(src)
	r.NoError(err)
	return doc
}

func ParseReader(t *testing.T, cab fs.FS, rc io.ReadCloser) *Document {
	t.Helper()

	r := require.New(t)

	p := testParser(t, cab)

	doc, err := p.ParseReader(rc)
	r.NoError(err)
	return doc
}
