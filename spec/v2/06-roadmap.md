# 06 · Roadmap

Each milestone has an exit gate. Nothing in a milestone starts until the previous milestone has
exited. **At every milestone exit, the result is published as its own locked artifact.**

**Every exit gate, M0 to M7, also needs a cross-tier fidelity review** (D-055):
- Opus, Sonnet and Haiku reviewers, blind to each other, audit the spec **and the build** against
  the [vision ledger](10-vision-ledger.md).
- The milestone exits only if no VR comes back `missing` or `contradicted`, unless the owner has
  signed off on a deviation.
- Reviews are saved in `reviews/`.

| M | Name | Work | Exit gate |
|---|---|---|---|
| **M0** | Freeze, spike, pre-register | Spec review rounds → spec 1.0. **Spike:** fill in the surface × capability matrix (exec bit, architectures, subagents, headless browser, usage logs) and choose the eval harness (Q11). Draft schemas and `config/defaults.toml`. **Pre-register:** seal the eval briefs, rubrics, thresholds, held-out list, the 3 tasks per specialist, and the SME fallback per domain and per skill. Recruit SMEs. Resolve Q5 (judge smoke test). (Q11 and Q12 are resolved: D-051, D-047.) | Spec 1.0 frozen. Capability matrix recorded as decisions. `evals/sealed/` committed. |
| **M1** | Kernel | J1, G1, G2, G6, S1, O1, plus `internal/sheaf`: the Go port of the sheaf-core graph subset (D-038) with golden-fixture parity in `meta check`, using the per-component cut. **`meta check`** (a local command agents run in seconds; D-052): atom count = 34 across the files, the index and the README; every atom diff cites a D-###; SKILL.md ≤ 500 words; prompt lint; `runs/` gitignored; `go test` including sheaf parity. Remote CI runs the same command afterwards, as a witness. | Tests pass on fixtures recorded from real gather runs on **dev briefs**. No mocks. |
| **M2** | Intake + domain + **walking skeleton** | I1, I2, I3, G5, D1, D2, D3, L2, plus **minimal A1, A2 and R1** | The dev briefs produce sourced region graphs that pass G1/G2, and drift is measured. **Walking skeleton (D-058):** one dev brief runs Frame → Map → Make → Check end to end, with the See lane off: the v2−lens path, where L3 runs as its plain merge (consistent with "L3 always runs"), not a stub. It produces a real rendered artifact, published for the owner to look at. |
| **M3** | Lens | L1, L3 (with union-find coherence), L4, G3 (all three layers), G4 | Every operator can be switched on or off by flag, and the v2−lens arm runs. G3/G4 reports are produced on the dev briefs. G3 numeric layers match the oracle within tolerance. |
| **M4** | Artifact, review, record | A1, A2, R1 completed; G7, R2, S2 | End-to-end runs on the dev briefs pass G6/G7 and the stage-5 coverage checks. **Plug-and-play test (P14):** a non-technical person installs the `meta-suite-next` plugin from the released zip on the Claude app and gets a first artifact from meta-lens with zero help; time and any confusion are recorded. Then draw the held-out domain. |
| **M5** | Pilot eval (Phase A) | [Eval](05-eval.md) Phase A, then the power calculation | Go/no-go recorded. On a no-go, fix the weakest atom and rerun Phase A **once**. |
| **M5b** | Specialists | The specialists on the kernel, in this build order (D-065, D-071): meta-review, meta-prompting (P1, P2), meta-design (DS1), meta-planning (PL1, PL2), meta-mvp (MV1; after design and planning), meta-craft (CR1), meta-compound (CP1), meta-context (CX1). meta-research is rebuilt as an alias of meta-lens in brief form (D-072). Each skill's deep and `--fast` modes. | Each skill runs end to end on 2 dev tasks in both modes and passes its own gates. Artifact-producing skills meet rubric L2 on every applicable dimension. **Plug-and-play test (P14)** repeated for one specialist (meta-mvp). |
| **M6** | Full eval + suite release | Phases B and C, plus the per-skill comparisons. Release the **v2 suite** (`meta-suite-next` becomes `meta-suite` 2.0, D-071) if the meta-lens ship criteria hold; each specialist ships as v2 or falls back to v1 under its per-skill rule. | Results page published, and v2.0 tagged or explicitly declined. |
| **M7** | Watch & alerts (D-074) | Still v2, before any v3 work: re-run a saved brief on a cadence, detect changed facts and gate results from the ledger and saved run state, notify the user. Its atom budget is set by a decision at M7 entry. | A watched brief re-runs unattended, a planted change in a source raises exactly one alert, and an unchanged run raises none. |

## Working rules

- **One atom per PR.** Each PR links its atom file, shows its check passing, and cites the D-###
  entry for any change to the spec.
- **Feedback during the build** is recorded as decision entries in the repo, not left in chat.
