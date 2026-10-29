# 08 · Open questions

Each question has an owner and the milestone by which it must be resolved. A resolved question
becomes a [decision](09-decisions.md).

| # | Question | Owner | By |
|---|---|---|---|
| Q1 | Does a Go binary execute from an uploaded skill folder in the Claude app sandbox, in Claude Code and in Cowork? What size limits apply to skill zips holding binaries for 3 platforms? | Claude (spike) | M0 |
| Q2 | Final name for the pipeline skill. `meta-lens` is a working name. | Manu | M0 |
| Q3 | Source for the Grant & Wang gluing construction behind `meta-score`. A citation is needed before that skill is scheduled. | Manu | before v2.x |
| Q4 | The "up to 95% token reduction" from progressive disclosure and memory banks. It will be measured as M6 against the v1 baseline rather than assumed. | eval | M6 |
| Q5 | Which judge family is used for the cross-model blind review (R1) and the pairwise judge (M2), and is it available on every surface? | Claude | M1 |
| Q6 | Recruiting the SME panel (M3): at least 1 expert per domain. | Manu | M5 |
| Q7 | How the rubrics are authored without expert reference reports. The proposal: a separate agent drafts from independent sources and an SME validates a sample. | Claude + SME | M5 |
| Q8 | Which v1 path is the v1 arm? The proposal is `meta-mvp`, since it is v1's end-to-end pipeline. | Manu | M0 |
| Q9 | The v1.2.0 git tag. This session's proxy blocks creating tags, so it needs to be run manually: `git tag -a v1.2.0 8cdbf73 -m "v1 baseline" && git push origin v1.2.0`. | Manu | M0 |
| Q10 | Is a headless browser available per surface for G7's render-dependent predicates? | Claude (spike) | M0 |
