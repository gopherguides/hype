package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gopherguides/hype/blog"
	"github.com/markbates/cleo"
	"github.com/markbates/plugins"
)

var _ plugins.Needer = &Blog{}

type Blog struct {
	cleo.Cmd

	Timeout time.Duration
	Verbose bool

	flags *flag.FlagSet
	mu    sync.RWMutex
}

func (cmd *Blog) WithPlugins(fn plugins.FeederFn) error {
	if cmd == nil {
		return fmt.Errorf("blog is nil")
	}

	if fn == nil {
		return fmt.Errorf("fn is nil")
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	cmd.Feeder = fn

	return nil
}

func (cmd *Blog) Flags(stderr io.Writer) (*flag.FlagSet, error) {
	usage := `
Usage: hype blog <command> [options]

Commands:
    init <name>     Create a new blog project
    build           Build the static site to public/
    serve           Start a local preview server (default: localhost:3000)
    new <slug>      Create a new article scaffold
    theme           Manage themes (add, list, remove)

Examples:
    hype blog init mysite
    hype blog init mysite --theme developer
    hype blog build
    hype blog serve
    hype blog new hello-world
    hype blog theme list
    hype blog theme add suspended
`

	if err := cmd.validate(); err != nil {
		return nil, err
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	if cmd.flags != nil {
		return cmd.flags, nil
	}

	cmd.flags = flag.NewFlagSet("blog", flag.ContinueOnError)
	cmd.flags.SetOutput(stderr)
	cmd.flags.DurationVar(&cmd.Timeout, "timeout", DefaultTimeout, "timeout for execution, defaults to 30 seconds (30s)")
	cmd.flags.BoolVar(&cmd.Verbose, "v", false, "enable verbose output")

	cmd.flags.Usage = func() {
		fmt.Fprintf(stderr, "Usage of %s:\n", os.Args[0])
		cmd.flags.PrintDefaults()
		fmt.Fprintln(stderr, usage)
	}

	return cmd.flags, nil
}

func (cmd *Blog) Main(ctx context.Context, pwd string, args []string) error {
	cmd.mu.Lock()
	to := cmd.Timeout
	if to == 0 {
		to = DefaultTimeout
		cmd.Timeout = to
	}
	cmd.mu.Unlock()

	if err := (&cmd.Cmd).Init(); err != nil {
		return err
	}

	flags, err := cmd.Flags(cmd.Stderr())
	if err != nil {
		return err
	}

	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	remaining := flags.Args()
	if len(remaining) == 0 {
		flags.Usage()
		return nil
	}

	subCmd := remaining[0]
	subArgs := remaining[1:]

	switch subCmd {
	case "init":
		return cmd.runInit(ctx, pwd, subArgs)
	case "build":
		return cmd.runBuild(ctx, pwd, subArgs)
	case "serve":
		return cmd.runServe(ctx, pwd, subArgs)
	case "new":
		return cmd.runNew(ctx, pwd, subArgs)
	case "theme":
		return cmd.runTheme(ctx, pwd, subArgs)
	default:
		return fmt.Errorf("unknown subcommand: %s", subCmd)
	}
}

func (cmd *Blog) runInit(ctx context.Context, pwd string, args []string) error {
	var theme string

	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.StringVar(&theme, "theme", "suspended", "theme to use (built-in name or git URL)")
	fs.StringVar(&theme, "t", "", "theme to use (shorthand)")
	fs.Usage = func() {
		fmt.Fprintln(cmd.Stdout(), `Usage: hype blog init <name> [options]

Create a new blog project with the given name.

Options:
    --theme, -t    Theme to use (default: "suspended")
                   Can be a built-in theme name or a git URL

Built-in themes:
    suspended      Minimal, typography-focused theme (default)
    developer      Code-focused with enhanced syntax highlighting
    cards          Card-based layout for visual blogs

Arguments:
    name    Name of the blog directory to create

Examples:
    hype blog init mysite
    hype blog init mysite --theme developer
    hype blog init mysite --theme https://github.com/user/my-hype-theme`)
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	if fs.NArg() == 0 {
		fs.Usage()
		return fmt.Errorf("missing required argument: name")
	}

	name := fs.Arg(0)
	dir := filepath.Join(pwd, name)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	themeName := theme
	if strings.HasPrefix(theme, "https://") || strings.HasPrefix(theme, "git@") {
		themeName = extractThemeName(theme)
		if err := cloneTheme(theme, filepath.Join(dir, "themes", themeName)); err != nil {
			return fmt.Errorf("failed to clone theme: %w", err)
		}
	} else if blog.IsBuiltinTheme(theme) {
		if err := blog.CopyBuiltinTheme(theme, filepath.Join(dir, "themes", theme)); err != nil {
			return fmt.Errorf("failed to copy theme: %w", err)
		}
	} else {
		return fmt.Errorf("unknown theme: %s (use a built-in name or git URL)", theme)
	}

	configContent := fmt.Sprintf(`title: "%s"
description: "A blog powered by hype"
baseURL: "https://example.com"
author:
  name: "Your Name"
  email: ""
  twitter: ""
theme: "%s"
highlight:
  style: "monokai"
  lineNumbers: false
seo:
  ogImage: "/images/og-default.png"
  twitterCard: "summary_large_image"
contentDir: "content"
outputDir: "public"
`, name, themeName)

	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create config.yaml: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "content"), 0755); err != nil {
		return fmt.Errorf("failed to create content directory: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "static", "images"), 0755); err != nil {
		return fmt.Errorf("failed to create static/images directory: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "layouts"), 0755); err != nil {
		return fmt.Errorf("failed to create layouts directory: %w", err)
	}

	faviconSVG := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32">
  <rect width="32" height="32" rx="6" fill="#1e293b"/>
  <path d="M8 8v16M24 8v16M8 16h16" stroke="#60a5fa" stroke-width="3" stroke-linecap="round"/>
</svg>`
	if err := os.WriteFile(filepath.Join(dir, "static", "favicon.svg"), []byte(faviconSVG), 0644); err != nil {
		return fmt.Errorf("failed to create favicon.svg: %w", err)
	}

	gitignoreContent := `public/
.DS_Store
`
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	fmt.Fprintf(cmd.Stdout(), "Created new blog at %s\n", dir)
	fmt.Fprintf(cmd.Stdout(), "Theme: %s\n", themeName)
	fmt.Fprintf(cmd.Stdout(), "\nNext steps:\n")
	fmt.Fprintf(cmd.Stdout(), "  cd %s\n", name)
	fmt.Fprintf(cmd.Stdout(), "  hype blog new hello-world\n")
	fmt.Fprintf(cmd.Stdout(), "  hype blog build\n")

	return nil
}

func extractThemeName(url string) string {
	url = strings.TrimSuffix(url, ".git")
	parts := strings.Split(url, "/")
	name := parts[len(parts)-1]
	name = strings.TrimPrefix(name, "hype-theme-")
	return name
}

func cloneTheme(url, destDir string) error {
	if err := os.MkdirAll(filepath.Dir(destDir), 0755); err != nil {
		return err
	}
	cmd := exec.Command("git", "clone", "--depth=1", url, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (cmd *Blog) runTheme(ctx context.Context, pwd string, args []string) error {
	if len(args) == 0 {
		return cmd.runThemeList(ctx, pwd, args)
	}

	subCmd := args[0]
	subArgs := args[1:]

	switch subCmd {
	case "add":
		return cmd.runThemeAdd(ctx, pwd, subArgs)
	case "list":
		return cmd.runThemeList(ctx, pwd, subArgs)
	case "remove":
		return cmd.runThemeRemove(ctx, pwd, subArgs)
	default:
		return fmt.Errorf("unknown theme command: %s (use: add, list, remove)", subCmd)
	}
}

func (cmd *Blog) runThemeAdd(ctx context.Context, pwd string, args []string) error {
	fs := flag.NewFlagSet("theme add", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(cmd.Stdout(), `Usage: hype blog theme add <name-or-url>

Add a theme to the current project.

Arguments:
    name-or-url    Built-in theme name or git URL

Examples:
    hype blog theme add developer
    hype blog theme add https://github.com/user/hype-theme-custom`)
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	if fs.NArg() == 0 {
		fs.Usage()
		return fmt.Errorf("missing required argument: name-or-url")
	}

	theme := fs.Arg(0)
	themesDir := filepath.Join(pwd, "themes")

	if strings.HasPrefix(theme, "https://") || strings.HasPrefix(theme, "git@") {
		themeName := extractThemeName(theme)
		destDir := filepath.Join(themesDir, themeName)
		if _, err := os.Stat(destDir); err == nil {
			return fmt.Errorf("theme %s already exists at %s", themeName, destDir)
		}
		if err := cloneTheme(theme, destDir); err != nil {
			return fmt.Errorf("failed to clone theme: %w", err)
		}
		fmt.Fprintf(cmd.Stdout(), "Added theme %s from %s\n", themeName, theme)
		fmt.Fprintf(cmd.Stdout(), "Update config.yaml to use this theme: theme: \"%s\"\n", themeName)
	} else if blog.IsBuiltinTheme(theme) {
		destDir := filepath.Join(themesDir, theme)
		if _, err := os.Stat(destDir); err == nil {
			return fmt.Errorf("theme %s already exists at %s", theme, destDir)
		}
		if err := blog.CopyBuiltinTheme(theme, destDir); err != nil {
			return fmt.Errorf("failed to copy theme: %w", err)
		}
		fmt.Fprintf(cmd.Stdout(), "Added built-in theme %s\n", theme)
		fmt.Fprintf(cmd.Stdout(), "Update config.yaml to use this theme: theme: \"%s\"\n", theme)
	} else {
		return fmt.Errorf("unknown theme: %s (use a built-in name or git URL)", theme)
	}

	return nil
}

func (cmd *Blog) runThemeList(ctx context.Context, pwd string, args []string) error {
	fmt.Fprintln(cmd.Stdout(), "Built-in themes:")
	for _, t := range blog.ListBuiltinThemes() {
		desc := ""
		switch t {
		case "suspended":
			desc = "Minimal, typography-focused theme"
		case "developer":
			desc = "Code-focused with enhanced syntax highlighting"
		case "cards":
			desc = "Card-based layout for visual blogs"
		}
		fmt.Fprintf(cmd.Stdout(), "  %-12s %s\n", t, desc)
	}

	themesDir := filepath.Join(pwd, "themes")
	if entries, err := os.ReadDir(themesDir); err == nil && len(entries) > 0 {
		fmt.Fprintln(cmd.Stdout(), "\nInstalled themes:")
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			themeDir := filepath.Join(themesDir, entry.Name())
			info, _ := blog.LoadThemeInfo(themeDir)
			desc := info.Description
			if desc == "" {
				desc = "(no description)"
			}
			fmt.Fprintf(cmd.Stdout(), "  %-12s %s\n", entry.Name(), desc)
		}
	}

	return nil
}

func (cmd *Blog) runThemeRemove(ctx context.Context, pwd string, args []string) error {
	fs := flag.NewFlagSet("theme remove", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(cmd.Stdout(), `Usage: hype blog theme remove <name>

Remove an installed theme from the project.

Arguments:
    name    Name of the theme to remove

Example:
    hype blog theme remove custom-theme`)
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	if fs.NArg() == 0 {
		fs.Usage()
		return fmt.Errorf("missing required argument: name")
	}

	theme := fs.Arg(0)
	themeDir := filepath.Join(pwd, "themes", theme)

	if _, err := os.Stat(themeDir); os.IsNotExist(err) {
		return fmt.Errorf("theme %s not found at %s", theme, themeDir)
	}

	if err := os.RemoveAll(themeDir); err != nil {
		return fmt.Errorf("failed to remove theme: %w", err)
	}

	fmt.Fprintf(cmd.Stdout(), "Removed theme %s\n", theme)
	return nil
}

func (cmd *Blog) runBuild(ctx context.Context, pwd string, args []string) error {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(cmd.Stdout(), `Usage: hype blog build

Build the static site from content/ to public/.

The build process:
    1. Reads config.yaml for site settings
    2. Discovers articles in content/ directory
    3. Processes markdown with hype (code execution, includes, etc.)
    4. Generates HTML with syntax highlighting
    5. Creates RSS feed, sitemap, and robots.txt

Template lookup order:
    1. layouts/ (your project overrides)
    2. themes/<theme>/layouts/ (theme templates)
    3. Built-in defaults (fallback)

Example:
    hype blog build`)
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	b, err := blog.New(pwd)
	if err != nil {
		return err
	}

	if err := b.Build(ctx); err != nil {
		return err
	}

	fmt.Fprintf(cmd.Stdout(), "Built %d articles to %s/\n", len(b.Articles), b.Config.OutputDir)
	return nil
}

func (cmd *Blog) runServe(ctx context.Context, pwd string, args []string) error {
	var addr string
	var watch bool

	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.StringVar(&addr, "addr", ":3000", "address to serve on (default :3000)")
	fs.StringVar(&addr, "a", ":3000", "address to serve on (shorthand)")
	fs.BoolVar(&watch, "watch", false, "watch for file changes and rebuild")
	fs.BoolVar(&watch, "w", false, "watch for file changes and rebuild (shorthand)")
	fs.Usage = func() {
		fmt.Fprintln(cmd.Stdout(), `Usage: hype blog serve [options]

Start a local HTTP server to preview the built site.

If public/ doesn't exist, the site will be built first.

Options:
    -addr, -a    Address to serve on (default ":3000")
    -watch, -w   Watch for file changes and rebuild automatically

Example:
    hype blog serve
    hype blog serve -watch
    hype blog serve -addr :8080 -w`)
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	b, err := blog.New(pwd)
	if err != nil {
		return err
	}

	publicDir := filepath.Join(pwd, b.Config.OutputDir)
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		fmt.Fprintf(cmd.Stdout(), "Building site first...\n")
		if err := b.Build(ctx); err != nil {
			return err
		}
	}

	finalAddr, triedPorts := findAvailablePort(addr)
	if len(triedPorts) > 0 {
		fmt.Fprintf(cmd.Stdout(), "Ports in use: %s\n", strings.Join(triedPorts, ", "))
	}

	fmt.Fprintf(cmd.Stdout(), "Serving %s at http://localhost%s\n", publicDir, finalAddr)
	if watch {
		fmt.Fprintf(cmd.Stdout(), "Watching for changes...\n")
	}
	fmt.Fprintf(cmd.Stdout(), "Press Ctrl+C to stop\n")

	server := &http.Server{
		Addr:    finalAddr,
		Handler: http.FileServer(http.Dir(publicDir)),
	}

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	var watcher *fsnotify.Watcher
	if watch {
		watcher, err = cmd.startWatcher(ctx, pwd, b)
		if err != nil {
			return fmt.Errorf("failed to start file watcher: %w", err)
		}
		defer watcher.Close()
	}

	go func() {
		<-quit
		fmt.Fprintf(cmd.Stdout(), "\nShutting down server...\n")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Fprintf(cmd.Stderr(), "Server shutdown error: %v\n", err)
		}
		close(done)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	<-done
	return nil
}

func (cmd *Blog) startWatcher(ctx context.Context, pwd string, b *blog.Blog) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	contentDir := filepath.Join(pwd, b.Config.ContentDir)
	staticDir := filepath.Join(pwd, "static")
	layoutsDir := filepath.Join(pwd, "layouts")
	themesDir := filepath.Join(pwd, "themes")

	if err := addWatchRecursive(watcher, contentDir); err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to watch content directory: %w", err)
	}

	if _, err := os.Stat(staticDir); err == nil {
		if err := addWatchRecursive(watcher, staticDir); err != nil {
			watcher.Close()
			return nil, fmt.Errorf("failed to watch static directory: %w", err)
		}
	}

	if _, err := os.Stat(layoutsDir); err == nil {
		if err := addWatchRecursive(watcher, layoutsDir); err != nil {
			watcher.Close()
			return nil, fmt.Errorf("failed to watch layouts directory: %w", err)
		}
	}

	if _, err := os.Stat(themesDir); err == nil {
		if err := addWatchRecursive(watcher, themesDir); err != nil {
			watcher.Close()
			return nil, fmt.Errorf("failed to watch themes directory: %w", err)
		}
	}

	configPath := filepath.Join(pwd, "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		if err := watcher.Add(configPath); err != nil {
			watcher.Close()
			return nil, fmt.Errorf("failed to watch config.yaml: %w", err)
		}
	}

	configPathYml := filepath.Join(pwd, "config.yml")
	if _, err := os.Stat(configPathYml); err == nil {
		if err := watcher.Add(configPathYml); err != nil {
			watcher.Close()
			return nil, fmt.Errorf("failed to watch config.yml: %w", err)
		}
	}

	go cmd.watchLoop(ctx, watcher, pwd)

	return watcher, nil
}

