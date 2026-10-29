# 04 · Atom index (v2.0 = these 25, no more)

**What an atom is:** one unit of functionality with:
- a typed contract (in → out)
- a **verdict owner**: the layer that issues its pass/fail
- at least one machine check
- a milestone

**Who does the work vs who rules on it.** Operators (LLM) produce things. The kernel judges them.

| Layer | Meaning |
|---|---|
| operator | The LLM produces the output; the verdict comes from the gates it feeds. |
| kernel | Deterministic, or aggregated over recorded judgments through [J1](atoms/J1-llm-sample.md). |
| operator→kernel | The LLM produces; the kernel issues this atom's own verdict. |

| ID | Atom | Layer | `meta` commands | M | One line |
|---|---|---|---|---|---|
| [I1](atoms/I1-brief-classify.md) | brief.classify | operator | sample | M2 | Rough query → typed `Brief` (1–2 question types) |
| [I2](atoms/I2-spec-expand.md) | spec.expand | operator | validate | M2 | `Brief` → fullest faithful `Spec` |
| [I3](atoms/I3-spec-roundtrip.md) | spec.roundtrip | operator | stage | M2 | Fresh agent reconstructs intent from `Spec` alone |
| [D1](atoms/D1-domain-cover.md) | domain.cover | operator | validate | M2 | Question tree + overlapping cover + relation vocabulary |
| [D2](atoms/D2-domain-gather.md) | domain.gather | operator | claim add, stage | M2 | Parallel per-region claims through typed writes |
| [D3](atoms/D3-domain-graph.md) | domain.graph | operator→kernel | graph build/shape | M2 | Per-region subgraphs + candidate interfaces + shape metrics |
| [L1](atoms/L1-lens-yoneda.md) | lens.yoneda | operator→kernel | sample | M3 | Define an entity by its relations; prove the profile discriminates |
| [L2](atoms/L2-lens-pullback.md) | lens.pullback | operator→kernel | gate G2 | M2 | Verified only when independent paths agree; counterexamples → exceptions |
| [L3](atoms/L3-lens-glue.md) | lens.glue | operator→kernel | interface declare/merge | M3 | The only merge, and only along declared interfaces |
| [L4](atoms/L4-lens-naturality.md) | lens.naturality | operator→kernel | sample | M3 | Two views: where translation fails = tension |
| [J1](atoms/J1-llm-sample.md) | llm.sample | kernel | sample issue/record/aggregate | M1 | Every LLM judgment or sampled answer: issued, recorded, aggregated |
| [G1](atoms/G1-gate-schema.md) | gate.schema | kernel | validate | M1 | Every typed object validates; claims rejected at write time |
| [G2](atoms/G2-gate-provenance.md) | gate.provenance | kernel | gate G2 | M1 | Sources, independence, no orphans, sides sourced |
| [G3](atoms/G3-gate-sheaf.md) | gate.sheaf | kernel | gate G3 | M3 | Overlap-conflict and coverage report + conflict components |
| [G4](atoms/G4-gate-operadic.md) | gate.operadic | kernel | gate G4 | M3 | Composed vs collapsed answers disagree = uncertainty (advisory in v2.0) |
| [G5](atoms/G5-gate-roundtrip.md) | gate.roundtrip | kernel | gate G5 | M2 | Drift between original and reconstructed intent |
| [G6](atoms/G6-gate-vocab.md) | gate.vocab | kernel | gate G6 | M1 | No internal formal vocabulary reaches the reader |
| [G7](atoms/G7-gate-render.md) | gate.render | kernel | gate G7 | M4 | Accessibility and layout predicates, computed |
| [R1](atoms/R1-review-blind.md) | review.blind | operator | stage | M4 | Fresh-context critic; outputs only an edit script |
| [R2](atoms/R2-review-expert-delta.md) | review.expert_delta | operator→kernel | delta, sample | M4 | Sourced insights a plain answer misses (deep tier + eval) |
| [A1](atoms/A1-artifact-select.md) | artifact.select | operator→kernel | validate | M4 | Question type(s) × graph shape × medium → form + layout |
| [A2](atoms/A2-artifact-compose.md) | artifact.compose | operator | gate G6/G7 | M4 | Domain structure → artifact structure; "how this was checked" |
| [S1](atoms/S1-ledger-run.md) | ledger.run | kernel | ledger | M1 | Append-only run record: hashes, usage, verdicts |
| [S2](atoms/S2-playbook-delta.md) | playbook.delta | operator→kernel | playbook | M4 | Record learned-rule candidates (no promotion in v2.0) |
| [O1](atoms/O1-loop-controller.md) | loop.controller | kernel | loop, stage | M1 | Stage state machine: next stage, caps, the single ask rule |

The parked domain-residue atom (formerly S3) is in the [parking lot](01-scope.md#parking-lot).

## Lane map (D-037)

| Lane | What the user sees | Atoms |
|---|---|---|
| **Frame** | The rough question becomes a checked brief | I1, I2, I3, G5 |
| **Map** | The domain is researched and mapped, with sources | D1, D2, D3, L2 |
| **See** | Structure, tensions and contested points are found | L1, L3, L4, G3, G4 |
| **Make** | The artifact is built | A1, A2, G6, G7 |
| **Check** | Blind review, and the record | R1, R2, S1, S2 |
| *(kernel, all lanes)* | — | J1, G1, G2, O1 |

**L2 sits in Map, not See.** Verifying claims against independent sources belongs to the harness,
not the lens (D-021).

## Atom file template

Every atom file uses these headings, in this order:
1. **Purpose**
2. **Internal construct** (never user-facing; "—" if none)
3. **In / Out**
4. **Contract**
5. **Check**
6. **Guards against**
7. **Evidence** (IDs in [evidence](07-evidence.md); "—" if none)

Any literal number in an atom is a **default** and names its config key.
