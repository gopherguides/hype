# Hype Content Calendar Scaffold (Q2 2026)

## Goals
- Grow awareness among Go developers, docs-first teams, and developer advocates.
- Convert interest into trial installs and docs visits.
- Build proof through practical demos and before/after docs examples.

## Cadence
- **2 posts/week** on X + LinkedIn
- **1 technical blog post/week** on `hypemd.dev/blog`
- **1 demo artifact/week** (template, snippet pack, or short walkthrough)

## Pillars
1. **Product Education** — how Hype works and where it saves time.
2. **Workflow Integrations** — CI, docs repos, blog pipelines, and developer tooling.
3. **Proof & Outcomes** — real examples, performance/reliability wins.
4. **Community** — user showcases, tips, and contribution opportunities.

## Weekly Scaffold (first 4 weeks)
| Week | Product Education | Workflow Integrations | Proof & Outcomes | Community |
|---|---|---|---|---|
| W1 | "What is Hype" explainer post | GitHub Actions publish flow | Before/after docs quality comparison | Ask for examples from users |
| W2 | Parser/CLI quick tips | Docs + blog in one repo | Build-time and error reporting walkthrough | Share template call for contributors |
| W3 | Metadata + SEO support | Social card and OG automation | Traffic uplift case study format | Feature request spotlight |
| W4 | Export/validation deep dive | Dokploy deployment guide | Reliability checklist post | Office-hours / AMA prompt |

## Two Concrete Drafts (PR-ready)

### Draft 01 — X + LinkedIn (Week 1, Tue)
- **Goal:** Introduce core problem/solution with a clear CTA.
- **Primary CTA:** `https://hypemd.dev`
- **Secondary CTA:** `https://github.com/gopherguides/hype`

**X version (<=280 chars)**
Most docs/blog stacks break at the seams: one pipeline for docs, another for content, glue scripts everywhere.

Hype keeps it markdown-native in one workflow:
• docs + blog from the same repo
• pre-deploy validation
• cleaner SEO/social metadata

Start: https://hypemd.dev

**LinkedIn version**
If your team maintains separate docs and blog pipelines, you’re paying a hidden tax in review overhead and deployment risk.

Hype is built around one markdown-native workflow so teams can:
- publish docs + blog from the same repo
- validate links/content before deployment
- ship cleaner metadata for SEO and social cards

If you already run on GitHub + CI, this fits naturally.

Start: https://hypemd.dev
Repo: https://github.com/gopherguides/hype

### Draft 02 — X + LinkedIn (Week 1, Thu)
- **Goal:** Position reliability + low-ops value for small teams.
- **Primary CTA:** Ask an engagement question to drive replies.

**X version (<=280 chars)**
Small docs teams don’t need more tooling—they need fewer failure points.

Hype’s direction is intentionally boring:
• markdown-native authoring
• deterministic build/export paths
• lightweight deployment model

What would you automate first: validation, metadata, or deploy checks?

**LinkedIn version**
Most documentation incidents come from pipeline complexity, not writing quality.

Hype aims to remove operational drag with:
- markdown-native authoring
- deterministic export/build behavior
- lightweight deployment patterns

The outcome for small teams: fewer runtime surprises and faster iteration.

What would you automate first in your docs workflow: validation, metadata, or deploy checks?

## Asset Dependencies
- [ ] Attach one screenshot/GIF per draft (quickstart + generated output).
- [ ] Add UTM tracking links for X and LinkedIn variants.
- [ ] Draft first matching technical blog post for W1 proof link.
