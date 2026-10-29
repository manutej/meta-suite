# I2 · spec.expand

**Purpose:** Expand an unfinished query into the fullest typed `Spec` that stays faithful to the
user's intent: the methodology, slots, constraints and success criteria the rest of the run is
held to.

**Internal construct:** the expanding half of an adjunction (the free construction). I3 provides
the compressing half, and G5 measures what the round trip lost.

**In:** `Brief`
**Out:** `Spec` (content-hashed)

**Contract**
- Every slot and constraint carries `origin` ∈ {`stated`, `inferred`, `default`}.
- No domain facts go in the spec. Facts come only from D2 gather.
- Success criteria are checkable. Each names the finding kind or gate that will satisfy it.
- The spec is an XML contract: typed slots, rules and success criteria, with no worked examples.

**Check:** G1. G5 runs after I3.

**Guards against:** over-expansion (requirements the user never implied) and under-expansion
(a vague target that lets the run drift).

**Evidence:** E1 (typed slots plus a validator), E22 (right altitude: expand intent, not boilerplate).
