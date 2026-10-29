# D1 · domain.cover

**Purpose:** Decide how the domain is split up. The same split is reused for gathering, for
recombining answers, and for gluing checks.

**Internal construct:**
- a **decomposition tree** of sub-questions, whose composition is measured by G4
- a **cover** of a declared universe by regions, whose computed overlaps and triple overlaps give
  the cell complex G3 works on

**In:** `Spec`
**Out:** `Cover`, containing:
- `universe` (the sub-question nodes plus expected entity kinds)
- `tree`
- `regions[]`, each a subset of the universe
- `relation_vocab`

`overlaps[]` and `triple_overlaps[]` are **computed by the kernel**, never declared by the model.

**Contract**
- **Covering condition.** Every element of `universe` belongs to at least one region. The kernel
  rejects a cover that leaves anything out.
- **Overlaps are computed.** An overlap is a pair of regions whose intersection is non-empty, so
  "adjacent" simply means "has a non-empty overlap." The overlap graph must be connected; if it
  isn't, D1 must add a bridging element.
- **Triple overlaps** are computed for G3's 2-cells and L3's coherence check.
- **Region count** is at most `modes.<mode>.max_regions`.
- **Relation vocabulary.**
  - `relation_vocab` is controlled.
  - Each relation declares `value_type` ∈ {`qualitative`, `numeric`, `relative`}. Numeric and
    relative relations also declare a unit dimension. This is what lets G3's numeric layers run.
- **Insider probes (VR-03).** Each region carries questions aimed at what practitioners argue about, where the rules don't hold, how things are actually enforced, and which practitioner sources (not just top search results) to read. Probes are questions, never answers.
- **No pre-written answers.** No region or node is defined by naming the answer it expects.

**Check:** G1, plus the kernel's checks that the cover is covering and its overlap graph connected.

**Guards against:**
- siloed research that can't be cross-checked
- gaps in coverage that nobody notices
- sprawling relation vocabularies

**Evidence:** E4, E5, E9, E26.
