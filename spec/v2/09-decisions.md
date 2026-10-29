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
| D-036 | 2026-10-29 | **The pipeline skill is named `meta-lens`.** `meta-atlas` and `meta-forge` were rejected. | The name says what sets the skill apart (how it sees). The rejected names collide with existing skills (`/atlas`, the Forge plugins) and emphasize the output instead. | Q2 closed |
| D-037 | 2026-10-29 | **Five user-facing lanes inside the one skill: Frame → Map → See → Make → Check.** The lanes are stages, not separate skills. The user only types a question; the type is inferred by I1. When the v1 skills are ported (v2.x), they become specialists behind `meta-lens` as the front door. | Keeps the plugin simple to use, and lets the internal structure be explained in plain words. | 03 lane column, 04 lane map |

## Review round 2: sheaf-core integration audit (2026-10-29)

An independent agent, given the sheaf-core repo and the v0.2 spec, reviewed G3/L3/D1 against working,
tested sheaf code.

| ID | Date | Decision | Why | Consequence |
|---|---|---|---|---|
| D-038 | 2026-10-29 | **The Go kernel embeds a port of the sheaf-core graph subset**: δ⁰ with per-edge maps, components, per-component eigen cut, H⁰/H¹ by counting, energy/discord, Hodge split, Gaussian null, relative H⁰ + dense harmonic extension, and optionally δ¹ on triangles. It is tested in CI against golden reports from sheaf-core. Rust is only the oracle and is never shipped. **Committing the port was gated on Q12, which D-047 has resolved.** | One binary per platform, no version drift, and the Go prototype already matched 11/11 within 2.4e-12. Refines D-003. | `kernel/internal/sheaf`, `testdata/sheaf`, roadmap M1 |
| D-039 | 2026-10-29 | **G3 gets three sheaf layers**: numeric key (deterministic), relational cycle (deterministic, Hodge harmonic share), and judgment balance (aggregated, reusing J1 votes). Free-text claims stay on pairwise judgments, and embedding stalks are banned. | Real computation wherever the claims have linear structure, and honesty where they don't. | G3, 03 G3Report |
| D-040 | 2026-10-29 | **The two sides of an obstruction come from the layer-3 balance section.** A frustrated component becomes an obstruction with no sides, tagged contested. | Conflict components are not two-sided: A–B and B–C pull in C even when A and C agree. | G3, Finding |
| D-041 | 2026-10-29 | **L3 merge checks coherence with union-find** against J1 "distinct" constraints, and refuses closures that join objects judged distinct. | Colimits over 3+ regions merge transitively. | L3 |
| D-042 | 2026-10-29 | **The Cover declares a universe and must satisfy the covering condition.** Overlaps and triple overlaps are computed, not declared, and adjacency means a non-empty overlap. | "Open cover" was empty language, and "adjacent" was defined in terms of overlaps. | D1, Cover |
| D-043 | 2026-10-29 | **Claims carry `value_type`, `unit` and `as_of`, and the kernel normalizes units.** Units live in the data, never in the restriction maps. | No numeric layer is possible without it. Unit factors inside the maps produced a spurious "not unique" result in the toy. | Claim, D1 vocab |
| D-044 | 2026-10-29 | **Amends D-025.** Linear relational cycle conflicts (3+ claims) are **in scope** through G3 layer 2. Three-way conflicts among free-text claims stay parked. | Layer 2 brings them back cheaply. | scope parking lot |
| D-045 | 2026-10-29 | G3 reports `judged_pairs / candidate_pairs`, samples stratified by key when capped, and drops "locality" as a separate axiom: for claim sets it reduces to entity resolution, which is D3's job. D1's tree is called a decomposition tree, and G4 measures its composition. L4 states that the functor condition is unchecked in v2.0. | Undefined denominator, sampling bias, and terms the math didn't support. | G3, D1, L4 |
| D-046 | 2026-10-29 | The v2−lens arm switches off **all three G3 layers**. M5 reports numeric and judgment obstructions separately. | Ablation hygiene. | eval |
| D-047 | 2026-10-29 | **Resolves Q12. The graph-sheaf subset is open-sourced under MIT inside meta-suite** (`kernel/internal/sheaf`, plus its tests and synthetic fixtures). sheaf-core stays private and is the reference oracle. The rest of sheaf-core (m1, FFI, Lanczos, CG, mod-p) stays proprietary. The canonical Go package is first developed in sheaf-core (`go/`) with parity tests, then copied in at M1. | The math is standard, and it is a key **inbound lead magnet**: working, tested gluing code people can run. | D-038 un-gated; `go/` must be self-contained |
| D-048 | 2026-10-29 | **The `meta sheaf analyze` command is exposed as a standalone tool**, so the gluing check can be used without the pipeline. It is a CLI surface on G3's engine, **not a new atom**, so the budget stays at 25. Docs and examples ship with it; landing-page and marketing copy stay in the separate articles track. | Lead-magnet value at near-zero scope cost. | 03 CLI table |

