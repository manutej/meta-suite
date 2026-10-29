# CP1 · compound.map

**Purpose:** Show how this project's work can make the next unit of work cheaper: what to capture, where it lives, how it is retrieved and refreshed.

**Internal construct:** —

**In:** `Spec`, plus optionally S1 ledgers from past runs
**Out:** `CompoundMap`: the loop, the durable anchors, the highest-leverage investments, and the capture/retrieve/refresh policy

**Contract**
- Every investment names the repeated cost it removes and how the saving will be measured.
- It points at S2 and the playbook for capture. Loading learned material across runs waits for v2.1 (D-063).
- **Modes:** deep by default; `--fast` loads only the core op (D-066).

**Check:** G1; every investment has a measurable saving.

**Guards against:** vague "build a knowledge base" advice.

**Evidence:** E15, E16
