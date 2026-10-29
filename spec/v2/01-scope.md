# 01 · Scope

## In v2.0

- **One new pipeline skill**, working name `meta-lens` (see [open questions](08-open-questions.md)).
  It takes a rough query and produces a verified research view or artifact.
- **Three question types:**

  | Type | Asks | Core operators |
  |---|---|---|
  | `profile` | What is X, really? | Yoneda + pullback |
  | `compare` | X vs Y, or view A vs view B | naturality + pullback |
  | `synthesize` | What does the field add up to? | glue + pullback |

- **The 25 atoms** in the [atom index](04-atoms.md), and nothing more.
- **A shared Go kernel** shipped as a single static `meta` binary. It handles schemas, gates,
  ledger, consistency metrics and the loop controller.
- **A progressive-disclosure layout** inside the skill folder (L0–L4).
- **A pre-registered evaluation** of v1 vs v2 vs v2-without-lens vs plain baseline, across models.
- **A results page per milestone.** Each is published as its own locked artifact.

## Non-goals for v2.0

- Porting the nine v1 skills. They stay as they are until v2.0 passes eval.
- Stitch integration, watch/alerting, a hosted service, or a public API.
- Swapping prompts in the middle of a stage. Refactoring happens between stages only.
- Rust. New question types (`path`, `tradeoff`, `analogy`, `monitor`).
- Custom UI beyond the host's artifact types and the existing design skills.
- `meta-score` integration, pending source verification.
- Articles and marketing collateral. These are a separate track that uses this spec as input.

## Parking lot

| Item | Why parked | Target |
|---|---|---|
| `path` question type + composition-gap operator | not needed to prove the lens | v2.1 |
| `tradeoff` question type + adjunction operator (beyond the intake round-trip) | same | v2.1 |
| `analogy` question type + functor operator (preserved vs broken relations) | same | v2.2 |
| Port v1 skills onto the kernel | only after v2.0 eval passes | v2.x |
| `meta-score` (paraphrase invariance of meta-prompt atoms) | its source construction is unverified | v2.x |
| Watch mode: re-run on a cadence, alert on changed facts | belongs with Stitch | v3 |
| Harness/middleware API for any agent | product layer | v3 |
| Learned orchestrator (trained conductor) | research bet | unscheduled |

## Change control

- **Atom budget = 25.** Adding an atom requires a decision entry that removes, merges or defers
  another one.
- **Parking-lot items move into scope only through a decision entry**, and only after M5 (pilot).
- **"While we're at it" changes are rejected by default.**
