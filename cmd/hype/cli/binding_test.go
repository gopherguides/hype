package cli

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/flect"
	"github.com/gopherguides/hype"
	"github.com/stretchr/testify/require"
)

func Test_NewBindingNode(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	root := "testdata/whole/simple"
	cab := os.DirFS(root)

	whole, err := WholeFromPath(root, "book", "chapter")
	whole.Name = flect.New("My Big Book")

	r.NoError(err)

	tcs := []struct {
		in  string
		out string
		err bool
	}{
		{in: "<binding></binding>", err: true},
		{in: `<binding part="two"></binding>`, out: `"Chapter 2: Two"`},
		{in: `<binding whole></binding>`, out: `book`},
		{in: `<binding whole="title"></binding>`, out: `My Big Book`},
		{in: `<binding part></binding>`, out: `chapter`},
	}

	for _, tc := range tcs {
		name := fmt.Sprintf("in: %q", tc.in)

		t.Run(name, func(t *testing.T) {

			r := require.New(t)

			p := hype.NewParser(cab)
			p.NodeParsers[hype.Atom("binding")] = NewBindingNodes(whole)

			doc, err := p.Parse(strings.NewReader(tc.in))
			if tc.err {
				r.Error(err)
				return
			}

			r.NoError(err)

			bindings := hype.ByType[*Binding](doc.Children())
			r.Len(bindings, 1)

			binding := bindings[0]

			act := binding.String()

			r.Equal(tc.out, act)

		})

	}

}