## M0 spike results (2026-10-29)

| ID | Date | Decision | Why | Consequence |
|---|---|---|---|---|
| D-049 | 2026-10-29 | The skill invokes the kernel only as `sh bin/meta …`. The launcher restores the exec bit and, on a noexec mount, falls back to a `$TMPDIR` copy. It reports `kernel=absent` only when both fail. Windows gets `meta.cmd`. | Spike: Python `zipfile` strips the exec bit (exit 126 with `./bin/meta`). The fallback was verified on a noexec tmpfs. | 03 degraded mode |
| D-050 | 2026-10-29 | Every kernel command emits deterministically ordered output. | Spike: the JSON Schema library returned errors in random order, so a deterministic gate gave different bytes on each run. | all kernel commands |
| D-051 | 2026-10-29 | **Resolves Q11.** The eval harness is `claude --bare -p --output-format json` (or the Agent SDK) with tools locked down. Tokens come from `modelUsage`. Dollars are labelled estimated unless reconciled against the usage API. An open-weight model joins Phase C only if it passes a pre-registered tool-use smoke test. **Resolves Q10 for the cloud sandbox:** Chrome for Testing and Playwright are present. | Spike, docs fetched. | 05 Phase C, M6 |
| D-052 | 2026-10-29 | **Local checks are the gate; remote CI is an asynchronous witness.** `meta check` runs in seconds and is what agents run before every push, so it is the gate that blocks progress. GitHub CI is **optional**: it runs the same command on PRs and merges as independent, non-self-reported evidence, with `continue-on-error`. It is never a required status check and never blocks a merge or an agent's iteration loop. Cross-builds also run locally. | Remote CI is too slow for an agentic build loop, but the only evidence an agent can't fake is a check it didn't run itself. | 06 M1; supersedes "CI checks" wording elsewhere |
| D-053 | 2026-10-29 | **P14 plug and play.** End users never see setup: there is no terminal, keys, config or git, and degraded mode still delivers the artifact. M4's exit gate includes a first-run test by a non-technical person. Probes, tags and parity tests are maintainer-only. | meta-suite's audience includes non-technical professionals; it is a lead magnet. | P14, M4 |
| D-054 | 2026-10-29 | **Q1 resolved for the Claude app.** Installed plugins land in a writable folder that keeps exec bits, and they are not mounted noexec. `meta-probe` passed there with exec bits stripped, and on a read-only mount through the temp-copy fallback. A fresh **Haiku** agent following only SKILL.md also passed every step with no help (plug-and-play evidence for P14). Still open: Cowork desktop and a real Mac. | Spike evidence from 2026-10-29. | 03 surface matrix, Q1 |

## Cross-tier fidelity review round 1 (2026-10-29)

