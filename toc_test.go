package hype

import (
	"context"
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_GenerateToC(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/toc")
	p.Section = 42

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "hype.md")
	r.NoError(err)

	toc, err := GenerateToC(p, doc.Children())
	r.NoError(err)

	r.Len(toc, 1)

	b, err := fs.ReadFile(p.FS, "hype.gold")
	r.NoError(err)

	exp := string(b)
	exp = strings.TrimSpace(exp)

	act := toc.String()
	act = strings.TrimSpace(act)

	r.Equal(exp, act)
}

func Test_ToC_MarshalJSON(t *testing.T) {
	t.Parallel()

	toc := &ToC{
		Element: NewEl("toc", nil),
	}

	testJSON(t, "toc", toc)
}
