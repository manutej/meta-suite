# I1 · brief.classify

**Purpose:** Turn a rough query and its attachments into a typed `Brief`. This includes choosing the
question type(s), which decides which lens operators run.

**Internal construct:** the brief as an object, Domain × Audience × Purpose × Question type(s).

**In:** `raw: text`, `attachments[]`
**Out:** `Brief` (see [core types](../03-architecture.md#core-types-json-schema-in-kernelschemas))

**Contract**
- `question_types` holds 1 or 2 values from {`profile`, `compare`, `synthesize`}. With two types,
  the operators are the union of both sets.
- **Unsupported types** (path, tradeoff, analogy) map to the nearest supported type and are logged
  as `unsupported_type`. The artifact must also state the mapping as an assumption.
- The classification is a J1 `judge` task, k samples (config: `tiers.<tier>.k`). Any disagreement
  is recorded in `assumptions[]`.
- **Routing (front door).** I1 also decides whether the request is really a research question or a specialist task: a meta-prompt, plan, design, bundle, review, craft standard, compounding map or context strategy. If it is specialist work, it routes there ([suite map](../11-suite.md)) and records the routing as an assumption the user can override.
- **Mode.** Deep is the default. `--fast` or "quick" selects fast mode (D-066).
- I1 also extracts:
  - `targets[]`: the subjects for L1
  - `views[2]`: the two sides for L4
  - `vocab_allowlist[]`: needed only when the brief is itself about formal methods (see G6)
- When to ask the user is decided by the single ask rule in [O1](O1-loop-controller.md). I1 only
  flags candidate assumptions with their stakes.

**Check:** G1. The classification agreement rate is logged to S1.

**Guards against:** running the wrong operators, and silently forcing a question into a type it
doesn't fit.

**Evidence:** E10 (graph structure helps synthesis and hurts lookup).
