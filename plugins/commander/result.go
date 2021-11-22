package commander

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/markbates/hepa"
	"github.com/markbates/hepa/filters"
	"github.com/mattn/go-shellwords"
)

type Result struct {
	Duration time.Duration
	Err      error  // error from running command
	ExitCode int    // exit code
	Root     string // directory where the command was run
	Pwd      string // where it was actually run
	args     []string
	stderr   []byte
	stdout   []byte
}

func (r Result) Args() []string {
	args := make([]string, len(r.args))
	copy(args, r.args)

	return args
}

func (r Result) CmdString() string {
	if len(r.args) == 0 {
		return ""
	}

	return fmt.Sprintf("$ %s", strings.Join(r.args, " "))
}

func (r Result) Stdout() io.Reader {
	return bytes.NewReader(r.stdout)
}

func (r Result) Stderr() io.Reader {
	return bytes.NewReader(r.stderr)
}

func (r Result) MarshalJSON() ([]byte, error) {
	x := resultJSON{
		Args:     strings.Join(r.args, " "),
		Exit:     r.ExitCode,
		Root:     r.Root,
		Duration: r.Duration,
		Stderr:   r.stderr,
		Stdout:   r.stdout,
	}

	if r.Err != nil {
		x.Error = r.Err.Error()
	}
	return json.MarshalIndent(x, "", "  ")
}

func (r *Result) UnmarshalJSON(data []byte) error {
	var x resultJSON

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	args, _ := shellwords.Parse(x.Args)

	r.Err = fmt.Errorf(x.Error)
	r.ExitCode = x.Exit
	r.Root = x.Root
	r.args = args
	r.stderr = x.Stderr
	r.stdout = x.Stdout
	r.Duration = time.Duration(x.Duration)

	return nil
}

func (r Result) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r Result) Out(ats Attributes, data Data) (string, error) {
	bb := &bytes.Buffer{}

	args := r.Args()

	if len(args) > 0 && !ats.HasKeys("hide-cmd") {
		fmt.Fprintf(bb, "%s\n\n", r.CmdString())
	}

	if len(r.stdout) > 0 && !ats.HasKeys("hide-stdout") {
		// fmt.Fprintf(bb, "-------\n")
		// line := fmt.Sprintf("STDOUT:\n\n%s\n", r.stdout)
		line := fmt.Sprintf("%s\n", r.stdout)
		fmt.Fprint(bb, line)
	}

	if len(r.stderr) > 0 && !ats.HasKeys("hide-stderr") {
		fmt.Fprintf(bb, "-------\n")
		line := fmt.Sprintf("STDERR:\n\n%s\n", r.stderr)
		fmt.Fprint(bb, line)
	}

	if len(data) > 0 && !ats.HasKeys("hide-data") {
		fmt.Fprintf(bb, "-------\n")

		keys := make([]string, 0, len(data))
		for k := range data {
			if !ats.HasKeys("hide-" + k) {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)

		w := tabwriter.NewWriter(bb, 0, 0, 1, ' ', tabwriter.StripEscape)

		for _, k := range keys {
			fmt.Fprintf(w, "%s:\t%s\n", k, data[k])
		}

		w.Flush()
	}

	pure := &hepa.Purifier{}
	pure = hepa.With(pure, filters.PWD())
	pure = hepa.With(pure, filters.Secrets())
	pure = hepa.With(pure, filters.Golang())

	b, err := pure.Clean(bb)
	if err != nil {
		return "", err
	}

	// b := bb.Bytes()
	s := strings.TrimSpace(string(b))

	return s, nil
}

type resultJSON struct {
	Args     string        `json:"args"`
	Duration time.Duration `json:"elasped"`
	Error    string        `json:"error"`
	Exit     int           `json:"exit"`
	Root     string        `json:"root"`
	Stderr   []byte        `json:"stderr"`
	Stdout   []byte        `json:"stdout"`
}
