# 03 · Architecture

## Layers

```
┌────────────────────────────────────────────────────────────┐
│ Skill (L0–L3 markdown)  meta-lens/SKILL.md → index → ops/   │  what the model reads
├────────────────────────────────────────────────────────────┤
│ Operators (LLM work)    intake · domain · lens · artifact   │  model + subagents
├────────────────────────────────────────────────────────────┤
│ Kernel (Go binary)      schemas · gates · ledger · loop     │  deterministic, fast
├────────────────────────────────────────────────────────────┤
│ State (files)           runs/ · playbook/ · residue/        │  survives sessions
└────────────────────────────────────────────────────────────┘
```

**The model does judgment; the kernel does bookkeeping and verdicts.**
- The model never computes a gate result. It calls `meta gate <id>` and reads the verdict.
- The kernel never makes a domain judgment.

## Pipeline (one run)

| Stage | Atoms | What happens | Exit check |
|---|---|---|---|
| 0 intake | I1 → I2 → I3 | Classify the brief, expand it into a typed spec, round-trip the spec back to intent. | G5 drift ≤ threshold (≤ 2 re-expands) |
| 1 decompose | D1 | Build the decomposition tree and the region cover. | G1 |
| 2 gather | D2 (parallel) | One subagent per region writes claims through typed tools. | G1, G2 |
| 3 assemble | D3 | Build the domain graph from the claims. | G1, G2 (no orphans, every morphism sourced) |
| 4 lens | L1–L4 per question type | Run the operators that match the question type. | G3 sheaf, G4 operadic |
| 5 artifact | A1 → A2 | Choose the form, map domain structure to artifact structure, render. | G6 vocab, G7 render |
| 6 review | R1, R2 | Blind critic plus expert-delta check. The output is an edit script. | edit script empty, or pass cap reached |
| 7 residue | S1, S2, S3 | Write the ledger, playbook candidates and domain residue. | always runs |

**Loop rule:** an edit script sends the run back to the **earliest stage it touches**, never to the
start. The number of passes is capped by the budget tier.

**Operator selection by question type:**

| Question type | Operators |
|---|---|
| `profile` | L1 Yoneda + L2 pullback |
| `compare` | L4 naturality + L2 pullback (+ L1 per side) |
| `synthesize` | L3 glue + L2 pullback |

**Budget tiers** (values in `config/defaults.toml`, not in code):

| Tier | Review passes | Max regions | Consistency samples |
|---|---|---|---|
| quick | 1 | 3 | 2 |
| standard | 2 | 6 | 3 |
| deep | 3 | 12 | 5 |

## Core types (JSON Schema in `kernel/schemas/`)

| Type | Key fields |
|---|---|
| `Brief` | `raw`, `question_type`, `domain`, `audience`, `medium?`, `stakes`, `assumptions[]` |
| `Spec` | `objective`, `slots[]` (name, type, required), `constraints[]`, `success_criteria[]`, `hash` |
| `Cover` | `tree` (question decomposition), `regions[]` (id, scope, sub_questions[]) |
| `Claim` | `id`, `subject`, `relation`, `object`, `value?`, `source` (locator, retrieved_at, independence_key), `region`, `agent`, `model`, `prompt_hash` |
| `DomainGraph` | `objects[]` (id, kind, labels), `morphisms[]` (id, src, dst, relation_type, claim_ids[]), `relation_vocab[]` |
| `Finding` | `id`, `kind` ∈ {verified, single_source, tension, obstruction, profile_trait, gap}, `evidence[]`, `plain_text` |
| `EditScript` | `edits[]` (target_stage, target_id, op, rationale), `reviewer`, `blind: true` |
| `RunRecord` | `run_id`, `brief_hash`, `tier`, `model_map`, `prompt_hashes{}`, `gate_results[]`, `tokens`, `wall_ms`, `seed` |
| `PlaybookEntry` | `id`, `rule`, `trigger`, `evidence_runs[]`, `elevation` ∈ {1_observation … 4_canon}, `ablation?` |

A source's `independence_key` (publisher or origin) is what lets [L2](atoms/L2-lens-pullback.md)
tell two independent paths apart from one source cited twice.

## Kernel: the `meta` binary (Go)

```
meta validate <file>                 # G1 against kernel/schemas
meta gate <G2..G7> <run_dir>         # deterministic verdict + report JSON
meta consistency <run_dir> --k N     # G4 metric
meta ledger append|show <run_dir>    # S1
meta playbook propose|promote        # S2 (promote requires ablation result)
meta loop next <run_dir>             # O1: which stage runs next, given state + edits
```

- **Why Go:** fast concurrent I/O and one static binary per platform, so the skill needs no Python
  environment. It is also the same language as the Stitch core.
- **Honest note on speed:** LLM latency dominates wall time, so most of the speed comes from
  elsewhere:
  - parallel gather
  - prompt caching (stable prefix first)
  - batched programmatic tool calls
  - a small-model gather tier with top-tier judges
  - gates that exit early

## Repository layout (v2 additions)

```
kernel/                       Go module: cmd/meta, internal/{schema,gates,ledger,loop,consistency,sheaf}
  schemas/*.json              core types
  config/defaults.toml        every threshold and tier value
skills/meta-lens/
  SKILL.md                    L0 trailhead, ≤500 words
  index.md                    L1: question types → operators → atom files
  ops/                        L2: one operator spec per file (procedure-as-program + XML contract)
  reference/                  L3: plain-language explainers, artifact patterns
  bin/                        built `meta` binaries (linux-amd64, linux-arm64, darwin-arm64)
evals/
  briefs/  rubrics/  runs/  report/
spec/v2/                      this spec
```

## Prompt contract format

Each operator file in `ops/` has:
- an XML contract (typed slots in, typed slots out, gates that must pass)
- a procedure written as a short program or pseudocode the model executes step by step through
  `meta` calls

The XML is the typed schema the loop reads and writes. The program is the procedure. Neither
contains worked examples of domain content.
