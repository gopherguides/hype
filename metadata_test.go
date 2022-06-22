package hype

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewMetadata(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := os.DirFS("testdata/metadata/simple")

	p := NewParser(cab)
	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	mds := ByType[*Metadata](doc.Children())
	r.Len(mds, 1)

	md := mds[0]
	class, ok := md.Get("class")
	r.True(ok)
	r.Equal("center, middle, inverse", class)

	act := doc.String()
	act = strings.TrimSpace(act)

	r.NotContains(act, `<metadata>`)
	r.NotContains(act, `duration: 45m`)
}

func Test_NewMetadata_Pages(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := os.DirFS("testdata/metadata/pages")

	p := NewParser(cab)
	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	mds := ByType[*Metadata](doc.Children())
	r.Len(mds, 2)

	pages := ByType[*Page](doc.Children())
	r.Len(pages, 2)

	mds = ByType[*Metadata](pages[0].Children())
	r.Len(mds, 1)

	md := mds[0]
	class, ok := md.Get("class")
	r.True(ok)
	r.Equal("center, middle, inverse", class)

	mds = ByType[*Metadata](pages[1].Children())
	r.Len(mds, 1)

	md = mds[0]
	ov, ok := md.Get("duration")
	r.True(ok)
	r.Equal("1h", ov)
}

func Test_NewMetadata_Multiple_Erro0r(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := os.DirFS("testdata/metadata/multi")

	p := NewParser(cab)
	_, err := p.ParseFile("module.md")
	r.Error(err)
}
