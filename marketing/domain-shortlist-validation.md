# Domain Shortlist Validation (2026-03-20)

## Current status
- **Registered + live primary:** `hypemd.dev`
- **Current website:** https://hypemd.dev
- **Decision baseline:** keep one canonical brand domain and avoid defensive registrations unless abuse/confusion appears.

## Validation criteria + weights
| Criterion | Weight | Why it matters |
|---|---:|---|
| Brand fit | 30% | Must map naturally to “Hype + Markdown docs.” |
| Memorability/clarity | 25% | Easy recall and low typing friction improve direct traffic. |
| Collision/confusion risk | 20% | Reduce support burden and typo-domain ambiguity. |
| Strategic expandability | 15% | Useful if product scope expands beyond docs pages. |
| Cost/overhead | 10% | Keep renewals/admin low until justified. |

## Weighted decision matrix
Scoring scale: **1 (weak) → 5 (strong)**. Collision risk is inverse-value (5 = low risk).

| Domain | Brand fit (30) | Clarity (25) | Collision (20) | Strategic (15) | Cost (10) | Weighted score (/5) | Recommendation |
|---|---:|---:|---:|---:|---:|---:|---|
| `hypemd.dev` | 5 | 5 | 4 | 5 | 5 | **4.80** | **Primary (keep)** |
| `hypedocs.dev` | 4 | 4 | 4 | 3 | 3 | 3.80 | Monitor only |
| `hypecli.dev` | 4 | 4 | 4 | 3 | 3 | 3.80 | Monitor only |
| `tryhype.dev` | 3 | 4 | 3 | 2 | 3 | 3.05 | Skip |
| `usehype.dev` | 3 | 4 | 3 | 2 | 3 | 3.05 | Skip |
| `hypemarkdown.dev` | 4 | 2 | 5 | 2 | 3 | 3.35 | Skip (too long) |
| `hype-docs.dev` | 2 | 2 | 4 | 1 | 3 | 2.35 | Skip (hyphen penalty) |

## Practical validation notes
- Canonical domain already resolves publicly and is live in production (`hypemd.dev`), which removes launch-risk from rebranding decisions.
- Secondary names are not required for near-term GTM execution; no immediate abuse signal justifies extra renewals.
- Best defensive candidates (if needed later): `hypedocs.dev` and `hypecli.dev`.

## Trigger-based defensive registration policy
Register defensive domains only if one of these occurs:
1. Brand impersonation/phishing report tied to Hype naming.
2. Material traffic leakage from typo/confusion in support/community channels.
3. Planned paid campaigns where typo capture materially reduces CAC.

## Decision
1. Keep `hypemd.dev` as canonical domain.
2. No immediate defensive registrations.
3. Re-evaluate quarterly (or immediately on abuse signal).

## Registrar + conflict spot-check log
| Timestamp (UTC) | Check | Result | Notes |
|---|---|---|---|
| 2026-03-20 22:22 | Public DNS resolution for `hypemd.dev` | Pass | Domain resolves and serves current site. |
| 2026-03-20 22:22 | GitHub/org collision scan for `hypedocs`/`hypecli` naming | No blocking collision found | Continue monitoring if paid campaigns begin. |
| 2026-03-20 22:22 | Trademark quick screen (`Hype` + docs tooling context) | No obvious immediate blocker | Full legal review only if broad commercial expansion starts. |

## Decision matrix issue update payload
Use this summary in issue/PR updates:
- Canonical domain remains `hypemd.dev` (weighted score 4.80/5).
- Defensive domains deferred by policy unless abuse/confusion signals appear.
- Next review checkpoint: end of Q2 2026 or sooner on incident trigger.

### Suggested issue mapping
- Primary update target: https://github.com/gopherguides/hype/issues/83
- Cross-reference from website architecture thread when domain implications arise: https://github.com/gopherguides/hype/issues/84

## Single-comment issue update template (de-duplicated)
Use this exact block for the next issue update to avoid correction-spam:

```md
Domain workstream update:
- Canonical domain remains `hypemd.dev` (weighted score: 4.80/5)
- Defensive registrations remain deferred by trigger policy (abuse/confusion only)
- Next checkpoint: end of Q2 2026 (or immediately on impersonation/traffic leakage)

Artifact:
- `marketing/domain-shortlist-validation.md`
```

## Anti-spam update rule
- Post **at most one** domain update comment per 24h unless there is a material decision change.
- If formatting breaks in a comment, edit locally first and only repost once with a "supersedes" note.
