package commander

import (
	"fmt"
	"os"
	"testing"

	"github.com/gopherguides/hype"
	"github.com/stretchr/testify/require"
)

func Test_Tags_SetSource(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	dir := "testdata/set-sources"
	cab := os.DirFS(dir)

	p, err := hype.NewParser(cab)
	r.NoError(err)

	p.Root = dir

	err = Register(p)
	r.NoError(err)

	doc, err := p.ParseFile("module.md")
	r.NoError(err)
	r.NotNil(doc)

	fmt.Println(doc.String())

	t.Fail()

}
