# 06 · Roadmap

Each milestone has an exit gate. Nothing in a milestone starts until the previous milestone has
exited. **At every milestone exit, the result is published as its own locked artifact.**

**Every exit gate also needs a cross-tier fidelity review** (D-055):
- Opus, Sonnet and Haiku reviewers, blind to each other, audit the spec **and the build** against
  the [vision ledger](10-vision-ledger.md).
- The milestone exits only if no VR comes back `missing` or `contradicted`, unless the owner has
  signed off on a deviation.
- Reviews are saved in `reviews/`.

| M | Name | Work | Exit gate |
|---|---|---|---|
| **M0** | Freeze, spike, pre-register | Spec review rounds → spec 1.0. **Spike:** fill in the surface × capability matrix (exec bit, architectures, subagents, headless browser, usage logs) and choose the eval harness (Q11). Draft schemas and `config/defaults.toml`. **Pre-register:** seal the eval briefs, rubrics, thresholds and held-out list. Recruit SMEs. **Decide Q12** (licence for the sheaf port). | Spec 1.0 frozen. Capability matrix recorded as decisions. `evals/sealed/` committed. |
| **M1** | Kernel | J1, G1, G2, G6, S1, O1, plus `internal/sheaf`: the Go port of the sheaf-core graph subset (D-038) with golden-fixture parity in CI, using the per-component cut. **`meta check`** (a local command agents run in seconds; D-052): atom count = 25 across the files, the index and the README; every atom diff cites a D-###; SKILL.md ≤ 500 words; prompt lint; `runs/` gitignored; `go test` including sheaf parity. Remote CI runs the same command afterwards, as a witness. | Tests pass on fixtures recorded from real gather runs on **dev briefs**. No mocks. |
| **M2** | Intake + domain + **walking skeleton** | I1, I2, I3, G5, D1, D2, D3, L2, plus **minimal A1, A2 and R1** | The dev briefs produce sourced region graphs that pass G1/G2, and drift is measured. **Walking skeleton (D-058):** one dev brief runs Frame → Map → Make → Check end to end, with the See lane off (that is the v2−lens path, not a stub). It produces a real rendered artifact, published for the owner to look at. |
| **M3** | Lens | L1, L3 (with union-find coherence), L4, G3 (all three layers), G4 | Every operator can be switched on or off by flag, and the v2−lens arm runs. G3/G4 reports are produced on the dev briefs. G3 numeric layers match the oracle within tolerance. |
| **M4** | Artifact, review, record | A1, A2, R1 completed; G7, R2, S2 | End-to-end runs on the dev briefs pass G6/G7 and the stage-5 coverage checks. **Plug-and-play test (P14):** a non-technical person installs the plugin from the released zip on the Claude app and gets a first artifact with zero help; time and any confusion are recorded. Then draw the held-out domain. |
| **M4b** | Specialists | The nine skills on the kernel, in this order (D-062): meta-research, meta-review, meta-prompting (P1, P2), meta-mvp (MV1), meta-design (DS1), meta-planning (PL1, PL2), meta-craft (CR1), meta-compound (CP1), meta-context (CX1). Each skill's deep and `--fast` modes. | Each skill runs end to end on 2 dev tasks in both modes and passes its own gates. The rubric minimum bar is met for the artifact-producing skills. |
| **M5** | Pilot eval (Phase A) | [Eval](05-eval.md) Phase A, then the power calculation | Go/no-go recorded. On a no-go, fix the weakest atom and rerun Phase A **once**. |
| **M6** | Full eval + release | Phases B and C, plus the per-skill comparisons. Release the **v2 suite** if the meta-lens ship criteria hold; each specialist ships as v2 or falls back to v1 under its per-skill rule. | Results page published, and v2.0 tagged or explicitly declined. |

## Working rules

- **One atom per PR.** Each PR links its atom file, shows its check passing, and cites the D-###
  entry for any change to the spec.
- **Feedback during the build** is recorded as decision entries in the repo, not left in chat.