func addWatchRecursive(watcher *fsnotify.Watcher, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && path != dir {
				return filepath.SkipDir
			}
			return watcher.Add(path)
		}
		return nil
	})
}

func (cmd *Blog) watchLoop(ctx context.Context, watcher *fsnotify.Watcher, pwd string) {
	var debounceTimer *time.Timer
	var mu sync.Mutex
	var changedFiles []string
	debounceDelay := 500 * time.Millisecond

	rebuild := func() {
		mu.Lock()
		files := changedFiles
		changedFiles = nil
		mu.Unlock()

		fmt.Fprintf(cmd.Stdout(), "\n")
		for _, f := range files {
			relPath, err := filepath.Rel(pwd, f)
			if err != nil {
				relPath = f
			}
			fmt.Fprintf(cmd.Stdout(), "Changed: %s\n", relPath)
		}
		fmt.Fprintf(cmd.Stdout(), "Rebuilding site...\n")

		b, err := blog.New(pwd)
		if err != nil {
			fmt.Fprintf(cmd.Stderr(), "Rebuild error: %v\n", err)
			return
		}
		if err := b.Build(ctx); err != nil {
			fmt.Fprintf(cmd.Stderr(), "Rebuild error: %v\n", err)
			return
		}
		fmt.Fprintf(cmd.Stdout(), "Rebuilt %d articles\n", len(b.Articles))
	}

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				if strings.Contains(event.Name, "public") {
					continue
				}
				if strings.HasPrefix(filepath.Base(event.Name), ".") {
					continue
				}

				mu.Lock()
				found := false
				for _, f := range changedFiles {
					if f == event.Name {
						found = true
						break
					}
				}
				if !found {
					changedFiles = append(changedFiles, event.Name)
				}
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.AfterFunc(debounceDelay, rebuild)
				mu.Unlock()

				if event.Has(fsnotify.Create) {
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						watcher.Add(event.Name)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(cmd.Stderr(), "Watcher error: %v\n", err)
		}
	}
}

