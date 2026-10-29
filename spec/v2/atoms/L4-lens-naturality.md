# L4 · lens.naturality

**Purpose:** Compare two views of the same domain and find exactly where they can't be reconciled.
Examples: buyer vs vendor, regulator vs practitioner, source family A vs B, then vs now.

**Internal construct:** a natural transformation between two functors (views). A square that fails
to commute is a tension.

**In:** the merged `DomainGraph` and `Brief.views[2]`. Each view is a mapping from objects and
relations to that side's terms and claims.
**Out:** `Finding[kind=tension]` with `sides[2]`, and `commute_rate`

**Contract**
- **Translation table.** The model writes a translation table from view A to view B, one row per
  object.
- **Square check.** For each relation X→Y in the shared domain, check that "translate, then relate
  in B" equals "relate in A, then translate." This is a J1 `judge` task with k samples.
- **Tensions.** Each failing square becomes a tension with `sides[0]` and `sides[1]`, each holding at
  least 1 claim.
- **Limit, stated honestly.** v2.0 checks the squares only. It does not check that each view respects how relations compose (the functor condition).
- **Comparing two entities.** When the brief compares two entities rather than two views, run L1 on
  each one, then L4 over the probes they share.
- **Cost cap.** At most `lens.naturality.max_squares` squares per run, highest-degree objects first.

**Check:** G2 confirms every tension is sourced on both sides. `commute_rate` is logged.

**Guards against:** false balance, and comparisons that list features but never reach the real
conflict.

**Evidence:** E3, E11. This construct has the least published validation, so its individual
ablation is mandatory ([eval](../05-eval.md)).
