# Hype Domain Shortlist Validation (2026-03-20)

## Objective
Select a primary marketing domain and two defensive registrations for the Hype project.

## Decision Matrix
Scoring: 1 (poor) to 5 (excellent). Weighted total = score × weight.

| Criteria | Weight | hype.sh | hypemd.dev | hype.run | hypecli.dev |
|---|---:|---:|---:|---:|---:|
| Brand clarity (Hype markdown tooling) | 5 | 3 | 5 | 3 | 4 |
| Memorability | 4 | 5 | 4 | 4 | 3 |
| Developer trust / credibility | 4 | 3 | 5 | 4 | 4 |
| SEO keyword relevance | 3 | 2 | 5 | 2 | 4 |
| Expansion flexibility | 2 | 4 | 4 | 4 | 3 |
| Typo/confusion risk | 3 | 3 | 4 | 3 | 4 |
| **Weighted total** |  | **53** | **85** | **58** | **67** |

## Validation Notes

### DNS/Resolution Snapshot (2026-03-21 UTC)

| Domain | Resolves now? | Notes |
|---|---|---|
| `hypemd.dev` | Yes (`129.212.149.192`) | Confirms active domain in use. |
| `hypecli.dev` | No A/CNAME observed | Candidate appears unconfigured; verify registrar availability. |
| `hype.run` | Yes (`76.223.54.146`, `13.248.169.48`) | Likely already registered by third party; may require alternate defensive pick. |
| `hype.sh` | Yes (`217.92.164.52`) | Registered/active by third party; not practical as primary. |

### `hypemd.dev` (Recommended primary)
- Already live and recognized in architecture notes.
- Highest semantic fit: includes "md" and naturally signals markdown/doc tooling.
- `.dev` adds technical trust and enforces HTTPS by default in modern browsers.

### `hypecli.dev` (Recommended defensive)
- Useful redirect target for CLI-first audience/search intent ("hype cli").
- Lower brand elegance than `hypemd.dev`, but clear and practical.
- Next step: registrar availability + purchase window.

### `hype.run`
- Good campaign phrase, but current DNS indicates it is already in use.
- Treat as aspirational; do not block launch plan on this domain.

### `hype.sh`
- Short and memorable, but too generic and currently occupied.
- Better as a convenience redirect only if ever acquirable.

## Recommendation
1. Keep `hypemd.dev` as canonical primary domain.
2. Prioritize acquisition attempt for `hypecli.dev`.
3. Replace `hype.run` defensive slot with a likely-available fallback candidate (`hypemd.run` or `usehype.dev`) unless acquisition path opens.
4. Route all alternates to canonical URLs with 301 redirects to consolidate SEO authority.

## Implementation Checklist
- [ ] Verify registrar availability + pricing for `hypecli.dev`, `hypemd.run`, and `usehype.dev`.
- [ ] Register first available fallback set and enable registrar lock.
- [ ] Add DNS + TLS config in Dokploy/edge layer.
- [ ] Configure 301 redirects to canonical `https://hypemd.dev` routes.
- [ ] Add canonical tags and update sitemap host references.

## Issue Update Snippet (for #83)
- Decision matrix refreshed with concrete DNS validation.
- Primary remains `hypemd.dev`.
- `hype.run` and `hype.sh` appear occupied; defensive domain recommendation updated to `hypecli.dev` + one fallback (`hypemd.run`/`usehype.dev`) pending registrar check.
