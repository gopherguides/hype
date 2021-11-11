package golang

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/flect"
	"golang.org/x/mod/sumdb/dirhash"
)

type Envelope struct {
	Body string
	Cmd  *exec.Cmd
	Doc  string
	Exit int
	Root string
}

func (e Envelope) Hash() (string, error) {
	s, err := dirhash.HashDir(e.Root, "golang-hype", dirhash.DefaultHash)
	if err != nil {
		return s, err
	}

	return strings.TrimSuffix(s, "="), nil
}

func (e Envelope) FilePath() (string, error) {
	s, err := e.Hash()
	if err != nil {
		return s, err
	}

	fp := runtime.Version()

	if len(e.Doc) > 0 {
		d := strings.ReplaceAll(e.Doc, "#", "/")
		fp = filepath.Join(fp, d)
	}

	if e.Cmd != nil {
		a := strings.Join(e.Cmd.Args, " ")
		a = flect.Underscore(a)
		fp = filepath.Join(fp, a)
	}

	fp = filepath.Join(fp, s)

	// m := map[string]interface{}{
	// 	"cmd":  e.Cmd.Args,
	// 	"doc":  e.Doc,
	// 	"root": e.Root,
	// }

	// b, err := json.Marshal(m)
	// if err != nil {
	// 	return s, err
	// }

	// s = hex.EncodeToString(b)
	return fp, nil
}

func (e Envelope) String() string {
	bb := &bytes.Buffer{}

	if len(e.Body) > 0 {
		fmt.Fprintln(bb, strings.TrimSpace(e.Body))
		fmt.Fprintf(bb, "\n--------------------\n")
	}

	fmt.Fprintf(bb, "go version:\t%s\n", runtime.Version())

	if e.Cmd != nil {
		fmt.Fprintf(bb, "command:\t$ %s\n", strings.Join(e.Cmd.Args, " "))
		if e.Exit != 0 {
			fmt.Fprintf(bb, "exit code:\t%d\n", e.Exit)
		}
	}

	if e.Root != "" {
		fmt.Fprintf(bb, "root:\t\t%s\n", e.Root)
	}

	if e.Doc != "" {
		fmt.Fprintf(bb, "docs:\t\t<a href=\"https://pkg.go.dev/%[1]s\" target='_blank'>%[1]s</a>\n", e.Doc)
	}

	return strings.TrimSuffix(bb.String(), "\n")
}
