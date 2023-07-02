package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gopherguides/hype"
	"github.com/markbates/iox"
	"github.com/stretchr/testify/require"
)

type echoTag struct {
	*hype.Element
}

func (e *echoTag) Execute(ctx context.Context, doc *hype.Document) error {
	if e == nil {
		return fmt.Errorf("echoTag is nil")
	}

	e.Lock()
	defer e.Unlock()

	s := e.Nodes.String()
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	e.Nodes = hype.Nodes{
		hype.Text(s),
	}

	return nil
}

func Test_Encode_JSON(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	root := "testdata/encode/json"
	cab := os.DirFS(root)

	p, err := NewParser(cab, root, root)
	r.NoError(err)

	p.NodeParsers["echo"] = func(p *hype.Parser, el *hype.Element) (hype.Nodes, error) {
		return hype.Nodes{
			&echoTag{el},
		}, nil
	}

	tcs := []struct {
		name string
		in   any
		args []string
		exp  string
	}{
		{name: "execute file", exp: "success/execute-file.json", args: []string{"module.md"}},
		{name: "parse file", exp: "success/parse-file.json", args: []string{"-p", "module.md"}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			cmd := &Encode{
				Parser: p,
			}

			bb := iox.Buffer{}
			cmd.IO.Out = &bb.Out
			cmd.IO.Err = &bb.Err

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := cmd.Main(ctx, root, tc.args)
			r.NoError(err)

			b, err := fs.ReadFile(cab, tc.exp)
			r.NoError(err)

			exp := string(b)
			exp = strings.TrimSpace(exp)

			act := bb.Out.String()
			act = strings.TrimSpace(act)

			// f, err := os.Create(filepath.Join(root, tc.exp))
			// r.NoError(err)

			// _, err = f.WriteString(act)
			// r.NoError(err)
			// r.NoError(f.Close())

			r.Equal(exp, act)

		})
	}

}
