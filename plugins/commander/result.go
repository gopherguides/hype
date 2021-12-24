package commander

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/markbates/hepa"
	"github.com/markbates/hepa/filters"
)

// Result is the result of running a command.
type Result struct {
	Duration time.Duration // time to run the command
	Err      error         // error from running command
	ExitCode int           // exit code
	Pwd      string        // where it was actually run
	Root     string        // directory where the command was run
	args     []string
	stderr   []byte
	stdout   []byte
}

// Args returns a copy of the args used to run the command.
func (r Result) Args() []string {
	args := make([]string, len(r.args))
	copy(args, r.args)

	return args
}

// CmdString returns a string representation of the command that was run.
//
// Example:
//	$ go run main.go
func (r Result) CmdString() string {
	if len(r.args) == 0 {
		return ""
	}

	return fmt.Sprintf("$ %s", strings.Join(r.args, " "))
}

// Stdout returns the results from stdout as a Reader.
func (r Result) Stdout() io.Reader {
	return bytes.NewReader(r.stdout)
}

// Stderr returns the results from stderr as a Reader.
func (r Result) Stderr() io.Reader {
	return bytes.NewReader(r.stderr)
}

// String returns the CmdString.
func (r Result) String() string {
	return r.CmdString()
}

func (r Result) sep(w io.Writer) {
	for i := 0; i < 80; i++ {
		fmt.Fprint(w, "-")
	}
	fmt.Fprintln(w)
}

// ToHTML returns the results as an HTML tag.
func (r Result) ToHTML(ats Attributes, data Data) (string, error) {
	bb := &bytes.Buffer{}

	args := r.Args()

	if len(args) > 0 && !ats.HasKeys("hide-cmd") {
		fmt.Fprintf(bb, "%s\n\n", r.CmdString())
	}

	if len(r.stdout) > 0 && !ats.HasKeys("hide-stdout") {
		line := fmt.Sprintf("%s\n", r.stdout)
		fmt.Fprint(bb, line)
	}

	if len(r.stderr) > 0 && !ats.HasKeys("hide-stderr") {
		r.sep(bb)
		line := fmt.Sprintf("STDERR:\n\n%s\n", r.stderr)
		fmt.Fprint(bb, line)
	}

	pd := map[string]string{}
	if !ats.HasKeys("hide-data") {
		for k, v := range data {
			if !ats.HasKeys("hide-" + k) {
				pd[k] = v
			}
		}
	}

	if len(pd) > 0 {
		r.sep(bb)
		keys := make([]string, 0, len(data))
		for k := range pd {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		w := tabwriter.NewWriter(bb, 0, 0, 1, ' ', tabwriter.StripEscape)

		for _, k := range keys {
			fmt.Fprintf(w, "%s:\t%s\n", k, pd[k])
		}

		w.Flush()
	}

	root := r.Root
	if root, err := filepath.Abs(root); err == nil {
		r.Root = root
	}

	pure := &hepa.Purifier{}
	pure = hepa.With(pure, filters.Replace(r.Root, "."))
	pure = hepa.With(pure, filters.PWD())
	pure = hepa.With(pure, filters.Secrets())
	pure = hepa.With(pure, filters.Golang())
	pure = hepa.With(pure, filters.Replace("willmark/", ""))

	b, err := pure.Clean(bb)
	if err != nil {
		return "", err
	}

	s := strings.TrimSpace(string(b))

	return s, nil
}
