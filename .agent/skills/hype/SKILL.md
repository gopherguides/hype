---
name: hype
description: Write and maintain documentation using the Hype content generation tool with dynamic code execution, snippets, and includes
license: MIT
compatibility:
  - claude-code
  - openai-codex
  - gemini-cli
  - cursor
  - github-copilot
  - vscode
tags:
  - documentation
  - markdown
  - go
  - code-execution
---

# Hype Documentation Tool

Hype is a content generation tool that extends Markdown with dynamic features for creating rich, automated documentation that stays in sync with your code.

## Core Principles

- **Packages**: Keep content in small, reusable units with relative links
- **Reuse**: Write documentation once, use everywhere (blog, book, README)
- **Includes**: Compose larger documents from smaller partials
- **Validation**: Verify code samples compile and run correctly
- **Asset validation**: Ensure local assets like images exist

## Tag Reference

### `<code>` - Display Source Code

Display source code from files with syntax highlighting.

```html
<!-- Entire file -->
<code src="main.go"></code>

<!-- Named snippet -->
<code src="main.go" snippet="example"></code>

<!-- Alternative snippet syntax -->
<code src="main.go#example"></code>

<!-- Line range (0-indexed) -->
<code src="main.go" range="10:20"></code>

<!-- Override language -->
<code src="config.txt" language="yaml"></code>

<!-- Escape HTML entities -->
<code src="template.html" esc></code>
```

### `<go>` - Execute Go Commands

Run Go commands and display output. Supports all `go` subcommands.

```html
<!-- Run a Go file -->
<go run="main.go"></go>

<!-- Run with source directory -->
<go src="myapp" run="."></go>

<!-- Build and show output -->
<go build="."></go>

<!-- Run tests -->
<go test="-v ./..."></go>

<!-- Show documentation for a symbol -->
<go doc="fmt.Println"></go>

<!-- Cross-compilation -->
<go build="." goos="linux" goarch="amd64"></go>

<!-- With environment variables -->
<go run="main.go" environ="DEBUG=true,LOG_LEVEL=info"></go>

<!-- Show code then run it -->
<go src="examples" code="main.go" run="."></go>

<!-- Expected non-zero exit -->
<go run="fail.go" exit="1"></go>

<!-- Custom timeout -->
<go run="slow.go" timeout="60s"></go>
```

### `<cmd>` - Execute Shell Commands

Run arbitrary shell commands and capture output.

```html
<!-- Simple command -->
<cmd exec="echo Hello World"></cmd>

<!-- Command in a directory -->
<cmd exec="ls -la" src="mydir"></cmd>

<!-- Expect failure -->
<cmd exec="false" exit="1"></cmd>

<!-- Any non-zero exit -->
<cmd exec="might-fail" exit="-1"></cmd>

<!-- With environment -->
<cmd exec="printenv FOO" environ="FOO=bar"></cmd>

<!-- Custom timeout (default 30s) -->
<cmd exec="long-process" timeout="120s"></cmd>
```

### `<include>` - Include Other Documents

Compose documents from partials.

```html
<!-- Include another markdown file -->
<include src="docs/intro.md"></include>

<!-- Include from subdirectory -->
<include src="chapters/getting-started/index.md"></include>
```

Included files maintain their relative paths for assets and links.

### `<youtube>` - Embed YouTube Videos

Embed YouTube videos in documentation.

```html
<!-- Basic embed -->
<youtube id="dQw4w9WgXcQ"></youtube>

<!-- With custom title -->
<youtube id="dQw4w9WgXcQ" title="Introduction Video"></youtube>
```

The `id` must be exactly 11 alphanumeric characters (the video ID from YouTube URLs).

## Snippet System

Snippets let you extract specific portions of code files. Mark regions with comments:

### Go, JavaScript, TypeScript

```go
// snippet:example
func Example() {
    fmt.Println("This is the snippet content")
}
// snippet:example
```

### HTML, Markdown

```html
<!-- snippet:header -->
<header>Navigation here</header>
<!-- snippet:header -->
```

### Ruby, Shell, YAML

```yaml
# snippet:config
database:
  host: localhost
  port: 5432
# snippet:config
```

### Supported Extensions

| Extension | Comment Format |
|-----------|----------------|
| `.go` | `// snippet:name` |
| `.js`, `.ts` | `// snippet:name` |
| `.html`, `.md` | `<!-- snippet:name -->` |
| `.rb` | `# snippet:name` |
| `.sh` | `# snippet:name` |
| `.yaml`, `.yml` | `# snippet:name` |
| `.env`, `.envrc` | `# snippet:name` |

## CLI Commands

### Export to Markdown

```bash
# Generate README from hype document
hype export -format=markdown -f hype.md > README.md

# Export with custom output
hype export -format=markdown -f docs/guide.md -o output.md
```

### Export to HTML

```bash
hype export -format=html -f document.md > output.html
```

## Best Practices

### Directory Structure

```
project/
├── .hype/
│   ├── hype.md           # Main document
│   └── docs/
│       ├── intro.md
│       ├── examples/
│       │   ├── basic.go
│       │   └── advanced.go
│       └── images/
│           └── diagram.png
├── README.md             # Generated output
└── src/
    └── ...
```

### Document Organization

1. Use `<include>` to break large documents into manageable pieces
2. Keep code examples in dedicated directories near your documentation
3. Use snippets to show relevant portions of larger files
4. Set appropriate timeouts for slow-running commands
5. Use `exit` attribute when demonstrating error cases

### Code Examples

1. Always test that code examples compile/run before documenting
2. Use snippets to keep examples focused and maintainable
3. Include enough context for examples to be understandable
4. Use `code` attribute on `<go>` tags to show source before output

## Error Handling

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `missing src attribute` | `<code>` or `<include>` without `src` | Add `src="path/to/file"` attribute |
| `unclosed snippet` | Snippet comment not closed | Add closing `// snippet:name` comment |
| `snippet not found` | Referenced snippet doesn't exist | Check snippet name matches exactly |
| `duplicate snippet` | Same snippet name used twice | Use unique names for each snippet |
| `invalid YouTube video ID` | ID not 11 alphanumeric chars | Use valid 11-char video ID |
| `exit code mismatch` | Command exit differs from expected | Set correct `exit` attribute or fix command |
| `timeout exceeded` | Command took too long | Increase `timeout` or optimize command |

### Debugging Tips

1. Run `hype export` with verbose output to see processing steps
2. Check that all source files exist and are readable
3. Verify snippet names match exactly (case-sensitive)
4. Test commands manually before adding to documentation

## File References

For complete attribute details, see:
- [Tag Reference](references/tag-reference.md) - All attributes for each tag
- [Troubleshooting](references/troubleshooting.md) - Common errors and solutions
