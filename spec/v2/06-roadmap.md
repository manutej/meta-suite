# 06 · Roadmap

Each milestone has an exit gate. Nothing in a milestone starts until the previous milestone has
exited. **At every milestone exit, the result is published as its own locked artifact.**

| M | Name | Work | Exit gate |
|---|---|---|---|
| **M0** | Freeze, spike, pre-register | Spec review rounds → spec 1.0. **Spike:** fill in the surface × capability matrix (exec bit, architectures, subagents, headless browser, usage logs) and choose the eval harness (Q11). Draft schemas and `config/defaults.toml`. **Pre-register:** seal the eval briefs, rubrics, thresholds and held-out list. Recruit SMEs. **Decide Q12** (licence for the sheaf port). | Spec 1.0 frozen. Capability matrix recorded as decisions. `evals/sealed/` committed. |
| **M1** | Kernel | J1, G1, G2, G6, S1, O1, plus `internal/sheaf`: the Go port of the sheaf-core graph subset (D-038) with golden-fixture parity in CI, using the per-component cut. **`meta check`** (a local command agents run in seconds; D-052): atom count = 25 across the files, the index and the README; every atom diff cites a D-###; SKILL.md ≤ 500 words; prompt lint; `runs/` gitignored; `go test` including sheaf parity. Remote CI runs the same command afterwards, as a witness. | Tests pass on fixtures recorded from real gather runs on **dev briefs**. No mocks. |
| **M2** | Intake + domain | I1, I2, I3, G5, D1, D2, D3, L2 | The dev briefs produce sourced region graphs that pass G1/G2, and drift is measured. |
| **M3** | Lens | L1, L3 (with union-find coherence), L4, G3 (all three layers), G4 | Every operator can be switched on or off by flag, and the v2−lens arm runs. G3/G4 reports are produced on the dev briefs. G3 numeric layers match the oracle within tolerance. |
| **M4** | Artifact, review, record | A1, A2, G7, R1, R2, S2 | End-to-end runs on the dev briefs pass G6/G7 and the stage-5 coverage checks. **Plug-and-play test (P14):** a non-technical person installs the plugin from the released zip on the Claude app and gets a first artifact with zero help; time and any confusion are recorded. Then draw the held-out domain. |
| **M5** | Pilot eval (Phase A) | [Eval](05-eval.md) Phase A, then the power calculation | Go/no-go recorded. On a no-go, fix the weakest atom and rerun Phase A **once**. |
| **M6** | Full eval + release | Phases B and C. Release if every ship criterion holds. | Results page published, and v2.0 tagged or explicitly declined. |

## Working rules

- **One atom per PR.** Each PR links its atom file, shows its check passing, and cites the D-###
  entry for any change to the spec.
- **Feedback during the build** is recorded as decision entries in the repo, not left in chat.
