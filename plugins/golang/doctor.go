package golang

import (
	"context"
	"fmt"
	"net/url"
	"path"
)

func (d *Godoc) Package(src string) (string, error) {
	pkg, err := url.Parse("//" + src)
	if err != nil {
		return "", err
	}

	x := path.Join(pkg.Host, pkg.Path)
	return x, nil
}

func (d *Godoc) goGet(ctx context.Context, src string) error {
	if d == nil {
		return fmt.Errorf("doctor is nil")
	}

	pkg, err := d.Package(src)
	if err != nil {
		return err
	}

	return execute(ctx, &StdIO{}, "get", pkg)
}

func (d *Godoc) tidy(ctx context.Context) error {
	if d == nil {
		return fmt.Errorf("doctor is nil")
	}

	return execute(ctx, &StdIO{}, "mod", "tidy")
}
