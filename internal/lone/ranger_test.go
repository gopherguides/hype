package lone

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
		{in: "-1:5", start: -1, end: 5},
		{in: "1:", start: 1, end: 0},
		{in: "1:", start: 1, end: 42, ranger: &Ranger{End: 42}},
		{in: "1:2", start: 1, end: 2},
		{in: "2:0", start: 2, end: 0},
		{in: ":", start: 0, end: 0},
		{in: ":-5", start: 0, end: -5},
		{in: ":2", start: 0, end: 2},
	}

	for _, tc := range tcs {
		name := tc.in
		if len(name) == 0 {
			name = "empty"
		}

		t.Run(name, func(t *testing.T) {

			r := require.New(t)

			rg := tc.ranger
			if rg == nil {
				rg = &Ranger{}
			}

			err := rg.Parse(tc.in)

			if tc.err {
				r.Error(err)
				return
			}

			r.True(rg.IsRange(tc.in))

			r.Equal(tc.start, rg.Start)
			r.Equal(tc.end, rg.End)

		})
	}

}

func Test_Range(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		in    string
		start int
		end   int
		err   bool
	}{
		{in: "", start: 0, end: 0, err: true},
		{in: "-1:5", start: -1, end: 5},
		{in: "1:", start: 1, end: 0},
		{in: "1:2", start: 1, end: 2},
		{in: "2:0", start: 2, end: 0},
		{in: ":", start: 0, end: 0},
		{in: ":-5", start: 0, end: -5},
		{in: ":2", start: 0, end: 2},
		{in: "100:200", start: 100, end: 200},
	}

	for _, tc := range tcs {
		name := tc.in
		if len(name) == 0 {
			name = "empty"
		}

		t.Run(name, func(t *testing.T) {

			r := require.New(t)

			start, end, err := Range(tc.in)

			if tc.err {
				r.Error(err)
				return
			}

			r.Equal(tc.start, start)
			r.Equal(tc.end, end)

		})
	}

}
