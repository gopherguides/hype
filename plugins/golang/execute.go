package golang

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Commander struct {
	IO   *StdIO
	Name string
	Args []string
}

func NewCommander(name string, args ...string) *Commander {
	return &Commander{
		Name: name,
		Args: args,
	}
}

func (r *Commander) Cmd(ctx context.Context) *exec.Cmd {
	args := append([]string{r.Name}, r.Args...)
	c := exec.CommandContext(ctx, "go", args...)
	c.Stdout = r.IO.Out()
	c.Stderr = r.IO.Err()
	c.Stdin = r.IO.In()
	return c
}

func (r *Commander) String() string {
	return fmt.Sprintf("go %s %s", r.Name, strings.Join(r.Args, " "))
}

func (r *Commander) Run(ctx context.Context, root string) error {
	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)

	if ext := filepath.Ext(root); len(ext) > 0 {
		root = filepath.Dir(root)
	}
	os.Chdir(root)

	bb := &bytes.Buffer{}
	r.IO = WithErr(r.IO, io.MultiWriter(bb, r.IO.Err()))

	c := r.Cmd(ctx)
	err := c.Run()
	// fmt.Println(">", c.String())
	if err != nil {
		s := strings.TrimSpace(bb.String())
		pwd, _ := os.Getwd()
		err = fmt.Errorf("$ %s: %w\n%s", r, err, s)
		return fmt.Errorf("%s: %w", pwd, err)
	}

	return nil
}

func execute(ctx context.Context, std *StdIO, name string, args ...string) error {
	r := NewCommander(name, args...)
	r.IO = std

	c := r.Cmd(ctx)

	bb := &bytes.Buffer{}

	c.Stdout = bb

	err := c.Run()

	if err != nil {
		return fmt.Errorf("$ %s: %w", r, err)
	}

	s := html.EscapeString(bb.String())
	io.Copy(std.Out(), strings.NewReader(s))

	return nil
}
