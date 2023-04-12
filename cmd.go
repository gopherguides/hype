package hype

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/markbates/clam"
	"github.com/mattn/go-shellwords"
)

// Cmd is a tag representing a command to be executed.
type Cmd struct {
	*Element

	Args         []string
	Env          []string
	ExpectedExit int
	Timeout      time.Duration

	res *CmdResult
}

func (c *Cmd) MD() string {
	if c == nil {
		return ""
	}

	return c.Children().MD()
}

// Result returns the result of executing the command.
func (c *Cmd) Result() *CmdResult {
	c.RLock()
	defer c.RUnlock()
	return c.res
}

// Execute the command.
func (c *Cmd) Execute(ctx context.Context, doc *Document) error {
	if c == nil {
		return ErrIsNil("cmd")
	}

	if c.Element == nil {
		return ErrIsNil("element")
	}

	if doc == nil {
		return ErrIsNil("document")
	}

	if c.Timeout == 0 {
		c.Timeout = time.Second * 30
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	cmd := &clam.Cmd{
		Env: c.Env,
	}

	src, ok := c.Get("src")
	if ok {
		dir := filepath.Join(doc.Root, src)
		cmd.Dir = dir
	}

	res, err := cmd.Run(ctx, c.Args...)
	if err != nil && res.Exit != c.ExpectedExit {
		return err
	}

	if res.Exit != c.ExpectedExit {
		return fmt.Errorf("unexpected exit code: %d", res.Exit)
	}

	cres, err := NewCmdResult(doc.Parser, c, res)
	if err != nil {
		return err
	}

	c.Lock()
	c.res = cres
	c.Nodes = Nodes{cres}
	c.Unlock()

	return nil
}

func NewCmd(el *Element) (*Cmd, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	c := &Cmd{
		Element: el,
		Timeout: time.Second * 30,
	}

	ex, err := el.ValidAttr("exec")
	if err != nil {
		return nil, err
	}

	args, err := shellwords.Parse(ex)
	if err != nil {
		return nil, err
	}

	if len(args) == 0 {
		return nil, c.WrapErr(fmt.Errorf("no command specified"))
	}

	c.Args = args

	if en, ok := el.Get("environ"); ok {
		c.Env = append(c.Env, strings.Split(en, ",")...)
	}

	if ee, ok := el.Get("exit"); ok {
		c.ExpectedExit, err = strconv.Atoi(ee)
		if err != nil {
			return nil, err
		}
	}

	if to, ok := el.Get("timeout"); ok {
		c.Timeout, err = time.ParseDuration(to)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func NewAttrCode(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	var nodes Nodes

	code, ok := el.Get("code")
	if !ok {
		return nodes, nil
	}

	el.Delete("code")

	src, ok := el.Get("src")
	if ok {
		code = filepath.Join(src, code)
	}

	ats := &Attributes{}

	if err := ats.Set("src", code); err != nil {
		return nil, err
	}

	cel := NewEl("code", el.Parent)
	cel.Attributes = ats
	codes, err := NewCodeNodes(p, cel)

	if err != nil {
		return nil, err
	}

	nodes = append(nodes, codes...)

	nodes = append(nodes, NewEl("hr", nil))

	return nodes, nil
}

func NewCmdNodes(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	nodes, err := NewAttrCode(p, el)
	if err != nil {
		return nil, err
	}

	cmd, err := NewCmd(el)
	if err != nil {
		return nil, err
	}

	nodes = append(nodes, cmd)
	return nodes, nil
}
