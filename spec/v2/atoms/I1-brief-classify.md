# I1 · brief.classify

**Purpose:** Turn a rough query, plus any attachments, into a typed `Brief`. Above all, pick the
question type, because that choice decides which lens operators run.

**Internal construct:** classify the brief as an object, Domain × Audience × Purpose × Question type.

**In:** `raw: text`, `attachments?: file[]`
**Out:** `Brief` (see [architecture](../03-architecture.md#core-types-json-schema-in-kernelschemas))

**Contract**
- `question_type` ∈ {`profile`, `compare`, `synthesize`}.
- An out-of-scope type (path, tradeoff, analogy) maps to the nearest supported type. It is also
  logged as `unsupported_type` in the ledger, which is evidence for the parking lot.
- Classification is sampled 3 times (configurable). Majority wins. Disagreement is recorded in
  `assumptions[]`.
- Ask the user only if the samples disagree **and** `stakes = high` (P9). Otherwise state the
  assumption and proceed.

**Check:** G1 (schema). Classification agreement rate is logged to S1.

**Guards against:** running the wrong operators because the question type was misread.

**Evidence:** E10 (graph structure helps synthesis and hurts lookup, so type selection matters).
