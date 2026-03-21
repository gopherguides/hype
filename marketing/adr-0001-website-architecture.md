# ADR 0001: Website Architecture for Hype Docs + Marketing

- **Status:** Proposed
- **Date:** 2026-03-20
- **Owners:** Hype maintainers
- **Related issue:** https://github.com/gopherguides/hype/issues/84

## Context

Hype needs a website architecture that supports two jobs without splitting the user journey:

1. **Marketing narrative** (what Hype is, why it matters, social proof, CTAs)
2. **Technical docs** (quickstart, command reference, examples)

Current constraints:
- Team bandwidth is limited; architecture should be simple to operate.
- Content velocity matters more than custom frontend complexity.
- Website should reinforce Hype itself as the engine (dogfooding).
- We want low hosting cost, fast page loads, and straightforward deploys.

## Decision drivers

1. Keep build/deploy system maintainable by a small team.
2. Optimize for content publishing speed (not bespoke web app development).
3. Keep docs and marketing in one information architecture.
4. Preserve future optionality for richer docs/search.
5. Favor static hosting + CDN for reliability and cost.

## Considered options

### Option A — Hype-generated static site (single repo) on DigitalOcean Static Sites
- Use Hype as primary generator for docs + lightweight marketing pages.
- Deploy static artifacts to DigitalOcean-managed static hosting/CDN.

### Option B — Hugo/Astro marketing site + separate docs stack
- Split responsibilities: modern framework for marketing, separate docs pipeline.
- Potentially better design flexibility, but higher cognitive load.

### Option C — Fully custom web app (Next.js) with dynamic content pipeline
- Maximum flexibility and integrations.
- Highest ops complexity and maintenance burden.

## Decision matrix

Scoring scale: 1 (weak) to 5 (strong)

| Criterion | Weight | A: Hype + DO Static | B: Split stack | C: Custom app |
|---|---:|---:|---:|---:|
| Team maintainability | 30% | 5 | 3 | 2 |
| Publishing velocity | 25% | 5 | 3 | 2 |
| Cost + ops simplicity | 20% | 5 | 3 | 2 |
| UX/design flexibility | 15% | 3 | 4 | 5 |
| Long-term extensibility | 10% | 4 | 4 | 5 |
| **Weighted total (/5)** |  | **4.60** | **3.20** | **2.60** |

## Decision

Adopt **Option A**: Hype-generated unified static site deployed to DigitalOcean Static Sites.

### Rationale

- Fastest route to consistent docs + marketing publishing.
- Lowest operational burden for current team size.
- Strong product signal: Hype powers its own public content.
- Keeps future migration path open if richer app behavior becomes necessary.

## Consequences

### Positive
- Single content workflow and fewer moving parts.
- Lower hosting and maintenance overhead.
- Easy rollback model via static deploys.

### Trade-offs / risks
- Less out-of-the-box interactivity than framework-heavy sites.
- Search and analytics integrations may need incremental additions.
- Design system sophistication depends on internal template investment.

## Guardrails + follow-up tasks

1. Define site IA v1: Home, Docs, Examples, Blog, Changelog.
2. Add deploy checklist (build, link check, smoke test, publish).
3. Define baseline analytics events (CTA clicks, docs entry points, quickstart conversions).
4. Revisit architecture after 90 days or at >50k monthly sessions.

## Rollout plan (v1)

- Week 1: finalize IA + templates for core pages.
- Week 2: migrate key docs/quickstart pages to unified nav.
- Week 3: connect analytics + publish content cadence from marketing calendar.
- Week 4: review metrics and backlog improvements.

## v1 information architecture (concrete)
- `/` Home (value proposition + primary CTA to quickstart)
- `/docs` Documentation hub
- `/docs/quickstart` First successful run path
- `/examples` Executable examples gallery
- `/blog` Narrative + release posts
- `/changelog` Product updates and release notes

## Measurement baseline (must-implement)
| Event | Trigger | Why it matters |
|---|---|---|
| `cta_quickstart_click` | Click from `/` to `/docs/quickstart` | Measures marketing -> docs conversion |
| `quickstart_copy_command` | Copy action on first command block | Measures trial intent |
| `repo_click` | Click to GitHub repo | Measures OSS contribution funnel |
| `docs_search_used` | First search interaction (when enabled) | Signals docs discoverability gaps |
