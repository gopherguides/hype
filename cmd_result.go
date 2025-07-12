package hype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/gobuffalo/flect"
	"github.com/markbates/clam"
)

// CmdResult is the result of executing a command.
type CmdResult struct {
	*Element
	*clam.Result
}

func (c *CmdResult) MarshalJSON() ([]byte, error) {
	if c == nil {
		return nil, ErrIsNil("cmd result")
	}

	c.RLock()
	defer c.RUnlock()

	m, err := c.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = toType(c)

	if c.Result == nil {
		return json.MarshalIndent(m, "", "  ")
	}

	m["args"] = c.Args
	m["dir"] = c.Dir
	m["duration"] = c.Duration.String()
	m["err"] = errForJSON(c.Err)
	m["exit"] = c.Exit
	m["stderr"] = string(c.Stderr)
	m["stdout"] = string(c.Stdout)

	return json.MarshalIndent(m, "", "  ")
}

func (c *CmdResult) MD() string {
	if c == nil {
		return ""
	}

	return c.Children().MD()
}

func (c *CmdResult) String() string {
	if c == nil || c.Element == nil {
		return ""
	}
	return c.Children().String()
}

func NewCmdResult(p *Parser, c *Cmd, res *clam.Result) (*CmdResult, error) {
	if res == nil {
		return nil, c.WrapErr(ErrIsNil("result"))
	}

	cmd := &CmdResult{
		Result:  res,
		Element: NewEl("result", nil),
	}

	cmd.Parent = c

	ats := c.Attrs()

	lang := "shell"
	lang = Language(ats, lang)

	var lines []string

	_, hcmd := c.Get("hide-cmd")
	if !hcmd {
		lines = append(lines, res.CmdString())
	}

	if len(res.Stdout) > 0 {
		s := string(res.Stdout)
		s = strings.TrimSpace(s)
		lines = append(lines, s)
	}

	if len(res.Stderr) > 0 {
		s := string(res.Stderr)
		s = strings.TrimSpace(s)
		lines = append(lines, s)
	}

	// actual body content:
	body := strings.Join(lines, "\n\n")
	body, err := resultBody(res, c.Attrs(), body)
	if err != nil {
		return nil, c.WrapErr(err)
	}

	cel := &FencedCode{
		Element: NewEl("code", cmd),
	}

	if err := cel.Set("language", lang); err != nil {
		return nil, err
	}

	if err := cel.Set("class", "language-"+lang); err != nil {
		return nil, err
	}

	body = strings.TrimSpace(body)
	cel.Nodes = append(cel.Nodes, Text(body))

	type dt struct {
		key string
		val string
	}

	datum := []dt{}

	if _, ok := c.Get("show-duration"); ok {
		datum = append(datum, dt{
			key: "Duration",
			val: res.Duration.String(),
		})
	}

	ats.Range(func(k string, v string) bool {
		if !strings.HasPrefix(k, "data-") {
			return true
		}

		v = strings.TrimSpace(v)
		if len(v) == 0 {
			return true
		}

		k = strings.TrimPrefix(k, "data-")
		k = flect.Titleize(k)

		datum = append(datum, dt{
			key: k,
			val: v,
		})

		return true
	})

	if len(datum) > 0 {
		sort.Slice(datum, func(i, j int) bool {
			return datum[i].key < datum[j].key
		})

		bb := &bytes.Buffer{}
		tw := tabwriter.NewWriter(bb, 0, 0, 0, ' ', 0)

		for i := 0; i < 80; i++ {
			fmt.Fprint(tw, "-")
		}

		fmt.Fprintln(tw)

		for _, d := range datum {
			fmt.Fprintf(tw, "%s:\t %s\n", d.key, d.val)
		}

		if err := tw.Flush(); err != nil {
			return nil, cmd.WrapErr(err)
		}

		text := fmt.Sprintf("\n\n%s", bb.String())
		cel.Nodes = append(cel.Nodes, Text(text))
	}

	cmd.Nodes = Nodes{cel}

	return cmd, nil
}

func resultBody(res *clam.Result, ats *Attributes, body string) (string, error) {
	body = strings.TrimSpace(body)

	if len(res.Dir) > 0 {
		body = strings.ReplaceAll(body, res.Dir, ".")
	}

	if pwd, err := os.Getwd(); err == nil {
		fp := fmt.Sprintf("%s%s", pwd, string(filepath.Separator))
		body = strings.ReplaceAll(body, fp, "")
	}

	body = html.EscapeString(body)

	mo, ok := ats.Get("truncate")
	if !ok {
		return body, nil
	}

	max, err := strconv.ParseInt(mo, 0, 64)
	if err != nil {
		return "", err
	}

	lines := make([]string, 0, max)
	for _, l := range strings.Split(body, "\n") {
		if len(lines) >= int(max) {
			lines = append(lines, "...")
			break
		}
		lines = append(lines, l)
	}

	return strings.Join(lines, "\n"), nil
}
