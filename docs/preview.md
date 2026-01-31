# Live Preview

Hype includes a live preview server with automatic file watching and browser reload for a seamless documentation authoring experience.

## Basic Usage

```bash
# Start preview server on default port (3000)
hype preview -f hype.md

# Open browser automatically
hype preview -f hype.md -open

# Use a different port
hype preview -f hype.md -port 8080
```

The preview server watches for file changes and automatically rebuilds the document, pushing updates to connected browsers via WebSocket.

## Watch Configuration

### Watch Directories

By default, the preview server watches the current directory (`.`). Use `-w` to specify which directories to watch:

```bash
# Watch multiple directories (replaces default, so include . to watch source file)
hype preview -f hype.md -w . -w ./src -w ./images
```

**Note:** The `-w` flag replaces the default watch directory rather than adding to it. If your source file is in the current directory and you want to watch additional directories, include `.` in your watch list to ensure changes to the source file trigger rebuilds.

### File Extensions

Filter which file types trigger rebuilds:

```bash
# Only watch specific extensions
hype preview -f hype.md -e md,html,go,png,jpg
```

### Include/Exclude Patterns

Use glob patterns to fine-tune what files are watched:

```bash
# Include specific patterns
hype preview -f hype.md -i "**/*.md" -i "**/*.go"

# Exclude directories
hype preview -f hype.md -x "**/vendor/**" -x "**/tmp/**"

# Combine include and exclude
hype preview -f hype.md -i "**/*.md" -x "**/node_modules/**"
```

## Themes

The preview server supports the same themes as HTML export:

```bash
# List available themes
hype preview -themes

# Use a specific theme
hype preview -f hype.md -theme solarized-dark

# Use custom CSS
hype preview -f hype.md -css ./my-styles.css
```

## Advanced Options

### Debounce Delay

Control how long the server waits after a file change before rebuilding:

```bash
# Shorter delay for faster feedback (100ms)
hype preview -f hype.md -d 100ms

# Longer delay for busy file systems (500ms)
hype preview -f hype.md -debounce 500ms
```

### Execution Timeout

Set a timeout for document execution (useful for documents with long-running commands):

```bash
hype preview -f hype.md -timeout 60s
```

### Verbose Output

Enable verbose mode to see file change events:

```bash
hype preview -f hype.md -v
```

## Flags Reference

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `-f` | | `hype.md` | Source markdown file to preview |
| `-port` | | `3000` | Server port |
| `-w` | `-watch` | | Directories to watch (repeatable) |
| `-e` | `-ext` | | File extensions to watch (comma-separated) |
| `-i` | `-include` | | Glob patterns to include (repeatable) |
| `-x` | `-exclude` | | Glob patterns to exclude (repeatable) |
| `-d` | `-debounce` | `300ms` | Debounce delay before rebuild |
| `-v` | `-verbose` | `false` | Verbose output (log file changes) |
| `-open` | | `false` | Auto-open browser on start |
| `-theme` | | `github` | Preview theme name |
| `-css` | | | Path to custom CSS file (overrides -theme) |
| `-themes` | | | List available themes and exit |
| `-timeout` | | `0` | Execution timeout (0 = no timeout) |

## How It Works

1. The server starts an HTTP server on the specified port
2. A file watcher monitors the source file and watch directories
3. When changes are detected, the server rebuilds the document
4. Connected browsers receive a WebSocket message to reload
5. The browser automatically refreshes with the updated content

The preview uses the same rendering pipeline as `hype export -format=html`, ensuring what you see matches the final output.
