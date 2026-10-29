# D1 · domain.cover

**Purpose:** Decide how the domain is split up. The same split is reused for gathering, for
recombining answers, and for checking consistency.

**Internal construct:**
- An operadic decomposition tree: each sub-question has a typed answer slot, and a parent's answer
  is composed from its children's.
- An open cover of the domain by regions with overlaps, which is what makes gluing checkable.

**In:** `Spec`
**Out:** `Cover` (tree + regions) and an initial `relation_vocab[]`, the typed relation kinds for
this domain.

**Contract**
- Every pair of adjacent regions shares at least one object or sub-question. With no overlaps
  there is nothing to glue, so G3 cannot run.
- The region count is capped by the budget tier.
- `relation_vocab` is a controlled vocabulary. D2 must use it, or propose an extension, which is
  flagged and reviewed in D3.
- No region is defined by naming a specific answer in advance.

**Check:** G1. The kernel verifies that the overlap graph is connected.

**Guards against:** siloed research that can't be cross-checked, and relation sprawl (a hundred
synonyms for "depends on").

**Evidence:** E4, E5 (operadic decomposition), E9 (fixed schema at write time).
