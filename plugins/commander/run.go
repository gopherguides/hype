package commander

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Runner executes a command.
type Runner struct {
	Args    []string      // args to pass to the command
	Env     []string      // env vars to pass to the command
	Name    string        // name of the command
	Root    string        // root directory to run the command in
	Timeout time.Duration // timeout for the command
	sync.RWMutex
}

func (r Runner) CmdString() string {
	if len(r.Args) == 0 {
		return r.Name
	}

	return fmt.Sprintf("$ %s %s", r.Name, strings.Join(r.Args, " "))
}

func (r *Runner) Run(ctx context.Context, exp int) (Result, error) {
	runDir := r.Root
	env := r.Env
	name := r.Name
	args := r.Args

	_, err := exec.LookPath(name)
	if err != nil {
		return Result{}, err
	}

	if ext := filepath.Ext(runDir); len(ext) > 0 {
		runDir = filepath.Dir(runDir)
	}

	runDir, err = filepath.Abs(runDir)
	if err != nil {
		return Result{}, err
	}

	var cancel context.CancelFunc

	if r.Timeout == 0 {
		r.Timeout = time.Second * 5
	}

	ctx, cancel = context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	c := exec.CommandContext(ctx, name, args...)
	c.Dir = runDir
	c.Env = append(os.Environ(), env...)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	c.Stdout = stdout
	c.Stderr = stderr

	res := Result{
		Pwd:  runDir,
		Root: r.Root,
		args: c.Args,
	}

	start := time.Now()

	err = c.Run()

	res.Duration = time.Since(start)

	res.Err = err

	if c.ProcessState != nil {
		res.ExitCode = c.ProcessState.ExitCode()
	}

	res.stderr = stderr.Bytes()
	res.stdout = stdout.Bytes()

	sch := []byte(res.Pwd)
	rpl := []byte(".")

	res.stderr = bytes.ReplaceAll(res.stderr, sch, rpl)
	res.stdout = bytes.ReplaceAll(res.stdout, sch, rpl)

	if res.ExitCode != exp {

		return res, fmt.Errorf("expected exit code %d, got %d", exp, res.ExitCode)

	}

	return res, nil
}

// var moot sync.Mutex

// func Run(ctx context.Context, root string, env []string, name string, args ...string) (Result, error) {
// 	moot.Lock()
// 	defer moot.Unlock()

// 	var err error

// 	runDir := root
// 	if ext := filepath.Ext(root); len(ext) > 0 {
// 		runDir = filepath.Dir(root)
// 	}

// 	runDir, _ = filepath.Abs(runDir)

// 	c := exec.CommandContext(ctx, name, args...)
// 	c.Dir = runDir
// 	c.Env = append(os.Environ(), env...)

// 	stdout := &bytes.Buffer{}
// 	stderr := &bytes.Buffer{}

// 	c.Stdout = stdout
// 	c.Stderr = stderr

// 	r := Result{
// 		Pwd:  runDir,
// 		Root: root,
// 		args: c.Args,
// 	}

// 	start := time.Now()

// 	err = c.Run()

// 	r.Duration = time.Since(start)

// 	r.Err = err

// 	if c.ProcessState != nil {
// 		r.ExitCode = c.ProcessState.ExitCode()
// 	}

// 	r.stderr = stderr.Bytes()
// 	r.stdout = stdout.Bytes()

// 	sch := []byte(r.Pwd)
// 	rpl := []byte(".")

// 	r.stderr = bytes.ReplaceAll(r.stderr, sch, rpl)
// 	r.stdout = bytes.ReplaceAll(r.stdout, sch, rpl)

// 	return r, nil
// }
