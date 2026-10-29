# G3 · gate.sheaf

**Purpose:** Measure whether the local claims glue into one consistent picture across the cover,
and locate where they don't, with real sheaf computation wherever the claims have linear structure.

**Internal construct:** a conflict graph over the cover, plus three cellular-sheaf layers computed
by the kernel's embedded sheaf subset (D-038). The sheaf-core repo serves as the test oracle.

| Layer | What it covers | Kind | Cells | Stalks | Maps | Output |
|---|---|---|---|---|---|---|
| **1 · Numeric key** | Claims with numeric `value_type` | deterministic | 0-cells = regions asserting a (subject, relation) key; 1-cells = computed overlaps where both sides assert it; 2-cells = triple overlaps | ℝᵈ, values in the **canonical unit** (kernel-normalized) | identity, or a projection onto the shared coordinates | `h0` (glues or not), per-edge and per-region `discord`. Relative H⁰ + harmonic extension say whether an uncovered overlap is **determined** by its neighbours |
| **2 · Relational cycle** | Linear relative claims: ratios as logs, differences, time offsets | deterministic | vertices = entities; edges = relative claims | ℝ | identity | Hodge split: the harmonic share (against its null) is inconsistency that only shows around cycles of 3+ claims. It localizes to the **cycle**, not to a single edge |
| **3 · Judgment balance** | All claims on one key, using the J1 judgments already collected (no extra calls) | kernel-aggregated | vertices = claims; edges = judged pairs | ℝ | +1 `compatible`, −1 `conflict` (weight √vote margin); `unrelated` = no edge | `balance`. H⁰ = 1 means balanced, and the section gives the two camps → `sides[2]`. H⁰ = 0 means frustrated |

**Where the math stops.**
- Free-text qualitative claims have no linear structure. Embedding vectors are **not** used as
  stalks, because similarity is not consistency.
- Those claims stay on blocking + J1 pairwise judgments + layer 3.
- Layer 3 checks that the **judgments** cohere. It says nothing about whether they are true.

**In:** claims (with `value_type`, `unit`, `as_of`), the `Cover` (universe, computed overlaps,
triple overlaps), and conflict candidates from L2/L3
**Out:** the G3 report, `ConflictEdge[]`, and `Finding[kind=obstruction]`

**Report fields:**
- `h0`
- `discord`
- `harmonic_share` with `harmonic_null`
- `balance`
- `near_zero_ambiguous`
- `judged_pairs` / `candidate_pairs`

**Contract**
- **Pairwise candidates.**
  - Blocked by the (subject, relation) key, up to `tiers.<tier>.g3_pair_cap`.
  - Each pair is a J1 `judge` task with the verdict values `conflict`, `compatible` or `unrelated`.
    Vote counts are kept.
  - Coverage is reported as `judged_pairs / candidate_pairs`. When the cap bites, pairs are sampled
    stratified by key rather than ordered by centrality, so the reported rate is not biased.
- **Obstructions and sides.**
  - A balanced conflicting component yields one obstruction with `sides[2]` taken from the layer-3
    section.
  - A **frustrated** component yields an obstruction with no sides, tagged `contested`, with the note
    "the evidence does not split cleanly into two camps."
- **`near_zero_ambiguous`.** When an eigenvalue sits near the cut, the affected findings get
  confidence `uncertain` instead of a verdict.
- **Report only.** At stage 4 G3 writes reports and nothing else.
- **Coverage verdict at stage 5.** Every obstruction must map to a `contested` section in A2's
  `section_map`.
- **Locality.** For claim sets, the sheaf locality condition reduces to entity resolution, which is
  D3's job. G3 does not report it as a separate axiom.

**Check:** the stage-5 structural coverage check. Every layer's numeric output must match golden
reports from the sheaf-core oracle within tolerance in `meta check` (see [roadmap](../06-roadmap.md)).

**Guards against:**
- smoothing over real disagreement
- three-way inconsistencies that no pairwise check sees (layer 2)
- "two sides" framings the evidence doesn't support (layer 3)

**Evidence:** E6, E7, E8, E26 (sheaf-core oracle and the toy proof).
