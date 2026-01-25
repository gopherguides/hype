# Hype Troubleshooting Guide

Common errors and their solutions when working with hype documents.

## Code Tag Errors

### Missing src attribute

**Error**: `missing src attribute`

**Cause**: The `<code>` tag requires a `src` attribute to know which file to display.

**Solution**: Add the `src` attribute pointing to your source file.

```html
<!-- Wrong -->
<code></code>

<!-- Correct -->
<code src="main.go"></code>
```

### File not found

**Error**: `failed to read file "path/to/file": open path/to/file: no such file or directory`

**Cause**: The specified source file doesn't exist or the path is incorrect.

**Solution**:
1. Verify the file exists at the specified path
2. Paths are relative to the hype document, not the current working directory
3. Check for typos in the filename

```html
<!-- If hype document is in .hype/ and file is in .hype/examples/ -->
<code src="examples/main.go"></code>
```

## Snippet Errors

### Unclosed snippet

**Error**: `unclosed snippet: path/to/file: ["snippetname"]`

**Cause**: A snippet was opened but never closed with a matching comment.

**Solution**: Add the closing snippet comment with the same name.

```go
// snippet:example
func Example() {
    fmt.Println("Hello")
}
// snippet:example  // <-- Don't forget this!
```

### Snippet not found

**Error**: `snippet "name" not found in "file.go"`

**Cause**: The requested snippet name doesn't exist in the source file.

**Solution**:
1. Check that the snippet name matches exactly (case-sensitive)
2. Verify the snippet comments are in the correct format for the file type
3. Ensure the file has been saved with the snippet markers

```html
<!-- Requesting snippet "Example" -->
<code src="main.go" snippet="Example"></code>
```

```go
// In main.go - name must match exactly
// snippet:Example
func Example() { }
// snippet:Example
```

### Duplicate snippet

**Error**: `duplicate snippet: path/to/file#name`

**Cause**: The same snippet name is used more than once in a file.

**Solution**: Use unique names for each snippet in a file.

```go
// Wrong - same name used twice
// snippet:example
func One() {}
// snippet:example

// snippet:example  // Duplicate!
func Two() {}
// snippet:example

// Correct - unique names
// snippet:example-one
func One() {}
// snippet:example-one

// snippet:example-two
func Two() {}
// snippet:example-two
```

### Wrong comment format

**Cause**: Using the wrong comment format for the file type.

**Solution**: Use the correct comment format for your file extension.

| Extension | Correct Format |
|-----------|----------------|
| `.go`, `.js`, `.ts` | `// snippet:name` |
| `.html`, `.md` | `<!-- snippet:name -->` |
| `.rb`, `.sh`, `.yaml`, `.yml` | `# snippet:name` |

## Include Errors

### Missing src attribute (include)

**Error**: `missing src attribute`

**Cause**: The `<include>` tag requires a `src` attribute.

**Solution**: Add the `src` attribute with the path to the markdown file.

```html
<!-- Wrong -->
<include></include>

<!-- Correct -->
<include src="docs/intro.md"></include>
```

### Included file not found

**Error**: `open path/to/file.md: no such file or directory`

**Cause**: The included markdown file doesn't exist.

**Solution**:
1. Verify the file exists
2. Check the path is relative to the including document
3. Ensure the file has a `.md` extension

## Command Execution Errors

### Exit code mismatch

**Error**: `exit status N` (when N differs from expected)

**Cause**: The command exited with a different code than expected.

**Solution**:
1. If the command should fail, set the expected exit code:
   ```html
   <cmd exec="false" exit="1"></cmd>
   ```
2. If expecting any failure, use `-1`:
   ```html
   <cmd exec="might-fail" exit="-1"></cmd>
   ```
3. If the command should succeed, fix the underlying issue

### Timeout exceeded

**Error**: `context deadline exceeded`

**Cause**: The command took longer than the timeout (default 30s).

**Solution**: Increase the timeout or optimize the command.

```html
<!-- Increase timeout to 2 minutes -->
<cmd exec="slow-build" timeout="120s"></cmd>

<go test="./..." timeout="5m"></go>
```

### Command not found

**Error**: `exec: "command": executable file not found in $PATH`

**Cause**: The command doesn't exist or isn't in PATH.

**Solution**:
1. Verify the command is installed
2. Use the full path to the executable
3. Check that required tools are available in the build environment

## YouTube Errors

### Invalid YouTube video ID

**Error**: `invalid YouTube video ID "xyz": must be 11 alphanumeric characters`

**Cause**: The video ID is not in the correct format.

**Solution**: Use the 11-character video ID from the YouTube URL.

```html
<!-- Wrong - using full URL -->
<youtube id="https://youtube.com/watch?v=dQw4w9WgXcQ"></youtube>

<!-- Wrong - ID too short -->
<youtube id="dQw4w9"></youtube>

<!-- Correct - just the 11-char ID -->
<youtube id="dQw4w9WgXcQ"></youtube>
```

### Finding the video ID

From YouTube URLs:
- `https://www.youtube.com/watch?v=dQw4w9WgXcQ` → ID is `dQw4w9WgXcQ`
- `https://youtu.be/dQw4w9WgXcQ` → ID is `dQw4w9WgXcQ`
- `https://www.youtube.com/embed/dQw4w9WgXcQ` → ID is `dQw4w9WgXcQ`

## Go Tag Errors

### Go command failed

**Error**: Various Go compilation or runtime errors

**Cause**: The Go code has errors or dependencies aren't available.

**Solution**:
1. Test the Go code manually first: `go run main.go`
2. Ensure `go.mod` exists if using modules
3. Run `go mod tidy` to resolve dependencies
4. Check that the `src` directory is correct

### Cross-compilation errors

**Error**: Build errors when using `goos` or `goarch`

**Cause**: Some packages don't support all OS/architecture combinations.

**Solution**:
1. Verify the target is supported by your dependencies
2. Check for CGO dependencies (which can't cross-compile)
3. Test the cross-compilation manually first

```bash
GOOS=linux GOARCH=amd64 go build .
```

## General Tips

### Debugging hype documents

1. **Test incrementally**: Add one tag at a time and verify it works
2. **Check paths**: Most errors are path-related; verify files exist
3. **Run commands manually**: Test `<cmd>` and `<go>` commands in terminal first
4. **Check snippet names**: They're case-sensitive and must match exactly
5. **Validate timeouts**: Long-running commands need appropriate timeouts

### Common path mistakes

```html
<!-- Paths are relative to the hype document location -->

<!-- If document is .hype/hype.md and code is .hype/examples/main.go -->
<code src="examples/main.go"></code>  <!-- Correct -->
<code src=".hype/examples/main.go"></code>  <!-- Wrong -->
<code src="../examples/main.go"></code>  <!-- Wrong (unless examples is sibling to .hype) -->
```
