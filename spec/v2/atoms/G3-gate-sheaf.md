# G3 · gate.sheaf

**Purpose:** Check whether the local claims form a single consistent global picture. Report exactly
where they do not, axiom by axiom, and name the smallest sets of claims that cannot all be true.

**Internal construct:** the sheaf axioms over the region cover:
- **locality:** sections that agree on every region are equal
- **gluing:** compatible local sections extend to a global one

**In:** claims, cover overlaps, and the conflict candidates from L2 and L3
**Out:** a per-axiom report and `minimal_conflict_sets[]`. Each conflict becomes
`Finding[kind=obstruction]`.

**Contract**
- Whether two natural-language claims conflict is a semantic judgment. The LLM judges each pair k
  times, and the kernel treats those judgments as data.
- The kernel computes:
  - a per-axiom pass/fail
  - minimal conflicting subsets (sets that are jointly inconsistent, but each of whose proper
    subsets is consistent)
- Obstructions are **findings**, not run failures. In the artifact they become "contested."
- The gate fails only if an obstruction is missing from the artifact, or is presented as settled.

**Check:** the per-axiom report is written to S1. Each axiom is reported separately; there is no
single blended score.

**Guards against:** smoothing over real disagreement, and a single "coherence score" that hides
which kind of failure happened.

**Evidence:** E6 (axioms fail separately), E7 (minimal conflicting sets), E8.

**Open:** cost of pairwise judgments on large claim sets. The plan is to block by shared key first.
