# ADR 0001: Website Architecture for Hype Marketing + Docs

- **Status:** Draft
- **Date:** 2026-03-20
- **Owners:** Gopher Guides / Hype maintainers

## Context
Hype already has a live public site at `hypemd.dev` deployed via Dokploy. We need an architecture direction that supports:
- marketing pages and docs content
- fast iteration with low ops overhead
- straightforward CI/CD from GitHub
- SEO, social metadata, and reliability

## Decision
Adopt a **Git-based static site workflow with Dokploy deployment** as the primary architecture, backed by:
1. Source-controlled content + templates in the repo
2. Build artifact generation in CI
3. Dokploy-managed deploys to production
4. Canonical domain routing through `hypemd.dev`

## Rationale
- Existing live deployment already validates Dokploy as an execution path.
- Static output minimizes runtime complexity and operational risk.
- GitHub-driven flow aligns with contributor expectations and review process.
- Supports incremental improvement without framework migration churn.

## Information architecture (v1)
- Home
- Docs
- Tutorials
- Blog
- Templates/Examples
- Changelog/Release notes

## Deployment topology
- **Origin:** GitHub repository
- **Build:** CI workflow on push/merge to main
- **Deploy target:** Dokploy-managed service(s)
- **Domain:** `hypemd.dev` as canonical host
- **TLS:** managed at edge/reverse-proxy layer (Dokploy stack)

## CI/CD outline
1. PR opened -> lint/test/build checks
2. Merge to main -> production build
3. Dokploy deploy hook/image update
4. Post-deploy smoke check (`/`, `/docs`, key assets)

## Consequences
### Positive
- Low operational burden, fast publish cycle
- Easy rollback via Git commit history/deploy history
- Works with existing team workflows and infra

### Negative / trade-offs
- Dynamic personalization/search depth may be limited without additional services
- Requires discipline in content structure and release process

## Follow-up actions
- Document smoke test script for post-deploy verification
- Add architecture diagram to docs site operations page
- Define SLO for docs uptime + publish latency
