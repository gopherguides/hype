# Slides

Hype can generate web-based presentations from your markdown documents using the `slides` command.

## Basic Usage

```bash
# Start slides server on default port (3000)
hype slides -f presentation.md

# Or specify the file as an argument
hype slides presentation.md

# Use a different port
hype slides -port 8080 presentation.md
```

Once started, open your browser to `http://localhost:3000` to view your presentation.

## Creating Slides

Slides are created using the standard hype `<page>` element. Each `<page>` becomes a slide in your presentation:

```markdown
<page>

# Slide 1

Welcome to my presentation!

</page>

<page>

# Slide 2

## Key Points

- Point one
- Point two
- Point three

</page>

<page>

# Code Example

<go src="example" run></go>

</page>
```

## Features

- **Live Code Execution**: Code blocks with `run` attribute execute and display output
- **Syntax Highlighting**: Code blocks are automatically highlighted
- **Navigation**: Use keyboard arrows or click to navigate between slides
- **Web-based**: No additional software required - just a browser

## Flags Reference

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `3000` | Port for the slides server |

## Tips

1. **Keep slides focused**: One main idea per slide works best
2. **Use code examples**: Hype's ability to execute code makes live demos easy
3. **Test navigation**: Check that your slides flow well before presenting
4. **Local images**: Images are served from the document's directory
