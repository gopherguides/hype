package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SrcAttr(t *testing.T) {
	t.Parallel()

	good := Attributes{
		"src": "foo.go",
	}

	skip := Attributes{
		"src":      "foo.go",
		"skip-src": "true",
	}

	table := []struct {
		name string
		ats  Attributes
		src  string
		err  bool
	}{
		{name: "src found", ats: good, src: "foo.go"},
		{name: "src not found", ats: Attributes{}, err: true},
		{name: "src found, but skipped", ats: skip, err: true},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			r := require.New(t)

			act, ok := SrcAttr(tt.ats)
			r.Equal(!tt.err, ok)

			r.Equal(Source(tt.src), act)
		})
	}

}
