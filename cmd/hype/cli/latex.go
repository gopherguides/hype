package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/gopherguides/hype/hytex"
	"github.com/markbates/cleo"
	"github.com/markbates/plugins"
	"golang.org/x/sync/errgroup"
)

type Latex struct {
	cleo.Cmd

	// a folder containing all chapters of a book, for example
	ContextPath string
	OutputPath  string
	Timeout     time.Duration // default: 5s
	FolderName  bool

	flags *flag.FlagSet
}

func (cmd *Latex) Description() string {
	return "converts markdown to latex"
}

func (cmd *Latex) Flags() (flags *flag.FlagSet, err error) {
	if err := cmd.validate(); err != nil {
		return nil, err
	}

	defer func() {
		p := recover()
		if p == nil {
			return
		}

		switch t := p.(type) {
		case error:
			err = t
		case string:
			err = fmt.Errorf(t)
		default:
			err = fmt.Errorf("%v", t)
		}
	}()

	cmd.Lock()
	defer cmd.Unlock()

	if cmd.flags != nil {
		return cmd.flags, nil
	}

	cmd.flags = flag.NewFlagSet("latex", flag.ContinueOnError)
	cmd.flags.SetOutput(io.Discard)

	cmd.flags.BoolVar(&cmd.FolderName, "f", cmd.FolderName, "use the folder name as the output file name")
	cmd.flags.DurationVar(&cmd.Timeout, "t", DefaultTimeout(), "timeout for execution")
	cmd.flags.StringVar(&cmd.ContextPath, "c", cmd.ContextPath, "a folder containing all chapters of a book, for example")
	cmd.flags.StringVar(&cmd.OutputPath, "o", "latex", "the output path")

	return cmd.flags, nil
}

func (cmd *Latex) validate() error {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	cmd.Lock()
	defer cmd.Unlock()

	if cmd.Timeout == 0 {
		cmd.Timeout = DefaultTimeout()
	}

	return nil
}

func (cmd *Latex) Main(ctx context.Context, pwd string, args []string) error {
	if err := cmd.validate(); err != nil {
		return plugins.Wrap(cmd, err)
	}

	flags, err := cmd.Flags()
	if err != nil {
		return plugins.Wrap(cmd, err)
	}

	if err := flags.Parse(args); err != nil {
		return plugins.Wrap(cmd, err)
	}

	var fn string

	args = flags.Args()

	if len(args) > 0 {
		pwd = filepath.Dir(args[0])
		fn = filepath.Base(args[0])
		if len(filepath.Ext(fn)) == 0 {
			fn = ""
			pwd = args[0]
		}
	}

	cab := cmd.FS

	if cab == nil {
		cab = os.DirFS(pwd)
	}

	err = WithTimeout(ctx, cmd.Timeout, func(ctx context.Context) error {

		if len(fn) == 0 {
			if err := cmd.executeFolder(ctx, cab, pwd); err != nil {
				return plugins.Wrap(cmd, err)
			}
			return nil
		}

		man := manifest{
			cab: cab,
			pwd: pwd,
			op:  cmd.OutputPath,
			fn:  fn,
		}

		return cmd.executeFile(ctx, man)
	})

	if err != nil {
		return plugins.Wrap(cmd, err)
	}

	return nil
}

func (cmd *Latex) executeFolder(ctx context.Context, cab fs.FS, pwd string) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	wg := errgroup.Group{}

	err := fs.WalkDir(cab, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(path)

		if base != "module.md" {
			return nil
		}

		fmt.Fprintln(cmd.Stdout(), "Processing: ", path)
		dir := filepath.Dir(path)

		cab, err := fs.Sub(cab, dir)
		if err != nil {
			return err
		}

		op := filepath.Join(cmd.OutputPath, dir)

		man := manifest{
			cab: cab,
			dir: dir,
			op:  op,
			fn:  base,
			pwd: filepath.Join(pwd, dir),
		}

		wg.Go(func() error {
			return cmd.executeFile(ctx, man)
		})

		return filepath.SkipDir
	})

	if err != nil {
		return err
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	return nil
}

func (cmd *Latex) executeFile(ctx context.Context, man manifest) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	cab := man.cab

	p, err := NewParser(cab, cmd.ContextPath, man.pwd)
	p.Root = man.pwd

	if err != nil {
		return err
	}

	doc, err := p.ParseExecuteFile(ctx, man.fn)
	if err != nil {
		return err
	}

	tex, err := hytex.Convert(ctx, doc)
	if err != nil {
		return err
	}

	return cmd.dump(tex, man.op)
}

func (cmd *Latex) dump(tex fs.FS, op string) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	if tex == nil {
		return fmt.Errorf("tex is nil")
	}

	var files []string
	err := fs.WalkDir(tex, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Base(path) == "module.tex" {
			files = append(files, path)
			return nil
		}

		rx, err := regexp.Compile(`(^|/)assets/`)
		if err != nil {
			return err
		}

		if rx.MatchString(path) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	sort.Strings(files)

	if len(op) == 0 {
		op = "latex"
	}

	for _, file := range files {
		err = cmd.dumpFile(tex, op, file)

		if err != nil {
			return err
		}
	}

	return nil
}
func (cmd *Latex) dumpFile(cab fs.FS, op string, file string) error {
	src, err := cab.Open(file)
	if err != nil {
		return err
	}

	defer src.Close()

	if file == "module.tex" && cmd.FolderName {
		file = filepath.Base(op) + ".tex"
	}

	fp := filepath.Join(op, file)

	err = os.MkdirAll(filepath.Dir(fp), 0755)
	if err != nil {
		return err
	}

	des, err := os.Create(fp)
	if err != nil {
		return err
	}

	defer des.Close()

	_, err = io.Copy(des, src)

	return err
}

type manifest struct {
	cab fs.FS
	pwd string // "pwd" the cmd is executed in (c:/x/y/z)
	op  string // "op" the output path for this manifest
	dir string // "dir" the directory of this manifest ("a/b/c")
	fn  string // "fn" the filename of this manifest ("module.md")
}

func (m manifest) String() string {
	b, _ := json.Marshal(map[string]any{
		"pwd": m.pwd,
		"op":  m.op,
		"dir": m.dir,
		"fn":  m.fn,
	})

	return string(b)
}
