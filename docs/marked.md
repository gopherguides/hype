# Marked 2 Integration

Hype integrates with [Marked 2](https://marked2app.com/), a powerful Markdown preview and export application for macOS.

## Overview

The `marked` command outputs hype documents in a format compatible with Marked 2's custom preprocessor feature. This allows you to preview hype documents directly in Marked 2, including dynamic code execution and includes.

## Setup

1. Open Marked 2 Preferences
2. Go to the "Advanced" tab
3. In "Custom Processor", set the path to the hype binary and enable "Preprocessing"
4. Enter the path: `/path/to/hype marked`

Marked 2 will set the `MARKED_PATH` and `MARKED_ORIGIN` environment variables automatically, telling hype which file to process.

## Basic Usage

The command is designed to be called by Marked 2 automatically:

```bash
hype marked
```

For manual testing (run from your document's directory):

```bash
hype marked -f hype.md
```

## Flags Reference

| Flag | Default | Description |
|------|---------|-------------|
| `-f` | | File to process (if not provided, reads from stdin) |
| `-p` | `false` | Parse only mode - parse the file but don't execute commands |
| `-timeout` | `30s` | Timeout for command execution |
| `-context` | | A folder containing all chapters of a book, for example |
| `-section` | `0` | Target section number |
| `-v` | `false` | Enable verbose output for debugging |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `MARKED_PATH` | Set by Marked 2 - used for file context and relative path resolution |
| `MARKED_ORIGIN` | Set by Marked 2 - the directory of the file being previewed |

## How It Works

1. Marked 2 detects a file change and calls hype as a preprocessor
2. Marked 2 pipes the file contents to hype via stdin
3. Hype uses `MARKED_PATH` for context and relative path resolution
4. Hype processes all includes, executes code blocks, and renders the document
5. The processed Markdown is output to stdout
6. Marked 2 renders the processed Markdown

## Page Breaks

Hype inserts page break comments between pages (`&lt;!--BREAK--&gt;`), which Marked 2 can use for pagination in exported documents.

## Troubleshooting

**Document not updating:**
- Ensure `MARKED_PATH` is being set correctly
- Try running with `-v` flag for verbose output
- Check that hype is installed and accessible

**Timeout errors:**
- Increase the timeout with `-timeout 30s` for documents with slow-running commands
- Use `-p` to test parsing without execution

**Includes not resolving:**
- Verify all include paths are relative to the document
- Check that the source files exist
