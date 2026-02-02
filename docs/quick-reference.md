## Quick Reference

> **Hype** is a Markdown content generator with dynamic code execution, includes, and validation.

### Install

```bash
go install github.com/gopherguides/hype/cmd/hype@latest
```

Or via Homebrew: `brew install gopherguides/hype/hype-md`

### Common Commands

| Command | Description |
|---------|-------------|
| `hype export -format=markdown -f doc.md` | Export to markdown (stdout) |
| `hype export -format=html -f doc.md -o doc.html` | Export to HTML file |
| `hype preview -f doc.md -open` | Live preview with hot reload |
| `hype validate -f doc.md` | Validate document structure |

### Key Tags

| Tag | Purpose | Example |
|-----|---------|---------|
| `<include>` | Include another file | `<include src="other.md">` |
| `<code>` | Show file contents | `<code src="main.go">` |
| `<go>` | Run Go code, show output | `<go run="main.go">` |
| `<cmd>` | Run shell command | `<cmd exec="ls -la">` |
| `<img>` | Include image | `<img src="diagram.png">` |

### AI Assistants

For detailed skill documentation, see [`.agent/skills/hype/`](.agent/skills/hype/).
