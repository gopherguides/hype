package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/markbates/cleo"
)

type App = cleo.Cmd

func New() *App {
	app := &App{
		FS: os.DirFS("."),
	}

	app.Add("preview", &Marked{})
	app.Add("marked", &Marked{})
	return app
}

func SectionFromPath(mp string) (int, error) {
	dir := filepath.Dir(mp)
	base := filepath.Base(dir)
	rx, err := regexp.Compile(`^(\d+)-.+`)
	if err != nil {
		return 0, err
	}

	match := rx.FindStringSubmatch(base)
	if len(match) < 2 {
		return 0, fmt.Errorf("could not find section: %q", mp)
	}

	sec, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}
	return sec, nil
}
