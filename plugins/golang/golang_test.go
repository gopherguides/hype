package golang

import (
	"io/fs"
	"os"
	"testing"

	"github.com/gopherguides/hype"
)

var testdata = os.DirFS("testdata")

func testParser(t testing.TB, cab fs.FS, root string) *hype.Parser {
	t.Helper()

	p, err := hype.NewParser(cab)
	if err != nil {
		t.Fatal(err)
	}

	Register(p, root)
	return p
}
