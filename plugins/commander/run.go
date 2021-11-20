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

type Runner struct {
	Args []string
	Env  []string
	Name string
	Root string
}

func (r *Runner) Run(ctx context.Context) (Result, error) {
	moot.Lock()
	defer moot.Unlock()
	var err error

	root := r.Root
	if ext := filepath.Ext(root); len(ext) > 0 {
		root = filepath.Dir(root)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	c := exec.CommandContext(ctx, r.Name, r.Args...)
	c.Dir = root
	c.Env = append(os.Environ(), r.Env...)
	c.Stdout = stdout
	c.Stderr = stderr

	res := Result{
		Env:  r.Env,
		Root: root,
		args: c.Args,
	}

	start := time.Now()
	err = c.Run()
	res.Duration = time.Since(start)

	res.Err = err
	res.ExitCode = c.ProcessState.ExitCode()
	res.stderr = stderr.Bytes()
	res.stdout = stdout.Bytes()

	pwd := root
	res.Pwd = pwd
	base := filepath.Base(pwd)

	res.stderr = bytes.ReplaceAll(res.stderr, []byte(pwd), []byte(base))
	res.stdout = bytes.ReplaceAll(res.stdout, []byte(pwd), []byte(base))

	return res, nil
}

var moot sync.Mutex
