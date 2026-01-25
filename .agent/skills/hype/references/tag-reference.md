# Hype Tag Reference

Complete attribute reference for all hype tags.

## `<code>` Tag

Display source code with syntax highlighting.

| Attribute | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `src` | string | Yes | - | Path to source file, relative to document |
| `snippet` | string | No | - | Name of snippet to extract from file |
| `range` | string | No | - | Line range in format `start:end` (0-indexed) |
| `language` | string | No | auto | Override detected language for syntax highlighting |
| `esc` | flag | No | - | HTML-escape the code content |

### Fragment Syntax

The `src` attribute supports fragment syntax: `src="file.go#snippetname"` as shorthand for `src="file.go" snippet="snippetname"`.

### Range Format

The `range` attribute accepts:
- `10:20` - Lines 10-20 (0-indexed, exclusive end)
- `:10` - First 10 lines
- `10:` - From line 10 to end

## `<go>` Tag

Execute Go commands and display output. Inherits all `<cmd>` attributes.

### Go Command Attributes

Each Go subcommand is an attribute. The attribute value becomes the command argument.

| Attribute | Example | Resulting Command |
|-----------|---------|-------------------|
| `run` | `run="main.go"` | `go run main.go` |
| `build` | `build="."` | `go build .` |
| `test` | `test="-v ./..."` | `go test -v ./...` |
| `doc` | `doc="fmt.Println"` | `go doc fmt.Println` |
| `fmt` | `fmt="."` | `go fmt .` |
| `vet` | `vet="./..."` | `go vet ./...` |
| `mod` | `mod="tidy"` | `go mod tidy` |
| `get` | `get="github.com/pkg"` | `go get github.com/pkg` |
| `install` | `install="."` | `go install .` |
| `list` | `list="-m all"` | `go list -m all` |
| `version` | `version=""` | `go version` |
| `env` | `env="GOPATH"` | `go env GOPATH` |
| `generate` | `generate="./..."` | `go generate ./...` |
| `clean` | `clean="-cache"` | `go clean -cache` |
| `bug` | `bug=""` | `go bug` |
| `fix` | `fix="."` | `go fix .` |
| `help` | `help="build"` | `go help build` |
| `tool` | `tool="pprof"` | `go tool pprof` |

### Cross-Compilation Attributes

| Attribute | Type | Default | Description |
|-----------|------|---------|-------------|
| `goos` | string | current | Target operating system (linux, darwin, windows, etc.) |
| `goarch` | string | current | Target architecture (amd64, arm64, 386, etc.) |

### Additional Attributes

| Attribute | Type | Default | Description |
|-----------|------|---------|-------------|
| `code` | string | - | Show source file before executing command |
| `sym` | string | - | Show Go doc with `-cmd -u -src -short` flags |
| `src` | string | - | Working directory for command execution |
| `exit` | int | 0 | Expected exit code (-1 for any non-zero) |
| `timeout` | duration | 30s | Maximum execution time |
| `environ` | string | - | Comma-separated environment variables (`KEY=val,KEY2=val2`) |

## `<cmd>` Tag

Execute arbitrary shell commands.

| Attribute | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `exec` | string | Yes | - | Command to execute |
| `src` | string | No | - | Working directory for command |
| `exit` | int | No | 0 | Expected exit code |
| `timeout` | duration | No | 30s | Maximum execution time |
| `environ` | string | No | - | Comma-separated environment variables |

### Exit Code Handling

| Value | Behavior |
|-------|----------|
| `0` | Command must succeed (exit 0) |
| `1`, `2`, etc. | Command must exit with this specific code |
| `-1` | Command must fail (any non-zero exit) |

### Timeout Format

Timeouts use Go duration format: `30s`, `5m`, `1h30m`, `500ms`.

## `<include>` Tag

Include content from another markdown file.

| Attribute | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `src` | string | Yes | - | Path to markdown file to include |

### Path Resolution

- Paths are relative to the including document
- Included files have their own relative paths adjusted automatically
- Assets (images, links) in included files resolve correctly

## `<youtube>` Tag

Embed YouTube videos.

| Attribute | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `id` | string | Yes | - | YouTube video ID (11 alphanumeric characters) |
| `title` | string | No | "YouTube video player" | Title for the iframe (accessibility) |

### Video ID Format

The video ID is the 11-character string from YouTube URLs:
- `https://www.youtube.com/watch?v=dQw4w9WgXcQ` → `id="dQw4w9WgXcQ"`
- `https://youtu.be/dQw4w9WgXcQ` → `id="dQw4w9WgXcQ"`

Valid characters: `a-z`, `A-Z`, `0-9`, `-`, `_`
