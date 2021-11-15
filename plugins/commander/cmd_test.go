package commander

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cmd_Tag(t *testing.T) {
	t.Parallel()
	// t.Skip()
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
