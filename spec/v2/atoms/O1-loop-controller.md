# O1 · loop.controller

**Purpose:** The deterministic state machine that runs the pipeline. It decides what happens next,
so the model doesn't have to.

**In:** run state, gate verdicts, edit scripts, the tier's budgets
**Out:** the next stage to run (`meta loop next`), or `done`, or `ask_user`

**Contract**
- **Stage order** is defined in the [architecture](../03-architecture.md#pipeline-one-run).
- **After a gate fails or an edit script arrives,** route back to the earliest stage touched, never
  to the start.
- **Hard caps** (in config, per tier): review passes, re-gathers, re-expands, total tokens.
- **Early exit:** if every gate passes and the edit script is empty, the run is done.
- **Ask policy (P9):** return `ask_user` only when an assumption is flagged high-stakes **and**
  irreversible. In an unattended run, record the assumption and continue instead.
- **Refactor between stages only (P10).** Promoted playbook entries are loaded at stage boundaries.
- **When a cap is hit:** ship the best artifact so far, with the open issues listed in "how this was
  checked."

**Check:** recorded state transitions are replayed in tests.

**Guards against:** runaway loops, restarting from zero, and the model deciding when it's "done."

**Evidence:** E1 (proposer/validator/executor split), E14 (measured search loop).
