package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Tags_MetaData(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	doc, err := p.ParseFile("html5.html")
	r.NoError(err)
	r.NotNil(doc)

	md := doc.MetaData()
	r.Equal("website", md["og:type"])
}
