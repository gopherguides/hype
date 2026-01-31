# CLI Reference

Hype provides several commands for working with dynamic markdown documents.

## Commands Overview

| Command | Description |
|---------|-------------|
| `export` | Export documents to different formats (markdown, HTML) |
| `preview` | Start a live preview server with auto-reload |
| `marked` | Integration with Marked 2 app |
| `slides` | Web-based presentation server |
| `blog` | Static blog generator |

---

## export

Export hype documents to markdown or HTML.

```bash
hype export [options]
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `-f` | `hype.md` | Input file to process |
| `-format` | `markdown` | Output format: `markdown` or `html` |
| `-o` | stdout | Output file path |
| `-theme` | `github` | Theme for HTML export |
| `-css` | | Path to custom CSS file |
| `-no-css` | `false` | Output raw HTML without styling |
| `-themes` | | List available themes and exit |
| `-timeout` | `30s` | Execution timeout |
| `-v` | `false` | Verbose output |

### Examples

```bash
# Export to markdown (default)
hype export -f hype.md > README.md

# Export to HTML
hype export -f docs.md -format html > docs.html

# Export with a theme
hype export -f docs.md -format html -theme solarized-dark

# Export with custom CSS
hype export -f docs.md -format html -css ./styles.css

# Export raw HTML (no styling)
hype export -f docs.md -format html -no-css

# List available themes
hype export -themes

# Output directly to file
hype export -f hype.md -format markdown -o README.md
```

---

## preview

Start a live preview server with file watching and auto-reload.

```bash
hype preview [options]
```

### Options

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `-f` | | `hype.md` | Source file to preview |
| `-port` | | `3000` | Server port |
| `-w` | `-watch` | | Additional directories to watch (repeatable) |
| `-e` | `-ext` | | File extensions to watch (comma-separated) |
| `-i` | `-include` | | Glob patterns to include (repeatable) |
| `-x` | `-exclude` | | Glob patterns to exclude (repeatable) |
| `-d` | `-debounce` | `300ms` | Debounce delay before rebuild |
| `-v` | `-verbose` | `false` | Verbose output |
| `-open` | | `false` | Auto-open browser on start |
| `-theme` | | `github` | Preview theme |
| `-css` | | | Custom CSS file path |
| `-themes` | | | List available themes |
| `-timeout` | | `0` | Execution timeout |

### Examples

```bash
# Basic preview
hype preview -f hype.md

# Open browser automatically
hype preview -f hype.md -open

# Watch additional directories
hype preview -f hype.md -w ./src -w ./images

# Filter by extension
hype preview -f hype.md -e md,go,html

# Use a dark theme
hype preview -f hype.md -theme solarized-dark
```

---

## marked

Integration with [Marked 2](https://marked2app.com/) for macOS.

```bash
hype marked [options]
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `-f` | | Input file (uses `MARKED_PATH` if not set) |
| `-p` | `false` | Parse only (no execution) |
| `-timeout` | `5s` | Execution timeout |
| `-context` | | Context folder path |
| `-section` | `0` | Target section number |
| `-v` | `false` | Verbose output |

### Environment Variables

- `MARKED_PATH` - Set by Marked 2 to the current file path
- `MARKED_ORIGIN` - Set by Marked 2 to the file's directory

---

## slides

Web-based presentation server.

```bash
hype slides [options] [file]
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `3000` | Server port |

### Examples

```bash
# Start slides server
hype slides presentation.md

# Use a different port
hype slides -port 8080 presentation.md
```

---

## blog

Static blog generator with theming support.

```bash
hype blog <command> [options]
```

### Subcommands

| Command | Description |
|---------|-------------|
| `init <name>` | Create a new blog project |
| `build` | Build static site to `public/` |
| `serve` | Start local preview server |
| `new <slug>` | Create a new article |
| `theme` | Manage themes (add, list, remove) |

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `-timeout` | `30s` | Execution timeout |
| `-v` | `false` | Verbose output |

### Examples

```bash
# Create a new blog
hype blog init mysite

# Create with a theme
hype blog init mysite --theme developer

# Build the site
hype blog build

# Start preview server
hype blog serve

# Create a new article
hype blog new hello-world

# List available themes
hype blog theme list

# Add a theme
hype blog theme add suspended
```

---

## Common Options

These options are available across most commands:

| Flag | Description |
|------|-------------|
| `-f` | Input file path |
| `-timeout` | Execution timeout for code blocks |
| `-v` | Enable verbose/debug output |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error |

---

## Getting Help

```bash
# Show available commands
hype

# Show help for a specific command
hype export --help
hype preview --help
hype blog --help
```
