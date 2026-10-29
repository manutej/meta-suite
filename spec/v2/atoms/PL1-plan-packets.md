# PL1 · plan.packets

**Purpose:** Turn a mapped task into work packets that fresh agents can execute without the planner present.

**Internal construct:** a decomposition whose composition must close: packets compose back into the goal.

**In:** `Spec`, the D1 `Cover` (domain mapping)
**Out:** `Plan`: packets (goal, inputs, outputs, typed done-criteria, context bundle, tools), dependencies, audit points

**Contract**
- Each packet is context-complete. A fresh agent needs nothing outside the packet.
- Done-criteria are checkable: each names a gate or a test.
- Audit points are placed where an error would be expensive. A blind R1 runs there.
- **Modes:** deep by default; `--fast` loads only the core op (D-066).

**Check:** G1; PL2.

**Guards against:** plans that only work if the planner stays in the room, and "done" that nobody can check.

**Evidence:** E14, E16
