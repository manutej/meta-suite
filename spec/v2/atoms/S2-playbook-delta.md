# S2 · playbook.delta

**Purpose:** Capture candidate learned rules as small delta entries, so later versions can test and
promote them. v2.0 records rules but never promotes them.

**Internal construct:** an accumulating edit log over the playbook, in the Writer style.

**In:** finished runs from S1, R1 edit scripts and gate failures
**Out:** `PlaybookEntry` candidates at `1_observation`. An entry rises to `2_field_note` when it
recurs in at least `playbook.recur_min` runs.

**Contract**
- **At most `playbook.max_per_run` proposals per run** (default 3). Each states a rule, the trigger
  that fired it, and the runs that evidence it.
- **The playbook only grows.** It changes by delta edits and is never rewritten wholesale.
- **No promotion in v2.0.** Promotion to `3_shared_insight` or `4_canon` needs ablations at eval
  scale, which belong to v2.1. Until then, entries are **not loaded** into run context.
- **Off during eval.** Playbook writes are disabled and `playbook_hash` is recorded per run, so eval
  runs stay independent of each other.

**Check:** G1. `meta playbook` has no promote command in v2.0.

**Guards against:** context collapse, and self-generated rules that make things worse without
anyone noticing (−1.3 points in SkillsBench).

**Evidence:** E15, E16, E17.
