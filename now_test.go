package hype

import (
	"context"
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_NowNodes(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	tm, err := time.Parse(time.RFC822, "24 Aug 76 12:34 UTC")

	r.NoError(err)

	p := testParser(t, "testdata/now")
	p.NowFn = func() time.Time {
		return tm
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "hype.md")
	r.NoError(err)

	act := doc.String()
	act = strings.TrimSpace(act)

	// fmt.Println(act)

	b, err := fs.ReadFile(p.FS, "hype.gold")
	r.NoError(err)

	exp := string(b)
	exp = strings.TrimSpace(exp)

	r.Equal(exp, act)
}

func Test_Now_MarshalJSON(t *testing.T) {
	t.Parallel()

	n := &Now{
		Element: NewEl("now", nil),
	}
	n.Nodes = append(n.Nodes, Text("1/2/2006"))

	testJSON(t, "now", n)
}
