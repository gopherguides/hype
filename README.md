
[<img alt="Release" src="https://img.shields.io/github/release/goreleaser/goreleaser.svg"></img>](https://github.com/gopherguides/hype/releases/latest)
[<img alt="Go Build Status" src="https://github.com/gopherguides/hype/actions/workflows/tests.yml/badge.svg"></img>](https://github.com/gopherguides/hype/actions)
[<img alt="Go Reference" src="https://pkg.go.dev/badge/github.com/goherguides/hype.svg"></img>](https://pkg.go.dev/github.com/gopherguides/hype)
[<img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/gopherguides/hype"></img>](https://goreportcard.com/report/github.com/gopherguides/hype)
[<img alt="Slack" src="https://img.shields.io/badge/Slack-hype-brightgreen"></img>](https://gophers.slack.com/archives/C05SKNHQY3U)

---

# Hype

Hype is a content generation tool that use traditional Markdown syntax, and allows it to be extended for almost any use to create dynamic, rich, automated output that is easily maintainable and reusable.

## Created with Hype

This README was created with hype.  Here was the command we used to create it:

From the `.hype` directory, run:

`hype export -format=markdown -f module.md > ../README.md
`

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

Notice the use of the `snippet` comment.  The format for the comment is:

`// snippet: <snippet_name_here>
`

You must have a beginning and an ending snippet for the code to work.

The output of including that tag will be as follows:

```go
func main() {
	fmt.Println("Hello World")
}
```

A `snippet` is not required in your `code` tag.  They default behavior of a `code` tag is to include the entire source file.

If we leave the tag out, it will result in the following code being included:

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}


```

Notice that none of the `snippet` comments were in the output?  This is because hype recognizes them as directives for the document, and will not show them in the actual output.

# Go Specific Commands

There are a number of [Go](https://go.dev/) specific commands you can run as well. Anything from executing the code and showing the output, to including go doc (from the standard library or your own source code), etc.

Let's look at how we use the `go` tag.

Here is the source code of the Go file we are going to include.  Notice the use of the `snippet` comments to identify the area of code we want included.  We'll see how to specify that in the next section when we include it in our markdown.

# Running Go Code

The following command will include the go source code, run it, and include the output of the program as well:

`<go src="src/hello" run="."></go>
`

Here is the result that will be included in your document from the above command:

```shell
$ go run .

Hello World

--------------------------------------------------------------------------------
Go Version: go1.22.0

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

---

```shell
$ go run .

Hello World

--------------------------------------------------------------------------------
Go Version: go1.22.0

```

## Snippets with Go

You can also specify the snippet in a `go` tag as well.  The result is that it will only include the code snippet in the included source:

`<go src="src/hello" run="." code="main.go#example"></go>
`

You can see now that only the snippet is included, but the output is still the same:

```go
func main() {
	fmt.Println("Hello World")
}
```

---

```shell
$ go run .

Hello World

--------------------------------------------------------------------------------
Go Version: go1.22.0

```

## Invalid Code

What if you want to include an example of code that does not compile?  We still want the code to be parsed and included, even though the code doesn't compile.  For this, we can state the expected output of the program.

`<go src="src/broken" run="." code="main.go#example" exit="1"></go>
`

The result now includes the snippet, and the error output from trying to compile the invalid source code.

```go
func main() {
	fmt.Prin("Hello World")
}
```

---

```shell
$ go run .

# github.com/gopherguides/hype/.hype/.
./main.go:7:6: undefined: fmt.Prin

--------------------------------------------------------------------------------
Go Version: go1.22.0

```

### GoDoc

While there are a number of `godoc` commands that will allow you to put your documentation from your code directly into your articles as well.  Here are some of the commands.

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
Go Version: go1.22.0

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
Go Version: go1.22.0

```

For more examples, see the [hype repo](https://www.github.com/gopherguides/hype).

# Arbitrary Commands

You can also use the `cmd` tag and the `exec` attribute to run arbitrary commands and include them in your documentation.  Here is the command to run the `tree` command and include it in our documentation:

`<cmd exec="tree" src="."></cmd>
`

Here is the output:

```shell
$ tree

.
├── includes.md
├── module.md
└── src
    ├── broken
    │   └── main.go
    └── hello
        └── main.go

3 directories, 4 files
```

# The Export Command

There are several options for running the `hype` command.  Most notable is the `export` option:

`$ hype export -h

Usage of hype:
  -f string
    	optional file name to preview, if not provided, defaults to module.md (default "module.md")
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

---

# README Source

You can view the source for this entire readme in the [.hype](https://github.com/gopherguides/corp/tree/main/.hype) directory.

Here is the current structure that we are using to create this readme:

```shell
$ tree

.
├── badges.md
├── license.md
├── module.md
└── quickstart
    ├── includes.md
    ├── module.md
    └── src
        ├── broken
        │   └── main.go
        └── hello
            └── main.go

4 directories, 7 files
```
---

# License

[Hype](https://github.com/gopherguides/hype) by [Gopher Guides LLC](https://github.com/gopherguides) is licensed under [Attribution-NonCommercial-ShareAlike 4.0 International<img src="https://mirrors.creativecommons.org/presskit/icons/cc.svg?ref=chooser-v1" style="height:22px!important;margin-left:3px;vertical-align:text-bottom;"></img><img src="https://mirrors.creativecommons.org/presskit/icons/by.svg?ref=chooser-v1" style="height:22px!important;margin-left:3px;vertical-align:text-bottom;"></img><img src="https://mirrors.creativecommons.org/presskit/icons/nc.svg?ref=chooser-v1" style="height:22px!important;margin-left:3px;vertical-align:text-bottom;"></img><img src="https://mirrors.creativecommons.org/presskit/icons/sa.svg?ref=chooser-v1" style="height:22px!important;margin-left:3px;vertical-align:text-bottom;"></img>](http://creativecommons.org/licenses/by-nc-sa/4.0/?ref=chooser-v1)

