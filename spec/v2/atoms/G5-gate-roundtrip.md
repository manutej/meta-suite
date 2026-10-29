# G5 · gate.roundtrip

**Purpose:** Turn intent drift into a number with a threshold.

**In:** `Brief`, `Brief'` (from I3)
**Out:** a field-level agreement table, a `drift` score, and a verdict

**Contract**
- **Exact fields:** `question_type` and `audience`.
- **Semantic fields:** `objective` and each constraint. Each is judged k times by an LLM, and the
  kernel aggregates.
- Weights and the threshold live in config.
- **If the gate fails:** return the per-field diff to I2. After 2 re-expands, proceed anyway and
  log the residual drift as an assumption shown in the artifact.

**Check:** `drift ≤ threshold`.

**Guards against:** a beautifully executed answer to the wrong question.

**Evidence:** E20.
