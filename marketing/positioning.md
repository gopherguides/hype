# Hype Positioning & Messaging Framework

## Overview

Hype is a free, open-source (MIT) Markdown content generator that keeps documentation accurate through dynamic code execution, includes, and validation. Built by [Gopher Guides](https://www.gopherguides.com).

**Website:** [hypemd.dev](https://hypemd.dev)
**X/Twitter:** [@hype_markdown](https://x.com/hype_markdown)
**License:** MIT

---

## Positioning Statement

For developers and technical content teams who maintain code-heavy documentation, Hype is a Markdown content generator that keeps docs accurate by executing and validating code samples automatically. Unlike static doc tools like Docusaurus, MkDocs, or mdBook, Hype treats documentation like code — with reusable packages, executable examples, and built-in validation — so your docs are never out of date.

## Tagline

**Your docs are never out of date.**

---

## Ideal Community Profiles (ICPs)

### 1. OSS Maintainers

Maintain READMEs with code examples that break silently as the API evolves. Need docs that stay correct without manual effort.

**Pain:** "Our README examples haven't worked since v2.3 and nobody noticed until a new contributor filed an issue."

### 2. DevRel & Developer Advocates

Produce tutorials, blog posts, and conference materials with code samples that must work across versions. Reuse content across multiple channels.

**Pain:** "I copy-paste the same setup snippet into five blog posts, then the API changes and I miss two of them."

### 3. Training & Education Teams

Build courses and workshops where every code example must compile and run. A broken sample derails an entire classroom.

**Pain:** "A student hit a broken example mid-workshop and I lost 20 minutes debugging live."

### 4. Docs-as-Code Teams

Treat documentation as a first-class engineering artifact with CI, reviews, and automation. Want the same rigor for docs that they have for code.

**Pain:** "We review docs in PRs but still ship broken examples because there's no way to test them."

---

## Key Differentiators

### 1. Executable Documentation

Hype runs your code samples and includes real output in the rendered docs. If the code doesn't compile or the output changes, you know immediately. Static tools like Docusaurus, MkDocs, and mdBook render whatever you paste in — broken or not.

### 2. Reusable Content Packages

Write a code example or explanation once, include it everywhere — README, blog, book, workshop. Hype's include and partial system works like code imports. Other tools require copy-pasting between separate doc sites, with no single source of truth.

### 3. Built-in Validation

Hype validates that code samples execute, local assets exist, and links resolve. It's like a test suite for your documentation. Alternatives rely on manual review or external linting tools bolted on after the fact.

---

## Messaging Framework

### Pain → Outcome → Proof

| Pain | Outcome | Proof |
|------|---------|-------|
| Code examples rot silently — broken docs erode trust and waste contributors' time | Every code sample is executed and validated on every build, so docs stay accurate automatically | Hype runs `go run`, `go build`, or any command on your fenced code blocks and fails the build if they break |
| Same content copy-pasted across README, blog, slides — one update, five places to fix | Write once, include everywhere — a single source of truth for every code example | Hype's `<include>` and partial system lets you compose docs from reusable Markdown packages |
| No way to "test" documentation — broken examples ship because nobody catches them | Docs get the same validation rigor as code — CI catches doc issues before users do | Run `hype export` in CI to validate all code samples, assets, and links on every PR |
| Locked into one output format — can't reuse tutorial content as slides or a blog post | One source, multiple formats — export to Markdown, HTML, or slides from the same content | `hype export -format=markdown`, `hype export -format=html`, `hype slides` |

### Short-Form Messaging

**One-liner (social bio / package description):**
Markdown content generator with dynamic code execution, includes, and validation. Your docs are never out of date.

**Elevator pitch (30 seconds):**
Hype is a free, open-source Markdown tool that executes your code samples and validates them automatically. Write a code example once, include it in your README, blog, and workshop — and know it actually works. If your API changes and an example breaks, the build fails before your users notice.

**Why Hype? (website hero):**
Documentation with code examples goes stale the moment your API changes. Hype fixes this by executing code samples, validating output, and failing the build when something breaks. Write once, include everywhere, and never ship broken docs again.

---

## Competitive Landscape

| Capability | Hype | Docusaurus | MkDocs | mdBook | Quarto |
|------------|------|------------|--------|--------|--------|
| Code execution & validation | Yes | No | No | No | Yes (notebooks) |
| Reusable content includes | Yes | Limited (MDX imports) | No | No | Yes (includes) |
| Asset validation | Yes | No | No | No | No |
| Multiple export formats | Markdown, HTML, Slides | HTML | HTML | HTML, PDF | HTML, PDF, Word |
| Blog generator | Yes (3 themes) | Yes | Plugin | No | Yes |
| Live preview | Yes | Yes | Yes | Yes | Yes |
| Free & open source | MIT | MIT | BSD | MPL-2.0 | GPL-2.0 |
| Language-agnostic execution | Yes (any command) | No | No | No | Partial (kernels) |

**Key advantage over Quarto:** Hype runs any shell command — not just notebook kernels. No Jupyter dependency, no language-specific setup. Fence a code block, tell Hype how to run it, done.

---

## How to Use This Document

This messaging framework is the source of truth for all public-facing copy:

- **Website hero & landing page** → Use "Why Hype?" section and tagline
- **Social bio (X/Twitter, GitHub)** → Use one-liner
- **Blog post intros** → Use elevator pitch or pain→outcome→proof rows
- **Conference talk abstracts** → Combine positioning statement + differentiators
- **README description** → Use one-liner + elevator pitch