func findAvailablePort(addr string) (string, []string) {
	var triedPorts []string

	port := 3000
	if strings.HasPrefix(addr, ":") {
		if p, err := strconv.Atoi(addr[1:]); err == nil {
			port = p
		}
	}

	maxAttempts := 100
	for i := 0; i < maxAttempts; i++ {
		testAddr := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", testAddr)
		if err == nil {
			ln.Close()
			return testAddr, triedPorts
		}
		triedPorts = append(triedPorts, strconv.Itoa(port))
		port++
	}

	return fmt.Sprintf(":%d", port), triedPorts
}

func (cmd *Blog) runNew(ctx context.Context, pwd string, args []string) error {
	fs := flag.NewFlagSet("new", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(cmd.Stdout(), `Usage: hype blog new <slug>

Create a new article scaffold with the given slug.

Arguments:
    slug    URL-friendly identifier for the article (e.g., my-first-post)

Creates:
    content/<slug>/module.md    Article content file
    content/<slug>/src/         Directory for source code files

Example:
    hype blog new my-first-post
    hype blog new go-concurrency-patterns`)
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	if fs.NArg() == 0 {
		fs.Usage()
		return fmt.Errorf("missing required argument: slug")
	}

	slug := fs.Arg(0)
	articleDir := filepath.Join(pwd, "content", slug)

	if err := os.MkdirAll(articleDir, 0755); err != nil {
		return fmt.Errorf("failed to create article directory: %w", err)
	}

	today := time.Now().Format("01/02/2006")
	moduleContent := fmt.Sprintf(`# Article Title

<details>
slug: %s
published: %s
author: Your Name
seo_description: Brief description of the article for SEO (150-160 chars)
tags: tag1, tag2
</details>

Write your article content here.

## Section 1

Your content...

## Section 2

More content...
`, slug, today)

	modulePath := filepath.Join(articleDir, "module.md")
	if err := os.WriteFile(modulePath, []byte(moduleContent), 0644); err != nil {
		return fmt.Errorf("failed to create module.md: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(articleDir, "src"), 0755); err != nil {
		return fmt.Errorf("failed to create src directory: %w", err)
	}

	fmt.Fprintf(cmd.Stdout(), "Created new article at %s\n", articleDir)
	fmt.Fprintf(cmd.Stdout(), "\nEdit %s to add your content.\n", modulePath)

	return nil
}

func (cmd *Blog) validate() error {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	if cmd.Timeout == 0 {
		cmd.Timeout = DefaultTimeout
	}

	return nil
}
