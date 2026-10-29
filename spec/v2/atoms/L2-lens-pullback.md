# L2 · lens.pullback

**Purpose:** A fact counts as verified only when independent paths agree on it. The same move also
answers "who or what satisfies both constraints."

**Internal construct:** a pullback, the fibre product over a shared key.

**In:** `Claim[]`, `DomainGraph`, keys (subject + relation + object, or a value)
**Out:** `Finding[kind ∈ {verified, single_source}]`, plus disagreement pairs passed to G3

**Contract**
- A claim is `verified` if at least 2 claims with **different `independence_key`s** agree on the key.
- A claim with one source is `single_source`. It is never stated as settled in the artifact.
- Paths that disagree on the same key are passed to G3 as conflict candidates, not silently
  resolved.
- **Constraint mode:** the intersection of two typed constraints over the graph (for example,
  segments that have both budget and pain).

**Check:** G2 confirms that independence keys really differ.

**Guards against:** hallucination, citation laundering (the same source counted twice), and
confident single-source claims.

**Evidence:** E3 (pullback verifier ablation: +1.8 to +2.8, modest but real).
