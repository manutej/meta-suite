# G1 · gate.schema

**Purpose:** Every typed object in a run validates against its JSON Schema. Claims are checked at
the moment they are written.

**In:** any typed object, or a run directory
**Out:** a verdict (`pass` / `fail`) plus a list of errors. A failed write returns the error to the
calling agent.

**Contract**
- The schemas live in `kernel/schemas/` and are versioned. Each run records the schema version it
  used.
- The gate is fully deterministic and makes no LLM calls.
- `meta claim add` calls G1 inline and refuses to store an invalid claim.

**Check:** unit tests against fixtures recorded from real runs (P6: no mocked data).

**Guards against:** loose, untyped intermediate data, which is the main way structure decays over a
long run.

**Evidence:** E9.
