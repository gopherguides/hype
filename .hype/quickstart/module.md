# Quick Start Guide

For more in depth examples, you can read our quick start guide
[here](https://www.gopherguides.com/articles/golang-hype-quickstart).

# The Basics

This is the syntax to include a code sample in your document:

```
<code src="src/hello/main.go" snippet="example"></code>
```

The above code snippet does the following:

- Includes the code snippet specified in the source code
- Validates that the code compiles

Here is the source file:

```go
package main

import "fmt"

// snippet: example
func main() {
	fmt.Println("Hello World")
}

// snippet: example
```

Notice the use of the `snippet` comment. The format for the comment is:

```
// snippet: <snippet_name_here>
```

You must have a beginning and an ending snippet for the code to work.

The output of including that tag will be as follows:

<code src="src/hello/main.go" snippet="example"></code>

A `snippet` is not required in your `code` tag. They default behavior of a `code` tag is to include the entire source file.

If we leave the tag out, it will result in the following code being included:

<code src="src/hello/main.go"></code>

Notice that none of the `snippet` comments were in the output? This is because hype recognizes them as directives for the document, and will not show them in the actual output.

# Go Specific Commands

There are a number of [Go](https://go.dev/) specific commands you can run as well. Anything from executing the code and showing the output, to including go doc (from the standard library or your own source code), etc.

Let's look at how we use the `go` tag.

Here is the source code of the Go file we are going to include. Notice the use of the `snippet` comments to identify the area of code we want included. We'll see how to specify that in the next section when we include it in our markdown.

# Running Go Code

The following command will include the go source code, run it, and include the output of the program as well:

```
<go src="src/hello" run="."></go>
```

Here is the result that will be included in your document from the above command:

<go src="src/hello" run="."></go>

## Running and Showing the Code

If you want to both run and show the code with the same tag, you can add the `code` attribute to the tag:

```
<go src="src/hello" run="." code="main.go"></go>
```

Now the source code is includes, as well as the output of the program:

<go src="src/hello" run="." code="main.go"></go>

## Snippets with Go

You can also specify the snippet in a `go` tag as well. The result is that it will only include the code snippet in the included source:

```
<go src="src/hello" run="." code="main.go#example"></go>
```

You can see now that only the snippet is included, but the output is still the same:

<go src="src/hello" run="." code="main.go#example"></go>

## Invalid Code

What if you want to include an example of code that does not compile? We still want the code to be parsed and included, even though the code doesn't compile. For this, we can state the expected output of the program.

```
<go src="src/broken" run="." code="main.go#example" exit="1"></go>
```

The result now includes the snippet, and the error output from trying to compile the invalid source code.

<go src="src/broken" run="." code="main.go#example" exit="1"></go>

### GoDoc

While there are a number of `godoc` commands that will allow you to put your documentation from your code directly into your articles as well. Here are some of the commands.

Here is the basic usage first:

```
<go doc="-short context"></go>
```

Here is the output for the above command:

<go doc="-short context"></go>

You can also be more specific.

```
<go doc="-short context.WithCancel"></go>
```

Here is the output for the above command:
<go doc="-short context.WithCancel"></go>

For more examples, see the [hype repo](https://www.github.com/gopherguides/hype).

# Arbitrary Commands

You can also use the `cmd` tag and the `exec` attribute to run arbitrary commands and include them in your documentation. Here is the command to run the `tree` command and include it in our documentation:

```
<cmd exec="tree" src="."></cmd>
```

Here is the output:

<cmd exec="tree" src="."></cmd>

# The Export Command

There are several options for running the `hype` command. Most notable is the `export` option:

```
$ hype export -h

Usage of hype:
  -f string
    	optional file name to preview, if not provided, defaults to hype.md (default "hype.md")
  -format string
    	content type to export to: markdown, html (default "markdown")
  -timeout duration
    	timeout for execution, defaults to 30 seconds (30s) (default 5s)
  -v	enable verbose output for debugging

Usage: hype export [options]

Examples:
	hype export -format html
	hype export -f README.md -format html
	hype export -f README.md -format markdown -timeout=10s
```

This allows you to see your compiled document either as a single markdown, or as an html document that you can preview in the browser.

# Including Markdown

To include a markdown file, use the include tag. This will run that markdown file through the hype.Parser being used and append the results to the current document.

The paths specified in the src attribute of the include are relative to the markdown file they are used in. This allows you to move entire directory structures around in your project without having to change references within the documents themselves.

The following code will parse the code/code.md and sourceable/sourceable.md documents and append them to the end of the document they were included in.

<code src="includes.md"></code>
