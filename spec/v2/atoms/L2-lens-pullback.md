# L2 · lens.pullback

**Purpose:** Treat a fact as verified only when independent paths agree on it. Use the same move to
find which things satisfy two constraints at once, and the counterexamples to a general rule.

**Internal construct:** a pullback, the fibre product over a shared key.

**In:** `Claim[]`, the merged `DomainGraph`, and keys (subject + relation + object, or a value)
**Out:**
- `Finding[kind ∈ {verified, single_source, exception}]`
- candidate conflict pairs, passed to G3

**Contract**
- **verified:** at least `lens.pullback.min_independent` claims (default 2) with different
  `independence_key`s agree on the key.
- **single_source:** fewer agreeing sources. Never stated as settled.
- **Disagreement on the same key** is passed to G3 as a conflict candidate. It is never resolved
  silently.
- **Constraint mode:** the intersection of two typed constraints. Objects that satisfy one
  constraint and are sourced as violating an expected co-constraint become `exception` findings.
  These are the insider exceptions experts know about.
- **L2 is part of the harness, not the lens.** It stays on in the v2−lens arm (see [eval](../05-eval.md)).

**Check:** G2 confirms the independence keys really differ.

**Guards against:** hallucination, citation laundering, and confident claims that rest on one source.

**Evidence:** E3 (pullback verifier: +1.8 to +2.8).
