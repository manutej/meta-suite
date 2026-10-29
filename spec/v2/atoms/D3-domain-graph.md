# D3 · domain.graph

**Purpose:** Turn each region's claims into a typed subgraph, and propose where regions describe the
same things. This produces the hypothesis the lens will test.

**Internal construct:** a category presentation per region. Objects are entities; morphisms are
typed relations, each backed by claims.

**In:** `Claim[]`, `relation_vocab`, `Cover.overlaps`
**Out:** a `RegionGraph` for each region, `Interface[]` candidates, and the vocabulary decisions

**Contract**
- **Entity resolution happens only within a region.**
  - Two objects in the same region merge when their relation profiles match at or above
    `graph.merge_threshold`.
  - Objects that look alike across regions are **not** merged here. They are proposed as interface
    candidates for L3.
- **Every morphism references at least one claim.** Vocabulary extensions proposed by D2 are either
  accepted or mapped onto an existing relation, and the decision is logged.
- **Shape metrics** are computed by `meta graph shape` on the merged graph after L3. They are hub
  concentration, longest chain, two-sidedness and temporal density.

**Check:** G1. G2: no orphan objects, no unsourced morphisms.

**Guards against:** duplicate entities that split the evidence, relations that exist only in the
model's head, and merges no one declared.

**Evidence:** E3, E10.
