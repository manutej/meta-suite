# D2 · domain.gather

**Purpose:** Collect local evidence with one subagent per region, running in parallel.

**Internal construct:** local sections. Each claim holds only for its region and source.

**In:** a region from `Cover`, `relation_vocab`, a stage token
**Out:** `Claim[]`, appended to `claims/<agent>.jsonl` only through `meta claim add`

**Contract**
- **Typed writes only.**
  - `meta claim add` validates the claim (G1) and returns either a kernel-issued ULID or a schema
    error the agent can act on.
  - Run, region, agent, model and prompt hash come from the **kernel-issued stage token**
    (`meta stage begin`). They are never typed in by the model.
- **Every claim names its source.**
  - A locator (URL or document reference). No long quotes.
  - An `independence_key`: the publisher or origin.
- **Model tier.** Configurable (`roles.gather`). Regions that D1 flags as contested, or that a review sends back, use `roles.gather_contested`, a stronger model by default. The default for other regions is a fast, cheap model; judges use a
  top-tier model.
- **Gather is not synthesis.** No summarising; synthesis happens in L3.
- **Targeted re-gather.** L1, R1 or R2 can trigger a re-gather of a single region. Re-entry to D2
  is limited by `loop.d2_reentry_max`.

**Check:** G1 at write time, then G2 after the stage. The kernel merges shards at the stage boundary.

**Guards against:**
- hallucinated facts entering as global truth
- citations that can't be traced
- forged provenance
- bloating the main thread's context

**Evidence:** E9, E22.
