package commander

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Run(ctx context.Context, root string, name string, args ...string) (Result, error) {
	// fmt.Println("root:", root)
	// fmt.Println("name:", name)
	// fmt.Printf("args:%q\n", args)
	owd, _ := os.Getwd()
	defer os.Chdir(owd)

	nwd := root
	if ext := filepath.Ext(root); len(ext) > 0 {
		nwd = filepath.Dir(root)
	}
	// fmt.Println("nwd:", nwd)
	os.Chdir(nwd)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	c := exec.CommandContext(ctx, name, args...)
	// fmt.Printf("TODO >> run.go:28 c.Args %[1]T %[1]v\n", c.Args)
	c.Stdout = stdout
	c.Stderr = stderr

	r := Result{
		Root: root,
		args: c.Args,
	}

	start := time.Now()
	err := c.Run()
	r.Duration = time.Since(start)

	r.Err = err
	r.ExitCode = c.ProcessState.ExitCode()
	r.stderr = stderr.Bytes()
	r.stdout = stdout.Bytes()

	pwd, _ := os.Getwd()
	// fmt.Printf("TODO >> run.go:44 pwd %[1]T %[1]v\n", pwd)
	base := filepath.Base(pwd)

	r.stderr = bytes.ReplaceAll(r.stderr, []byte(pwd), []byte(base))
	r.stdout = bytes.ReplaceAll(r.stdout, []byte(pwd), []byte(base))

	return r, nil
}