Opus, Sonnet and Haiku reviewers, blind to each other, audited the spec against the
[vision ledger](10-vision-ledger.md). **All three returned "drifting, recoverable."** The engine is
faithful, but the product the owner described was under-served: artifact quality, speed,
self-adaptation, VR-00 breadth, and the fidelity loop itself. The full results are in
[reviews/2026-10-29-fidelity-r1.md](reviews/2026-10-29-fidelity-r1.md).

| ID | Date | Decision | Why | Consequence |
|---|---|---|---|---|
| D-055 | 2026-10-29 | **Cross-tier fidelity review at every milestone exit.** Opus, Sonnet and Haiku, blind to each other, audit the spec and the build against the vision ledger. Any VR `missing` or `contradicted` blocks the exit unless the owner signs off. Results are saved in `reviews/`. | VR-22. The ledger cited this decision before it existed, and all three reviewers flagged that. | 06 exit gates |
| D-056 | 2026-10-29 | **The artifact is measured and gated.** M8 is a first-look rating of the rendered artifact against v1-mvp and v1. A2 has an artifact bar using the owner's exemplars, and R1 reviews the rendered output. Ship criterion 6 is "v2 ≥ v1-mvp on M8". | VR-02/03. Text-only grading (D-020) could not see "stunning". D-020 is kept for M1. | 05, A2, R1, P11 |
| D-057 | 2026-10-29 | **Amends D-005: adapt during the run at stage boundaries.** O1 rewrites the next stage's prompt from edit scripts and gate results. Each version is hashed and logged, never mid-stage. Frozen in eval except in a `v2-adaptive` arm. | VR-10/08. The eval-hygiene reason only applies to eval arms. | P10, O1, 05 |
| D-058 | 2026-10-29 | **Walking skeleton at M2.** One dev brief runs Frame → Map → Make → Check end to end with a real rendered artifact, published for the owner. | VR-24. The first visible output came too late (M4). | 06 M2 |
| D-059 | 2026-10-29 | **Opus-native execution:** programmatic tool-calling loops, batched J1 (`meta sample run`), prompt caching, per-stage budgets, and wall time and turns in the ledger. **Speed is a ship criterion** (v2-quick ≤ v1 p50). | VR-06. Speed had gone backwards and was unmeasured. | 03, J1, 05 |
| D-060 | 2026-10-29 | **Depth where nuance is made:** D1 insider probes per region, a stronger gather model for contested regions, and R2-lite in the standard tier. **Amends D-032.** | VR-03. Every atom verified, but nothing pushed toward insider knowledge. | D1, D2, R2 |
| D-061 | 2026-10-29 | **v1 inventory** (`11-suite.md`): every v1 element is mapped to a v2.0 atom, a v2.x port, or dropped (none). It sets a port order and a plan B if Phase A is a no-go. **"Works best with Opus 5.5"** may be claimed only if Opus beats Sonnet and Haiku on M1 and M8. v1 stays the default install until v2.0 passes. M9 (human input) and M10 (speed) added. A1 takes the reasoning sequence as input. | VR-00, 23, 14, 07. | 11, 05, A1 |

## Owner decisions on the open trade-offs (2026-10-29)

| ID | Date | Decision | Why | Consequence |
|---|---|---|---|---|
| D-062 | 2026-10-29 | **Owner approved the port order in [suite map](11-suite.md):** meta-lens first. Then review, prompting and mvp in v2.1, and design, planning, craft, compound and context in v2.2. v1 stays the default install meanwhile. | Signed off by the owner (VR-00 deviation accepted). | 11 |
| D-063 | 2026-10-29 | **Owner chose: cross-run learning (playbook loading and domain memory) arrives in v2.1, after proof.** Only rules that pass a with/without test are used. In-run stage-boundary adaptation (D-057) stays in v2.0. | Signed off by the owner (VR-10/12 phasing accepted). | S2, parking lot |
| D-064 | 2026-10-29 | **Owner chose depth as the default.** Standard tier is the default run, and the quick tier is one tap away. The speed criterion (D-059) still applies. | Signed off by the owner. | O1 default tier, I1 |

