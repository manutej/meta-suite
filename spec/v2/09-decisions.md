# 09 · Decision log

**Format:** ID · date · decision · why · consequence.
- Entries are append-only.
- To change a decision, add a new entry that supersedes the old one.

| ID | Date | Decision | Why | Consequence |
|---|---|---|---|---|
| D-001 | 2026-10-29 | The category-theory lens is v2's core engine. It is mandatory internally and invisible externally. | It forces relational domain mapping and a hypothesis-and-test structure, and it is what worked in v1. | P1, G6, operators L1–L4 |
| D-002 | 2026-10-29 | v1 is frozen at `8cdbf73` (1.2.0) on `main`. v2 is developed on branch `v2`, and the spec is tracked in the repo. | Comparisons need a stable baseline. | Q9 tag |
| D-003 | 2026-10-29 | The kernel is written in Go. Rust is rejected for v2.0. | LLM latency dominates. Go gives one static binary per platform and shares a language with the Stitch core. | Q1 spike |
| D-004 | 2026-10-29 | No Python in the v2 kernel or runtime. | Fewer environments and one execution path. | v1 skills keep their own scripts |
| D-005 | 2026-10-29 | Refactor between stages only. Prompts are content-hashed. | Swapping prompts mid-run makes results impossible to compare and adds drift. | P10, S1, S2 |
| D-006 | 2026-10-29 | The v2.0 operators are Yoneda, pullback, glue and naturality. Composition gaps, tradeoff (adjunction) and analogy (functor) operators are parked. | Published evidence plus scope control. | [scope](01-scope.md) |
| D-007 | 2026-10-29 | The v2.0 question types are `profile`, `compare` and `synthesize`. | Each maps to proven operators. | I1 |
| D-008 | 2026-10-29 | Operadic consistency is the headline label-free gate. The sheaf gate reports each axiom separately. | Strongest empirical CT signal (E4); axioms fail independently (E6). | G3, G4 |
| D-009 | 2026-10-29 | Verification is blind and preferably cross-model. Nothing grades its own work. | E20, E21 | R1, I3 |
| D-010 | 2026-10-29 | The playbook changes only by delta edits, with elevation levels. Promotion to canon requires an ablation. | E15, E17 | S2 |
| D-011 | 2026-10-29 | The eval is pre-registered and staged: pilot first, then full, then portability. | Avoids tuning to the test, and controls cost. | [eval](05-eval.md) |
| D-012 | 2026-10-29 | v2.0 has an atom budget of 25. Additions must displace or defer something else. | Scope creep is the main project risk. | [scope](01-scope.md#change-control) |
| D-013 | 2026-10-29 | No client names or material in the public repo. Eval domains are public. | Client sensitivity. | P13 |
| D-014 | 2026-10-29 | v2.0 ships as one new pipeline skill alongside the v1 skills. The rest are ported after eval. | Proves the kernel before migrating. | [scope](01-scope.md) |
| D-015 | 2026-10-29 | The ship criteria and thresholds in [eval](05-eval.md) are the initial values. They are frozen at M0. | Pre-registration. | [eval](05-eval.md) |
| D-016 | 2026-10-29 | Intake expands rough queries into a typed spec, and the expansion is checked by an adjunction round trip (I2/I3/G5). | Underspecified intent is the real input. Minimal-prompt guidance applies to boilerplate, not to intent. | I2, I3, G5 |
