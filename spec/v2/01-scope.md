# 01 · Scope

## In v2.0

- **`meta-lens`, the suite's front door** (D-036, D-065). It turns a rough query into a verified
  research view or artifact, or routes it to the right specialist, through five lanes: **Frame → Map → See → Make → Check** (D-037). The user
  only types a question; the question type is inferred (I1), and the lanes run on their own.
- **Three research question types** (for the meta-lens and meta-research lanes; the specialists add their own task types). A brief can carry one or two of them; with two, the operator sets are
  combined.

  | Type | Asks | Core operators |
  |---|---|---|
  | `profile` | What is X, really? | Yoneda + pullback (+ glue along interfaces) |
  | `compare` | X vs Y, or view A vs view B | naturality + pullback (+ glue along interfaces) |
  | `synthesize` | What does the field add up to? | glue + pullback |

- **The whole new suite** ([suite map](11-suite.md)): `meta-lens` plus the nine rebuilt skills, on **34 atoms** (25 core + 9 specialist) in the [atom index](04-atoms.md), and nothing else.
- **Two modes for every skill:** deep (default) and `--fast` (D-066).
- **An artifact rubric** whose minimum bar is the owner's recent work ([rubric](12-artifact-rubric.md), D-067).
- **A shared Go kernel**, shipped as a single static `meta` binary for 5 platforms. It owns schemas,
  gates, judgment sampling, the ledger and the loop controller. It holds no model credentials.
- **Degraded mode** when the binary can't run: every gate is reported `unchecked`, and no verdict is
  ever simulated.
- **A progressive-disclosure layout** inside the skill folder (L0–L4).
- **A pre-registered evaluation:** v1 vs v2 vs v2-without-lens vs v2-quick vs a plain baseline,
  across models.
- **One results page per milestone**, each published as its own locked artifact.

## Not in v2.0

- **v2 in the non-Claude bundles** (`dist/` for ChatGPT, Grok and plain-prompt hosts).
- **Product and service layers:** Stitch integration, watch mode and alerting, a hosted service, a
  public API.
- **Swapping prompts mid-stage.** Adaptation happens at stage boundaries (D-057).
- **Rust in the shipped kernel.** Rust is used only as the test oracle (sheaf-core, D-038).
- **New question types** (`path`, `tradeoff`, `analogy`, `monitor`).
- **Custom UI** beyond the host's artifact types and the existing design skills.
- **`meta-score` integration**, until its source is verified.
- **Articles and marketing collateral.** These are a separate track that uses this spec as input.

## Parking lot

| Item | Why parked | Target |
|---|---|---|
| `path` question type + composition-gap operator (second-order effects) | not needed to prove the lens | v2.1 |
| Domain residue (formerly atom S3): saved domain graph reused as a prior | creates dependence between eval runs | v2.1 |
| Playbook promotion and loading (S2 levels 3–4) | needs ablations at eval scale | v2.1 |
| `tradeoff` question type + adjunction operator (beyond the intake round trip) | not needed to prove the lens | v2.1 |
| Minimal conflicting sets; three-way conflicts among **free-text** claims (linear relational cycles are in scope via G3 layer 2, D-044) | exponential judgment cost | v2.1 |
| `analogy` question type + functor operator | not needed to prove the lens | v2.2 |
| v2 in non-Claude bundles | the binary can't run there | v2.x |
| `meta-score` | source unverified | v2.x |
| Watch mode: re-run on a cadence, alert on changed facts | belongs with Stitch | v3 |
| Harness/middleware API for any agent | product layer | v3 |
| Learned orchestrator | research bet | unscheduled |

## Change control

- **Atom budget = 34.** `meta check` enforces it from M1 (D-052): it checks that the atom files, the index rows and
  the README all agree. Adding an atom requires a decision entry that removes, merges or defers
  another.
- **Every atom diff cites a D-### entry.** `meta check` enforces this from M1.
- **Parking-lot items enter scope only through a decision**, and only after M5.
- **"While we're at it" changes are rejected by default.**
