# HTML Export

Hype can export your documents to styled HTML with built-in CSS themes.

## Basic Usage

```bash
# Export with default GitHub theme
hype export -format html -f hype.md > output.html

# Export to a file
hype export -format html -f hype.md -o output.html
```

## Themes

Hype includes 7 built-in themes:

| Theme | Description |
|-------|-------------|
| `github` | **Default**. Auto light/dark based on system preference |
| `github-dark` | GitHub dark mode only |
| `solarized-light` | Warm light tones |
| `solarized-dark` | Solarized dark variant |
| `swiss` | Minimalist Swiss typography |
| `air` | Clean, centered layout |
| `retro` | Nostalgic/vintage style |

### List Available Themes

```bash
hype export -themes
```

### Select a Theme

```bash
hype export -format html -theme solarized-dark -f hype.md -o output.html
```

## Custom CSS

Use your own CSS file instead of a built-in theme:

```bash
hype export -format html -css ./my-styles.css -f hype.md -o output.html
```

Your custom CSS should style the `.markdown-body` class which wraps the document content.

## Raw HTML (No Styling)

To get raw HTML without any CSS (the previous default behavior):

```bash
hype export -format html -no-css -f hype.md
```

## Flags Reference

| Flag | Description |
|------|-------------|
| `-format html` | Export as HTML |
| `-theme <name>` | Select a built-in theme (default: `github`) |
| `-css <path>` | Use a custom CSS file |
| `-no-css` | Output raw HTML without styling |
| `-themes` | List available themes and exit |
| `-o <path>` | Write output to file instead of stdout |
