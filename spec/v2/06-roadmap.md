# 06 · Roadmap

Each milestone has an exit gate. Nothing in a milestone starts until the previous milestone has
exited. Each milestone's result is published as its own locked artifact.

| M | Name | Atoms / work | Exit gate |
|---|---|---|---|
| **M0** | Freeze + spike | Spec review rounds → spec 1.0. Spike: does a Go `meta` binary run from a skill folder on Claude Code, the Claude app (uploaded skill) and Cowork (linux-amd64/arm64, darwin-arm64)? Is a headless browser available? Draft the schemas. | Spec 1.0 frozen. Spike results recorded as decisions. |
| **M1** | Kernel | G1, G2, G6, S1, O1. Schemas, `config/defaults.toml`, CI (word-count and prompt lint). | Tests pass on fixtures recorded from real gather runs. No mocks. |
| **M2** | Intake + domain | I1, I2, I3, G5, D1, D2, D3, L2 | 4 pilot briefs produce sourced domain graphs that pass G1/G2. Drift is measured. |
| **M3** | Lens | L1, L3, L4, G3, G4 | Each operator can be switched on or off by a flag. Per-axiom and consistency reports are produced on the 4 pilot briefs. |
| **M4** | Artifact + review + learning | A1, A2, G7, R1, R2, S2, S3 | End-to-end runs on the 4 pilot briefs, with every artifact passing G6/G7. |
| **M5** | Pilot eval (Phase A) | Run [eval](05-eval.md) Phase A | Go/no-go recorded. On a no-go, fix the weakest atom and rerun Phase A once. |
| **M6** | Full eval + release | Phases B and C. Release v2.0 if the ship criteria hold. | Results page published, and v2.0 tagged or explicitly declined. |

## Working rules

- **One atom, one PR.** Each PR links its atom file and shows its check passing.
- **Pre-registration is locked before M5.** Briefs, rubrics and thresholds are committed before any
  v2 eval run.
- **Changes are tracked in the repo.** Feedback during the build is recorded as entries in the
  [decision log](09-decisions.md), not left in chat.
