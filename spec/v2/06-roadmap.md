# 06 · Roadmap

Each milestone has an exit gate. Nothing in a milestone starts until the previous milestone has
exited. **At every milestone exit, the result is published as its own locked artifact.**

| M | Name | Work | Exit gate |
|---|---|---|---|
| **M0** | Freeze, spike, pre-register | Spec review rounds → spec 1.0. **Spike:** fill in the surface × capability matrix (exec bit, architectures, subagents, headless browser, usage logs) and choose the eval harness (Q11). Draft schemas and `config/defaults.toml`. **Pre-register:** seal the eval briefs, rubrics, thresholds and held-out list. Recruit SMEs. | Spec 1.0 frozen. Capability matrix recorded as decisions. `evals/sealed/` committed. |
| **M1** | Kernel | J1, G1, G2, G6, S1, O1. CI checks: atom count = 25 across the files, the index and the README; every atom diff cites a D-###; SKILL.md ≤ 500 words; prompt lint; `runs/` gitignored. | Tests pass on fixtures recorded from real gather runs on **dev briefs**. No mocks. |
| **M2** | Intake + domain | I1, I2, I3, G5, D1, D2, D3, L2 | The dev briefs produce sourced region graphs that pass G1/G2, and drift is measured. |
| **M3** | Lens | L1, L3, L4, G3, G4 | Every operator can be switched on or off by flag, and the v2−lens arm runs. G3/G4 reports are produced on the dev briefs. |
| **M4** | Artifact, review, record | A1, A2, G7, R1, R2, S2 | End-to-end runs on the dev briefs pass G6/G7 and the stage-5 coverage checks. Then draw the held-out domain. |
| **M5** | Pilot eval (Phase A) | [Eval](05-eval.md) Phase A, then the power calculation | Go/no-go recorded. On a no-go, fix the weakest atom and rerun Phase A **once**. |
| **M6** | Full eval + release | Phases B and C. Release if every ship criterion holds. | Results page published, and v2.0 tagged or explicitly declined. |

## Working rules

- **One atom per PR.** Each PR links its atom file, shows its check passing, and cites the D-###
  entry for any change to the spec.
- **Feedback during the build** is recorded as decision entries in the repo, not left in chat.
