# I2 · spec.expand

**Purpose:** Expand an unfinished query into the fullest typed `Spec` that stays faithful to it. The
spec sets the methodology, slots, constraints and success criteria the rest of the run is held to.

**Internal construct:** the expanding half of an adjunction (the free construction). I3 is the
compressing half, and G5 measures what the round trip loses.

**In:** `Brief`, and on a re-expand, a `DriftReport`
**Out:** `Spec`, as JSON with a content hash. When shown to a model it is rendered as an XML contract.

**Contract**
- Every slot and constraint carries `origin` ∈ {`stated`, `inferred`, `default`}.
- No domain facts go in the spec. Facts come only from D2.
- Success criteria are checkable. Each names the finding kind or gate that will satisfy it.
- No worked examples of domain content.

**Check:** G1. G5 runs after I3.

**Guards against:** over-expansion (requirements the user never implied) and under-expansion (a
vague target that lets the run drift).

**Evidence:** E1, E22.
