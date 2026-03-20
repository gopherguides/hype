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

### `hypemd.dev` (Recommended primary)
- Already live and recognized in previous architecture notes.
- Highest semantic fit: includes "md" and naturally signals markdown/doc tooling.
- `.dev` adds technical trust and enforces HTTPS by default in modern browsers.

### `hypecli.dev` (Recommended defensive)
- Useful redirect target for CLI-first audience/search intent ("hype cli").
- Lower brand elegance than `hypemd.dev`, but clear and practical.

### `hype.run` (Recommended defensive)
- Strong campaign/CTA option ("run hype").
- Useful for short links and launch experiments.

### `hype.sh`
- Short and memorable, but over-generic and potentially confusing without context.
- Better as a convenience redirect than primary brand domain.

## Recommendation
1. Keep `hypemd.dev` as canonical primary domain.
2. Acquire `hypecli.dev` and `hype.run` as defensive/marketing redirect domains.
3. Route all alternates to canonical URLs with 301 redirects to consolidate SEO authority.

## Implementation Checklist
- [ ] Verify availability and pricing for `hypecli.dev` and `hype.run`.
- [ ] Register domains and enable registrar lock.
- [ ] Add DNS + TLS config in Dokploy/edge layer.
- [ ] Configure 301 redirects to canonical `https://hypemd.dev` routes.
- [ ] Add canonical tags and update sitemap host references.
