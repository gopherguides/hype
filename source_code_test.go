package hype

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_SourceCode_MD(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	root := filepath.Join("testdata", "to_md", "source_code")

	tcs := []struct {
		name string
	}{
		{name: "full"},
		{name: "snippet"},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			dir := filepath.Join(root, tc.name)
			cab := os.DirFS(dir)

			p := NewParser(cab)
			p.Root = filepath.Dir(dir)

			doc, err := p.ParseExecuteFile(ctx, "hype.md")
			r.NoError(err)

			act := doc.MD()
			act = strings.TrimSpace(act)

			// fmt.Println(act)

			b, err := os.ReadFile(filepath.Join(root, tc.name, "hype.gold"))
			r.NoError(err)

			exp := strings.TrimSpace(string(b))

			r.Equal(exp, act)

		})
	}

}
