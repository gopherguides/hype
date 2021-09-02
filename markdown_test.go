package hype

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parser_ParseMD_TagBug(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	in := "1. In the terminal type `git config --global user.name \"<your name>\"` and press \"enter\"."

	p := testParser(t, week01)
	doc, err := p.ParseMD([]byte(in))
	r.NoError(err)

	act := doc.String()
	exp := `TODO`

	fmt.Println(act)
	r.Equal(exp, act)

}
