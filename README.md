
[<img alt="Release" src="https://img.shields.io/github/release/goreleaser/goreleaser.svg"></img>](https://github.com/gopherguides/hype/releases/latest)
[<img alt="Go Build Status" src="https://github.com/gopherguides/hype/actions/workflows/tests.yml/badge.svg"></img>](https://github.com/gopherguides/hype/actions)
[<img alt="Go Reference" src="https://pkg.go.dev/badge/github.com/goherguides/hype.svg"></img>](https://pkg.go.dev/github.com/gopherguides/hype)
[<img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/gopherguides/hype"></img>](https://goreportcard.com/report/github.com/gopherguides/hype)
[<img alt="Slack" src="https://img.shields.io/badge/Slack-hype-brightgreen"></img>](https://gophers.slack.com/archives/C05SKNHQY3U)

---

# Hype

Hype is a content generation tool that use traditional Markdown syntax, and allows it to be extended for almost any use to create dynamic, rich, automated output that is easily maintainable and reusable.

Hype follows the same principals that we use for coding:


* packages (keep relevant content in small, reusable packages, with all links relative to the package)
* reuse - write your documentation once (even in your code), and use everywhere (blog, book, github repo, etc)
* partials/includes - support including documents into a larger document (just like code!)
* validation - like tests, but validate all your code samples are valid (or not if that is what you expect).
* asset validation - ensure local assets like images, etc actually exist


## Created with Hype

This README was created with hype. Here was the command we used to create it:

From the `.hype` directory, run:

`hype export -format=markdown -f hype.md > ../README.md
`

You can also use a [github action](#using-github-actions-to-update-your-readme) to automatically update your README as well.

---

# Quick Start Guide

For more in depth examples, you can read our quick start guide
[here](https://www.gopherguides.com/articles/golang-hype-quickstart).

# The Basics

This is the syntax to include a code sample in your document:

`<code src="src/hello/main.go" snippet="example"></code>
`

The above code snippet does the following:


* Includes the code snippet specified in the source code
* Validates that the code compiles


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

`// snippet: <snippet_name_here>
`

You must have a beginning and an ending snippet for the code to work.

The output of including that tag will be as follows:

```go
func main() {
	fmt.Println("Hello World")
}
```
> *source: docs/quickstart/src/hello/main.go:example*


A `snippet` is not required in your `code` tag. The default behavior of a `code` tag is to include the entire source file.

If we leave the tag out, it will result in the following code being included:

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}
```
> *source: docs/quickstart/src/hello/main.go*


Notice that none of the `snippet` comments were in the output? This is because hype recognizes them as directives for the document, and will not show them in the actual output.

# Go Specific Commands

There are a number of [Go](https://go.dev/) specific commands you can run as well. Anything from executing the code and showing the output, to including go doc (from the standard library or your own source code), etc.

Let's look at how we use the `go` tag.

Here is the source code of the Go file we are going to include. Notice the use of the `snippet` comments to identify the area of code we want included. We'll see how to specify that in the next section when we include it in our markdown.

# Running Go Code

The following command will include the go source code, run it, and include the output of the program as well:

`<go src="src/hello" run="."></go>
`

Here is the result that will be included in your document from the above command:

```shell
$ go run .

Hello World

--------------------------------------------------------------------------------
Go Version: go1.23.0

```

## Running and Showing the Code

If you want to both run and show the code with the same tag, you can add the `code` attribute to the tag:

`<go src="src/hello" run="." code="main.go"></go>
`

Now the source code is includes, as well as the output of the program:

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}
```
> *source: docs/quickstart/src/hello/main.go*


---

```shell
$ go run .

Hello World

--------------------------------------------------------------------------------
Go Version: go1.23.0

```

## Snippets with Go

You can also specify the snippet in a `go` tag as well. The result is that it will only include the code snippet in the included source:

`<go src="src/hello" run="." code="main.go#example"></go>
`

You can see now that only the snippet is included, but the output is still the same:

```go
func main() {
	fmt.Println("Hello World")
}
```
> *source: docs/quickstart/src/hello/main.go#example:example*


---

```shell
$ go run .

Hello World

--------------------------------------------------------------------------------
Go Version: go1.23.0

```

## Invalid Code

What if you want to include an example of code that does not compile? We still want the code to be parsed and included, even though the code doesn't compile. For this, we can state the expected output of the program.

`<go src="src/broken" run="." code="main.go#example" exit="1"></go>
`

The result now includes the snippet, and the error output from trying to compile the invalid source code.

```go
func main() {
	fmt.Prin("Hello World")
}
```
> *source: docs/quickstart/src/broken/main.go#example:example*


---

```shell
$ go run .

# github.com/gopherguides/hype/.
./main.go:7:6: undefined: fmt.Prin

--------------------------------------------------------------------------------
Go Version: go1.23.0

```

### GoDoc

While there are a number of `godoc` commands that will allow you to put your documentation from your code directly into your articles as well. Here are some of the commands.

Here is the basic usage first:

`<go doc="-short context"></go>
`

Here is the output for the above command:

```shell
$ go doc -short context

var Canceled = errors.New("context canceled")
var DeadlineExceeded error = deadlineExceededError{}
func AfterFunc(ctx Context, f func()) (stop func() bool)
func Cause(c Context) error
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithCancelCause(parent Context) (ctx Context, cancel CancelCauseFunc)
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)
func WithDeadlineCause(parent Context, d time.Time, cause error) (Context, CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
func WithTimeoutCause(parent Context, timeout time.Duration, cause error) (Context, CancelFunc)
type CancelCauseFunc func(cause error)
type CancelFunc func()
type Context interface{ ... }
    func Background() Context
    func TODO() Context
    func WithValue(parent Context, key, val any) Context
    func WithoutCancel(parent Context) Context

--------------------------------------------------------------------------------
Go Version: go1.23.0

```

You can also be more specific.

`<go doc="-short context.WithCancel"></go>
`

Here is the output for the above command:
```shell
$ go doc -short context.WithCancel

func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
    WithCancel returns a copy of parent with a new Done channel. The returned
    context's Done channel is closed when the returned cancel function is called
    or when the parent context's Done channel is closed, whichever happens
    first.

    Canceling this context releases resources associated with it, so code should
    call cancel as soon as the operations running in this Context complete.

--------------------------------------------------------------------------------
Go Version: go1.23.0

```

For more examples, see the [hype repo](https://www.github.com/gopherguides/hype).

# Arbitrary Commands

You can also use the `cmd` tag and the `exec` attribute to run arbitrary commands and include them in your documentation. Here is the command to run the `tree` command and include it in our documentation:

`<cmd exec="tree" src="."></cmd>
`

Here is the output:

```shell
$ tree

.
├── hype.md
├── includes.md
└── src
    ├── broken
    │   └── main.go
    └── hello
        └── main.go

3 directories, 4 files
```

# The Export Command

There are several options for running the `hype` command. Most notable is the `export` option:

`$ hype export -h

Usage of hype:
  -f string
     optional file name to preview, if not provided, defaults to hype.md (default "hype.md")
  -format string
     content type to export to: markdown, html (default "markdown")
  -timeout duration
     timeout for execution, defaults to 30 seconds (30s) (default 5s)
  -v enable verbose output for debugging

Usage: hype export [options]

Examples:
 hype export -format html
 hype export -f README.md -format html
 hype export -f README.md -format markdown -timeout=10s
`

This allows you to see your compiled document either as a single markdown, or as an html document that you can preview in the browser.

# Including Markdown

To include a markdown file, use the include tag. This will run that markdown file through the hype.Parser being used and append the results to the current document.

The paths specified in the src attribute of the include are relative to the markdown file they are used in. This allows you to move entire directory structures around in your project without having to change references within the documents themselves.

The following code will parse the code/code.md and sourceable/sourceable.md documents and append them to the end of the document they were included in.

```md
<include src="code/code.md"></include>

<include src="sourceable/sourceable.md"></include>
```
> *source: docs/quickstart/includes.md*


---

# README Source

You can view the source for this entire readme in the [.hype](https://github.com/gopherguides/corp/tree/main/.hype) directory.

Here is the current structure that we are using to create this readme:

```shell
$ tree

.
├── Makefile
├── README.md
├── atom.go
├── atomx
│   ├── atoms.go
│   ├── atoms.ts
│   ├── atomx.go
│   ├── atomx_test.go
│   └── gen.go
├── attributes.go
├── attributes_test.go
├── binding
│   ├── errors.go
│   ├── part.go
│   ├── part_test.go
│   ├── testdata
│   │   ├── toc
│   │   │   ├── 01-one
│   │   │   │   ├── assets
│   │   │   │   │   └── foo.png
│   │   │   │   ├── hype.md
│   │   │   │   ├── hype.tex.gold
│   │   │   │   ├── simple
│   │   │   │   │   ├── assets
│   │   │   │   │   │   └── foo.png
│   │   │   │   │   ├── simple.md
│   │   │   │   │   └── src
│   │   │   │   │       └── greet
│   │   │   │   │           ├── go.mod
│   │   │   │   │           └── main.go
│   │   │   │   └── src
│   │   │   │       └── greet
│   │   │   │           ├── go.mod
│   │   │   │           └── main.go
│   │   │   ├── 02-two
│   │   │   │   ├── assets
│   │   │   │   │   └── foo.png
│   │   │   │   ├── hype.md
│   │   │   │   ├── hype.tex.gold
│   │   │   │   ├── simple
│   │   │   │   │   ├── assets
│   │   │   │   │   │   └── foo.png
│   │   │   │   │   ├── simple.md
│   │   │   │   │   └── src
│   │   │   │   │       └── greet
│   │   │   │   │           ├── go.mod
│   │   │   │   │           └── main.go
│   │   │   │   └── src
│   │   │   │       └── greet
│   │   │   │           ├── go.mod
│   │   │   │           └── main.go
│   │   │   └── 03-three
│   │   │       ├── assets
│   │   │       │   └── foo.png
│   │   │       ├── hype.md
│   │   │       ├── hype.tex.gold
│   │   │       ├── simple
│   │   │       │   ├── assets
│   │   │       │   │   └── foo.png
│   │   │       │   ├── simple.md
│   │   │       │   └── src
│   │   │       │       └── greet
│   │   │       │           ├── go.mod
│   │   │       │           └── main.go
│   │   │       └── src
│   │   │           └── greet
│   │   │               ├── go.mod
│   │   │               └── main.go
│   │   └── whole
│   │       └── simple
│   │           ├── 01-one
│   │           │   └── hype.md
│   │           ├── 02-two
│   │           │   └── hype.md
│   │           └── 03-three
│   │               └── hype.md
│   ├── whole.go
│   └── whole_test.go
├── body.go
├── body_test.go
├── cmd
│   └── hype
│       ├── cli
│       │   ├── binding.go
│       │   ├── binding_test.go
│       │   ├── cli.go
│       │   ├── cli_darwin.go
│       │   ├── commander.go
│       │   ├── encode.go
│       │   ├── encode_test.go
│       │   ├── env.go
│       │   ├── env_test.go
│       │   ├── export.go
│       │   ├── marked.go
│       │   ├── marked_test.go
│       │   ├── parser.go
│       │   ├── pwd.go
│       │   ├── slides.go
│       │   ├── testdata
│       │   │   ├── encode
│       │   │   │   └── json
│       │   │   │       ├── hype.md
│       │   │   │       └── success
│       │   │   │           ├── execute-file.json
│       │   │   │           └── parse-file.json
│       │   │   ├── latex
│       │   │   │   ├── file
│       │   │   │   │   ├── assets
│       │   │   │   │   │   └── foo.png
│       │   │   │   │   ├── hype.tex.gold
│       │   │   │   │   ├── index.md
│       │   │   │   │   ├── simple
│       │   │   │   │   │   ├── assets
│       │   │   │   │   │   │   └── foo.png
│       │   │   │   │   │   ├── simple.md
│       │   │   │   │   │   └── src
│       │   │   │   │   │       └── greet
│       │   │   │   │   │           ├── go.mod
│       │   │   │   │   │           └── main.go
│       │   │   │   │   └── src
│       │   │   │   │       └── greet
│       │   │   │   │           ├── go.mod
│       │   │   │   │           └── main.go
│       │   │   │   ├── multi
│       │   │   │   │   ├── one
│       │   │   │   │   │   ├── assets
│       │   │   │   │   │   │   └── foo.png
│       │   │   │   │   │   ├── hype.md
│       │   │   │   │   │   ├── hype.tex.gold
│       │   │   │   │   │   ├── simple
│       │   │   │   │   │   │   ├── assets
│       │   │   │   │   │   │   │   └── foo.png
│       │   │   │   │   │   │   ├── simple.md
│       │   │   │   │   │   │   └── src
│       │   │   │   │   │   │       └── greet
│       │   │   │   │   │   │           ├── go.mod
│       │   │   │   │   │   │           └── main.go
│       │   │   │   │   │   └── src
│       │   │   │   │   │       └── greet
│       │   │   │   │   │           ├── go.mod
│       │   │   │   │   │           └── main.go
│       │   │   │   │   ├── three
│       │   │   │   │   │   ├── assets
│       │   │   │   │   │   │   └── foo.png
│       │   │   │   │   │   ├── hype.md
│       │   │   │   │   │   ├── hype.tex.gold
│       │   │   │   │   │   ├── simple
│       │   │   │   │   │   │   ├── assets
│       │   │   │   │   │   │   │   └── foo.png
│       │   │   │   │   │   │   ├── simple.md
│       │   │   │   │   │   │   └── src
│       │   │   │   │   │   │       └── greet
│       │   │   │   │   │   │           ├── go.mod
│       │   │   │   │   │   │           └── main.go
│       │   │   │   │   │   └── src
│       │   │   │   │   │       └── greet
│       │   │   │   │   │           ├── go.mod
│       │   │   │   │   │           └── main.go
│       │   │   │   │   └── two
│       │   │   │   │       ├── assets
│       │   │   │   │       │   └── foo.png
│       │   │   │   │       ├── hype.md
│       │   │   │   │       ├── hype.tex.gold
│       │   │   │   │       ├── simple
│       │   │   │   │       │   ├── assets
│       │   │   │   │       │   │   └── foo.png
│       │   │   │   │       │   ├── simple.md
│       │   │   │   │       │   └── src
│       │   │   │   │       │       └── greet
│       │   │   │   │       │           ├── go.mod
│       │   │   │   │       │           └── main.go
│       │   │   │   │       └── src
│       │   │   │   │           └── greet
│       │   │   │   │               ├── go.mod
│       │   │   │   │               └── main.go
│       │   │   │   └── simple
│       │   │   │       ├── assets
│       │   │   │       │   └── foo.png
│       │   │   │       ├── hype.md
│       │   │   │       ├── hype.tex.gold
│       │   │   │       ├── simple
│       │   │   │       │   ├── assets
│       │   │   │       │   │   └── foo.png
│       │   │   │       │   ├── simple.md
│       │   │   │       │   └── src
│       │   │   │       │       └── greet
│       │   │   │       │           ├── go.mod
│       │   │   │       │           └── main.go
│       │   │   │       └── src
│       │   │   │           └── greet
│       │   │   │               ├── go.mod
│       │   │   │               └── main.go
│       │   │   └── whole
│       │   │       └── simple
│       │   │           ├── 01-one
│       │   │           │   └── hype.md
│       │   │           ├── 02-two
│       │   │           │   └── hype.md
│       │   │           └── 03-three
│       │   │               └── hype.md
│       │   ├── toc.go
│       │   ├── toc_test.go
│       │   └── vscode.go
│       └── main.go
├── cmd.go
├── cmd_error.go
├── cmd_error_test.go
├── cmd_result.go
├── cmd_result_test.go
├── cmd_test.go
├── code.go
├── code_test.go
├── comment.go
├── comment_test.go
├── dist
│   ├── CHANGELOG.md
│   ├── artifacts.json
│   ├── config.yaml
│   ├── hype_0.1.0_checksums.txt
│   ├── hype_Darwin_arm64.tar.gz
│   ├── hype_Darwin_x86_64.tar.gz
│   ├── hype_Linux_arm64.tar.gz
│   ├── hype_Linux_i386.tar.gz
│   ├── hype_Linux_x86_64.tar.gz
│   ├── hype_Windows_arm64.zip
│   ├── hype_Windows_i386.zip
│   ├── hype_Windows_x86_64.zip
│   ├── hype_darwin_amd64_v1
│   │   └── hype
│   ├── hype_darwin_arm64
│   │   └── hype
│   ├── hype_linux_386
│   │   └── hype
│   ├── hype_linux_amd64_v1
│   │   └── hype
│   ├── hype_linux_arm64
│   │   └── hype
│   ├── hype_windows_386
│   │   └── hype.exe
│   ├── hype_windows_amd64_v1
│   │   └── hype.exe
│   ├── hype_windows_arm64
│   │   └── hype.exe
│   └── metadata.json
├── docs
│   ├── badges.md
│   ├── license.md
│   └── quickstart
│       ├── hype.md
│       ├── includes.md
│       └── src
│           ├── broken
│           │   └── main.go
│           └── hello
│               └── main.go
├── document.go
├── document_test.go
├── element.go
├── element_test.go
├── empty.go
├── empty_test.go
├── errors.go
├── execute.go
├── execute_error.go
├── execute_error_test.go
├── execute_test.go
├── fenced_code.go
├── fenced_code_test.go
├── figcaption.go
├── figcaption_test.go
├── figure.go
├── figure_test.go
├── finders.go
├── finders_test.go
├── go.mod
├── go.sum
├── godoc.go
├── golang.go
├── golang_test.go
├── heading.go
├── heading_test.go
├── hype.go
├── hype.md
├── hype_test.go
├── image.go
├── image_test.go
├── include.go
├── include_test.go
├── inline_code.go
├── inline_code_test.go
├── internal
│   └── lone
│       ├── ranger.go
│       └── ranger_test.go
├── li.go
├── li_test.go
├── license.md
├── link.go
├── link_test.go
├── md.go
├── md_test.go
├── mdx
│   ├── parser.go
│   ├── parser_test.go
│   └── testdata
│       ├── assignment.md
│       ├── basics.md
│       ├── booleans.md
│       ├── constants.md
│       ├── hype.md
│       ├── numbers.md
│       ├── src
│       │   ├── constants
│       │   │   ├── const
│       │   │   │   └── main.go
│       │   │   ├── const-err
│       │   │   │   └── main.go
│       │   │   ├── const-infer
│       │   │   │   └── main.go
│       │   │   └── const_type
│       │   │       └── main.go
│       │   ├── go.mod
│       │   ├── numbers
│       │   │   ├── maxuint8
│       │   │   │   └── main.go
│       │   │   ├── maxuint8-overflow
│       │   │   │   └── main.go
│       │   │   └── maxuint8-saturation
│       │   │       └── main.go
│       │   ├── utf8
│       │   │   ├── utf8
│       │   │   │   └── main.go
│       │   │   ├── utf8-loop
│       │   │   │   └── main.go
│       │   │   ├── utf8-range
│       │   │   │   └── main.go
│       │   │   └── utf8-rune
│       │   │       └── main.go
│       │   └── variables
│       │       ├── multiple
│       │       │   └── main.go
│       │       └── zero
│       │           └── main.go
│       ├── strings.md
│       ├── utf8.md
│       └── variables.md
├── metadata.go
├── metadata_test.go
├── node.go
├── node_test.go
├── now.go
├── now_test.go
├── ol.go
├── ol_test.go
├── page.go
├── page_test.go
├── paragraph.go
├── paragraph_test.go
├── parse_error.go
├── parse_error_test.go
├── parser.go
├── parser_test.go
├── post_execute.go
├── post_execute_error.go
├── post_execute_error_test.go
├── post_execute_test.go
├── post_parse_error.go
├── post_parse_error_test.go
├── post_parser.go
├── post_parser_test.go
├── pre_execute.go
├── pre_execute_error.go
├── pre_execute_error_test.go
├── pre_execute_test.go
├── pre_parse_error.go
├── pre_parse_error_test.go
├── pre_parser.go
├── pre_parser_test.go
├── ref.go
├── ref_processor.go
├── ref_processor_test.go
├── ref_test.go
├── references_test.go
├── restripe_figures.go
├── revive.toml
├── slides
│   ├── app.go
│   └── templates
│       ├── assets
│       │   ├── app.css
│       │   └── app.js
│       └── slides.html
├── snippet.go
├── snippet_test.go
├── source_code.go
├── source_code_test.go
├── table.go
├── table_test.go
├── tag.go
├── td.go
├── td_test.go
├── testdata
│   ├── auto
│   │   ├── blockquote
│   │   │   ├── html
│   │   │   │   ├── hype.gold
│   │   │   │   ├── hype.md
│   │   │   │   └── shine.txt
│   │   │   └── md
│   │   │       ├── hype.gold
│   │   │       ├── hype.md
│   │   │       └── shine.txt
│   │   ├── commands
│   │   │   ├── greet
│   │   │   │   ├── hype.gold
│   │   │   │   ├── hype.md
│   │   │   │   └── src
│   │   │   │       ├── go.mod
│   │   │   │       └── main.go
│   │   │   ├── results
│   │   │   │   ├── data
│   │   │   │   │   ├── hype.gold
│   │   │   │   │   └── hype.md
│   │   │   │   └── truncate
│   │   │   │       ├── hype.gold
│   │   │   │       ├── hype.md
│   │   │   │       └── src
│   │   │   │           ├── go.mod
│   │   │   │           └── main.go
│   │   │   ├── side-by-side
│   │   │   │   ├── hype.gold
│   │   │   │   ├── hype.md
│   │   │   │   └── values
│   │   │   │       ├── _string.md
│   │   │   │       ├── assets
│   │   │   │       │   └── string-keys.svg
│   │   │   │       ├── src
│   │   │   │       │   └── string-keys
│   │   │   │       │       ├── go.mod
│   │   │   │       │       └── main.go
│   │   │   │       └── values.md
│   │   │   └── timeout
│   │   │       ├── go.mod
│   │   │       └── main.go
│   │   ├── metadata
│   │   │   └── simple
│   │   │       ├── hype.gold
│   │   │       └── hype.md
│   │   ├── parser
│   │   │   └── hello
│   │   │       ├── hype.gold
│   │   │       ├── hype.md
│   │   │       └── second
│   │   │           ├── second.md
│   │   │           └── src
│   │   │               ├── go.mod
│   │   │               └── main.go
│   │   ├── refs
│   │   │   ├── fenced
│   │   │   │   ├── hype.gold
│   │   │   │   └── hype.md
│   │   │   ├── figure-styles
│   │   │   │   ├── hype.gold
│   │   │   │   └── hype.md
│   │   │   ├── images
│   │   │   │   ├── assets
│   │   │   │   │   └── nodes.svg
│   │   │   │   ├── hype.gold
│   │   │   │   └── hype.md
│   │   │   ├── includes
│   │   │   │   ├── assets
│   │   │   │   │   └── foo.png
│   │   │   │   ├── hype.gold
│   │   │   │   ├── hype.md
│   │   │   │   ├── simple
│   │   │   │   │   ├── assets
│   │   │   │   │   │   └── foo.png
│   │   │   │   │   ├── simple.md
│   │   │   │   │   └── src
│   │   │   │   │       └── greet
│   │   │   │   │           ├── go.mod
│   │   │   │   │           └── main.go
│   │   │   │   └── src
│   │   │   │       └── greet
│   │   │   │           ├── go.mod
│   │   │   │           └── main.go
│   │   │   └── simple
│   │   │       ├── assets
│   │   │       │   └── foo.png
│   │   │       ├── hype.gold
│   │   │       ├── hype.md
│   │   │       └── src
│   │   │           └── greet
│   │   │               ├── go.mod
│   │   │               └── main.go
│   │   ├── snippets
│   │   │   ├── range
│   │   │   │   ├── hype.gold
│   │   │   │   ├── hype.md
│   │   │   │   └── src
│   │   │   │       ├── go.mod
│   │   │   │       └── main.go
│   │   │   └── simple
│   │   │       ├── hype.gold
│   │   │       ├── hype.md
│   │   │       └── src
│   │   │           ├── go.mod
│   │   │           └── main.go
│   │   ├── toc
│   │   │   ├── hype.gold
│   │   │   └── hype.md
│   │   └── vars
│   │       ├── details
│   │       │   ├── hype.gold
│   │       │   └── hype.md
│   │       ├── metadata
│   │       │   ├── hype.gold
│   │       │   └── hype.md
│   │       └── var_tag
│   │           ├── hype.gold
│   │           └── hype.md
│   ├── commands
│   │   └── bad-exit
│   │       ├── go.mod
│   │       └── main.go
│   ├── doc
│   │   ├── execution
│   │   │   ├── failure
│   │   │   │   └── hype.md
│   │   │   ├── nested_failure
│   │   │   │   ├── hype.md
│   │   │   │   └── second
│   │   │   │       ├── second.md
│   │   │   │       └── src
│   │   │   │           ├── go.mod
│   │   │   │           └── main.go
│   │   │   └── success
│   │   │       └── hype.md
│   │   ├── pages
│   │   │   ├── hype.md
│   │   │   └── second
│   │   │       ├── second.md
│   │   │       └── src
│   │   │           ├── go.mod
│   │   │           └── main.go
│   │   ├── simple
│   │   │   └── hype.md
│   │   ├── snippets
│   │   │   ├── hype.md
│   │   │   └── src
│   │   │       └── main.ts
│   │   └── to_md
│   │       ├── basics
│   │       │   ├── basics.md
│   │       │   └── src
│   │       │       └── background
│   │       │           ├── empty
│   │       │           │   ├── go.mod
│   │       │           │   └── main.go
│   │       │           └── implementation
│   │       │               ├── go.mod
│   │       │               └── main.go
│   │       ├── cancellation
│   │       │   ├── assets
│   │       │   │   ├── cancellation.svg
│   │       │   │   └── nodes.svg
│   │       │   ├── cancellation.md
│   │       │   └── src
│   │       │       ├── basic
│   │       │       │   ├── go.mod
│   │       │       │   └── main.go
│   │       │       └── cancelling
│   │       │           ├── go.mod
│   │       │           └── main.go
│   │       ├── errors
│   │       │   ├── errors.md
│   │       │   └── src
│   │       │       ├── canceled
│   │       │       │   ├── go.mod
│   │       │       │   └── main.go
│   │       │       └── deadline
│   │       │           ├── go.mod
│   │       │           └── main.go
│   │       ├── graffles
│   │       │   └── context.graffle
│   │       ├── hype.gold
│   │       ├── hype.md
│   │       ├── nodes
│   │       │   ├── assets
│   │       │   │   └── nodes.svg
│   │       │   ├── nodes.md
│   │       │   └── src
│   │       │       └── node-tree
│   │       │           ├── go.mod
│   │       │           ├── go.sum
│   │       │           ├── main.go
│   │       │           └── stdout.txt
│   │       ├── rules
│   │       │   └── rules.md
│   │       ├── signals
│   │       │   ├── signals.md
│   │       │   └── src
│   │       │       ├── signals
│   │       │       │   ├── go.mod
│   │       │       │   └── main.go
│   │       │       └── testing
│   │       │           ├── go.mod
│   │       │           ├── go.sum
│   │       │           ├── signals.go
│   │       │           ├── signals_test.go
│   │       │           └── stdout.txt
│   │       ├── timeouts
│   │       │   ├── src
│   │       │   │   ├── timeout
│   │       │   │   │   ├── go.mod
│   │       │   │   │   └── main.go
│   │       │   │   ├── with-deadline
│   │       │   │   │   ├── go.mod
│   │       │   │   │   ├── go.sum
│   │       │   │   │   └── main.go
│   │       │   │   └── with-timeout
│   │       │   │       ├── go.mod
│   │       │   │       ├── go.sum
│   │       │   │       └── main.go
│   │       │   └── timeouts.md
│   │       └── values
│   │           ├── _securing.md
│   │           ├── _strings.md
│   │           ├── assets
│   │           │   └── string-keys.svg
│   │           ├── src
│   │           │   ├── basic
│   │           │   │   ├── go.mod
│   │           │   │   ├── go.sum
│   │           │   │   └── main.go
│   │           │   ├── custom-const
│   │           │   │   ├── go.mod
│   │           │   │   └── main.go
│   │           │   ├── custom-keys
│   │           │   │   ├── go.mod
│   │           │   │   └── main.go
│   │           │   ├── keys
│   │           │   │   ├── go.mod
│   │           │   │   ├── main.go
│   │           │   │   └── stdout.txt
│   │           │   ├── malicious
│   │           │   │   ├── bar
│   │           │   │   │   └── bar.go
│   │           │   │   ├── foo
│   │           │   │   │   └── foo.go
│   │           │   │   ├── go.mod
│   │           │   │   └── main.go
│   │           │   ├── resolution
│   │           │   │   ├── go.mod
│   │           │   │   ├── go.sum
│   │           │   │   ├── main.go
│   │           │   │   └── stdout.txt
│   │           │   ├── secured
│   │           │   │   ├── bar
│   │           │   │   │   └── bar.go
│   │           │   │   ├── foo
│   │           │   │   │   └── foo.go
│   │           │   │   ├── go.mod
│   │           │   │   └── main.go
│   │           │   └── string-keys
│   │           │       ├── go.mod
│   │           │       └── main.go
│   │           └── values.md
│   ├── golang
│   │   └── sym
│   │       ├── cmd
│   │       │   ├── go.mod
│   │       │   └── main.go
│   │       ├── go.mod
│   │       └── sym.go
│   ├── includes
│   │   ├── broken
│   │   │   └── hype.md
│   │   ├── sublevel
│   │   │   ├── below
│   │   │   │   └── b.md
│   │   │   └── hype.md
│   │   └── toplevel
│   │       └── hype.md
│   ├── json
│   │   ├── body.json
│   │   ├── cmd.json
│   │   ├── cmd_error.json
│   │   ├── cmd_result.json
│   │   ├── comment.json
│   │   ├── document.json
│   │   ├── element.json
│   │   ├── execute_error.json
│   │   ├── fenced_code.json
│   │   ├── figcaption.json
│   │   ├── figure.json
│   │   ├── heading.json
│   │   ├── image.json
│   │   ├── include.json
│   │   ├── inline_code.json
│   │   ├── li.json
│   │   ├── link.json
│   │   ├── metadata.json
│   │   ├── now.json
│   │   ├── ol.json
│   │   ├── p.json
│   │   ├── page.json
│   │   ├── parse_error.json
│   │   ├── parser.json
│   │   ├── post_execute_error.json
│   │   ├── post_parse_error.json
│   │   ├── pre_execute_error.json
│   │   ├── pre_parse_error.json
│   │   ├── ref.json
│   │   ├── table.json
│   │   ├── td.json
│   │   ├── th.json
│   │   ├── thead.json
│   │   ├── toc.json
│   │   ├── tr.json
│   │   ├── ul.json
│   │   └── var.json
│   ├── markdown
│   │   └── unknown-atom
│   │       ├── _included.md
│   │       └── hype.md
│   ├── metadata
│   │   ├── multi
│   │   │   └── hype.md
│   │   └── pages
│   │       └── hype.md
│   ├── now
│   │   ├── hype.gold
│   │   └── hype.md
│   ├── parser
│   │   ├── errors
│   │   │   ├── execute
│   │   │   │   └── hype.md
│   │   │   ├── folder
│   │   │   │   ├── 01-one
│   │   │   │   │   ├── assets
│   │   │   │   │   │   └── foo.png
│   │   │   │   │   ├── hype.md
│   │   │   │   │   ├── hype.tex.gold
│   │   │   │   │   ├── simple
│   │   │   │   │   │   ├── assets
│   │   │   │   │   │   │   └── foo.png
│   │   │   │   │   │   ├── simple.md
│   │   │   │   │   │   └── src
│   │   │   │   │   │       └── greet
│   │   │   │   │   │           ├── go.mod
│   │   │   │   │   │           └── main.go
│   │   │   │   │   └── src
│   │   │   │   │       └── greet
│   │   │   │   │           ├── go.mod
│   │   │   │   │           └── main.go
│   │   │   │   ├── 02-two
│   │   │   │   │   ├── assets
│   │   │   │   │   │   └── foo.png
│   │   │   │   │   ├── hype.md
│   │   │   │   │   ├── hype.tex.gold
│   │   │   │   │   ├── simple
│   │   │   │   │   │   ├── assets
│   │   │   │   │   │   │   └── foo.png
│   │   │   │   │   │   ├── simple.md
│   │   │   │   │   │   └── src
│   │   │   │   │   │       └── greet
│   │   │   │   │   │           ├── go.mod
│   │   │   │   │   │           └── main.go
│   │   │   │   │   └── src
│   │   │   │   │       └── greet
│   │   │   │   │           ├── go.mod
│   │   │   │   │           └── main.go
│   │   │   │   └── 03-three
│   │   │   │       ├── assets
│   │   │   │       │   └── foo.png
│   │   │   │       ├── hype.md
│   │   │   │       ├── hype.tex.gold
│   │   │   │       ├── simple
│   │   │   │       │   ├── assets
│   │   │   │       │   │   └── foo.png
│   │   │   │       │   ├── simple.md
│   │   │   │       │   └── src
│   │   │   │       │       └── greet
│   │   │   │       │           ├── go.mod
│   │   │   │       │           └── main.go
│   │   │   │       └── src
│   │   │   │           └── greet
│   │   │   │               ├── go.mod
│   │   │   │               └── main.go
│   │   │   ├── post_execute
│   │   │   │   └── hype.md
│   │   │   ├── post_parse
│   │   │   │   └── hype.md
│   │   │   ├── pre_execute
│   │   │   │   └── hype.md
│   │   │   └── pre_parse
│   │   │       └── hype.md
│   │   └── folder
│   │       ├── 01-one
│   │       │   ├── assets
│   │       │   │   └── foo.png
│   │       │   ├── hype.md
│   │       │   ├── hype.tex.gold
│   │       │   ├── simple
│   │       │   │   ├── assets
│   │       │   │   │   └── foo.png
│   │       │   │   ├── simple.md
│   │       │   │   └── src
│   │       │   │       └── greet
│   │       │   │           ├── go.mod
│   │       │   │           └── main.go
│   │       │   └── src
│   │       │       └── greet
│   │       │           ├── go.mod
│   │       │           └── main.go
│   │       ├── 02-two
│   │       │   ├── assets
│   │       │   │   └── foo.png
│   │       │   ├── hype.md
│   │       │   ├── hype.tex.gold
│   │       │   ├── simple
│   │       │   │   ├── assets
│   │       │   │   │   └── foo.png
│   │       │   │   ├── simple.md
│   │       │   │   └── src
│   │       │   │       └── greet
│   │       │   │           ├── go.mod
│   │       │   │           └── main.go
│   │       │   └── src
│   │       │       └── greet
│   │       │           ├── go.mod
│   │       │           └── main.go
│   │       └── 03-three
│   │           ├── assets
│   │           │   └── foo.png
│   │           ├── hype.md
│   │           ├── hype.tex.gold
│   │           ├── simple
│   │           │   ├── assets
│   │           │   │   └── foo.png
│   │           │   ├── simple.md
│   │           │   └── src
│   │           │       └── greet
│   │           │           ├── go.mod
│   │           │           └── main.go
│   │           └── src
│   │               └── greet
│   │                   ├── go.mod
│   │                   └── main.go
│   ├── snippets
│   │   ├── snip.txt
│   │   ├── snippets.go
│   │   ├── snippets.html
│   │   ├── snippets.js
│   │   └── snippets.rb
│   ├── table
│   │   ├── data
│   │   │   └── hype.md
│   │   ├── headless
│   │   │   ├── hype.gold
│   │   │   └── hype.md
│   │   ├── md_in_html
│   │   │   ├── hype.gold
│   │   │   └── hype.md
│   │   └── md_in_md
│   │       ├── hype.gold
│   │       └── hype.md
│   ├── to_md
│   │   └── source_code
│   │       ├── full
│   │       │   ├── hype.gold
│   │       │   ├── hype.md
│   │       │   └── src
│   │       │       └── main.go
│   │       └── snippet
│   │           ├── hype.gold
│   │           ├── hype.md
│   │           └── src
│   │               └── main.go
│   └── toc
│       ├── hype.gold
│       └── hype.md
├── text.go
├── th.go
├── th_test.go
├── thead.go
├── thead_test.go
├── time.go
├── title.go
├── title_test.go
├── tmpl.go
├── tmpl_test.go
├── toc.go
├── toc_test.go
├── tr.go
├── tr_test.go
├── type.go
├── ul.go
├── ul_test.go
├── unwrap.go
├── var.go
└── var_test.go

322 directories, 581 files
```
---

# Using Github Actions to update your README

This repo uses the action to keep the README up to date.

## Requirements

For this action to work, you need to either configure your repo with specific permissions, or use a `personal access token`.

### Repo Permissions

You need to give permission to your GitHub Actions to create a pull request in your GitHub repo settings _(Settings -> Actions -> General)_.

Under `Workflow Permissions`


* Check `Allow GitHub Actions to create and approve pull requests`.
* Check `Read and write permissions`


### Personal Access Token

Alternately, you can use tokens to give permission to your action.

It is recommend to use a GitHub [Personnal Acces Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-fine-grained-personal-access-token) like: `${{secrets.PAT}}` instead of using `${{secrets.GITHUB_TOKEN}}` in GitHub Actions.

## The Action

The current action is set to only generate the readme on a pull request and commit it back to that same pull request.  You can modify this to your own needs.

```yml
name: Generate README with Hype
on: [pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          repository: ${{ github.event.pull_request.head.repo.full_name }}
          ref: ${{ github.event.pull_request.head.ref }}
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: "1.22.x"
          cache-dependency-path: subdir/go.sum
      - name: Install hype
        run: go install github.com/gopherguides/hype/cmd/hype@latest
      - name: Run hype
        run: hype export -format=markdown -f hype.md > README.md
      - name: Commit README back to the repo
        run: |-
          git rev-parse --abbrev-ref HEAD
          git config user.name 'GitHub Actions'
          git config user.email 'actions@github.com'
          git diff --quiet || (git add README.md && git commit -am "Updated README")
          git push origin ${{github.event.pull_request.head.ref}}
```
> *source: .github/workflows/hype.yml*


---

# Issues

There are several issues that still need to be worked on. Please see the issues tab if you are interested in helping.

---

# License

[Hype](https://github.com/gopherguides/hype) by [Gopher Guides LLC](https://github.com/gopherguides) is licensed under [Attribution-NonCommercial-ShareAlike 4.0 International<img src="https://mirrors.creativecommons.org/presskit/icons/cc.svg?ref=chooser-v1" style="height:22px!important;margin-left:3px;vertical-align:text-bottom;"></img><img src="https://mirrors.creativecommons.org/presskit/icons/by.svg?ref=chooser-v1" style="height:22px!important;margin-left:3px;vertical-align:text-bottom;"></img><img src="https://mirrors.creativecommons.org/presskit/icons/nc.svg?ref=chooser-v1" style="height:22px!important;margin-left:3px;vertical-align:text-bottom;"></img><img src="https://mirrors.creativecommons.org/presskit/icons/sa.svg?ref=chooser-v1" style="height:22px!important;margin-left:3px;vertical-align:text-bottom;"></img>](http://creativecommons.org/licenses/by-nc-sa/4.0/?ref=chooser-v1)

