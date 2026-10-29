# O1 · loop.controller

**Purpose:** A deterministic state machine that runs the pipeline. It decides what runs next, so
the model doesn't.

**Internal construct:** —

**In:** run state, gate verdicts, edit scripts, mode budgets, usage
**Out:** the next stage (from `meta loop next`), `done`, or `ask_user`. It also issues stage tokens
(`meta stage begin`).

**Contract**
- **Stage order** follows the [architecture](../03-architecture.md#pipeline-one-run).
- **Routing back.** A failed gate or an edit script routes the run to the earliest stage it touches.
  - D2 can be re-entered at most `loop.d2_reentry_max` times.
  - Re-expansion is capped by `intake.re_expand_max`.
  - Review passes are capped by `modes.<mode>.review_passes`.
- **Token cap.** `modes.<mode>.token_cap` is enforced only when `usage.source` is `api` or
  `harness`. Otherwise the run is capped by a call budget.
- **The single ask rule (P9).** Return `ask_user` only when an assumption is flagged **high-stakes
  and either expensive to redo or irreversible**. In an unattended run, record the assumption and
  continue. I1 and every other atom defer to this rule.
- **Stage tokens.** Each stage gets a token that stamps run, stage, region, agent, model and prompt
  hash onto everything written during that stage.
- **Stage-boundary adaptation (D-057, D-070).** On by default. Before each stage starts, O1 (deterministic) decides whether that stage's op prompt is rewritten, from the edit scripts and gate results so far, and records why. The new text is produced by a J1 `answer` task on the host model; the kernel hashes it, logs the version and the diff, and never applies it mid-stage. Only the `v2-frozen` ablation arm turns it off.
- **Human-input log (M9).** Every `ask_user`, user interruption and retry is counted in S1.
- **Early exit.** When every gate passes and the edit script is empty, the run is done.
- **Hitting a cap.** Ship the best artifact so far, with the open issues listed under "how this was
  checked."

**Check:** replay tests on state transitions recorded from real runs.

**Guards against:** runaway loops, restarting from zero, the model deciding when it is done, and
forged provenance.

**Evidence:** E1, E14.
