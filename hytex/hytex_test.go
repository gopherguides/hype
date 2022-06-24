package hytex

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gopherguides/hype"
	"github.com/stretchr/testify/require"
)

func Test_Covert(t *testing.T) {
	t.Parallel()

	t.Skip()
	r := require.New(t)

	const root = "testdata/convert"

	cab := os.DirFS(root)
	p := hype.NewParser(cab)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "module.md")
	r.NoError(err)

	fs, err := Convert(ctx, doc)
	r.NoError(err)

	r.NotNil(fs)
}
