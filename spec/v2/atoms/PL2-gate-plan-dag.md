# PL2 · gate.plan_dag

**Purpose:** Deterministically check that a plan is executable.

**Internal construct:** —

**In:** `Plan`
**Out:** a verdict, the parallel groups, the critical path, and defects

**Contract**
- Dependencies form an acyclic graph, and every input is produced by some packet or by the user.
- Parallel groups and the critical path are **computed**, never asserted by the model.
- Every packet has typed outputs and at least one done-criterion.
- **Modes:** deep by default; `--fast` loads only the core op (D-066).

**Check:** Zero defects.

**Guards against:** circular plans, orphan inputs, and fake parallelism.

**Evidence:** —
