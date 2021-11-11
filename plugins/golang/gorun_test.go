package golang

import (
	"fmt"
	"os"
	"testing"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

func Test_NewGoRun(t *testing.T) {
	t.Parallel()
	t.Skip()
	table := []struct {
		name string
		src  string
		exp  string
		err  bool
	}{
		{name: "specific .go file", src: "./testdata/cmd/main.go"},
		{name: "directory", src: "./testdata/cmd", exp: "TODO"},
		{name: "unknown", src: "???", err: true},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			r := require.New(t)
			hn := htmx.AttrNode(GORUN.String(), htmx.Attributes{
				"src": tt.src,
			})

			node := hype.NewNode(hn)

			gr, err := NewGoRun(node, "")

			if tt.err {
				r.Error(err)
				return
			}
			r.NoError(err)

			r.NoError(err)
			r.NotNil(gr)

			act := gr.String()
			r.Equal(tt.exp, act)
		})

		// p := testParser(t, testdata, "testdata")

		// doc, err := p.ParseFile("run.md")
		// r.NoError(err)

		// act := doc.String()
		// exp := `TODO`

		// fmt.Println(act)
		// r.Equal(exp, act)
	}
}

func Test_NewGoRun_Parser(t *testing.T) {
	t.Skip()
	t.Parallel()
	r := require.New(t)

	os.RemoveAll(".testdata/cmd/output")
	p := testParser(t, testdata, "testdata")
	doc, err := p.ParseFile("run.md")
	r.NoError(err)

	act := doc.String()
	exp := `TODO`

	fmt.Println(act)
	r.Equal(exp, act)

}
