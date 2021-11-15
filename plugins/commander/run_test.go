package commander

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Run(t *testing.T) {
	t.Parallel()

	t.Skip()
	ctx := context.Background()

	table := []struct {
		name string
		args []string
		cmd  string
		exp  *Result
		dir  string
		err  bool
	}{
		{
			name: "echo good",
			cmd:  "echo",
			args: []string{"Hello, World"},
			exp: &Result{
				stdout: []byte("Hello, World\n"),
			},
		},
		{
			name: "unknown command",
			cmd:  "unknown",
			exp: &Result{
				ExitCode: -1,
			},
		},
		{
			name: "bad go run",
			dir:  "testdata/cmd",
			cmd:  "go",
			args: []string{"run", "-tags", "sad", "."},
			exp: &Result{
				ExitCode: 1,
				stderr:   []byte("boom!\nexit status 255\n"),
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			res, err := Run(ctx, tt.dir, tt.cmd, tt.args...)

			// fmt.Println(res)

			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(res)

			exit := res.ExitCode
			r.Equal(tt.exp.ExitCode, exit)

			if exit != 0 {
				r.NotNil(res.Err)

				assertReaders(t, res.Stderr(), tt.exp.Stderr())
				return
			}

			assertReaders(t, res.Stdout(), tt.exp.Stdout())
		})
	}

}
