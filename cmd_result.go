package hype

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gopherguides/hype/atomx"
	"github.com/markbates/clam"
)

// CmdResult is the result of executing a command.
type CmdResult struct {
	*Element
	*clam.Result
}

func NewCmdResult(p *Parser, c *Cmd, res *clam.Result) (*CmdResult, error) {
	if res == nil {
		return nil, ErrIsNil("result")
	}

	cmd := &CmdResult{
		Result:  res,
		Element: NewEl("result", nil),
	}

	cmd.Parent = c

	lang := "text"
	lang = Language(c.Attrs(), lang)

	var lines []string

	_, hcmd := c.Get("hide-cmd")
	if !hcmd {
		// lines = append(lines, fmt.Sprintf("```%s\n%s", lang, res.CmdString()))
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
		return nil, err
	}

	pre := NewEl(atomx.Pre, cmd)
	cel := NewEl(atomx.Code, pre)
	cel.Set("language", lang)
	cel.Set("class", "language-"+lang)
	cel.Nodes = append(cel.Nodes, TextNode(body))

	pre.Nodes = append(pre.Nodes, cel)

	cmd.Nodes = Nodes{pre}

	return cmd, nil
}

func resultBody(res *clam.Result, ats *Attributes, body string) (string, error) {
	body = strings.TrimSpace(body)

	if len(res.Dir) > 0 {
		body = strings.ReplaceAll(body, res.Dir, ".")
	}

	if pwd, err := os.Getwd(); err == nil {
		body = strings.ReplaceAll(body, fmt.Sprintf("%s%s", pwd, string(filepath.Separator)), "")
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
