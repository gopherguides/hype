
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
Go Version: go1.25.6

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
Go Version: go1.25.6

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
Go Version: go1.25.6

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
Go Version: go1.25.6

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
Go Version: go1.25.6

```

You can also be more specific.

`<go doc="-short context.WithCancel"></go>
`

Here is the output for the above command:
```shell
$ go doc -short context.WithCancel

func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
    WithCancel returns a derived context that points to the parent context but
    has a new Done channel. The returned context's Done channel is closed when
    the returned cancel function is called or when the parent context's Done
    channel is closed, whichever happens first.

    Canceling this context releases resources associated with it, so code should
    call cancel as soon as the operations running in this Context complete.

--------------------------------------------------------------------------------
Go Version: go1.25.6

```

For more examples, see the [hype repo](https://www.github.com/gopherguides/hype).

# Arbitrary Commands

You can also use the `cmd` tag and the `exec` attribute to run arbitrary commands and include them in your documentation. Here is the command to run the `tree` command and include it in our documentation:

```html
<cmd exec="tree" src="."></cmd>

```

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

4 directories, 4 files
```

# Embedding YouTube Videos

You can embed YouTube videos directly in your document using the `youtube` tag:

```html
<youtube id="VIDEO_ID"></youtube>

```

Where `VIDEO_ID` is the 11-character video ID from the YouTube URL. For example, from `https://www.youtube.com/watch?v=dQw4w9WgXcQ`, the video ID is `dQw4w9WgXcQ`.

You can also add an optional title for accessibility:

```html
<youtube id="dQw4w9WgXcQ" title="Introduction to Error Handling"></youtube>

```

The `youtube` tag renders a responsive iframe embed with proper security attributes:

```html
<div class="youtube-embed">
  <iframe src="https://www.youtube.com/embed/VIDEO_ID"
    title="Video Title"
    frameborder="0"
    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
    referrerpolicy="strict-origin-when-cross-origin"
    allowfullscreen></iframe>
</div>

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
$ tree ./docs

./docs
├── badges.md
├── license.md
└── quickstart
    ├── hype.md
    ├── includes.md
    └── src
        ├── broken
        │   └── main.go
        └── hello
            └── main.go

5 directories, 6 files
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
on:
  pull_request:
    types: [opened, synchronize, reopened]
jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      pull-requests: write
    steps:
      - uses: actions/checkout@v4
        with:
          ref: ${{github.event.pull_request.head.ref}}
          repository: ${{github.event.pull_request.head.repo.full_name}}
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version-file: 'go.mod'
          cache-dependency-path: subdir/go.sum
      - name: Install hype
        run: go install github.com/gopherguides/hype/cmd/hype@latest
      - name: Run hype
        run: hype export -format=markdown -f hype.md -o README.md
      - name: Check for changes
        id: git-check
        run: |
          git diff --quiet README.md || echo "changed=true" >> $GITHUB_OUTPUT
      - name: Commit README changes if any
        if: steps.git-check.outputs.changed == 'true'
        run: |
          git config --local user.email "github-actions[bot]@users.noreply.github.com"
          git config --local user.name "github-actions[bot]"
          git add README.md
          git commit -m "Update README.md with latest Hype changes"
          git push
```
> *source: .github/workflows/hype.yml*


---

# AI Assistant Integration

