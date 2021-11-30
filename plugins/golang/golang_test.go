package golang

import (
	"io/fs"
	"os"
	"testing"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

var testdata = os.DirFS("testdata")

func testParser(t testing.TB, cab fs.FS, root string) *hype.Parser {
	t.Helper()

	p, err := hype.NewParser(cab)
	if err != nil {
		t.Fatal(err)
	}

	Register(p)
	return p
}

func Test_NewGo(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		ats  hype.Attributes
		err  bool
		exp  string
	}{
		{
			name: "go run",
			exp:  "<cmd exec=\"go run main.go\" src=\"demo\"><pre class=\"code-block\"><code class=\"language-plain\" language=\"plain\"></code></pre></cmd>",
			ats: hype.Attributes{
				"run": "main.go",
				"src": "demo",
			},
		},
		{
			name: "go test",
			exp:  "<cmd exec=\"go test ./...\" hide-duration=\"true\" src=\"demo\"><pre class=\"code-block\"><code class=\"language-plain\" language=\"plain\"></code></pre></cmd>",
			ats: hype.Attributes{
				"test": "./...",
				"src":  "demo",
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			node := hype.NewNode(htmx.AttrNode("go", tt.ats))

			tag, err := NewGo(node)
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)

			act := tag.String()
			// fmt.Println(act)
			r.Equal(tt.exp, act)
		})
	}
}
