package commander

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Result_Tag(t *testing.T) {
	t.Parallel()
	t.Skip()

	res := Result{
		stderr: []byte("bad output"),
		stdout: []byte("good output"),
	}

	cargs := []string{"echo", "Hello, World"}
	cs := strings.Join(cargs, " ")

	table := []struct {
		name string
		args []string
		ats  Attributes
		exp  string
		data Data
	}{
		{
			name: "happy path",
			exp:  "happy.html",
			ats: Attributes{
				"exec": cs,
			},
			args: cargs,
		},
		{
			name: "happy path with data",
			exp:  "happy-data.html",
			ats: Attributes{
				"exec": cs,
			},
			data: Data{
				"go":       "vX.X.X",
				"duration": "1.30s",
			},
			args: cargs,
		},
		{
			name: "hide-cmd",
			exp:  "hide-cmd.html",
			ats: Attributes{
				"exec":     cs,
				"hide-cmd": "",
			},
			args: cargs,
		},
		{
			name: "hide-stdout",
			exp:  "hide-stdout.html",
			ats: Attributes{
				"exec":        cs,
				"hide-stdout": "",
			},
			args: cargs,
		},
		{
			name: "hide-stderr",
			exp:  "hide-stderr.html",
			ats: Attributes{
				"exec":        cs,
				"hide-stderr": "",
			},
			args: cargs,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			res.args = tt.args

			tag := res.Tag(tt.ats, tt.data)
			r.NotNil(t, tag)

			act := tag.String()
			// fmt.Println(act)
			assertExp(t, tt.exp, act)
		})
	}
}
