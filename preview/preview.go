package preview

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/fsnotify/fsnotify"
	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/internal/portutil"
	"github.com/gopherguides/hype/themes"
)

var defaultExtensions = []string{
	"md", "html", "go", "css", "png", "jpg", "jpeg", "gif", "svg", "webp",
}

var defaultExcludes = []string{
	".git",
	".hg",
	".svn",
	".idea",
	".vscode",
	"node_modules",
	"vendor",
	"tmp",
	"__pycache__",
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
	File          string
	Port          int
	WatchDirs     []string
	Extensions    []string
	IncludeGlobs  []string
	ExcludeGlobs  []string
	DebounceDelay time.Duration
	Verbose       bool
	OpenBrowser   bool
	Theme         string
	CustomCSS     string
	Timeout       time.Duration
	OnReady       func(port int)
}

func DefaultConfig() Config {
	return Config{
		File:          "hype.md",
		Port:          3000,
		WatchDirs:     []string{"."},
		Extensions:    defaultExtensions,
		ExcludeGlobs:  defaultExcludes,
		DebounceDelay: 300 * time.Millisecond,
		Theme:         themes.DefaultTheme,
	}
}

type Server struct {
	config     Config
	parser     *hype.Parser
	httpServer *http.Server
	watcher    *fsnotify.Watcher
	liveReload *LiveReload

	currentHTML string
	pwd         string
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

func (s *Server) Port() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config.Port
}

func (s *Server) Run(ctx context.Context, pwd string) error {
	s.liveReload = NewLiveReload()
	s.pwd = pwd

	if err := s.build(ctx, pwd); err != nil {
		return fmt.Errorf("initial build failed: %w", err)
	}

	watcher, err := s.startWatcher(ctx, pwd)
	if err != nil {
		return fmt.Errorf("failed to start watcher: %w", err)
	}
	s.watcher = watcher
	defer func() { _ = watcher.Close() }()

	port, triedPorts := portutil.FindAvailablePort(s.config.Port)
	if len(triedPorts) > 0 {
		s.stdout("Ports in use: ")
		for i, p := range triedPorts {
			if i > 0 {
				s.stdout(", ")
			}
			s.stdout("%d", p)
		}
		s.stdout("\n")
	}
	s.config.Port = port

	mux := http.NewServeMux()
	mux.HandleFunc("/_livereload", s.liveReload.HandleWebSocket)
	mux.HandleFunc("/", s.handleRequest)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	serverErr := make(chan error, 1)
	go func() {
		s.stdout("Starting preview server at http://localhost:%d\n", port)
		s.stdout("Watching for changes... (Press Ctrl+C to stop)\n")
		if s.config.OnReady != nil {
			s.config.OnReady(port)
		}
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

	execCtx := ctx
	if s.config.Timeout > 0 {
		var cancel context.CancelFunc
		execCtx, cancel = context.WithTimeout(ctx, s.config.Timeout)
		defer cancel()
	}

	if err := doc.Execute(execCtx); err != nil {
		return fmt.Errorf("execute error: %w", err)
	}

	html, err := s.wrapHTML(doc.String())
	if err != nil {
		return fmt.Errorf("wrap HTML error: %w", err)
	}
	s.currentHTML = html
	return nil
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "" {
		s.handlePreview(w, r)
		return
	}

	urlPath := strings.TrimPrefix(r.URL.Path, "/")
	cleanPath := filepath.Clean(urlPath)

	if strings.HasPrefix(cleanPath, "..") || strings.Contains(cleanPath, ".."+string(filepath.Separator)) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	filePath := filepath.Join(s.pwd, cleanPath)

	relPath, err := filepath.Rel(s.pwd, filePath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		s.handlePreview(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
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
		relPath = filepath.ToSlash(relPath)

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
	relPath = filepath.ToSlash(relPath)

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
		extLower := strings.ToLower(ext)
		for _, e := range s.config.Extensions {
			if extLower == strings.ToLower(e) {
				return true
			}
		}
		return false
	}

	return true
}

func (s *Server) wrapHTML(content string) (string, error) {
	css, err := s.getCSS()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Hype Preview</title>
    <style>
%s
    </style>
</head>
<body>
    <article class="markdown-body">
%s
    </article>
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
</html>`, css, content), nil
}

func (s *Server) getCSS() (string, error) {
	if s.config.CustomCSS != "" {
		return themes.LoadCustomCSS(s.config.CustomCSS)
	}
	return themes.GetCSS(s.config.Theme)
}
