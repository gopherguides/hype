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

	runDir := root
	if ext := filepath.Ext(root); len(ext) > 0 {
		runDir = filepath.Dir(root)
	}

	runDir, _ = filepath.Abs(runDir)

	c := exec.CommandContext(ctx, name, args...)
	c.Dir = runDir
	c.Env = append(os.Environ(), env...)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	c.Stdout = stdout
	c.Stderr = stderr

	r := Result{
		Pwd:  runDir,
		Root: root,
		args: c.Args,
	}

	start := time.Now()

	err = c.Run()

	r.Duration = time.Since(start)

	r.Err = err

	if c.ProcessState != nil {
		r.ExitCode = c.ProcessState.ExitCode()
	}

	r.stderr = stderr.Bytes()
	r.stdout = stdout.Bytes()

	sch := []byte(r.Pwd)
	rpl := []byte(".")

	r.stderr = bytes.ReplaceAll(r.stderr, sch, rpl)
	r.stdout = bytes.ReplaceAll(r.stdout, sch, rpl)

	return r, nil
}
