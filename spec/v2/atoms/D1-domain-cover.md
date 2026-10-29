# D1 · domain.cover

**Purpose:** Decide how the domain is split up. The same split is reused for gathering, for
recombining answers and for consistency checks.

**Internal construct:**
- An operadic decomposition tree: each node is a sub-question with a typed answer slot, and a
  parent's answer is composed from its children's.
- An open cover of the domain by overlapping regions, which is what makes gluing checkable.

**In:** `Spec`
**Out:** `Cover`, containing the tree, regions, overlaps and `relation_vocab`

**Contract**
- **Overlaps.** Every pair of adjacent regions shares at least one tree node or expected object,
  declared in `overlaps[]`.
- **Region count.** At most `tiers.<tier>.max_regions`.
- **Relation vocabulary.** `relation_vocab` is controlled. D2 must use it, or propose an extension;
  D3 decides on extensions and logs the decision.
- **No pre-written answers.** No region or node is defined by naming the answer it expects.

**Check:** G1. The kernel verifies that the overlap graph is connected.

**Guards against:** siloed research that can't be cross-checked, and sprawling relation vocabularies.

**Evidence:** E4, E5, E9.
