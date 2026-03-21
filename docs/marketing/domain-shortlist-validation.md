# Hype Domain Shortlist Validation (2026-03-20)

## Objective
Select a primary marketing domain and two defensive registrations for the Hype project.

## Decision Matrix
Scoring: 1 (poor) to 5 (excellent). Weighted total = score × weight.

| Criteria | Weight | hype.sh | hypemd.dev | hype.run | hypecli.dev | hypemd.run | usehype.dev |
|---|---:|---:|---:|---:|---:|---:|---:|
| Brand clarity (Hype markdown tooling) | 5 | 3 | 5 | 3 | 4 | 5 | 4 |
| Memorability | 4 | 5 | 4 | 4 | 3 | 3 | 4 |
| Developer trust / credibility | 4 | 3 | 5 | 4 | 4 | 4 | 4 |
| SEO keyword relevance | 3 | 2 | 5 | 2 | 4 | 4 | 3 |
| Expansion flexibility | 2 | 4 | 4 | 4 | 3 | 4 | 4 |
| Typo/confusion risk | 3 | 3 | 4 | 3 | 4 | 4 | 3 |
| **Weighted total** |  | **53** | **85** | **58** | **67** | **73** | **69** |

## Validation Notes

### DNS/Resolution Snapshot (2026-03-21 UTC)

| Domain | Resolves now? | Notes |
|---|---|---|
| `hypemd.dev` | Yes (`129.212.149.192`) | Confirms active domain in use. |
| `hypecli.dev` | No A/CNAME observed | Candidate appears unconfigured; verify registrar availability. |
| `hype.run` | Yes (`76.223.54.146`, `13.248.169.48`) | Already registered by third party; not suitable for near-term launch plan. |
| `hype.sh` | Yes (`217.92.164.52`) | Registered/active by third party; not practical as primary. |
| `hypemd.run` | No A/CNAME observed | Strong fallback candidate for redirects/campaign links. |
| `usehype.dev` | No A/CNAME observed | Good fallback if `hypecli.dev` unavailable. |

## Decision (current)
1. **Primary canonical:** `hypemd.dev`
2. **Defensive target #1:** `hypecli.dev`
3. **Defensive target #2 fallback order:** `hypemd.run`, then `usehype.dev`

## Recommendation
1. Keep `hypemd.dev` as canonical primary domain.
2. Prioritize acquisition attempt for `hypecli.dev`.
3. Replace `hype.run` defensive slot with likely-available fallback (`hypemd.run` then `usehype.dev`).
4. Route all alternates to canonical URLs with 301 redirects to consolidate SEO authority.

## Implementation Checklist
- [ ] Verify registrar availability + pricing for `hypecli.dev`, `hypemd.run`, and `usehype.dev`.
- [ ] Register first available fallback set and enable registrar lock.
- [ ] Add DNS + TLS config in Dokploy/edge layer.
- [ ] Configure 301 redirects to canonical `https://hypemd.dev` routes.
- [ ] Add canonical tags and update sitemap host references.
- [ ] Add domain ownership + renewal reminder to ops checklist.

## Issue Update Snippet (for #83)
- Decision matrix expanded with fallback candidates and weighted scoring.
- Primary remains `hypemd.dev`; defensive priority is now `hypecli.dev` then `hypemd.run`.
- `hype.run` and `hype.sh` remain occupied and are removed from launch-critical path.
