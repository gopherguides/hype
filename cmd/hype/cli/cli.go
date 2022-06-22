package cli

import (
	"os"

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
