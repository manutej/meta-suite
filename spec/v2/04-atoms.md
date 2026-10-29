# 04 · Atom index (v2.0 = these 25, no more)

**What an atom is:** one unit of functionality with a typed contract (in → out), exactly one owner
layer, at least one machine check, and a milestone.

**Two layers:**
- **Operators** (LLM) produce things.
- **Gates** (kernel, deterministic) judge them.

Where a gate needs semantic judgment (for example, whether two claims conflict), the LLM produces
sampled local judgments, and the kernel aggregates them and issues the verdict.

| ID | Atom | Layer | Milestone | One line |
|---|---|---|---|---|
| [I1](atoms/I1-brief-classify.md) | brief.classify | operator | M2 | Rough query → typed `Brief` with question type |
| [I2](atoms/I2-spec-expand.md) | spec.expand | operator | M2 | `Brief` → the fullest faithful typed `Spec` |
| [I3](atoms/I3-spec-roundtrip.md) | spec.roundtrip | operator | M2 | Fresh agent reconstructs intent from `Spec` alone |
| [D1](atoms/D1-domain-cover.md) | domain.cover | operator | M2 | Question tree + overlapping region cover + relation vocabulary |
| [D2](atoms/D2-domain-gather.md) | domain.gather | operator | M2 | Parallel per-region claims, written through typed tools |
| [D3](atoms/D3-domain-graph.md) | domain.graph | operator+kernel | M2 | Claims → `DomainGraph` (objects, typed relations, sources) |
| [L1](atoms/L1-lens-yoneda.md) | lens.yoneda | operator | M3 | Define an entity by its relations; prove the profile discriminates |
| [L2](atoms/L2-lens-pullback.md) | lens.pullback | operator+kernel | M2 | A fact is verified only when independent paths agree |
| [L3](atoms/L3-lens-glue.md) | lens.glue | operator+kernel | M3 | Merge local models only along a declared shared interface |
| [L4](atoms/L4-lens-naturality.md) | lens.naturality | operator+kernel | M3 | Two views of one domain: where translation fails = tension |
| [G1](atoms/G1-gate-schema.md) | gate.schema | kernel | M1 | Every typed object validates; claims rejected at write time |
| [G2](atoms/G2-gate-provenance.md) | gate.provenance | kernel | M1 | Sources, independence, no orphans, no unsourced relations |
| [G3](atoms/G3-gate-sheaf.md) | gate.sheaf | kernel | M3 | Per-axiom consistency report + minimal conflicting claim sets |
| [G4](atoms/G4-gate-operadic.md) | gate.operadic | kernel | M3 | Recompose answers at different depths; disagreement = uncertainty |
| [G5](atoms/G5-gate-roundtrip.md) | gate.roundtrip | kernel | M2 | Drift score between original intent and reconstructed intent |
| [G6](atoms/G6-gate-vocab.md) | gate.vocab | kernel | M1 | No internal formal vocabulary in user-facing output |
| [G7](atoms/G7-gate-render.md) | gate.render | kernel | M4 | Accessibility and layout predicates, computed |
| [R1](atoms/R1-review-blind.md) | review.blind | operator | M4 | Fresh-context, preferably cross-model critic; outputs only an edit script |
| [R2](atoms/R2-review-expert-delta.md) | review.expert_delta | operator+kernel | M4 | Count sourced insights a plain answer would miss |
| [A1](atoms/A1-artifact-select.md) | artifact.select | operator+kernel | M4 | Question type × graph shape × medium → artifact form |
| [A2](atoms/A2-artifact-compose.md) | artifact.compose | operator | M4 | Domain structure → artifact structure; render; "how this was checked" |
| [S1](atoms/S1-ledger-run.md) | ledger.run | kernel | M1 | Append-only run record with prompt hashes and gate results |
| [S2](atoms/S2-playbook-delta.md) | playbook.delta | operator+kernel | M4 | Small learned rules, promoted only by ablation |
| [S3](atoms/S3-residue-domain.md) | residue.domain | kernel | M4 | Persist the domain graph so the next nearby run starts warm |
| [O1](atoms/O1-loop-controller.md) | loop.controller | kernel | M1 | Stage state machine: next stage, caps, budgets, ask policy |

## Atom file template

Each atom file has:
- **Purpose**
- **Internal construct** (never user-facing)
- **In / Out** types
- **Contract** (rules)
- **Check**
- **Failure modes it guards**
- **Evidence** (IDs in [evidence](07-evidence.md))
- **Open** (anything unresolved)
