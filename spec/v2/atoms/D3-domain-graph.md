# D3 · domain.graph

**Purpose:** Assemble the claims into one typed `DomainGraph`. This graph is the hypothesis the lens
tests.

**Internal construct:** a category presentation. Objects are entities, morphisms are typed relations,
and each morphism is backed by claims.

**In:** `Claim[]`, `relation_vocab`
**Out:** `DomainGraph`

**Contract**
- Two objects merge only when their relation profiles match above a threshold set in config. This
  resolves entities by their relations, not by name similarity alone.
- Every morphism references at least one claim. Any proposed vocabulary extensions are accepted or
  mapped onto the existing vocabulary, and the decision is logged.
- The graph's shape metrics are computed here for A1:
  - hub concentration
  - longest chain
  - two-sidedness
  - temporal density

**Check:** G1, G2 (no orphan objects, no unsourced morphisms).

**Guards against:** duplicate entities splitting the evidence, and relations that exist only in
the model's head.

**Evidence:** E3 (structured evidence before synthesis), E10.
