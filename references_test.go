package hype

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_References_Simple(t *testing.T) {
	t.Parallel()

	table := []string{
		"fenced",
		"includes",
		"simple",
	}

	for _, tc := range table {
		t.Run(tc, func(t *testing.T) {

			r := require.New(t)

			root := filepath.Join("testdata/refs", tc)
			cab := os.DirFS(root)

			p := NewParser(cab)
			p.Root = root

			doc, err := p.ParseFile("module.md")
			r.NoError(err)

			err = doc.Execute(context.Background())
			r.NoError(err)

			act := doc.String()
			act = strings.TrimSpace(act)

			// fmt.Println(act)
			compareOutputFile(t, cab, act, "module.gold")
		})
	}
}

func Test_References_Fenced(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/refs/fenced"
	cab := os.DirFS(root)

	p := NewParser(cab)
	p.Root = root

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	err = doc.Execute(context.Background())
	r.NoError(err)

	act := doc.String()
	act = strings.TrimSpace(act)

	// fmt.Println(act)
	compareOutputFile(t, cab, act, "module.gold")
}

func Test_Refs_Images(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/refs/images"
	cab := os.DirFS(root)

	p := NewParser(cab)
	p.Root = root

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	err = doc.Execute(context.Background())
	r.NoError(err)

	act := doc.String()
	act = strings.TrimSpace(act)

	// fmt.Println(act)
	compareOutputFile(t, cab, act, "module.gold")
}
