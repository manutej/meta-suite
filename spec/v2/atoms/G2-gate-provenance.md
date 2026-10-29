# G2 · gate.provenance

**Purpose:** No statement in the domain graph or in any finding exists without a traceable origin.

**Internal construct:** —

**In:** a run directory (claims, graphs, findings)
**Out:** a verdict, plus a report of:
- orphan objects
- unsourced morphisms
- fake independence
- unsourced sides

**Contract**
- Every claim has a locator and an `independence_key`. Its stage-token fields come from the kernel.
- Every morphism references at least 1 claim, and every object takes part in at least 1 morphism.
- A `verified` finding has at least `lens.pullback.min_independent` distinct independence keys.
- A `tension` finding has at least 1 claim on each of its `sides`.
- Deterministic, with no network calls. Whether a link still resolves is not checked inside a gate.

**Check:** unit tests on recorded fixtures.

**Guards against:** "according to research" with nothing behind it, and one source counted twice.

**Evidence:** E3, E21.
