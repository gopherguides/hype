package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Ranger_Parse(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		in     string
		ranger *Ranger
		start  int
		end    int
		err    bool
	}{
		{in: "", start: 0, end: 0, err: true},
		{in: "1:2", start: 1, end: 2},
		{in: "2:0", start: 2, end: 0},
		{in: ":2", start: 0, end: 2},
		{in: ":", start: 0, end: 0},
		{in: "1:", start: 1, end: 42, ranger: &Ranger{End: 42}},
	}

	for _, tc := range tcs {
		name := tc.in
		if len(name) == 0 {
			name = "empty"
		}

		t.Run(name, func(t *testing.T) {

			r := require.New(t)

			lone := tc.ranger
			if lone == nil {
				lone = &Ranger{}
			}

			err := lone.Parse(tc.in)

			if tc.err {
				r.Error(err)
				return
			}

			r.True(lone.IsRange(tc.in))

			r.Equal(tc.start, lone.Start)
			r.Equal(tc.end, lone.End)

		})
	}

}
