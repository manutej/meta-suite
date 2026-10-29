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

## Review round 1 (2026-10-29)

Two blind reviewers looked at v0.1: one for scope and coherence, one for feasibility and eval.
Every P1 finding was accepted. The P2 and P3 findings were accepted except where noted below.

| ID | Date | Decision | Why | Consequence |
|---|---|---|---|---|
| D-017 | 2026-10-29 | Add **J1 llm.sample** as the single owner of every LLM judgment and sampled answer: k is odd, aggregation is by majority, the judge differs from the producer where possible, and every judgment is a ledger event. The kernel holds no model credentials. | "The LLM judges k times" appeared in 6 atoms with no owner, no rule and no CLI. | J1; G3, G4, G5, L1, L4, R2 rewritten |
| D-018 | 2026-10-29 | **Park S3 residue.domain** (to v2.1). **Cut S2 to recording only**: entries are not promoted or loaded in v2.0. The atom budget stays at 25. | Residue and playbook make eval runs dependent on each other, and promotion needs eval-scale ablations. | [scope](01-scope.md), S2 |
| D-019 | 2026-10-29 | Gates split into **deterministic** (G1, G2, G6) and **kernel-aggregated** (G3, G4, G5, R2). At stage 4, G3 and G4 only write reports. Their verdicts move to stage 5 as structural checks on A2's `section_map`. | The gates' verdicts depended on an artifact that didn't exist yet, and "deterministic" was overstated. | 03, G3, G4, A2 |
| D-020 | 2026-10-29 | The v1 arm is `meta-research`, and every arm's output is graded as normalized plain text. | `meta-mvp` produces a pitch bundle, so comparing against it is not like for like. | [eval](05-eval.md) |
| D-021 | 2026-10-29 | **v2−lens keeps L2** (verification is harness, not lens), replaces each lens step with a specified plain step, and is compute-matched to v2 within ±15% tokens. | Otherwise the "lens effect" could just be the effect of more compute. | [eval](05-eval.md) |
| D-022 | 2026-10-29 | **Dev briefs are kept apart from eval briefs.** Pre-registration happens at M0 by SHA-256. The held-out domain is drawn after M4 by a committed procedure. | The v0.1 pilot briefs doubled as development briefs, contaminating the eval. | 05, 06 |
| D-023 | 2026-10-29 | Ship criterion 5 becomes a **cost ceiling**: v2's median $/run ≤ multiplier × v1, with a 3× maximum, set from Phase A. M1 points per dollar are reported, with v2-quick as the cost-matched row. "Fewer tokens than v1" is dropped. | All the sampling and subagents make "tokens ≤ v1" implausible. | 05, [vision](00-vision.md) |
| D-024 | 2026-10-29 | Statistics: the brief is the unit; cluster bootstrap; power calculation from Phase A; Phase B grows to ≥24 briefs if power < 0.8; criterion 3 is pooled across models. | 12 briefs × 3 seeds is not 36 independent runs. | 05 |
| D-025 | 2026-10-29 | G3 operational form for v2.0: pairwise conflict edges, blocked by key and capped, then connected components. The sheaf conditions are measured as overlap-conflict rate (gluing) and uncovered-overlap count (coverage). Minimal unsatisfiable sets and conflicts visible only among three or more claims are deferred. | Exponential cost, and locality as originally stated could not be computed. | G3 |
| D-026 | 2026-10-29 | G4 is **advisory** until Phase A shows its correlation with M1 is at least the pre-registered ρ. | Its evidence comes from reasoning benchmarks, not open research. | G4 |
| D-027 | 2026-10-29 | **Degraded mode.** If the binary is missing, the run is `kernel=absent` and every gate is reported `unchecked`; the model never simulates a verdict. Binaries ship for 5 platforms. v2 is excluded from non-Claude bundles in `dist/` until v2.x. | Distribution risk on sandboxes, and hosts that cannot execute binaries. | 03 |
| D-028 | 2026-10-29 | Claims get kernel-issued ULIDs and per-agent shard files. A stage token stamps provenance. The ledger is written under a file lock. | Parallel writes could collide, and model-reported provenance could be forged. | D2, O1, S1 |
| D-029 | 2026-10-29 | **One ask rule**, owned by O1: ask only when an assumption is high-stakes **and** either expensive to redo or irreversible. | v0.1 had three conflicting ask rules. | O1, I1, P9 |
| D-030 | 2026-10-29 | Briefs can have **1–2 question types** (union of operators). The artifact form and layout tables are config defaults, with layout driven by shape metrics. An unsupported-type mapping is shown to the user as an assumption. | Guards against the design collapsing into templates. | I1, A1 |
| D-031 | 2026-10-29 | **"Second-order effects" is cut from the v2.0 north star** and parked with the composition operator. "Insider exceptions" are produced by L2's constraint-mode counterexamples (a new `exception` finding kind). | The north star promised output that no atom produces. | 00, L2 |
| D-032 | 2026-10-29 | **R2 runs only in the deep tier and in eval.** Judged findings are capped. Re-entry to D2 is limited to 1 per run. | Cost control. | R2, O1 |
| D-033 | 2026-10-29 | Every literal value in an atom is a config default with a named key. L1's probes come only from `relation_vocab`. | Hidden hardcoding (P6). | atoms |
| D-034 | 2026-10-29 | JSON is the canonical format for every typed object. XML is used only to render contracts inside prompts. | I2 and G1 disagreed on the format. | 03, I2 |
| D-035 | 2026-10-29 | **Supersedes D-015.** Thresholds, briefs and rubrics freeze at **M0** as the hashed config plus `evals/sealed/`. | Pre-registration was dated inconsistently across the spec. | 05, 06 |
