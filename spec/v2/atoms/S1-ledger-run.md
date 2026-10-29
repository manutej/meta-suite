# S1 · ledger.run

**Purpose:** An append-only record of each run. It is the substrate for evaluation, playbook
learning and audit.

**In:** events from every stage
**Out:** `runs/<run_id>/ledger.jsonl` and a `RunRecord` summary

**Contract**
- Each prompt is identified by a content hash. The same hash means the same prompt bytes.
- Each run records:
  - the model mapping per role (gather, judge, reviewer)
  - the seed and the tier
  - tokens and wall time per stage
  - every gate verdict and report
  - every edit script
  - the user's internal-view flag
- Append-only: corrections are new events, never rewrites.
- Runs stay in the user's workspace. They are never committed to the public repo (P13).

**Check:** G1 on every event.

**Guards against:** irreproducible claims about what improved, and the question "which prompt
version produced this?"

**Evidence:** E1 (edit log as an accumulating record), E15.
