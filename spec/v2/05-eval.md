# 05 · Evaluation (pre-registered at M0)

**Goal:** prove three things, rather than assume them:
1. v2 beats v1.
2. The lens is what makes the difference.
3. The results hold across models.

The following are committed, by SHA-256, in `evals/sealed/` at **M0 exit**, before any v2 code runs
against them:
- the briefs and their rubrics
- the thresholds (the hashed `config/defaults.toml` eval section)
- the held-out domain list and the procedure for drawing from it

Any later change is logged as a decision and flagged in the report.

## Briefs: development is kept apart from evaluation

- **Development set.** `evals/dev-briefs/` is used through M2–M4. It shares no brief, and no
  domain, with the eval set.
- **Eval set.** 12 briefs, made of 3 question types × 4 domains.
  - Three of the domains are named at M0.
  - The fourth is a **held-out** domain, drawn after M4 from the sealed list by the committed
    procedure. It never appears in development.
- **Who writes them.** Eval briefs and rubrics are authored blind to the builder (see Rubrics).
- **Form.** Each brief is a rough query, the way users actually ask. All domains are public (P13).

## Arms

| Arm | Definition |
|---|---|
| **v1** | meta-suite 1.2.0 `meta-research` (decompose → cross-verify → cite), its native output |
| **v2** | `meta-lens`, standard tier |
| **v2−lens** | v2 with L1, L3's interface discipline, L4, G4 and **all three G3 layers** off (see below) |
| **v2-quick** | v2 at quick tier; a secondary row that roughly matches v1's cost |
| **v2-adaptive** | v2 with stage-boundary adaptation **on** (D-057); every other arm freezes it |
| **v1-mvp** | meta-suite 1.2.0 `meta-mvp`, native artifact. Used **only for M8**, so v2's artifact is compared with v1's most visible output |
| **baseline** | same model, same tools, plain request |

**What v2−lens keeps and loses:**
- **Kept:** L2, because cross-source verification belongs to the harness, not the lens.
- **Replaced:**
  - D1's tree becomes a flat list of sub-questions.
  - L3 becomes a plain merge, with D3's cross-region candidates merged by name.
  - Lens vocabulary is stripped from the ops prompts.
- **Compute-matched.** It runs extra review passes and self-consistency samples until its median
  tokens are within ±15% of v2's.

The two contrasts:
- **v2 − v2−lens** is the lens effect, at equal compute.
- **v2 − baseline** is the harness effect.

**Every eval run:**
- uses a fresh workspace
- has playbook writes off (S2)
- records `config_hash` and `playbook_hash`
- uses R2 (eval mode) in the v2 arms only to measure, never to route the loop

## Models and phases

