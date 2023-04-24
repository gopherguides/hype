package hype

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Document_JSON(t *testing.T) {
	t.Skip()

	t.Parallel()
	r := require.New(t)

	root := "testdata/doc/to_md"

	p := testParser(t, root)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "module.md")
	// doc, err := p.ParseFile("module.md")
	r.NoError(err)

	b, err := json.MarshalIndent(doc, "", "  ")
	r.NoError(err)

	act := string(b)
	act = strings.TrimSpace(act)

	// fmt.Println(act)
	f, err := os.Create("testdata/doc/to_md/module.json")
	r.NoError(err)
	defer f.Close()

	_, err = f.WriteString(act)
	r.NoError(err)

	exp := ``

	r.Equal(exp, act)

}
