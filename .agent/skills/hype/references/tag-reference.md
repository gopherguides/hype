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
| `replace-N` | string | No | - | Regex pattern to match (N is 1, 2, 3, etc.) |
| `replace-N-with` | string | No | "" | Replacement text for matched pattern |

### Output Replacement

Commands often produce output containing dynamic content like timestamps, version numbers, or UUIDs. When you regenerate your documentation, this dynamic content changes even though your actual code hasn't—creating noise in your version control and making it impossible to see what really changed.

Use `replace-N` and `replace-N-with` attribute pairs to replace dynamic content with stable placeholders. This ensures predictable, reproducible output every time you regenerate—essential for blogs, READMEs, and any documentation you generate repeatedly.

```html
<!-- Replace Go version -->
<cmd exec="go version"
     replace-1="go1\.\d+\.\d+"
     replace-1-with="goX.X.X">
</cmd>

<!-- Replace timestamps -->
<cmd exec="./build.sh"
     replace-1="\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}"
     replace-1-with="YYYY-MM-DD HH:MM:SS">
</cmd>

<!-- Multiple replacements (applied in numeric order) -->
<cmd exec="go run main.go"
     replace-1="go1\.\d+\.\d+"
     replace-1-with="goX.X.X"
     replace-2="\d{4}-\d{2}-\d{2}"
     replace-2-with="[DATE]"
     replace-3="[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"
     replace-3-with="[UUID]">
</cmd>

<!-- Remove content (replace with empty) -->
<cmd exec="cat version.txt"
     replace-1=" \(build .*?\)"
     replace-1-with="">
</cmd>
```

Replacements are applied before HTML escaping, in numeric order (1, 2, 3, ...). Invalid regex patterns cause an error at export time.

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

## Mermaid Diagrams

Render Mermaid diagrams as ASCII art using fenced code blocks.

### Syntax

Use standard markdown fenced code blocks with `mermaid` as the language:

    ```mermaid
    graph LR
        A[Start] --> B{Decision}
        B -->|Yes| C[Action]
        B -->|No| D[End]
    ```

### Supported Diagram Types

| Type | Directive | Description |
|------|-----------|-------------|
| Flowchart | `graph LR`, `graph TD` | Left-to-right or top-down flowcharts |
| Flowchart | `flowchart LR`, `flowchart TD` | Alternative flowchart syntax |
| Sequence | `sequenceDiagram` | Interaction between participants |

### Graph Directions

| Direction | Description |
|-----------|-------------|
| `LR` | Left to Right |
| `TD` / `TB` | Top Down / Top to Bottom |

### Output Format

- **HTML export**: Rendered as `<pre><code class="language-plain">...</code></pre>`
- **Markdown export**: Rendered as plain fenced code block with ASCII art

### Limitations

The ASCII rendering is provided by [mermaid-ascii](https://github.com/AlexanderGrooff/mermaid-ascii). Not all Mermaid features are supported:

**Not supported:**
- Subgraph nesting
- Non-rectangular node shapes (diamonds render as rectangles)
- Class diagrams
- State diagrams
- Gantt charts
- Pie charts
- Diagonal arrows
