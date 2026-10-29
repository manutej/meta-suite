# I3 · spec.roundtrip

**Purpose:** Measure drift between what the user meant and what the spec now says.

**Internal construct:** the compressing half of the adjunction (forgetful). The unit of the
adjunction is what the round trip loses.

**In:** `Spec` only. The agent must never see the raw query or `Brief`.
**Out:** `Brief'`, the reconstructed intent.

**Contract**
- Runs in a fresh context (a subagent) with no access to the drafting trace.
- Reconstructs the question type, objective, audience and constraints from the spec alone.
- G5 then compares `Brief` with `Brief'`.

**Check:** G5 drift ≤ threshold. If it fails, I2 re-expands with the drift report (at most 2
times), then the run proceeds with the residual drift logged as an assumption.

**Guards against:** intent drift at the most expensive point to get wrong, before any research
budget is spent.

**Evidence:** E20 (self-correction needs external signal; a blind reconstructor is that signal).
