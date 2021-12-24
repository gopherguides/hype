package commander

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gopherguides/hype"
	"github.com/mattn/go-shellwords"
	"golang.org/x/net/html"
)

var _ hype.Tag = &Cmd{}
var _ hype.SetSourceable = &Cmd{}

// Cmd can be used to run a command.
// HTML Attributes:
// 	exec (required): The command to run.
// 		`<cmd exec="tree -a"></cmd>`
// 	src: The directory for the command to be run in.
// 		`<cmd src="./foo/bar"></cmd>`
// 	environ: A comma-separated list of environment variables to pass to the command.
// 		`<cmd environ="key=value,key2=value2"></cmd>`
// 	exit: The expected exit code. Defaults to 0.
// 	timeout: The timeout duration for the command. Defaults to infinity
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
	return c.work(p, c.Attrs()["src"])
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
func (c *Cmd) Validate(checks ...hype.ValidatorFn) error {
	checks = append(checks, hype.AtomValidator("cmd"))
	return c.Node.Validate(html.ElementNode, checks...)
}

// NewCmd returns a new Cmd from the given node.
// The command is *not* run until Finalize is called.
//
func NewCmd(node *hype.Node) (*Cmd, error) {
	cmd := &Cmd{
		Node: node,
	}

	if err := cmd.Validate(); err != nil {
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

	return cmd, cmd.Validate()
}

func (cmd *Cmd) work(p *hype.Parser, src string) error {
	if p == nil {
		return fmt.Errorf("parser is nil")
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

	cacheKey, err := cmd.CacheKey(root)
	if err != nil {
		return err
	}

	if !ats.HasKeys("no-cache") {
		if p.Cache != nil {

			b, err := p.Cache.Retrieve(root, cacheKey)
			if err == nil {
				cmd.Children = hype.Tags{hype.QuickText(string(b))}
				return cmd.Validate()
			}
		}
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

	runDir := filepath.Join(root, src)
	jog := &Runner{
		Args:    ag,
		Root:    runDir,
		Env:     cmd.Env,
		Name:    name,
		Timeout: timeout,
	}

	res, err := jog.Run(ctx)

	if err != nil {
		return err
	}

	if res.ExitCode != cmd.ExpectedExit {

		io.Copy(os.Stderr, res.Stderr())
		io.Copy(os.Stdout, res.Stdout())
		return fmt.Errorf("%s: exit code %d != %d", cmd.StartTag(), res.ExitCode, cmd.ExpectedExit)
	}

	data["duration"] = res.Duration.String()

	s, err := res.ToHTML(ats, data)
	if err != nil {
		return err
	}

	cmd.Children = hype.Tags{hype.QuickText(s)}

	if p.Cache == nil {
		return nil
	}

	if err := p.Cache.Store(root, cacheKey, []byte(s)); err != nil {
		return err
	}

	return nil
}

// CacheKey returns a unique key for the command.
func (cmd *Cmd) CacheKey(root string) (string, error) {

	h, err := hash(root)
	if err != nil {
		return "", fmt.Errorf("could not hash %s: %w", root, err)
	}

	tag := cmd.Node.StartTag()

	th := md5.New()
	fmt.Fprint(th, tag)
	hs := fmt.Sprintf("%x", th.Sum(nil))

	s := filepath.Join(h, hs)
	return s, nil
}
