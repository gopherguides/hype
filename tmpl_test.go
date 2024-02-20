package hype

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_GoTemplates(t *testing.T) {
	t.Skip()
	t.Parallel()
	r := require.New(t)

	root := "testdata/gotmpls"

	p := testParser(t, root)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "module.md")
	r.NoError(err)

	act := doc.String()

	// fmt.Println(act)

	exp := ``

	r.Equal(exp, act)

}
