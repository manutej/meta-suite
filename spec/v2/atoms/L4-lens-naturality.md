# L4 · lens.naturality

**Purpose:** Compare two views of the same domain and find exactly where they cannot be reconciled.
Examples: buyer vs vendor, regulator vs practitioner, source family A vs B, last year vs this year.

**Internal construct:** a natural transformation between two functors (views). A naturality square
that fails to commute is a tension.

**In:** `DomainGraph`, view A, view B (each a mapping of objects and relations to that view's terms
and claims)
**Out:** `Finding[kind=tension]`, `commute_rate`

**Contract**
- The model writes a translation table from view A to view B, one row per object.
- For each relation X→Y in the shared domain, check whether "translate, then relate in B" equals
  "relate in A, then translate." Each comparison is sampled k times; the kernel aggregates.
- Each failing square becomes a tension finding with at least one source **per side**.
- `compare` questions about two entities (rather than two views) run L1 on each side plus L4 over
  the shared probes.

**Check:** every tension is sourced on both sides (G2). `commute_rate` is logged.

**Guards against:** false balance, and comparisons that list features without finding where the
real conflict is.

**Evidence:** E3, E11 (analysis gap). This construct has the least published validation of the
four, so its ablation is mandatory in eval.
