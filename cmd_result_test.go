package hype

import (
	"io"
	"testing"
	"time"

	"github.com/markbates/clam"
)

func Test_CmdResult_MarshalJSON(t *testing.T) {
	t.Parallel()

	cr := &CmdResult{
		Element: NewEl("cmd", nil),
		Result: &clam.Result{
			Args:     []string{"echo", "hello"},
			Dir:      "/tmp",
			Duration: time.Second,
			Env:      []string{"FOO=bar", "BAR=baz"},
			Err:      io.EOF,
			Exit:     1,
			Stderr:   []byte("nothing"),
			Stdout:   []byte("foo\nbar\nbaz\n"),
		},
	}

	testJSON(t, "cmd_result", cr)

}
