package preview

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/fsnotify/fsnotify"
	"github.com/gopherguides/hype"
)

var defaultExtensions = []string{
	"md", "html", "go", "css", "png", "jpg", "jpeg", "gif", "svg", "webp",
}

var defaultExcludes = []string{
	"**/.git/**",
	"**/.hg/**",
	"**/.svn/**",
	"**/.idea/**",
	"**/.vscode/**",
	"**/node_modules/**",
	"**/vendor/**",
	"**/tmp/**",
	"**/__pycache__/**",
}

type Config struct {
	File           string
	Port           int
	WatchDirs      []string
	Extensions     []string
	IncludeGlobs   []string
	ExcludeGlobs   []string
	PollInterval   time.Duration
	DebounceDelay  time.Duration
	Verbose        bool
	OpenBrowser    bool
	Theme          string
}

func DefaultConfig() Config {
	return Config{
		File:          "hype.md",
		Port:          3000,
		WatchDirs:     []string{"."},
		Extensions:    defaultExtensions,
		ExcludeGlobs:  defaultExcludes,
		DebounceDelay: 300 * time.Millisecond,
		Theme:         "github",
	}
}

type Server struct {
	config     Config
	parser     *hype.Parser
	httpServer *http.Server
	watcher    *fsnotify.Watcher
	liveReload *LiveReload

	currentHTML string
	mu          sync.RWMutex

	stdout func(format string, args ...any)
	stderr func(format string, args ...any)
}

func New(cfg Config, parser *hype.Parser) *Server {
	return &Server{
		config: cfg,
		parser: parser,
		stdout: func(format string, args ...any) {
			_, _ = fmt.Fprintf(os.Stdout, format, args...)
		},
		stderr: func(format string, args ...any) {
			fmt.Fprintf(os.Stderr, format, args...)
		},
	}
}

func (s *Server) SetOutput(stdout, stderr func(format string, args ...any)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stdout = stdout
	s.stderr = stderr
}

func (s *Server) Run(ctx context.Context, pwd string) error {
	s.liveReload = NewLiveReload()

	if err := s.build(ctx, pwd); err != nil {
		return fmt.Errorf("initial build failed: %w", err)
	}

	watcher, err := s.startWatcher(ctx, pwd)
	if err != nil {
		return fmt.Errorf("failed to start watcher: %w", err)
	}
	s.watcher = watcher
	defer func() { _ = watcher.Close() }()

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handlePreview)
	mux.HandleFunc("/_livereload", s.liveReload.HandleWebSocket)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: mux,
	}

	serverErr := make(chan error, 1)
	go func() {
		s.stdout("Starting preview server at http://localhost:%d\n", s.config.Port)
		s.stdout("Watching for changes... (Press Ctrl+C to stop)\n")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(shutdownCtx)
	case err := <-serverErr:
		return err
	}
}

func (s *Server) build(ctx context.Context, pwd string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fileDir := filepath.Dir(s.config.File)
	fileName := filepath.Base(s.config.File)

	var parserFS = os.DirFS(pwd)
	if fileDir != "." && fileDir != "" {
		parserFS = os.DirFS(filepath.Join(pwd, fileDir))
	}

	p := s.parser
	if p == nil {
		p = hype.NewParser(parserFS)
	} else {
		p.FS = parserFS
	}
	p.Root = filepath.Join(pwd, fileDir)

	doc, err := p.ParseFile(fileName)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	if err := doc.Execute(ctx); err != nil {
		return fmt.Errorf("execute error: %w", err)
	}

	s.currentHTML = wrapHTML(doc.String(), s.config.Theme)
	return nil
}

func (s *Server) handlePreview(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	html := s.currentHTML
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, html)
}

func (s *Server) startWatcher(ctx context.Context, pwd string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watchDirs := s.config.WatchDirs
	if len(watchDirs) == 0 {
		watchDirs = []string{"."}
	}

	for _, dir := range watchDirs {
		absDir := dir
		if !filepath.IsAbs(dir) {
			absDir = filepath.Join(pwd, dir)
		}
		if err := s.addWatchRecursive(watcher, absDir); err != nil {
			_ = watcher.Close()
			return nil, fmt.Errorf("failed to watch %s: %w", dir, err)
		}
	}

	go s.watchLoop(ctx, watcher, pwd)

	return watcher, nil
}

