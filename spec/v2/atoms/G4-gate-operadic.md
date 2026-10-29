# G4 · gate.operadic

**Purpose:** A label-free quality signal. Answer each node of the question tree two ways and measure
how often the answers disagree:
- composed from its children's answers
- answered directly at a collapsed depth

**Internal construct:** operadic consistency. Composition should be associative across levels of the
decomposition tree.

**In:** `Cover.tree`, the answers at each node, and the tier's sample count k
**Out:** a disagreement rate per node and for the whole run, plus a list of flagged nodes

**Contract**
- For each internal node, compare the composed answer with the collapsed answer, k samples each.
  An LLM judges equivalence, and the kernel aggregates.
- Nodes above the disagreement threshold (set in config) are flagged as uncertain. They trigger a
  targeted re-gather (D2) for their regions, up to the tier cap, and the artifact marks them as
  lower confidence.
- The run-level rate is reported in "how this was checked" and logged as an eval metric.

**Check:** run-level disagreement ≤ threshold, **or** every flagged node is surfaced in the artifact
as uncertain.

**Guards against:** confident answers that only held together because of how the question was split.

**Evidence:** E4 (r = 0.86–0.94 with accuracy across 12 LLMs, label-free). This is the strongest
CT-derived signal found.
