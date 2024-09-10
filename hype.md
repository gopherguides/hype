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

# Issues

There are several issues that still need to be worked on. Please see the issues tab if you are interested in helping.

<include src="docs/license.md"></include>