func (s *Server) addWatchRecursive(watcher *fsnotify.Watcher, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(dir, path)
		if relPath == "" {
			relPath = "."
		}

		for _, pattern := range s.config.ExcludeGlobs {
			if matched, _ := doublestar.Match(pattern, relPath); matched {
				return filepath.SkipDir
			}
			if matched, _ := doublestar.Match(pattern, info.Name()); matched {
				return filepath.SkipDir
			}
		}

		return watcher.Add(path)
	})
}

func (s *Server) watchLoop(ctx context.Context, watcher *fsnotify.Watcher, pwd string) {
	var debounceTimer *time.Timer
	var mu sync.Mutex
	var changedFiles []string

	rebuild := func() {
		mu.Lock()
		files := changedFiles
		changedFiles = nil
		mu.Unlock()

		if s.config.Verbose {
			s.stdout("\n")
			for _, f := range files {
				relPath, err := filepath.Rel(pwd, f)
				if err != nil {
					relPath = f
				}
				s.stdout("Changed: %s\n", relPath)
			}
		}
		s.stdout("Rebuilding...\n")

		if err := s.build(ctx, pwd); err != nil {
			s.stderr("Build error: %v\n", err)
			return
		}

		s.stdout("Done. Reloading browser...\n")
		s.liveReload.Reload()
	}

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Create) {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					_ = s.addWatchRecursive(watcher, event.Name)
				}
			}

			if !s.shouldWatch(event.Name, pwd) {
				continue
			}

			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
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
				debounceTimer = time.AfterFunc(s.config.DebounceDelay, rebuild)
				mu.Unlock()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			s.stderr("Watcher error: %v\n", err)
		}
	}
}

func (s *Server) shouldWatch(path string, pwd string) bool {
	relPath, err := filepath.Rel(pwd, path)
	if err != nil {
		relPath = path
	}

	for _, pattern := range s.config.ExcludeGlobs {
		if matched, _ := doublestar.Match(pattern, relPath); matched {
			return false
		}
	}

	ext := filepath.Ext(path)
	if ext != "" {
		ext = ext[1:]
	}

	if len(s.config.IncludeGlobs) > 0 {
		for _, pattern := range s.config.IncludeGlobs {
			if matched, _ := doublestar.Match(pattern, relPath); matched {
				return true
			}
		}
		return false
	}

	if len(s.config.Extensions) > 0 {
		for _, e := range s.config.Extensions {
			if ext == e {
				return true
			}
		}
		return false
	}

	return true
}

func wrapHTML(content, theme string) string {
	themeCSS := getThemeCSS(theme)
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Hype Preview</title>
    <style>
        %s
        body {
            max-width: 900px;
            margin: 0 auto;
            padding: 2rem;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
            line-height: 1.6;
        }
        pre {
            background: #f6f8fa;
            padding: 1rem;
            border-radius: 6px;
            overflow-x: auto;
        }
        code {
            font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
            font-size: 0.9em;
        }
        img {
            max-width: 100%%;
            height: auto;
        }
    </style>
</head>
<body>
%s
<script>
(function() {
    var ws = new WebSocket('ws://' + location.host + '/_livereload');
    ws.onmessage = function(e) {
        if (e.data === 'reload') {
            location.reload();
        }
    };
    ws.onclose = function() {
        console.log('LiveReload disconnected. Attempting to reconnect...');
        setTimeout(function() {
            location.reload();
        }, 1000);
    };
    ws.onerror = function(err) {
        console.log('LiveReload error:', err);
    };
})();
</script>
</body>
</html>`, themeCSS, content)
}

func getThemeCSS(theme string) string {
	switch theme {
	case "github-dark":
		return `
        body { background: #0d1117; color: #c9d1d9; }
        pre { background: #161b22; }
        a { color: #58a6ff; }
        h1, h2, h3, h4, h5, h6 { color: #c9d1d9; border-bottom-color: #21262d; }
        `
	case "github":
		fallthrough
	default:
		return `
        body { background: #fff; color: #24292f; }
        pre { background: #f6f8fa; }
        a { color: #0969da; }
        h1, h2, h3, h4, h5, h6 { border-bottom: 1px solid #d0d7de; padding-bottom: 0.3em; }
        `
	}
}
