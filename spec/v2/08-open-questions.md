# 08 · Open questions

Each question has an owner and a milestone by which it must be resolved. A resolved question becomes
an entry in the [decision log](09-decisions.md).

| # | Question | Owner | By |
|---|---|---|---|
| Q1 | Can a skill-folder binary run, with its exec bit set, in the Claude app's uploaded-skill sandbox, in Claude Code and in Cowork? What size limits apply to skill zips that carry binaries for 5 platforms? | Claude (spike) | M0 |
| Q2 | **Resolved by D-036:** the skill is named `meta-lens`. | — | — |
| Q3 | Source for the Grant & Wang gluing construction behind `meta-score`. It is needed before that skill is scheduled. | Manu | before v2.x |
| Q4 | The claimed "up to 95% token reduction" from progressive disclosure and memory banks. It will be measured as M6 against v1, not assumed. | eval | M6 |
| Q5 | Which judge family to use for cross-family M1/M2 grading and rubric authoring in the eval harness. | Claude | M0 |
| Q6 | Recruit at least 2 SMEs per domain on the sealed list, the held-out candidates included. | Manu | M0 |
| Q7 | Rubric authoring protocol. Proposed: a cross-family agent drafts from independent sources; SMEs validate a sample for the named domains and every item for the held-out domain. | Claude + SMEs | M0 |
| Q8 | **Resolved by D-020:** the v1 arm is `meta-research`. | — | — |
| Q9 | Create the v1.2.0 git tag. This session's proxy blocks tag creation, so run it manually: `git tag -a v1.2.0 8cdbf73 -m "v1 baseline" && git push origin v1.2.0`. | Manu | M0 |
| Q10 | **Cloud sandbox resolved (D-051).** Other surfaces: `meta-probe` step 6. Is a headless browser available on each surface for G7's render predicates? | Claude (spike) | M0 |
| Q12 | **Resolved by D-047: publish the subset under MIT in meta-suite.** Original question, kept for the record: sheaf-core is private and proprietary (its ADR-010), while meta-suite is MIT and public. Porting the graph subset into `kernel/internal/sheaf` publishes that code. The subset is standard material (coboundary, Laplacian, Hodge split, H⁰/H¹ by counting), but publishing it is irreversible. Options: publish under MIT, or keep the port in a private module linked at build time. | Manu | M0 (before any port is committed) |
| Q11 | **Resolved by D-051.** Eval harness: which headless agent harness runs every arm with identical tools and API usage logs? Can it drive an open-weight model? If not, drop that model from Phase C. | Claude (spike) | M0 |
