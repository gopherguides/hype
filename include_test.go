package hype

import (
	"encoding/json"
	"io/fs"
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_NewInclude(t *testing.T) {
	t.Parallel()

	validInc := htmx.AttrNode(t, "include", map[string]string{
		"src": "html5.html",
	})
	fileMissing := htmx.AttrNode(t, "include", map[string]string{
		"src": "404.html",
	})
	srcMissing := htmx.ElementNode(t, "include")

	table := []struct {
		name string
		cab  fs.FS
		node *html.Node
		err  bool
	}{
		{name: "missing src attr", cab: testdata, node: srcMissing, err: true},
		{name: "missing src file", cab: testdata, node: fileMissing, err: true},
		{name: "nil all the way", err: true},
		{name: "non include tag", node: htmx.ElementNode(t, "p"), err: true},
		{name: "valid include", cab: testdata, node: validInc},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			p := testParser(t, tt.cab)
			i, err := p.NewInclude(NewNode(tt.node))

			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(i)
		})
	}

}

func Test_Include_JSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	validInc := htmx.AttrNode(t, "include", map[string]string{
		"src": "html5.html",
	})

	inc := &Include{
		Node: NewNode(validInc),
	}

	exp := `{"attributes":{"src":"html5.html"},"data":"include","type":"element"}`

	b, err := json.Marshal(inc)
	r.NoError(err)
	r.Equal(exp, string(b))
}
