package commander

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/htmx"
	"github.com/mattn/go-shellwords"
	"golang.org/x/net/html"
)

var _ hype.Tag = &Cmd{}
var _ hype.SetSourceable = &Cmd{}
var _ hype.Validatable = &Cmd{}

// Cmd can be used to run a command.
// HTML Attributes:
// 	exec (required): The command to run.
// 		`<cmd exec="tree -a"></cmd>`
// 	src: The directory for the command to be run in.
// 		`<cmd src="./foo/bar"></cmd>`
// 	environ: A comma-separated list of environment variables to pass to the command.
// 		`<cmd environ="key=value,key2=value2"></cmd>`
// 	exit: The expected exit code. Defaults to 0.
// 	timeout: The timeout duration for the command. Defaults to 30s
// 		`<cmd timeout="3s"></cmd>`
// 	no-cache: If true, the results of the command will not be cached.
// 		`<cmd no-cache></cmd>`
// 	hide-cmd: If true, the command will not be displayed.
// 		`<cmd hide-cmd></cmd>`
// 	hide-stdout: If true, the stdout of the command will not be displayed.
// 		`<cmd hide-stdout></cmd>`
// 	hide-stderr: If true, the stderr of the command will not be displayed.
// 		`<cmd hide-stderr></cmd>`
// 	hide-data: If true, no additional data will be displayed.
// 		`<cmd hide-data></cmd>`
type Cmd struct {
	*hype.Node
	ExpectedExit int      // Expect exit code.
	Args         []string // Arguments to pass to the command.
	Env          []string // Environment variables to pass to the command.
}

// Source returns the src attribute of the tag.
func (c *Cmd) Source() (hype.Source, bool) {
	return hype.SrcAttr(c.Attrs())
}

// SetSource sets the src attribute of the tag.
func (c *Cmd) SetSource(src string) {
	c.Set("src", src)
}

// Finalize runs the command and saves the output.
// This is called *after* the parsing has completed.
func (c *Cmd) Finalize(p *hype.Parser) error {
	return c.work(p)
}

// StartTag returns the start tag of the command.
func (c *Cmd) StartTag() string {
	return c.Node.StartTag() + `<pre class="code-block"><code class="language-plain" language="plain">`
}

// EndTag returns the end tag of the command.
func (c *Cmd) EndTag() string {
	return "</code></pre>" + c.Node.EndTag()
}

// String returns an HTML representation of the command.
func (c *Cmd) String() string {
	bb := &bytes.Buffer{}
	fmt.Fprint(bb, c.StartTag())
	fmt.Fprint(bb, c.Children.String())
	fmt.Fprint(bb, c.EndTag())
	return bb.String()
}

// Validate returns an error if the command is invalid.
func (c *Cmd) Validate(p *hype.Parser, checks ...hype.ValidatorFn) error {
	checks = append(checks, hype.AtomValidator("cmd"))
	return c.Node.Validate(p, html.ElementNode, checks...)
}

// NewCmd returns a new Cmd from the given node.
// The command is *not* run until Finalize is called.
// If the `code` attribute is set, then a `hype.Element`
// is created and returned containing the code and the
// command.
func NewCmd(cab fs.FS, node *hype.Node) (hype.Tag, error) {
	cmd := &Cmd{
		Node: node,
	}

	if err := cmd.Validate(nil); err != nil {
		return nil, err
	}

	ats := cmd.Attrs()

	if env, ok := ats["environ"]; ok {
		cmd.Env = strings.Split(env, ",")
	}

	if ex, ok := ats["exit"]; ok {
		i, err := strconv.Atoi(ex)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", cmd.StartTag(), err)
		}
		cmd.ExpectedExit = i
	}

	if _, ok := ats["code"]; ok {
		return newCodeCmd(cab, cmd)
	}

	return cmd, cmd.Validate(nil)
}

func newCodeCmd(cab fs.FS, cmd *Cmd) (*hype.Element, error) {
	if cmd == nil {
		return nil, errors.New("cmd is nil")
	}

	el := &hype.Element{
		Node: hype.NewNode(htmx.ElementNode("div")),
	}

	ats := cmd.Attrs()

	src, ok := ats["src"]
	if !ok {
		return nil, fmt.Errorf("missing src attribute")
	}

	code, ok := ats["code"]
	if !ok {
		return nil, errors.New("code is not set")
	}

	cs := strings.Split(code, ",")
	for i, c := range cs {
		cs[i] = filepath.Join(src, c)
	}

	ats["src"] = strings.Join(cs, ",")

	hn := hype.NewNode(htmx.AttrNode("code", ats))

	sc, err := hype.NewSourceCode(cab, hn, nil)
	if err != nil {
		return nil, err
	}
	el.Children = append(el.Children, sc)

	el.Children = append(el.Children, cmd)
	return el, nil
}

