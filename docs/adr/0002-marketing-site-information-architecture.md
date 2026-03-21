# ADR 0002: Marketing Site Information Architecture and Conversion Path

- **Status:** Draft
- **Date:** 2026-03-20
- **Owners:** Gopher Guides / Hype maintainers

## Context
We need a website structure that supports both discoverability (SEO/social) and conversion (trial/usage) without introducing a heavy application runtime.

Current state:
- Canonical domain: `hypemd.dev`
- Existing architecture direction: static-site workflow + Dokploy deployment (ADR 0001)

Gap:
- No explicit decision on page hierarchy, conversion funnel, and content ownership boundaries.

## Decision
Adopt a **docs-adjacent marketing IA** with explicit conversion path and reusable page templates.

### v1 page architecture
1. **Home (`/`)**
   - Value proposition
   - Primary CTA: Quickstart
   - Secondary CTA: GitHub repo
2. **Use Cases (`/use-cases`)**
   - Team scenarios: docs teams, dev advocates, OSS maintainers
3. **Docs (`/docs`)**
   - Reference + guides
4. **Blog (`/blog`)**
   - Product updates, tutorials, and proof posts
5. **Examples (`/examples`)**
   - Real snippets/templates and before/after artifacts
6. **Changelog (`/changelog`)**
   - Release updates mapped to capability changes

### Conversion path
Social post / search result -> Home or Use Case page -> Quickstart -> Repository stars/issues/discussions

## Options considered
| Option | Summary | Pros | Cons | Decision |
|---|---|---|---|---|
| A. Docs-adjacent marketing IA (selected) | Marketing pages and docs in one coherent structure | Shared templates, easier navigation, stronger SEO continuity | Requires editorial discipline | **Chosen** |
| B. Separate marketing microsite | Dedicated marketing site apart from docs | Design freedom for campaigns | Split ownership, duplicate content risk, SEO fragmentation | Rejected |
| C. Docs-only site with blog | Skip explicit marketing pages | Lowest effort initially | Weak narrative for first-time visitors; poorer conversion | Rejected |

## Rationale
- Keeps maintenance low while improving message clarity.
- Allows iterative content work without architecture churn.
- Aligns with a small-team operating model and repo-first workflows.

## Consequences
### Positive
- Clear discoverability and conversion path.
- Better editorial planning: each page type has explicit purpose.
- Reusable template system supports consistent branding.

### Trade-offs
- Requires ongoing content governance to avoid stale pages.
- Initial setup effort for page scaffolding and templates.

## Follow-up Actions
- [ ] Add route/page stubs for `/use-cases`, `/examples`, and `/changelog`.
- [ ] Define per-page metadata checklist (title, description, OG image, canonical).
- [ ] Add measurement plan (click-through to quickstart, docs depth, returning visitors).

## Measurement Plan (v1)
- **North-star conversion:** quickstart CTA click-through rate from home/use-cases pages.
- **Activation proxy:** unique visitors reaching `/docs/quickstart` within same session.
- **Retention proxy:** 7-day returning visitor ratio to docs/blog pages.
- **Channel attribution:** UTM-tagged social links for post-level performance.

## Owner Mapping
- Marketing docs and calendars: `docs/marketing/*`
- Architecture decisions: `docs/adr/*`
- Brand assets: `assets/brand/*`
- Product documentation source of truth: `docs/*`
