# I3 · spec.roundtrip

**Purpose:** Measure the drift between what the user meant and what the spec now says.

**Internal construct:** the compressing half of the adjunction (forgetful). What the round trip loses
is the drift.

**In:** `Spec` only
**Out:** `Brief'`, the reconstructed intent, using the same type as `Brief`

**Contract**
- Runs in a fresh subagent context with no access to the raw query, the `Brief` or the drafting
  trace.
- On surfaces without subagents, the reconstruction runs in the main thread and is recorded
  `blind=false`.
- Reconstructs the question type(s), objective, audience and constraints.

**Check:** G5. If G5 fails, I2 re-expands with the `DriftReport`, at most `intake.re_expand_max`
times (default 2). After that the run proceeds, and the residual drift is shown in the artifact as
an assumption.

**Guards against:** answering the wrong question with great care.

**Evidence:** E20.