## Owner correction: v2.0 is a new suite (2026-10-29)

| ID | Date | Decision | Why | Consequence |
|---|---|---|---|---|
| D-065 | 2026-10-29 | **Supersedes D-014 and D-062. v2.0 is a whole new suite.** `meta-lens` is the front door, and all nine v1 skills are rebuilt on the shared kernel and lanes, each as a thin recipe with at most 2 own atoms. The atom budget is **34** (25 core + 9 specialist). The port order becomes the build order inside M4b. Each specialist has a per-skill v1 comparison and a per-skill v1 fallback. | Owner: "I want a new suite, the new skill was just temporary" (VR-27). | 11 suite map, 04, 05, 06 M4b |
| D-066 | 2026-10-29 | **Amends D-064. Two modes for every skill: deep (default) and `--fast`.** Fast loads only L0 and the core op (lighter skill weight through progressive disclosure), uses k = 1, at most 3 regions, one review pass, no G4 and no R2. Deep scales with scope, up to full deep research. The speed criterion is measured on `--fast`. | Owner: "make deep default and --fast as a flag", like Sonnet vs Opus (VR-29). | 03 modes, 05, I1 |
| D-067 | 2026-10-29 | **Artifact rubric** ([12](12-artifact-rubric.md)). Its minimum bar (Level 2) is derived from the owner's recent top artifacts: visuals, colour, layering, mobile readiness, and more. Level 3 is "beyond what the owner imagined." Machine-checkable dimensions run as checks, and judged dimensions feed M8. | Owner request (VR-28). | A2, R1, M8, DS1 |

## Fidelity review round 2 (2026-10-29)

Opus, Sonnet and Haiku reviewers, blind to each other, audited v0.5 against the
[vision ledger](10-vision-ledger.md). Their consolidated findings are applied below.

