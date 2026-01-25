<include src="docs/badges.md"></include>

# Hype

Hype is a content generation tool that use traditional Markdown syntax, and allows it to be extended for almost any use to create dynamic, rich, automated output that is easily maintainable and reusable.

Hype follows the same principals that we use for coding:

- packages (keep relevant content in small, reusable packages, with all links relative to the package)
- reuse - write your documentation once (even in your code), and use everywhere (blog, book, github repo, etc)
- partials/includes - support including documents into a larger document (just like code!)
- validation - like tests, but validate all your code samples are valid (or not if that is what you expect).
- asset validation - ensure local assets like images, etc actually exist

## Created with Hype

This README was created with hype. Here was the command we used to create it:

From the `.hype` directory, run:

```
hype export -format=markdown -f hype.md > ../README.md
```

You can also use a [github action](#using-github-actions-to-update-your-readme) to automatically update your README as well.

<include src="docs/quickstart/hype.md"></include>

# README Source

You can view the source for this entire readme in the [.hype](https://github.com/gopherguides/corp/tree/main/.hype) directory.

Here is the current structure that we are using to create this readme:

<cmd exec="tree ./docs" src=".">

<include src=".github/workflows/hype.md"></include>

# AI Assistant Integration

Hype includes an [Agent Skill](https://agentskills.io) to help AI coding assistants write hype-compatible documentation. The skill is located in `.agent/skills/hype/`.

## One-Line Install (curl)

Install the hype skill globally for your preferred AI tool with a single command:

### Claude Code

```bash
mkdir -p ~/.claude/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=2 -C ~/.claude/skills hype-main/.agent/skills/hype
```

### OpenAI Codex

```bash
mkdir -p ~/.codex/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=2 -C ~/.codex/skills hype-main/.agent/skills/hype
```

### Gemini CLI

```bash
mkdir -p ~/.gemini/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=2 -C ~/.gemini/skills hype-main/.agent/skills/hype
```

### Cursor

```bash
mkdir -p ~/.cursor/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=2 -C ~/.cursor/skills hype-main/.agent/skills/hype
```

### GitHub Copilot

```bash
mkdir -p ~/.copilot/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=2 -C ~/.copilot/skills hype-main/.agent/skills/hype
```

### Universal (vendor-agnostic)

```bash
mkdir -p ~/.agent/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=2 -C ~/.agent/skills hype-main/.agent/skills/hype
```

## Project-Local Install

To install the skill for a specific project only:

```bash
mkdir -p .agent/skills && curl -sL https://github.com/gopherguides/hype/archive/main.tar.gz | tar -xz --strip-components=2 -C .agent/skills hype-main/.agent/skills/hype
```

## Using openskills

Alternatively, use [openskills](https://www.npmjs.com/package/openskills) for cross-tool installation:

```bash
npm install -g openskills
openskills install gopherguides/hype --universal
```

The skill activates automatically when working with hype documents.

# Issues

There are several issues that still need to be worked on. Please see the issues tab if you are interested in helping.

<include src="docs/license.md"></include>
