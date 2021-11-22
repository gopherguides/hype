package commander

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var moot sync.Mutex

func Run(ctx context.Context, root string, env []string, name string, args ...string) (Result, error) {
	moot.Lock()
	defer moot.Unlock()
	var err error
	// fmt.Println("root:", root)
	// fmt.Println("name:", name)
	// fmt.Printf("args:%q\n", args)

	nwd := root
	if ext := filepath.Ext(root); len(ext) > 0 {
		nwd = filepath.Dir(root)
	}

	nwd, _ = filepath.Abs(nwd)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	c := exec.CommandContext(ctx, name, args...)
	c.Stdout = stdout
	c.Stderr = stderr
	c.Env = append(os.Environ(), env...)
	c.Dir = nwd

	r := Result{
		Root: root,
		args: c.Args,
	}

	start := time.Now()
	// fmt.Println(">", c.Args)
	err = c.Run()
	r.Duration = time.Since(start)

	r.Err = err
	r.ExitCode = c.ProcessState.ExitCode()
	r.stderr = stderr.Bytes()
	r.stdout = stdout.Bytes()

	r.Pwd = nwd

	sch := []byte(r.Pwd)
	rpl := []byte(".")

	r.stderr = bytes.ReplaceAll(r.stderr, sch, rpl)
	r.stdout = bytes.ReplaceAll(r.stdout, sch, rpl)

	return r, nil
}
