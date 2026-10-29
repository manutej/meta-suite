# G5 · gate.roundtrip

**Purpose:** Turn intent drift into a number with a threshold.

**Internal construct:** —

**In:** `Brief`, `Brief'`
**Out:** `DriftReport`

**Contract**
- **Exact-match field:** `question_types`.
- **Semantic fields:** `audience`, `objective` and each constraint. Each is a J1 `judge` task.
- **Configuration:** field weights and the threshold are set in config: `gates.g5.*`.
- **On failure:** the per-field diff goes back to I2, within the limit of `intake.re_expand_max`
  re-expands.

**Check:** `drift ≤ gates.g5.threshold`.

**Guards against:** a beautifully executed answer to the wrong question.

**Evidence:** E20.
