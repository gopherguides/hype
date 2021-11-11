package golang

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/gopherguides/hype"
	"golang.org/x/net/html"
)

var _ hype.Tag = &GoRun{}

type GoRun struct {
	*hype.Node
	root string
	args []string
}

func (gd GoRun) Source() (hype.Source, bool) {
	return hype.SrcAttr(gd.Attrs())
}

func (gd GoRun) StartTag() string {
	bb := &bytes.Buffer{}
	fmt.Fprint(bb, gd.Node.StartTag())
	fmt.Fprint(bb, `<pre><code language="output" class="language-output">`)
	return bb.String()
}

func (gd GoRun) EndTag() string {
	bb := &bytes.Buffer{}
	fmt.Fprint(bb, "</code></pre>")
	fmt.Fprint(bb, gd.Node.EndTag())
	return bb.String()
}

func (gd GoRun) String() string {
	bb := &bytes.Buffer{}
	fmt.Fprint(bb, gd.StartTag())
	fmt.Fprint(bb, gd.Children.String())
	fmt.Fprint(bb, gd.EndTag())
	return bb.String()
}

func (gd *GoRun) Validate(checks ...hype.ValidatorFn) error {
	if gd == nil {
		return fmt.Errorf("run is nil")
	}

	_, ok := hype.TagSource(gd)
	if !ok {
		return fmt.Errorf("%s is not a tag source %v", gd.Atom(), gd)
	}

	checks = append(checks, hype.AtomValidator(GORUN))
	return gd.Node.Validate(html.ElementNode, checks...)
}

func NewGoRun(node *hype.Node, root string) (*GoRun, error) {

	gd := &GoRun{
		Node: node,
		root: root,
	}

	if err := gd.Validate(); err != nil {
		return nil, err
	}

	ctx := context.Background()

	if f, err := gd.Get("args"); err == nil {
		gd.args = CleanFlags(strings.Split(f, ",")...)
	}

	source, err := gd.Get("src")
	if err != nil {
		return nil, err
	}

	fp := filepath.Join(root, source)

	if _, err := os.Stat(fp); err != nil {
		return nil, err
	}

	paths := []string{"."}
	if len(filepath.Ext(fp)) > 0 {
		fp = filepath.Dir(fp)
	}

	if files, err := gd.Get("files"); err == nil {
		paths = strings.Split(files, ",")
	}

	args := append(gd.args, paths...)

	cmd := NewCommander("run", args...)

	key := gd.key(cmd, fp)
	os.MkdirAll(filepath.Dir(key), 0755)

	stdout, err := os.Create(key + ".stdout")
	if err != nil {
		return nil, err
	}
	defer stdout.Close()

	cmd.IO = WithOut(cmd.IO, stdout)

	stderr, err := os.Create(key + ".stderr")
	if err != nil {
		return nil, err
	}
	defer stderr.Close()

	cmd.IO = WithErr(cmd.IO, stderr)

	bb := &bytes.Buffer{}
	switch gd.Attrs()["out"] {
	case "stderr":
		cmd.IO = WithErr(cmd.IO, io.MultiWriter(stderr, bb))
	default:
		cmd.IO = WithOut(cmd.IO, io.MultiWriter(stdout, bb))
	}

	for _, w := range []io.Writer{cmd.IO.Out(), cmd.IO.Err()} {
		fmt.Fprintf(w, "$ %s\n", cmd.String())
		fmt.Fprintf(w, "\t%s\n", runtime.Version())
		fmt.Fprintf(w, "\tpwd: %s\n\n", strings.TrimPrefix(fp, root))
	}

	// cmd.IO = WithOut(cmd.IO, io.MultiWriter(stdout, bb))
	err = cmd.Run(ctx, fp)
	if err != nil {
		return nil, err
	}

	gd.Children = hype.Tags{
		hype.QuickText(bb.String()),
	}
	return gd, nil
}

func (gd *GoRun) key(cmd *Commander, root string) string {
	k := cmd.String()
	k = flect.Underscore(k)
	k = filepath.Join(root, "output", runtime.Version(), k)
	return k
}