| ID | Date | Decision | Why | Consequence |
|---|---|---|---|---|
| D-068 | 2026-10-29 | **Consistency sweep.** Atom count 34 and `meta check` (not CI) everywhere. The fast arm is `v2-fast`. One config scheme, `modes.<mode>.<key>` with mode ∈ {fast, deep_narrow, deep_broad}; the 03 modes table is the single source of mode values. New CLI rows: `meta check`, `meta rubric`, `meta palette`. G7 thresholds get config keys. PL2 runs identically in both modes. The M2 walking skeleton uses the v2−lens plain merge for L3. G6 allowlist entries can come from the brief's own domain vocabulary. v2.0 ships the `meta` CLI as an agent-callable harness core; a stable public API/SDK is v3. "Success, defined" gains the artifact, `--fast` and per-specialist criteria. **Supersede markers** (old rows stay as written): D-012 and D-018 (budget 25) are superseded by D-065 (34); D-037's "ported in v2.x" and D-062 (port order) are superseded by D-065; the tier wording ("quick", "standard", "deep tier", `v2-quick`) in D-023, D-032, D-059, D-060 and D-064 now reads as the D-066 modes `v2-fast`, `deep_narrow` and `deep_broad`. | Reviewers found stale counts, three naming schemes and conflicting numbers across files. | 00, 01, 03, 06, 08, 11, I1, J1, D1, G3, G6, G7, O1, R1, R2, S1, PL2 |
| D-069 | 2026-10-29 | **Rubric wired in.** G7 owns the rubric's machine checks as the config table `rubric.checks`, run by `meta rubric` (no new atom); a check without a browser is `unchecked` and can't support an L2 claim. A2's build bar is L2 on every applicable dimension, plus one L3 pass in deep mode. R1 scores every applicable dimension (L1 → P1 edit, L3 opportunity → P3). M8 uses the 60-second protocol. Ship criterion 6 is absolute and relative. Brand inputs override token-count and font numbers; structural dimensions don't bend. N/A per form. | The rubric existed but nothing enforced it. | 12, G7, A2, R1, DS1, P11, 05 |
| D-070 | 2026-10-29 | **Amends D-057: adaptation is the product.** Ship criteria are measured on `v2` with stage-boundary adaptation on; `v2-frozen` is the ablation arm. O1 decides when and records why; the new text is a J1 `answer` task; the kernel hashes and logs the version and diff. | VR-10: the shipped product must be the measured one. | P10, O1, 05 arms |
| D-071 | 2026-10-29 | **The suite is real in the architecture.** One folder per skill under `skills/`; the plugin ships one shared `bin/` at its root, and each per-skill zip in `dist/` bundles its own copy (about 6.7 MB). v2 ships as the preview plugin `meta-suite-next` until the ship criteria pass, then becomes `meta-suite` 2.0; v1 stays installable from 1.2.0. Core types for every specialist output. The P14 first-run test covers meta-lens and meta-mvp. | VR-27/21/23: the architecture still described one skill. | 03, 06 M4/M5b, P14 |
| D-072 | 2026-10-29 | **Lens in the specialists, v1 gaps closed.** CR1, CP1 and CX1 get internal constructs and a See-lane use; meta-review surfaces disagreement between lenses as a tension. Homes added for user-facing takeaways, blind-spot probes, DS1's three-role convergence and precedence, MV1's data model, meta-prompt and ambiguity panel, credibility triage, the lens catalog, and the dashboard and explainer forms. meta-research becomes an alias of meta-lens in brief form, with no separate per-skill eval. | VR-00/01: "nothing is dropped" was not yet true. | 11, A1, A2, R1, DS1, MV1, CR1, CP1, CX1 |
| D-073 | 2026-10-29 | **Eval right-sized.** Arms named per phase. Per-skill comparison: 8 specialists, Opus, 3 sealed tasks, 1 seed, deep (fast for speed only); a pooled suite test (CI excludes 50%) plus per-skill non-inferiority (≥ 50%), else v1 fallback. SME fallback pre-registered per domain and per skill. New metric M11 (context tokens loaded). M9 becomes ship criterion 9. About 480 runs (target ≤ ~550). Q5: a non-Claude judge if it passes a smoke test, otherwise "same-family, different-tier" stated as a limitation. Criterion 3 is tested on Opus. Roadmap: M4 → M5 pilot → M5b specialists → M6. | VR-13: the per-skill plan was under-powered and "all 5" was ambiguous. | 05, 06, 08 |
| D-074 | 2026-10-29 | **Owner: watch and alerts are in v2**, as milestone M7 after the M6 suite release and before any v3 work: re-run on a cadence, detect changed facts and gate results from the ledger and saved run state, notify the user. M7's atom budget is set by a decision at M7 entry; no atoms are added now. Resolves Q13. | Owner: "Alerts or monitoring should be a part of the product but we can push it in a later milestone within v2" (VR-14). | 00, 01, 06 M7, 08 |

## Build starts (2026-10-29)

| ID | Date | Decision | Why | Consequence |
|---|---|---|---|---|
| D-075 | 2026-10-29 | **Vertical slice first.** The owner wants something working to evaluate, so the M2 walking skeleton is built now as **slice 1**, in parallel with the remaining M0 work (freeze, round-3 review, pre-registration, SMEs). Slice 1 contains the kernel commands it needs, `meta sheaf analyze` (Go port vendored from sheaf-core under MIT, D-047), and meta-lens on the v2−lens path (See off, plain merge). It lives in `next/`, the `meta-suite-next` preview plugin (D-071), and v1 is untouched. Its first real run is published as a milestone artifact, next to v1 meta-research on the same question. **Amends the roadmap's strict milestone sequence for this slice only.** | Owner: "push towards getting something working that we can evaluate." VR-24. | next/, 06 |

