package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Marked_Errors(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		in   string
		exp  string
		dir  string
		args []string
		err  bool
	}{
		{name: "valid", in: `# Hello`, exp: `<h1>Hello</h1>`},
		{name: "invalid", in: `<code src="404.bad" snippet="oops"></code>`, err: true},
		{name: "timeout", in: `<cmd exec="sleep 1"></cmd>`, err: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			cmd := &Marked{}

			out := &bytes.Buffer{}

			cmd.In = strings.NewReader(tc.in)
			cmd.Out = out

			err := cmd.Main(ctx, tc.dir, tc.args)

			act := out.String()

			fmt.Println(act)

			if tc.err {
				r.Error(err)
				return
			}

			r.NoError(err)

			r.Contains(act, tc.exp)

		})

	}

}

func Test_Marked_Timeout(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cmd := &Marked{}

	cmd.In = strings.NewReader(`<cmd exec="sleep 1"></cmd>`)

	err := cmd.Main(context.Background(), "", []string{"-timeout", "1ms"})
	r.Error(err)
	r.True(errors.Is(err, context.DeadlineExceeded))
}

func Test_Marked_Path_Env(t *testing.T) {

	r := require.New(t)

	t.Setenv("MARKED_PATH", "42-foo/module.md")

	cmd := &Marked{}

	out := &bytes.Buffer{}
	cmd.Out = out

	cmd.In = strings.NewReader(`<figure id="foo"><figcaption>hello</figcaption></figure>`)

	err := cmd.Main(context.Background(), "", []string{})

	r.NoError(err)

	act := out.String()

	// fmt.Println(act)
	r.Contains(act, `Figure 42.1`)
}
