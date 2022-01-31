package commander

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gopherguides/hype"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Tag(t *testing.T) {
	t.Parallel()
	t.Skip()
	r := require.New(t)

	os.RemoveAll("~/.hype")
	p := testParser(t, testdata, "testdata")

	doc, err := p.ParseFile("run.md")
	r.NoError(err)
	r.NotNil(doc)

	act := doc.String()

	fmt.Println(act)
	assertExp(t, "cmd.exp.html", act)
}

func Test_Cmd_Run_Error(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cmd := cmdTag(t, hype.Attributes{
		"exec": "go run -tags sad .",
		"src":  "cmd",
		"exit": "1",
	})

	r.NotNil(cmd)

	p := testParser(t, testdata, "testdata")
	p.FileName = "run.md"

	err := cmd.work(p)
	r.Error(err)

	act := err.Error()
	act = strings.TrimSpace(act)

	fmt.Println(act)
	exp := `expected exit code 0, got 1:
<cmd exec="go run -tags sad ." exit="1" src="cmd">
file name:	"run.md"
command:	"$ go run -tags sad ."`

	r.True(strings.HasPrefix(act, exp))
}
