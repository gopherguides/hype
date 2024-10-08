package binding

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_PartFromPath(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		in   string
		num  int
		name string
		err  bool
	}{
		{in: `.`, err: true},
		{in: `01-foo/hype.md`, num: 1, name: "foo"},
		{in: `01-foo`, num: 1, name: "foo"},
		{in: `01-pkgs-mods-deps`, num: 1, name: "pkgs-mods-deps"},
		{in: `012-foo/hype.md`, num: 12, name: "foo"},
		{in: `012-foo`, num: 12, name: "foo"},
		{in: `01234-foo/hype.md`, num: 1234, name: "foo"},
		{in: `01234-foo`, num: 1234, name: "foo"},
		{in: `1-foo/hype.md`, num: 1, name: "foo"},
		{in: `1-foo`, num: 1, name: "foo"},
		{in: `12-foo/hype.md`, num: 12, name: "foo"},
		{in: `12-foo`, num: 12, name: "foo"},
		{in: `1234-foo/hype.md`, num: 1234, name: "foo"},
		{in: `1234-foo`, num: 1234, name: "foo"},
		{in: ``, err: true},
		{in: `foo.md`, err: true},
		{in: `src/simple/1-foo/hype.md`, num: 1, name: "foo"},
	}

	for _, tc := range tcs {
		name := tc.in
		if name == "" {
			name = "empty"
		}

		t.Run(name, func(t *testing.T) {

			r := require.New(t)

			cab := fstest.MapFS{}

			sec, err := PartFromPath(cab, tc.in)
			if tc.err {
				r.Error(err)
				return
			}

			r.NoError(err)

			r.Equal(tc.num, sec.Number)
			r.Equal(tc.name, sec.Name.String())

		})
	}

}
