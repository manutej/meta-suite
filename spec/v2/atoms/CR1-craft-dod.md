# CR1 · craft.dod

**Purpose:** Emit a craft constitution and a Definition of Done for this specific project, where every item is checkable.

**Internal construct:** —

**In:** `Spec` (project, medium, stakes)
**Out:** `Constitution`: rules (each with a reason), slop-vs-craft contrasts for this work, and DoD items (each naming its check)

**Contract**
- Every DoD item names how it is verified: a test, a gate, a measurement or a review lens.
- Rules are specific to the project. Generic platitudes are rejected; whether a rule is generic is a J1 `judge` task.
- **Modes:** deep by default; `--fast` loads only the core op (D-066).

**Check:** Every DoD item has a check.

**Guards against:** "best practices" lists nobody can verify.

**Evidence:** —
