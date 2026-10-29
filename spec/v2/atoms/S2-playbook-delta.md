# S2 · playbook.delta

**Purpose:** Self-improvement without context collapse. Small learned rules accumulate, and none is
trusted until it has been measured.

**Internal construct:** a Writer-style accumulating edit log over the playbook.

**In:** completed runs (from S1), R1 edit scripts, gate failures
**Out:** `PlaybookEntry` candidates, and promotions or demotions

**Contract**
- A run proposes at most 3 entries. Each has a rule, a trigger (when it applies), and its evidence
  runs.
- Promotion through elevation levels:

  | Level | Requirement |
  |---|---|
  | `1_observation` | proposed by a run |
  | `2_field_note` | recurs in at least N runs |
  | `3_shared_insight` | a reviewer agrees |
  | `4_canon` | an ablation on held-out briefs shows a gain and no regression |

- An entry whose later ablation fails is demoted.
- Only entries at level `3_shared_insight` or above are loaded into context, and only at stage
  boundaries (P10). Nothing is swapped mid-stage.
- The playbook is never rewritten wholesale. It changes only by delta edits.

**Check:** `meta playbook promote` refuses to promote to `4_canon` without an attached ablation result.

**Guards against:**
- brevity bias
- context collapse
- self-generated rules that make things worse (self-generated skills averaged −1.3 points in
  SkillsBench)

**Evidence:** E15, E16, E17.
