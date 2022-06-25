package cli

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_Latex_Timeout(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cmd := &Latex{}
	cmd.FS = fstest.MapFS{
		"module.md": &fstest.MapFile{
			Data: []byte(`<cmd exec="sleep 1"></cmd>`),
		},
	}

	args := []string{"-t", "1ms"}
	err := cmd.Main(context.Background(), "", args)
	r.Error(err)

	r.True(errors.Is(err, context.DeadlineExceeded))
}

func Test_Latex_Main(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	op, err := os.MkdirTemp("", "hype-cli-latex-test")
	r.NoError(err)
	defer os.RemoveAll(op)

	cmd := &Latex{}
	root := "testdata/latex/simple"
	cab := os.DirFS(root)
	cmd.FS = cab

	args := []string{"-o", op}
	err = cmd.Main(context.Background(), root, args)
	r.NoError(err)

	tex := os.DirFS(op)

	exp := []string{
		"assets/foo.png",
		"module.tex",
		"simple/assets/foo.png",
	}

	var act []string

	err = fs.WalkDir(tex, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		act = append(act, path)

		return nil
	})

	r.NoError(err)

	r.Equal(exp, act)

}

func Test_Latex_Main_File(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	op, err := os.MkdirTemp("", "hype-cli-latex-test")
	r.NoError(err)
	defer os.RemoveAll(op)

	cmd := &Latex{}
	root := "testdata/latex/file"
	cab := os.DirFS(root)
	cmd.FS = cab

	args := []string{"-o", op, "index.md"}
	err = cmd.Main(context.Background(), root, args)
	r.NoError(err)

	tex := os.DirFS(op)

	exp := []string{
		"assets/foo.png",
		"module.tex",
		"simple/assets/foo.png",
	}

	var act []string

	err = fs.WalkDir(tex, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		act = append(act, path)

		return nil
	})

	r.NoError(err)

	r.Equal(exp, act)

}

func Test_Latex_Main_Multiple(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	op, err := os.MkdirTemp("", "hype-cli-latex-multi-test")
	r.NoError(err)
	defer os.RemoveAll(op)

	cmd := &Latex{}
	root := "testdata/latex/multi"

	args := []string{"-o", op, root}
	err = cmd.Main(context.Background(), root, args)
	r.NoError(err)

	tex := os.DirFS(op)

	exp := []string{
		"one/assets/foo.png",
		"one/module.tex",
		"one/simple/assets/foo.png",
		"three/assets/foo.png",
		"three/module.tex",
		"three/simple/assets/foo.png",
		"two/assets/foo.png",
		"two/module.tex",
		"two/simple/assets/foo.png",
	}

	var act []string

	err = fs.WalkDir(tex, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		act = append(act, path)

		return nil
	})

	r.NoError(err)

	// fmt.Printf("%#v\n", act)
	r.Equal(exp, act)

}
