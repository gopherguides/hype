package cli

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/markbates/cleo"
	"github.com/markbates/iox"
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

	op, err := os.MkdirTemp("", "Test_Latex_Main")
	r.NoError(err)
	defer os.RemoveAll(op)

	oi := iox.Discard()

	cmd := &Latex{
		Cmd: cleo.Cmd{
			IO: oi,
		},
	}

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

	op, err := os.MkdirTemp("", "Test_Latex_Main_File")
	r.NoError(err)
	defer os.RemoveAll(op)

	cmd := &Latex{
		Cmd: cleo.Cmd{
			IO: iox.Discard(),
		},
	}

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

	op, err := os.MkdirTemp("", "Test_Latex_Main_Multiple")
	r.NoError(err)
	defer os.RemoveAll(op)

	cmd := &Latex{
		Cmd: cleo.Cmd{
			IO: iox.Discard(),
		},
	}

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

func Test_Latex_FolderName(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := fstest.MapFS{
		"01-foo/module.md": &fstest.MapFile{
			Data: []byte(`# Hello`),
		},
	}

	op, err := os.MkdirTemp("", "Test_Latex_FolderName")
	r.NoError(err)
	defer os.RemoveAll(op)

	cmd := &Latex{
		Cmd: cleo.Cmd{
			IO: iox.Discard(),
		},
	}

	cmd.FS = cab

	args := []string{"-o", op, "-f"}
	err = cmd.Main(context.Background(), "", args)
	r.NoError(err)

	tex := os.DirFS(op)

	exp := []string{
		"01-foo/01-foo.tex",
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
