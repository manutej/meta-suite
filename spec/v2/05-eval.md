# 05 · Evaluation (pre-registered)

**Goal:** prove, not assume, three things:
1. v2 beats v1.
2. The lens is what makes the difference.
3. The results hold across models.

Everything on this page is fixed **before** any v2 run. Any change after that is a logged decision
and gets flagged in the report.

## Arms

| Arm | What runs |
|---|---|
| **v1** | meta-suite 1.2.0 end-to-end path (`meta-mvp`, see open question Q8) |
| **v2** | `meta-lens` with all 25 atoms |
| **v2−lens** | v2 with L1–L4, G3 and G4 off. Everything else stays: intake, gather, gates G1/G2/G5/G6/G7, review, artifact. |
| **baseline** | the same model with the same tools, given a plain request |

v2 − v2−lens is the lens's contribution. v2 − baseline is the harness's contribution.

## Briefs

- **12 briefs:** 3 question types × 4 domains.
- **Domains:** three named in advance, plus one **held-out** domain drawn at random from a sealed
  list only after the spec is frozen. This guards against tuning to the domains you already know.
- **All domains are public.** No client material (P13).
- Each brief is written as a **rough query**, the way real users ask, not as a polished prompt.

## Models

| Phase | Models | Arms | Seeds |
|---|---|---|---|
| A · pilot | Opus 5.5 | all 4 | 3 (4 briefs, one per domain) |
| B · full | Opus 5.5, Sonnet 5.5 | all 4 | 3 (all 12 briefs) |
| C · portability | Haiku 5.5 + one open-weight model (OpenRouter) | v1, v2 | 1 (all 12 briefs) |

Phase A must pass its go/no-go before Phase B is run.

## Metrics

| ID | Metric | How |
|---|---|---|
| M1 | **Rubric: recall + analysis** (primary) | Binary rubric items written per brief **before** any run, by a separate agent from independent sources. Items are tagged recall, analysis or presentation. Score = % of recall + analysis items met. |
| M2 | Blind pairwise preference | A judge from a different model family, position-swapped, comparing v2 vs v1 and v2 vs v2−lens |
| M3 | SME panel | Human experts, at least 2 briefs per domain, blind to arm. They score depth, correctness and "would I learn something?" |
| M4 | Operadic consistency | from G4 |
| M5 | Contested-fact precision | share of G3 obstructions that a judge and an SME confirm are genuinely contested |
| M6 | Cost | tokens, wall time and dollars per run |
| M7 | Stability | variance across seeds for M1 and M2 |

## Ship criteria for v2.0 (all must hold)

1. **M1:** v2 − v1 ≥ +10 points mean, with a bootstrap 95% CI that excludes 0.
2. **M2:** v2 wins against v1 in at least 60% of comparisons.
3. **The lens is proven:** v2 − v2−lens on M1 has a CI that excludes 0, on Opus and on Sonnet.
4. **No regression:** no model in Phase C shows v2 more than 5 points below v1 on M1.
5. **M6:** v2's median tokens per run ≤ v1's median.
6. **M3:** the SME mean score for v2 > v1.

These thresholds are initial and are recorded in [decisions](09-decisions.md). They can change
before the freeze, and never afterwards.

## Outputs

- `evals/report/`: the tables, per-model breakdowns, and the ablation of each operator. L1–L4 are
  each switched off individually in Phase B, on Opus only.
- **One published results page per phase,** each a separate, locked artifact.
- **The report states the result even if v2 loses.** A negative result is itself a finding.
