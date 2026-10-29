# CX1 · context.scaffold

**Purpose:** Design a project's context strategy and emit a droppable L0–L4 substrate scaffold.

**Internal construct:** the context tiers form a cover; L0/L1 summaries must glue with the L2–L4 content beneath them.

**In:** `Spec` (project, team, surfaces)
**Out:** `ContextPlan` (what to keep, where, what loads when, compression, handoff) plus scaffold files (L0 trailhead ≤ 500 words, an L1 index, L2–L4 folders)

**Contract**
- It uses the same progressive-disclosure tiers as the suite itself.
- Every piece of kept knowledge has a load trigger and a refresh rule.
- **Summaries are supported.** Each claim in an L0/L1 summary is checked against the deeper layers (G3 layer 3 judgments); unsupported summary claims are edits.
- The scaffold is plain files, ready to drop into a repo or a Notion hub.
- **Modes:** deep by default; `--fast` loads only the core op (D-066).

**Check:** G1; L0 word limit; every item has a load trigger; G3 layer 3 finds no unsupported summary claim.

**Guards against:** context dumps that load everything every time.

**Evidence:** E22, E23
