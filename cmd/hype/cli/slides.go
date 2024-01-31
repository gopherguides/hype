package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/slides"
	"github.com/markbates/cleo"
	"github.com/markbates/fsx"
	"github.com/markbates/iox"
)

type Slides struct {
	cleo.Cmd

	// Env to be used by the app
	// If nil, os.Getenv will be used.
	*Env

	// Web app to be used by the app
	slides.App

	Parser *hype.Parser

	// Server to be used by the app
	// If nil, a default server will be created.
	Server *http.Server

	// Port to listen on. Defaults to 3000.
	Port int

	mu sync.RWMutex
}

// snippet: main

func (a *Slides) Main(ctx context.Context, pwd string, args []string) error {
	if a == nil {
		return fmt.Errorf("nil app")
	}

	flags := a.flags()
	err := flags.Parse(args)
	if err != nil {
		return err
	}
	if a.Parser != nil {
		a.App.Parser = a.Parser
	}
	a.App.PWD = pwd
	a.App.FileName = "module.md"

	if p := flags.Arg(0); len(p) > 0 {
		a.App.PWD = filepath.Dir(p)
		a.App.FileName = filepath.Base(p)
	}

	a.App.PWD, err = filepath.Abs(a.App.PWD)
	if err != nil {
		return err
	}
	return WithinDir(a.App.PWD, func() error {
		srv, err := a.server()
		if err != nil {
			return err
		}

		ctx, cause := context.WithCancelCause(ctx)
		defer cause(nil)

		srv.BaseContext = func(_ net.Listener) context.Context {
			return ctx
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil {
				cause(err)
			}
		}()

		srvCtx, srvCancel := context.WithCancel(ctx)
		defer srvCancel()

		go func() {
			<-ctx.Done()
			defer srvCancel()
			if err := srv.Shutdown(ctx); err != nil {
				cause(err)
			}
		}()

		<-ctx.Done()

		err = context.Cause(ctx)
		if err != nil && err != context.Canceled {
			return err
		}

		<-srvCtx.Done()

		return nil
	})
}

// snippet: main

func (a *Slides) SetIO(oi iox.IO) {
	if a == nil {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.IO = oi
}

func (a *Slides) Print(w io.Writer) error {
	if a == nil {
		return fmt.Errorf("nil app")
	}

	if w == nil {
		w = os.Stdout
	}

	flags := a.flags()
	flags.SetOutput(w)
	flags.Usage()

	return nil
}

func (a *Slides) Describe() string {
	return "launches a web server"
}

func (a *Slides) Getenv(key string) (s string) {
	if a == nil || a.Env == nil {
		return os.Getenv(key)
	}

	return a.Env.Getenv(key)
}

func (a *Slides) flags() *flag.FlagSet {

	flags := flag.NewFlagSet("server", flag.ContinueOnError)

	flags.SetOutput(a.Stderr())
	flags.IntVar(&a.Port, "port", 3000, "port to listen on")

	return flags
}

func (a *Slides) server() (*http.Server, error) {
	if a == nil {
		return nil, fmt.Errorf("nil app")
	}

	a.mu.RLock()
	srv := a.Server
	port := a.Port
	a.mu.RUnlock()

	mux := http.NewServeMux()
	mux.Handle("/", a)

	cab := &fsx.ArrayFS{}
	cab.Append(slides.AssetsFS)
	cab.Append(a.Parser.FS)

	mux.Handle("/templates/assets/", http.FileServer(http.FS(cab)))
	mux.Handle("/assets/", http.FileServer(http.FS(cab)))

	if srv != nil {
		if srv.Handler == nil {
			srv.Handler = mux
		}
		return srv, nil
	}

	if port == 0 {
		p := a.Getenv("PORT")
		pi, _ := strconv.Atoi(p)
		if pi == 0 {
			pi = 3000
		}
		port = pi
	}

	srv = &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	srv.Handler = mux

	return srv, nil
}
