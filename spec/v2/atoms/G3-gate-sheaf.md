# G3 · gate.sheaf

**Purpose:** Measure how well the local claims fit together where regions overlap, and name the
claims that conflict.

**Internal construct:** the sheaf conditions over the region cover, made operational for v2.0 as two
computable measures:
- **Gluing:** the overlap-conflict rate. It is the share of overlap keys that hold conflicting
  claims from different regions.
- **Coverage:** the uncovered-overlap count. It is the number of declared overlaps where one side
  has no claims at all.

**In:** claims, `Cover.overlaps`, and conflict candidates from L2 and L3
**Out:** the two measures, `ConflictEdge[]`, and conflict components. Each component becomes a
`Finding[kind=obstruction]`.

**Contract**
- **Blocking.** Candidate pairs are grouped by `(subject, relation)` key. Only pairs that share a key
  are judged, up to `tiers.<tier>.g3_pair_cap`, ordered by how central the subject is in the graph.
- **Each pair is a J1 `judge` task** with the verdict values `conflict`, `compatible` or
  `unrelated`.
- **Conflict components.** Connected components of the conflict edges.
  - For v2.0 these replace minimal conflicting sets, since computing true minimal sets over natural
    language costs exponentially many judgments.
  - Conflicts among three or more claims that no pair reveals are out of scope for v2.0.
- **Report only.** At stage 4, G3 writes its report and nothing else.
- **Coverage verdict at stage 5.** G3 passes only if every obstruction is mapped in A2's
  `section_map` to a section tagged `contested`.

**Check:** the stage-5 structural coverage check on `section_map`. Both measures are logged to S1.

**Guards against:** smoothing over real disagreement, and a single blended "coherence" score.

**Evidence:** E6 (report conditions separately), E7 (conflict sets), E8.
