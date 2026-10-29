# 02 · Principles (the v2 doctrine)

Each principle names the gate or mechanism that enforces it. If a principle has no enforcement, it
is a wish. Principles that have no gate say so.

| # | Principle | Enforced by |
|---|---|---|
| P1 | The lens is mandatory inside the machine and invisible to the reader. | [G6 vocab gate](atoms/G6-gate-vocab.md) |
| P2 | Every LLM output is a **local claim with provenance**. Global statements come only from gates that combine local claims. | [D2](atoms/D2-domain-gather.md), [G2](atoms/G2-gate-provenance.md) |
| P3 | Structure is enforced by code **at write time** (typed tools), not by instructions or after-the-fact checks. | [G1](atoms/G1-gate-schema.md), typed tool calls in D2 |
| P4 | Verification is blind and preferably cross-model. Nothing grades its own work. | [R1](atoms/R1-review-blind.md) |
| P5 | Every claim of improvement is measured: pre-registered, ablated, with variance reported. | [eval](05-eval.md) |
| P6 | No hardcoded domain content and no mocked data. Test fixtures are recorded from real runs. Thresholds live in config, not code. | code review + [G1](atoms/G1-gate-schema.md) |
| P7 | Prompts are pitched at the right altitude: no personas or chain-of-thought boilerplate. Procedures are written as programs where possible. | `meta check` prompt lint (M1) |
| P8 | Progressive disclosure: SKILL.md ≤ 500 words, and 2–3 modules loaded per run. Depth lives in nested folders. | `meta check` word count (M1) |
| P9 | Ask only under the single rule owned by O1: an assumption that is high-stakes **and** either expensive to redo or irreversible. Otherwise state the assumption and proceed. | [O1](atoms/O1-loop-controller.md) |
| P10 | **Adapt during the run, at stage boundaries** (D-057). Before each stage, O1 can rewrite that stage's prompt from the edit scripts and gate results so far. O1 decides when; the new text is a J1 `answer` task; the kernel hashes and logs each version and diff. A stage's prompt never changes mid-stage. Adaptation is on in the product and in the `v2` eval arm; `v2-frozen` is the ablation (D-070). | [O1](atoms/O1-loop-controller.md), [S1](atoms/S1-ledger-run.md) |
| P11 | Every artifact passes the accessibility gate before it is shown, **and** meets the artifact bar: Level 2 on every applicable dimension of the [rubric](12-artifact-rubric.md) (D-056, D-069). An SME would forward it, and the reader learns something in the first 60 seconds. | [G7](atoms/G7-gate-render.md) + `meta rubric`, [R1](atoms/R1-review-blind.md) |
| P12 | Show what the code did. Every artifact carries a plain-language "how this was checked" section built from the actual gate results. | [A2](atoms/A2-artifact-compose.md) |
| P13 | Nothing client-identifying goes in the public repo. Eval domains are public. | review checklist |
| P14 | **Plug and play for non-technical users.** Install is one action (add the plugin, or upload one zip). There is no terminal, no keys, no config and no git. The user types a question and gets the artifact. The kernel, gates and binaries are invisible. If the kernel can't run, the user still gets the artifact, with a plain note that some checks didn't run. All setup and verification steps belong to the maintainer, never the user. | M4 exit (meta-lens) and M5b exit (meta-mvp): a first-run test by a non-technical person (D-053, D-071) |

## Working with a model that hallucinates, drifts and is non-deterministic

- **Hallucination** is contained by P2 + [L2 pullback](atoms/L2-lens-pullback.md). A fact that only
  one source supports is flagged as single-source and never stated as settled.
- **Drift** is measured, not hoped against:
  - [I3](atoms/I3-spec-roundtrip.md) + [G5](atoms/G5-gate-roundtrip.md) round-trip the spec back to
    the original intent.
  - [L3 glue](atoms/L3-lens-glue.md) refuses to merge anything without a declared shared interface.
- **Non-determinism becomes a signal.** Sampling the same structure more than once and measuring
  agreement ([G4 operadic consistency](atoms/G4-gate-operadic.md)) shows where the model or the
  domain is uncertain.