func (cmd *Cmd) work(p *hype.Parser) error {
	if p == nil {
		return fmt.Errorf("parser is nil")
	}

	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	root := p.Root

	data := Data{}

	ex, err := cmd.Get("exec")
	if err != nil {
		return fmt.Errorf("%s: %w", cmd.StartTag(), err)
	}

	args, err := shellwords.Parse(ex)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("exec is empty")
	}

	cmd.Args = args

	ats := cmd.Attrs()
	src := ats["src"]

	e := strings.Join(cmd.Env, ",")
	e = strings.TrimSpace(e)
	if len(e) > 0 {
		data["env"] = e
	}

	name := cmd.Args[0]
	if x, ok := cmds[name]; ok {
		for k, v := range x {
			ats[k] = v
		}
	}

	if ats.HasKeys("data-go") || name == "go" {
		data["go"] = runtime.Version()
	}

	fp := filepath.Join(root, src)
	sum, err := hash(fp)
	if err != nil {
		return err
	}

	timeout := 5 * time.Second
	if to, ok := ats["timeout"]; ok {
		d, err := time.ParseDuration(to)
		if err != nil {
			return err
		}
		timeout = d
	}

	ctx := context.Background()

	var ag []string
	if len(cmd.Args) > 1 {
		ag = cmd.Args[1:]
	}

	jog := &Runner{
		Args:    ag,
		Root:    fp,
		Env:     cmd.Env,
		Name:    name,
		Timeout: timeout,
	}

	cache := &Cache{
		Command:   jog.CmdString(),
		Exit:      cmd.ExpectedExit,
		GoVersion: runtime.Version(),
		Src:       fp,
		Tag:       cmd.StartTag(),
		Sum:       sum,
	}

	if !ats.HasKeys("no-cache") && p.DB != nil {

		exists, err := func() (bool, error) {
			err := cache.Fetch(p.DB)
			if errors.Is(err, sql.ErrNoRows) {
				return false, nil
			}

			if err != nil {
				return false, err
			}

			return true, nil
		}()

		if err != nil {
			return err
		}

		if exists {
			cmd.Children = hype.Tags{hype.QuickText(cache.Body)}
			return cmd.Validate(p)
		}

	}

	res := Result{
		ExitCode: cmd.ExpectedExit,
		Root:     fp,
		Pwd:      src,
		Sum:      sum,
	}
	res, err = jog.Run(ctx, cmd.ExpectedExit)

	if err != nil {
		bb := &bytes.Buffer{}
		fmt.Fprintln(bb, cmd.Node.StartTag())
		fmt.Fprintln(bb, "\n-------")
		fmt.Fprintf(bb, "File Name:\t%q\n", p.FileName)
		fmt.Fprintf(bb, "Command:\t%q\n", res.CmdString())
		fmt.Fprintf(bb, "Duration:\t%q\n", res.Duration.String())
		fmt.Fprintf(bb, "Exit Code:\t%d\n", res.ExitCode)
		fmt.Fprintf(bb, "PWD:\t\t%q\n", res.Pwd)
		fmt.Fprintf(bb, "Root:\t\t%q\n", res.Root)
		if len(data) > 0 {
			fmt.Fprintln(bb, "Data:")
			for k, v := range data {
				fmt.Fprintf(bb, "\t%q:\t%q\n", k, v)
			}
		}
		fmt.Fprintln(bb, "\n-------")
		fmt.Fprintln(bb, "STDOUT:")
		io.Copy(bb, res.Stdout())
		fmt.Fprintln(bb, "\n-------")
		fmt.Fprintln(bb, "STDERR:")
		io.Copy(bb, res.Stderr())
		fmt.Fprintln(bb, "\n--------")
		fmt.Fprintln(bb, "ENVIRON:")
		for _, e := range os.Environ() {
			fmt.Fprintf(bb, "\t%q\n", e)
		}
		if len(cmd.Env) > 0 {
			fmt.Fprintln(bb, "CUSTOM ENV:")
			for _, e := range cmd.Env {
				fmt.Fprintf(bb, "\t%q\n", e)
			}
		}
		return fmt.Errorf("%w:\n%s", err, bb.String())
	}

	data["duration"] = res.Duration.String()

	s, err := res.ToHTML(ats, data)
	if err != nil {
		return err
	}

	cmd.Children = hype.Tags{hype.QuickText(s)}

	if p.DB == nil {
		return nil
	}

	cache.Body = s

	if err := cache.Insert(p.DB); err != nil {
		return err
	}

	return nil
}
