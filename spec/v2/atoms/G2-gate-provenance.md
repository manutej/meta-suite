# G2 · gate.provenance

**Purpose:** Make it impossible for any statement in the domain graph, or in a finding, to exist
without a traceable origin.

**In:** a run directory (claims, graph, findings)
**Out:** a verdict plus a report: orphan objects, morphisms with no source, and verified findings
whose sources are not actually independent

**Contract**
- Every claim must have a source locator and an `independence_key`.
- Every morphism must reference at least one claim, and every object must take part in at least one
  morphism.
- A `verified` finding must be backed by claims with at least 2 distinct `independence_key`s.
- A `tension` finding must have at least one source for each side.
- The gate is deterministic. It does not check whether links still resolve: there is no network
  call inside a gate.

**Check:** unit tests on fixtures recorded from real runs.

**Guards against:** "according to research" with nothing behind it, and a single source being cited
twice to look like corroboration.

**Evidence:** E3, E21 (don't trust the model to police itself).
