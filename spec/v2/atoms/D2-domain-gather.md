# D2 · domain.gather

**Purpose:** Collect local evidence, one subagent per region, running in parallel.

**Internal construct:** local sections. Each claim is true only of the region and source it came from.

**In:** `Cover.regions[i]`, `relation_vocab`, tier settings
**Out:** `Claim[]`, written only through `meta claim add`

**Contract**
- Claims are written through the typed tool. Writes that fail the schema are rejected immediately,
  with the error returned to the agent.
- Every claim carries:
  - a source locator, such as a URL or document reference. No long quotes.
  - an `independence_key`, the publisher or origin.
  - its region, agent, model and prompt hash.
- The gather model tier is configurable: a fast, cheap model for gathering, top-tier models for
  judging.
- The agent searches; it does not summarize. Synthesis happens later, in L3.
- A targeted re-gather can be triggered by L1, G4, R1 or R2 for a single region.

**Check:** G1 at write time. G2 after the stage.

**Guards against:** hallucinated facts entering as global truth, citations that can't be traced,
and context bloat in the main thread.

**Evidence:** E9 (+27 F1 when the schema is enforced during construction), E22 (sub-agents for
context hygiene).
