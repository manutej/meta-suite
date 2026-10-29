# R1 · review.blind

**Purpose:** An adversarial review by a critic that never saw the work being made.

**Internal construct:** —

**In:** the artifact, `Spec`, findings and gate reports. **No** drafting trace, chat or intermediate
reasoning.
**Out:** `EditScript` only. The critic cannot edit.

**Contract**
- **Isolation.** Fresh subagent context, using the reviewer model from the role map
  (`roles.reviewer`). On a surface without subagents the review is recorded `blind=false`.
- **Product default.** The reviewer is the host's best model in a fresh context. The eval measures
  this default configuration, not a cross-family reviewer that users would not have.
- **Edits.** Each edit names `target_stage`, `target_id`, `op`, `rationale` and a priority from P1
  to P3.
- **The reviewer looks at the rendered artifact**, not only its text: a screenshot or page render where available. The first lens is: would an SME forward this, and does the reader learn something in the first 60 seconds?
- **Rubric scoring (D-069).** R1 scores every applicable [rubric](../12-artifact-rubric.md) dimension on the render. Any L1 is a P1 edit. An L2 with a route to L3 is a P3 edit, executed only in deep mode's L3 pass.
- **Lenses** (catalog in config: `review.lenses`):
  - faithfulness to the spec
  - unsupported claims
  - expert concerns that are missing
  - generic content that could fit any domain
  - findings stated with more certainty than their evidence supports
- **Stop rule.** The loop stops when there are no P1 or P2 edits, or when the pass cap
  `modes.<mode>.review_passes` is reached.

**Check:** G1 on the edit script. The `blind` attestation is recorded in S1.

**Guards against:** self-grading, and mistaking polish for depth.

**Evidence:** E20, E21.
