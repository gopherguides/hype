package hype

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parser_Snippets(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	doc, err := p.ParseFile("html5.html")
	r.NoError(err)
	r.NotNil(doc)

	exp := node_string_exp
	fmt.Println(doc.String())
	r.Equal(exp, doc.String())
}