Hype includes an [Agent Skill](https://agentskills.io) to help AI coding assistants write hype-compatible documentation. The skill is located in `.agent/skills/hype/`.

## macOS / Linux

Install the hype skill globally for your preferred AI tool with a single command:

### Claude Code

```bash
mkdir -p ~/.claude/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=3 -C ~/.claude/skills hype-main/.agent/skills/hype

```

### OpenAI Codex

```bash
mkdir -p ~/.codex/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=3 -C ~/.codex/skills hype-main/.agent/skills/hype

```

### Gemini CLI

```bash
mkdir -p ~/.gemini/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=3 -C ~/.gemini/skills hype-main/.agent/skills/hype

```

### Cursor

```bash
mkdir -p ~/.cursor/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=3 -C ~/.cursor/skills hype-main/.agent/skills/hype

```

### GitHub Copilot

```bash
mkdir -p ~/.copilot/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=3 -C ~/.copilot/skills hype-main/.agent/skills/hype

```

### Universal (vendor-agnostic)

```bash
mkdir -p ~/.agent/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=3 -C ~/.agent/skills hype-main/.agent/skills/hype

```

## Windows (PowerShell)

Run these commands in PowerShell to install the hype skill:

### Claude Code

```powershell
$tmp="$env:TEMP\hype-skill"; Invoke-WebRequest -Uri "https://github.com/gopherguides/hype/archive/main.zip" -OutFile "$tmp.zip"; Expand-Archive -Path "$tmp.zip" -DestinationPath $tmp -Force; New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.claude\skills" | Out-Null; Copy-Item -Path "$tmp\hype-main\.agent\skills\hype" -Destination "$env:USERPROFILE\.claude\skills\hype" -Recurse -Force; Remove-Item -Path "$tmp.zip","$tmp" -Recurse -Force

```

### OpenAI Codex

```powershell
$tmp="$env:TEMP\hype-skill"; Invoke-WebRequest -Uri "https://github.com/gopherguides/hype/archive/main.zip" -OutFile "$tmp.zip"; Expand-Archive -Path "$tmp.zip" -DestinationPath $tmp -Force; New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.codex\skills" | Out-Null; Copy-Item -Path "$tmp\hype-main\.agent\skills\hype" -Destination "$env:USERPROFILE\.codex\skills\hype" -Recurse -Force; Remove-Item -Path "$tmp.zip","$tmp" -Recurse -Force

```

### Gemini CLI

```powershell
$tmp="$env:TEMP\hype-skill"; Invoke-WebRequest -Uri "https://github.com/gopherguides/hype/archive/main.zip" -OutFile "$tmp.zip"; Expand-Archive -Path "$tmp.zip" -DestinationPath $tmp -Force; New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.gemini\skills" | Out-Null; Copy-Item -Path "$tmp\hype-main\.agent\skills\hype" -Destination "$env:USERPROFILE\.gemini\skills\hype" -Recurse -Force; Remove-Item -Path "$tmp.zip","$tmp" -Recurse -Force

```

### Cursor

```powershell
$tmp="$env:TEMP\hype-skill"; Invoke-WebRequest -Uri "https://github.com/gopherguides/hype/archive/main.zip" -OutFile "$tmp.zip"; Expand-Archive -Path "$tmp.zip" -DestinationPath $tmp -Force; New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.cursor\skills" | Out-Null; Copy-Item -Path "$tmp\hype-main\.agent\skills\hype" -Destination "$env:USERPROFILE\.cursor\skills\hype" -Recurse -Force; Remove-Item -Path "$tmp.zip","$tmp" -Recurse -Force

```

### GitHub Copilot

```powershell
$tmp="$env:TEMP\hype-skill"; Invoke-WebRequest -Uri "https://github.com/gopherguides/hype/archive/main.zip" -OutFile "$tmp.zip"; Expand-Archive -Path "$tmp.zip" -DestinationPath $tmp -Force; New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.copilot\skills" | Out-Null; Copy-Item -Path "$tmp\hype-main\.agent\skills\hype" -Destination "$env:USERPROFILE\.copilot\skills\hype" -Recurse -Force; Remove-Item -Path "$tmp.zip","$tmp" -Recurse -Force

```

### Universal (vendor-agnostic)

```powershell
$tmp="$env:TEMP\hype-skill"; Invoke-WebRequest -Uri "https://github.com/gopherguides/hype/archive/main.zip" -OutFile "$tmp.zip"; Expand-Archive -Path "$tmp.zip" -DestinationPath $tmp -Force; New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.agent\skills" | Out-Null; Copy-Item -Path "$tmp\hype-main\.agent\skills\hype" -Destination "$env:USERPROFILE\.agent\skills\hype" -Recurse -Force; Remove-Item -Path "$tmp.zip","$tmp" -Recurse -Force

```

## Project-Local Install

To install the skill for a specific project only:

**macOS / Linux:**

```bash
mkdir -p .agent/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=3 -C .agent/skills hype-main/.agent/skills/hype

```

**Windows (PowerShell):**

```powershell
$tmp="$env:TEMP\hype-skill"; Invoke-WebRequest -Uri "https://github.com/gopherguides/hype/archive/main.zip" -OutFile "$tmp.zip"; Expand-Archive -Path "$tmp.zip" -DestinationPath $tmp -Force; New-Item -ItemType Directory -Force -Path ".agent\skills" | Out-Null; Copy-Item -Path "$tmp\hype-main\.agent\skills\hype" -Destination ".agent\skills\hype" -Recurse -Force; Remove-Item -Path "$tmp.zip","$tmp" -Recurse -Force

```

## Using openskills

Alternatively, use [openskills](https://www.npmjs.com/package/openskills) for cross-platform installation:

```bash
npm install -g openskills
openskills install gopherguides/hype --universal

```

The skill activates automatically when working with hype documents.

# Issues

There are several issues that still need to be worked on. Please see the issues tab if you are interested in helping.

---

# License

MIT License

Copyright © 2025 Gopher Guides LLC

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

