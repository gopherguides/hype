package hype

import (
	"encoding/json"
	"io/fs"
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_NewImage(t *testing.T) {
	t.Parallel()

	validImg := htmx.AttrNode("img", map[string]string{
		"src": "assets/foo.png",
	})

	fileMissing := htmx.AttrNode("img", map[string]string{
		"src": "assets/404.png",
	})
	srcMissing := htmx.ElementNode("img")

	table := []struct {
		name string
		cab  fs.FS
		node *html.Node
		err  bool
	}{
		{name: "missing src attr", cab: testdata, node: srcMissing, err: true},
		{name: "missing src file", cab: testdata, node: fileMissing, err: true},
		{name: "nil all the way", err: true},
		{name: "non image tag", node: htmx.ElementNode("p"), err: true},
		{name: "valid image", cab: testdata, node: validImg},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			i, err := NewImage(tt.cab, NewNode(tt.node))

			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(i)
		})
	}

}

func Test_Image_JSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	validImg := htmx.AttrNode("img", map[string]string{
		"src": "assets/foo.png",
	})

	img := &Image{Node: NewNode(validImg)}

	exp := `{"atom":"img","attributes":{"src":"assets/foo.png"},"data":"img","type":"element"}`
	b, err := json.Marshal(img)
	r.NoError(err)
	r.Equal(exp, string(b))
}
