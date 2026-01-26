package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/preview"
	"github.com/gopherguides/hype/themes"
	"github.com/markbates/cleo"
	"github.com/markbates/plugins"
)

var _ plugins.Needer = &Preview{}

type Preview struct {
	cleo.Cmd

	File          string
	Port          int
	WatchDirs     stringSlice
	Extensions    string
	IncludeGlobs  stringSlice
	ExcludeGlobs  stringSlice
	DebounceDelay time.Duration
	Verbose       bool
	OpenBrowser   bool
	Theme         string
	CustomCSS     string
	ListThemes    bool
	Timeout       time.Duration

	Parser *hype.Parser

	flags *flag.FlagSet
	mu    sync.RWMutex
}

type stringSlice []string

func (s *stringSlice) String() string {
	if s == nil {
		return ""
	}
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (cmd *Preview) WithPlugins(fn plugins.FeederFn) error {
	if cmd == nil {
		return fmt.Errorf("preview is nil")
	}

	if fn == nil {
		return fmt.Errorf("fn is nil")
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	cmd.Feeder = fn

	return nil
}

func (cmd *Preview) ScopedPlugins() plugins.Plugins {
	if cmd == nil {
		return nil
	}

	return cmd.Cmd.ScopedPlugins()
}

func (cmd *Preview) SetParser(p *hype.Parser) error {
	if cmd == nil {
		return fmt.Errorf("preview is nil")
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	cmd.Parser = p
	return nil
}

func (cmd *Preview) Flags(stderr io.Writer) (*flag.FlagSet, error) {
	usage := `
Usage: hype preview [options]

Starts a live preview server with file watching and auto-reload.

Available themes: ` + strings.Join(themes.ListThemes(), ", ") + `

Examples:
    hype preview -f hype.md
    hype preview -f hype.md -port 8080 -theme solarized-dark
    hype preview -f hype.md -css ./custom.css
    hype preview -f hype.md -w ./src -w ./images
    hype preview -f hype.md -e md,html,go,png,jpg
    hype preview -f hype.md -i "**/*.md" -i "**/*.go"
    hype preview -f hype.md -x "**/vendor/**" -x "**/tmp/**"
    hype preview -themes
`

	if err := cmd.validate(); err != nil {
		return nil, err
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	if cmd.flags != nil {
		return cmd.flags, nil
	}

	cmd.flags = flag.NewFlagSet("preview", flag.ContinueOnError)
	cmd.flags.SetOutput(stderr)

	cmd.flags.StringVar(&cmd.File, "f", "hype.md", "markdown file to preview")
	cmd.flags.IntVar(&cmd.Port, "port", 3000, "port for the preview server")
	cmd.flags.Var(&cmd.WatchDirs, "w", "directories to watch (repeatable)")
	cmd.flags.Var(&cmd.WatchDirs, "watch", "directories to watch (repeatable)")
	cmd.flags.StringVar(&cmd.Extensions, "e", "", "file extensions to watch (comma-separated, e.g., md,html,go)")
	cmd.flags.StringVar(&cmd.Extensions, "ext", "", "file extensions to watch (comma-separated)")
	cmd.flags.Var(&cmd.IncludeGlobs, "i", "glob patterns to include (repeatable)")
	cmd.flags.Var(&cmd.IncludeGlobs, "include", "glob patterns to include (repeatable)")
	cmd.flags.Var(&cmd.ExcludeGlobs, "x", "glob patterns to exclude (repeatable)")
	cmd.flags.Var(&cmd.ExcludeGlobs, "exclude", "glob patterns to exclude (repeatable)")
	cmd.flags.DurationVar(&cmd.DebounceDelay, "d", 300*time.Millisecond, "debounce delay before rebuild")
	cmd.flags.DurationVar(&cmd.DebounceDelay, "debounce", 300*time.Millisecond, "debounce delay before rebuild")
	cmd.flags.BoolVar(&cmd.Verbose, "v", false, "verbose output (log file changes)")
	cmd.flags.BoolVar(&cmd.Verbose, "verbose", false, "verbose output (log file changes)")
	cmd.flags.BoolVar(&cmd.OpenBrowser, "open", false, "auto-open browser on start")
	cmd.flags.StringVar(&cmd.Theme, "theme", themes.DefaultTheme, "preview theme name")
	cmd.flags.StringVar(&cmd.CustomCSS, "css", "", "path to custom CSS file (overrides -theme)")
	cmd.flags.BoolVar(&cmd.ListThemes, "themes", false, "list available themes and exit")
	cmd.flags.DurationVar(&cmd.Timeout, "timeout", 0, "timeout for document execution (0 = no timeout)")

	cmd.flags.Usage = func() {
		_, _ = fmt.Fprintf(stderr, "Usage of %s:\n", os.Args[0])
		cmd.flags.PrintDefaults()
		_, _ = fmt.Fprintln(stderr, usage)
	}

	return cmd.flags, nil
}

func (cmd *Preview) Main(ctx context.Context, pwd string, args []string) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	if err := (&cmd.Cmd).Init(); err != nil {
		return err
	}

	flags, err := cmd.Flags(cmd.Stderr())
	if err != nil {
		return err
	}

	if err := flags.Parse(args); err != nil {
		return err
	}

	if cmd.ListThemes {
		_, _ = fmt.Fprintln(cmd.Stdout(), "Available themes:")
		for _, t := range themes.ListThemes() {
			_, _ = fmt.Fprintf(cmd.Stdout(), "  %s\n", t)
		}
		return nil
	}

	if cmd.Theme != "" && cmd.CustomCSS == "" && !themes.IsBuiltinTheme(cmd.Theme) {
		return fmt.Errorf("unknown theme: %s (use -themes to list available themes)", cmd.Theme)
	}

	cfg := preview.DefaultConfig()
	cfg.File = cmd.File
	cfg.Port = cmd.Port
	cfg.Verbose = cmd.Verbose
	cfg.OpenBrowser = cmd.OpenBrowser
	cfg.Theme = cmd.Theme
	cfg.CustomCSS = cmd.CustomCSS
	cfg.DebounceDelay = cmd.DebounceDelay

	if len(cmd.WatchDirs) > 0 {
		cfg.WatchDirs = cmd.WatchDirs
	}

	if cmd.Extensions != "" {
		cfg.Extensions = strings.Split(cmd.Extensions, ",")
		for i, ext := range cfg.Extensions {
			cfg.Extensions[i] = strings.TrimSpace(ext)
		}
	}

	if len(cmd.IncludeGlobs) > 0 {
		cfg.IncludeGlobs = cmd.IncludeGlobs
	}

	if len(cmd.ExcludeGlobs) > 0 {
		cfg.ExcludeGlobs = append(cfg.ExcludeGlobs, cmd.ExcludeGlobs...)
	}

	if cmd.Timeout > 0 {
		cfg.Timeout = cmd.Timeout
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		_, _ = fmt.Fprintf(cmd.Stdout(), "\nShutting down preview server...\n")
		cancel()
	}()

	srv := preview.New(cfg, cmd.Parser)
	srv.SetOutput(
		func(format string, args ...any) {
			_, _ = fmt.Fprintf(cmd.Stdout(), format, args...)
		},
		func(format string, args ...any) {
			_, _ = fmt.Fprintf(cmd.Stderr(), format, args...)
		},
	)

	if cmd.OpenBrowser {
		go func() {
			time.Sleep(500 * time.Millisecond)
			_ = openBrowser(fmt.Sprintf("http://localhost:%d", cfg.Port))
		}()
	}

	return srv.Run(ctx, pwd)
}

func (cmd *Preview) validate() error {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	return nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}
