package hype

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TOC(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := os.DirFS("testdata/toc")

	p := testParser(t, cab)

	var docs []*Document

	err := fs.WalkDir(cab, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		if base != "module.md" {
			return nil
		}

		doc, err := p.ParseFile(path)
		r.NoError(err)
		r.NotNil(doc)
		docs = append(docs, doc)
		return nil
	})

	r.Len(docs, 2)
	r.NoError(err)

	toc, err := TocsFromDocs("My Big Book", docs...)
	r.NoError(err)
	r.NotNil(toc)

	exp := `My Big Book
	CHAPTER 0
		Chapter 0
			Section 0.1
				Subsection 0.1.1
			Section 0.2
				Subsection 0.2.1
					Subsubsection 0.2.1a
		Chapter 0.A
			Section 0.A.1
	CHAPTER 1
		Chapter 1
			Section 1.1
			Section 1.2
			Section 1.3
				Subsection 1.3.1
					Subsubsection 1.3.1a
				Subsection 1.3.2
		Chapter 1A
			Section 1.A.1
				Subsection 1.A.1.1
				Subsection 1.A.1.2
			Section 1.A.2
		Chapter 1B
			Section 1.B.1`

	act := toc.String()
	act = strings.TrimSpace(act)

	fmt.Println(act)

	r.Equal(exp, act)
}
