package hype

import (
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
	exp := `<li>In the terminal type <code>git config --global user.name &#34;&lt;your name&gt;&#34;</code> and press “enter”.</li>`

	// fmt.Println(act)
	r.Contains(act, exp)

}
