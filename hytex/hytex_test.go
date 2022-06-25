package hytex

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/gopherguides/hype"
	"github.com/stretchr/testify/require"
)

func Test_Covert(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	const root = "testdata/convert"

	cab := os.DirFS(root)
	p := hype.NewParser(cab)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "module.md")
	r.NoError(err)

	tex, err := Convert(ctx, doc)
	r.NoError(err)
	r.NotNil(tex)

	act, err := fs.ReadFile(tex, "module.tex")
	r.NoError(err)

	act = bytes.TrimSpace(act)

	fmt.Println(string(act))

	exp, err := fs.ReadFile(cab, "module.tex.gold")
	r.NoError(err)

	exp = bytes.TrimSpace(exp)
	r.Equal(string(exp), string(act))

	assets := []string{
		"assets/foo.png",
		"simple/assets/foo.png",
	}

	for _, asset := range assets {
		_, err = fs.Stat(tex, asset)
		r.NoError(err)
	}

}
