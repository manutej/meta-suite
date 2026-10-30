# PL1 · plan.packets

**Purpose:** Turn a mapped task into an executable plan. The default topology is **one implementer** (the whole spec in one context, compacting if needed), plus a **fresh reviewer** (sees only the spec and the diff), plus a **verifier** (sees only the user-flow matrix and makes live calls). Packets for fresh agents are the exception, emitted only when they pass the context-split test (D-084).

**Internal construct:** a decomposition whose composition must close: packets compose back into the goal.

**In:** `Spec`, the D1 `Cover` (domain mapping)
**Out:** `Plan`: packets (goal, inputs, outputs, typed done-criteria, context bundle, tools), dependencies, audit points

**Contract**
- **Context-split test, per unit.** Ask whether what the main session knows would *help* this unit, *hurt* it, or *not matter*.
  - **Help:** the same implementer does it.
  - **Hurt:** a fresh reviewer or verifier does it.
  - **Not matter:** a separate errand (for example, independent read-only research).
  - A packet is emitted only if it also passes the five parallel criteria and needs nothing another unit will learn.
- **Roles are named by function, not as personas:** implementer, reviewer, verifier, errand.
- Any emitted packet is context-complete: a fresh agent needs nothing outside it.
- Done-criteria are checkable: each names a gate or a test.
- Audit points are placed where an error would be expensive. A blind R1 runs there.
- **Modes:** deep by default; `--fast` loads only the core op (D-066).

**Check:** G1; PL2.

**Guards against:** plans that only work if the planner stays in the room, and "done" that nobody can check.

**Evidence:** E14, E16
