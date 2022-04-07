package golang

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/atomx"
	"github.com/gopherguides/hype/plugins/commander"
	"github.com/mattn/go-shellwords"
)

// go doc -cmd -u -src -short User.String

type Symbol struct {
	*hype.FencedCode
}

func (sym *Symbol) StartTag() string {
	return fmt.Sprintf("<pre>%s", sym.Node.StartTag())
}

func (sym *Symbol) EndTag() string {
	return fmt.Sprintf("%s</pre>", sym.Node.EndTag())
}

func (sym *Symbol) Source() (hype.Source, bool) {
	if sym == nil || sym.Node == nil {
		return "", false
	}

	src, err := sym.Get("src")
	if err != nil {
		return "", false
	}

	return hype.Source(src), true
}

func (sym *Symbol) SetSource(src string) {
	if sym == nil || sym.Node == nil {
		return
	}

	sym.Set("src", src)
}

func (src *Symbol) Finalize(p *hype.Parser) error {

	fp, err := src.Get("src")
	if err != nil {
		return err
	}

	fp = filepath.Join(p.Root, fp)

	ag := []string{"doc", "-cmd", "-u", "-src", "-short"}

	s, err := src.Get("sym")
	if err != nil {
		return err
	}

	words, err := shellwords.Parse(s)
	if err != nil {
		return err
	}

	ag = append(ag, words...)

	jog := &commander.Runner{
		Args:    ag,
		Root:    fp,
		Name:    "go",
		Timeout: time.Duration(time.Second),
	}

	res := commander.Result{
		Root: fp,
		Pwd:  fp,
	}

	res, err = jog.Run(context.Background(), 0)
	if err != nil || res.Err != nil {

		b, e := ioutil.ReadAll(res.Stderr())
		if e != nil {
			return fmt.Errorf("error reading stderr: %w", e)
		}

		if res.Err != nil {
			err = res.Err
		}

		err = fmt.Errorf("error running: %w\n%s", err, string(b))
		return err
	}

	b, err := ioutil.ReadAll(res.Stdout())
	if err != nil {
		return err
	}

	src.Children = hype.Tags{hype.QuickText(string(b))}

	fc, err := hype.NewFencedCode(src.Node)

	if err != nil {
		return err
	}

	src.FencedCode = fc

	return nil
}

func NewSymbol(node *hype.Node) (*Symbol, error) {
	if node == nil {
		return nil, fmt.Errorf("node is nil")
	}

	if _, err := node.Get("src"); err != nil {
		return nil, fmt.Errorf("missing src")
	}

	node.DataAtom = atomx.Code
	node.Set("language", "go")

	sym := &Symbol{
		FencedCode: &hype.FencedCode{
			Node: node,
		},
	}

	return sym, nil
}
