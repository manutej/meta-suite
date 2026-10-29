# S1 · ledger.run

**Purpose:** Keep an append-only record of each run. Evaluation, learning and audit all read from it.

**Internal construct:** an accumulating edit and event log, in the Writer style.

**In:** events from every stage, and judgments from J1
**Out:** `runs/<run_id>/ledger.jsonl` and a `RunRecord`

**Contract**
- **Prompts are identified by content hash.**
- **Every `RunRecord` captures:**
  - surface
  - `kernel` flag (present or absent)
  - model map
  - `config_hash`, `schema_version`, `playbook_hash`
  - seed and tier
  - `usage`, with its source (`api`, `harness` or `estimated`)
  - wall time per stage
  - every gate verdict, judgment and edit script
  - blind attestations
  - the internal-view flag
- **Append-only, under a file lock.** A correction is a new event.
- **Runs stay in the user's workspace.** They are never committed. CI checks that `runs/` is
  gitignored (P13).

**Check:** G1 on every event.

**Guards against:** improvement claims nobody can reproduce, and "which prompt produced this?"

**Evidence:** E1, E15.
