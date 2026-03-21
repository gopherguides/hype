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
| Week (2026) | Product Education | Workflow Integrations | Proof & Outcomes | Community |
|---|---|---|---|---|
| W1 (Mar 23) | "What is Hype" explainer | GitHub Actions publish flow | Before/after docs quality comparison | Ask for examples from users |
| W2 (Mar 30) | Parser/CLI quick tips | Docs + blog in one repo | Build-time and error reporting walkthrough | Share template call for contributors |
| W3 (Apr 6) | Metadata + SEO support | Social card and OG automation | Traffic uplift case study format | Feature request spotlight |
| W4 (Apr 13) | Export/validation deep dive | Dokploy deployment guide | Reliability checklist post | Office-hours / AMA prompt |

## Two Concrete Drafts (PR-ready)

### Draft 01 — X + LinkedIn (Week 1, Tue Mar 24)
- **Goal:** Introduce core problem/solution with a clear CTA.
- **Primary CTA:** `https://hypemd.dev/?utm_source=social&utm_medium=x&utm_campaign=q2_w1_intro`
- **Secondary CTA:** `https://github.com/gopherguides/hype?utm_source=social&utm_medium=linkedin&utm_campaign=q2_w1_intro`
- **Asset:** quickstart screenshot (`assets/marketing/q2-w1-intro.png`)

**X version (<=280 chars)**
Docs and blogs shouldn’t require 2+ pipelines and glue scripts.

Hype keeps publishing markdown-native in one workflow:
• docs + blog from one repo
• pre-deploy validation
• cleaner SEO/social metadata

Start: https://hypemd.dev/?utm_source=social&utm_medium=x&utm_campaign=q2_w1_intro

**LinkedIn version**
Many teams pay a hidden tax by splitting docs and blog pipelines.

Hype keeps publishing markdown-native so teams can:
- ship docs + blog from one repo
- validate links/content before deployment
- publish cleaner metadata for SEO and social cards

Start: https://hypemd.dev/?utm_source=social&utm_medium=linkedin&utm_campaign=q2_w1_intro
Repo: https://github.com/gopherguides/hype?utm_source=social&utm_medium=linkedin&utm_campaign=q2_w1_intro

### Draft 02 — X + LinkedIn (Week 1, Thu Mar 26)
- **Goal:** Position reliability + low-ops value for small teams.
- **Primary CTA:** Drive replies via concrete question.
- **Asset:** validation output screenshot (`assets/marketing/q2-w1-reliability.png`)

**X version (<=280 chars)**
Small docs teams don’t need more tooling—they need fewer failure points.

Hype focuses on boring reliability:
• markdown-native authoring
• deterministic build/export paths
• low-ops deployment patterns

What would you automate first: validation, metadata, or deploy checks?

**LinkedIn version**
Most docs incidents come from pipeline complexity, not writing quality.

Hype reduces operational drag with:
- markdown-native authoring
- deterministic export/build behavior
- lightweight deployment patterns

For your docs workflow, what would you automate first: validation, metadata, or deploy checks?

## Publishing Ops Checklist
- [ ] Create `assets/marketing/q2-w1-intro.png` and `assets/marketing/q2-w1-reliability.png`.
- [ ] Add first-week posts to scheduler with local timezone (America/Chicago).
- [ ] Create matching W1 technical blog draft with same UTM campaign tags.
- [ ] Capture baseline metrics (profile visits, link clicks, repo stars) before publishing.
