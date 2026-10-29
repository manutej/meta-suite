# P2 · prompt.sibling_test

**Purpose:** Prove a meta-prompt serves its category by running it on sibling instances it has never seen.

**Internal construct:** —

**In:** `MetaPrompt`, the category
**Out:** a sibling report: instances, outputs, per-instance pass/fail against the output contract, plus an edit script

**Contract**
- Generate `prompt.siblings` sibling instances (default 3 deep, 1 fast). They must differ from the original on its free variables.
- Run the meta-prompt on each sibling in a fresh context (J1 `answer` tasks), then judge contract conformance and quality (J1 `judge` tasks).
- A failing sibling produces an edit script back to P1. Iteration follows O1's caps.
- **Modes:** deep by default; `--fast` loads only the core op (D-066).

**Check:** Pass rate ≥ `prompt.sibling_pass_min`.

**Guards against:** overfitting to the user's one example, the most common meta-prompt failure.

**Evidence:** E1, E14
