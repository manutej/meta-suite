# R1 · review.blind

**Purpose:** An adversarial review by a critic that never saw the work being produced.

**In:** the artifact, the `Spec`, findings, and gate reports. The critic gets **no** drafting trace,
no chat history and no intermediate reasoning.
**Out:** an `EditScript` only. The critic cannot edit anything.

**Contract**
- The reviewer runs in a fresh context, on a different model family where one is configured. If it
  falls back to the same family, that fallback is logged.
- Each edit names `target_stage`, `target_id`, `op` and `rationale`, and is graded P1, P2 or P3.
- Lenses the reviewer applies:
  - faithfulness to the spec
  - unsupported claims
  - missing obvious expert concerns
  - generic content that could fit any domain
  - findings stated with more certainty than their evidence supports
- O1 applies the edits by routing back to the earliest stage they touch.
- The review loop stops when there are no P1 or P2 edits, or when the tier's pass cap is reached.

**Check:** the edit script is schema-valid (G1), and the `blind: true` attestation is recorded in S1.

**Guards against:** self-grading bias, and polish being mistaken for depth.

**Evidence:** E20, E21.
