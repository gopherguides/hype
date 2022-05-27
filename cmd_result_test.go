package hype

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CmdResult_DataAttrs(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/commands/results/data"

	cab := os.DirFS(root)

	p := NewParser(cab)

	doc, err := p.ParseExecuteFile(context.Background(), "module.md")
	r.NoError(err)

	act := doc.String()

	// fmt.Println(act)
	compareOutputFile(t, cab, act, "module.gold")

}