| Phase | Models | Arms | Briefs × seeds |
|---|---|---|---|
| A · pilot | Opus 5.5 | all 5 | 4 (one per domain, from the eval set's first draw) × 3 |
| B · full | Opus 5.5, Sonnet 5.5 | all 5 | 12 × 3, or more if the power calculation requires it (below) |
| C · portability | Haiku 5.5 (+ an open-weight model only if Q11 resolves) | v1, v2 | 12 × 1 |

**Harness.** Every arm runs headless on the same Claude agent harness, with the same tools and the
`meta` binary present. Usage comes from API logs.

**Phase A go/no-go:**
- v2 − v1 on M1 is at least +5 points mean, **and**
- infrastructure failures (kernel errors, schema rejections the run never recovered from) occur in
  fewer than 10% of runs.

## Metrics

| ID | Metric | How |
|---|---|---|
| M1 | **Rubric: recall + analysis** (primary) | Item-by-item binary grading on **normalized plain text**: formatting, arm markers and the "how this was checked" section are stripped. The grader is a cross-family judge, blind to arm. SMEs audit a 20% sample, and judge–SME agreement (κ) is reported. |
| M2 | Blind pairwise preference | Cross-family judge, position-swapped, on normalized text: v2 vs v1, and v2 vs v2−lens |
| M3 | SME panel | At least 2 SMEs per domain, blind to arm. They score depth, correctness, and "did I learn something?" |
| M4 | Operadic consistency | From G4. Its correlation with M1 decides whether G4 is promoted. |
| M5 | Contested-fact precision | The share of G3 obstructions a judge and an SME confirm as genuinely contested, reported **separately** for numeric-layer obstructions (layers 1–2) and judgment-layer obstructions (pairwise + layer 3) |
| M6 | Cost | Dollars, tokens and wall time per run, from usage logs only. M1 points per dollar are also reported. |
| M7 | Stability | Between-seed variance of M1 and M2 |
| M8 | **Artifact first look** (VR-02) | **Rendered artifacts, not stripped**, compared side by side, blind to arm. A cross-family judge plus SMEs rate four things: "would forward to a client", "learned something in the first 60 seconds", visual hierarchy and clarity, and "looks like it was made by an expert". Compared against v1-mvp and v1's native artifacts. |
| M9 | Human input per run (VR-14) | Count of clarifying questions, interventions and retries a user needed, logged by O1 |
| M10 | Speed (VR-06) | p50 and p90 wall time and turns per run, from harness logs |

## Rubrics

**Authorship:**
- Drafted by an agent from a **different model family from every arm's producer model**, working
  from independent sources.
- SMEs validate a sample of items for the named domains and **every item** for the held-out domain.
- SME-authored items are scored and reported separately as well.

**Format:** each item is tagged recall, analysis or presentation. M1 counts recall and analysis
only.

## Statistics

- **The unit of independence is the brief.** Confidence intervals come from a **cluster bootstrap**:
  resample briefs, then seeds within each brief.
- **Power.** A power calculation using Phase A variance is committed before Phase B. If power for
  criterion 1 is below 0.8, Phase B grows to at least 24 briefs at 2 seeds, still drawn by the
  sealed procedure.
- **Criterion 3** is a pooled test across Opus and Sonnet. Per-model results are reported alongside
  it.

## Ship criteria for v2.0 (all must hold)

1. **M1:** v2 − v1 ≥ +10 points mean, with a cluster-bootstrap 95% CI that excludes 0.
2. **M2:** v2 is preferred over v1 in at least 60% of comparisons.
3. **Lens:** v2 − v2−lens on M1 is greater than 0, with a pooled CI that excludes 0.
4. **No regression:** no model shows v2 more than 5 points below v1 on M1.
5. **Speed:** v2-quick p50 wall time ≤ v1's. v2 standard p50 ≤ `eval.wall_multiplier_max` × v1's (committed at M0).
6. **Artifact:** v2 ≥ v1-mvp and ≥ v1 on M8 (pooled), with SMEs agreeing in direction.
7. **Cost:** v2's median dollars per run is at most `eval.cost_multiplier_max` × v1's (the ceiling
   is 3×; the exact value is set from Phase A and committed before Phase B). M1 points per dollar
   for v2-quick is also reported against v1.
8. **SME:** v2 beats v1 on M3 in a paired sign test. If fewer than 2 SMEs per domain could be
   recruited, M3 is reported as a non-blocking signal, but only with the owner's explicit sign-off, logged as a decision.

## The "works best with Opus 5.5" claim (VR-23)

- **v1 stays the default install** until v2.0 meets every ship criterion.
- "Works best with Opus 5.5" may be said only if v2 on Opus beats v2 on Sonnet **and** on Haiku on M1 and M8, with pooled CIs that exclude 0.
- Otherwise, the framing states the measured per-tier results ("strongest on X, good on Y"). No claim goes beyond the data.

## Outputs

- `evals/report/` contains:
  - the tables
  - per-model and per-question-type breakdowns
  - **individual operator ablations** for L1, L3 and L4, run in Phase B on Opus only
- A results page is published at the end of each milestone, M5 and M6 included. Each page is
  locked as its own artifact.
- **A loss is reported as plainly as a win.**
