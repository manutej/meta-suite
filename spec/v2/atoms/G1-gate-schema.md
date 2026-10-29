# G1 · gate.schema

**Purpose:** Every typed object validates against its JSON Schema. Claims and judgments are checked
when they are written.

**Internal construct:** —

**In:** any typed object or run directory
**Out:** a verdict and a list of errors. A rejected write returns the error to the agent.

**Contract**
- Schemas live in `kernel/schemas/` and are versioned. Each run records its `schema_version`.
- Fully deterministic, with no LLM involvement.
- `meta claim add` and `meta sample record` call G1 inline and refuse invalid input.

**Check:** unit tests on fixtures recorded from real runs. No mocked data (P6).

**Guards against:** loose intermediate data, which is the main way structure decays over a long run.

**Evidence:** E9.
