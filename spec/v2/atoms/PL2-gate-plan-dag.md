# PL2 · gate.plan_dag

**Purpose:** Deterministically check that a plan is executable.

**Internal construct:** —

**In:** `Plan`
**Out:** a verdict, the parallel groups, the critical path, and defects

**Contract**
- Dependencies form an acyclic graph, and every input is produced by some packet or by the user.
- Parallel groups and the critical path are **computed**, never asserted by the model.
- Every packet has typed outputs and at least one done-criterion.
- **Seam count.** The number of cross-unit dependencies and contract requests is computed and reported. A plan whose packets exist only for parallelism, without a passing context-split result, is a defect.
- **Modes:** runs identically in both modes.

**Check:** Zero defects.

**Guards against:** circular plans, orphan inputs, and fake parallelism.

**Evidence:** —
