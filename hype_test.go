package hype

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	goVersion = func() string {
		return "go.test"
	}
}

type brokenReader struct{}

func (brokenReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("broken reader")
}

func testParser(t testing.TB, root string) *Parser {
	t.Helper()

	cab := os.DirFS(root)

	p := NewParser(cab)
	p.Root = root

	return p
}

func testModule(t testing.TB, root string) {
	t.Helper()

	r := require.New(t)

	p := testParser(t, root)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "module.md")
	r.NoError(err)

	act := doc.String()

	// fmt.Println(act)

	b, err := fs.ReadFile(p.FS, "module.gold")
	r.NoError(err)

	exp := string(b)

	if exp != act {
		fmt.Println(act)
		fp := filepath.Join("tmp", root)
		err = os.MkdirAll(fp, 0755)
		r.NoError(err)

		f, err := os.Create(filepath.Join(fp, "output.html"))
		r.NoError(err)
		defer f.Close()

		_, err = f.Write([]byte(act))
		r.NoError(err)

		err = f.Close()
		r.NoError(err)

		r.Equal(exp, act)
	}
}

func Test_Testdata_Auto_Modules(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	root := "testdata/auto"

	err := fs.WalkDir(os.DirFS(root), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(path)
		if base != "module.md" {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			dir := filepath.Dir(path)

			testModule(t, filepath.Join(root, dir))
		})

		return filepath.SkipDir
	})

	r.NoError(err)
}
