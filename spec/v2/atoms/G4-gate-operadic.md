# G4 · gate.operadic

**Purpose:** Get a quality signal without labels. Each internal node of the question tree is
answered two ways: composed from its children's answers, and answered directly in one collapsed
step. How often the two disagree measures uncertainty.

**Internal construct:** operadic consistency. Composition should hold across depths of the tree.

**In:** `Cover.tree`, the merged graph and claims
**Out:**
- per-node disagreement
- the run-level rate
- flagged nodes, each becoming a `Finding[kind=uncertain]`

**Contract**
- **Answers come from J1 `answer` tasks.** For each internal node, J1 produces k composed answers
  (from the children's answers) and k collapsed answers (straight from the node's claims).
- **Equivalence is a J1 `judge` task.**
- **Flagging.** A node is flagged when its disagreement exceeds `gates.g4.node_threshold`. Flagged
  nodes may trigger a targeted D2 re-gather, within the loop caps.
- **Report only at stage 4.** At stage 5 every flagged node must map to a section tagged `uncertain`
  in `section_map`.
- **Advisory in v2.0.** The run-level rate is reported, not enforced, until Phase A shows its
  correlation with M1 is at least `eval.g4_promote_rho` (pre-registered). Only then does it become a
  blocking gate.

**Check:** the stage-5 structural coverage check. The run-level rate is logged as eval metric M4.

**Guards against:** confident answers that only held together because of how the question happened
to be split.

**Evidence:** E4. Its r = 0.86–0.94 was measured on reasoning benchmarks, not open research, hence
the advisory status.
