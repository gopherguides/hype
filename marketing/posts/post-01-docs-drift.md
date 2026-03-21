# Post Draft 01 — Docs Drift

**Channel:** X + LinkedIn  
**Theme:** Problem framing  
**Goal:** Awareness + clickthrough to docs

## Final copy (X)
Your docs were correct when you wrote them.
They break quietly when APIs change.

That’s docs drift.

Hype executes + validates Markdown code samples during generation, so stale examples fail in CI before users hit them.

If docs are part of your product, they should be testable.

Start: https://hypemd.dev?utm_source=x&utm_medium=social&utm_campaign=q2_launch
Repo: https://github.com/gopherguides/hype

## Final copy (LinkedIn)
Most documentation failures are silent until users find them.

Hype helps teams catch docs drift earlier by executing and validating code examples directly from Markdown during doc generation.

That means stale snippets fail in CI instead of failing in front of customers.

If documentation is product surface area, test it like product code.

Start: https://hypemd.dev?utm_source=linkedin&utm_medium=social&utm_campaign=q2_launch
Repo: https://github.com/gopherguides/hype

## First comment/reply template
If you want to test this quickly, run `hype export` in CI on one docs folder and fail on broken examples.

## Media notes
- Asset file: `marketing/assets/post-01-docs-drift.png`
- Suggested asset: side-by-side screenshot (stale docs output vs CI failure from executable docs check).
- Alt text: “Comparison of a stale code example in static docs versus a failing CI check using executable Markdown validation.”

## Publishing metadata
- Scheduled slot (CT): 2026-03-23 09:00
- Owner: Hype maintainer
- Status: copy ready, awaiting image export
