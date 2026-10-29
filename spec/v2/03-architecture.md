# 03 · Architecture

## Layers

```
┌────────────────────────────────────────────────────────────┐
│ Skill (L0–L3 markdown)  meta-lens/SKILL.md → index → ops/   │  what the model reads
├────────────────────────────────────────────────────────────┤
│ Operators (LLM work)    intake · domain · lens · artifact   │  host model + subagents
├────────────────────────────────────────────────────────────┤
│ Kernel (Go binary)      schemas · gates · ledger · loop     │  bookkeeping + verdicts
├────────────────────────────────────────────────────────────┤
│ State (files)           runs/ · playbook/                   │  user workspace, never committed
└────────────────────────────────────────────────────────────┘
```

**The model does judgment; the kernel does bookkeeping and verdicts.**
- The kernel holds **no model credentials** in v2.0. Any LLM call it needs (a judgment, or a sampled
  answer) is issued as a **task file**.
- The host model or a subagent executes the task and records the result with
  `meta sample record`. The kernel then aggregates the results ([J1](atoms/J1-llm-sample.md)).

**Two kinds of gate:**

| Kind | Atoms | Verdict |
|---|---|---|
| Deterministic | G1, G2, G6, G3 layers 1–2 (numeric sheaf, relational cycles) | Pure computation, reproducible from run files alone. |
| Kernel-aggregated | G3 pairwise + layer 3 (judgment balance), G4, G5, R2 | Computed over recorded LLM judgments. Reproducible given the judgment files. |

**Canonical format:** every typed object is **JSON** (validated by G1). XML is only a rendering of a
contract inside prompts.

## Pipeline (one run)

| Lane | Stage | Atoms | What happens | Exit check |
|---|---|---|---|---|
| **Frame** | 0 intake | (kernel check) → I1 → I2 → I3 | `meta version` probe. Classify the brief, expand it into a spec, round-trip the spec. | G5 drift ≤ threshold, at most `re_expand_max` re-expands |
| **Map** | 1 decompose | D1 | Question tree, overlapping region cover, relation vocabulary | G1, cover overlaps connected |
| **Map** | 2 gather | D2 (parallel) | One subagent per region. Claims are written via `meta claim add`. | G1 at write time, G2 |
| **Map** | 3 assemble | D3 | Per-region subgraphs plus candidate interfaces | G1, G2 |
| **See** | 4 lens | operators for the question type(s); J1 | Lens operators run. G3 and G4 **emit reports only.** | reports written |
| **Make** | 5 artifact | A1 → A2 | Choose form and layout, compose, render | G6, G7, plus coverage: every obstruction and flagged node maps to a section tagged contested or uncertain in `section_map` |
| **Check** | 6 review | R1 (all tiers), R2 (deep tier and eval) | Blind critic, plus expert-delta in deep runs | no P1/P2 edits, or pass cap reached |
| **Check** | 7 record | S1, S2 | Ledger and playbook observations | always runs |

**Loop rule:** an edit script sends the run back to the **earliest stage it touches**. A return to
D2 is allowed at most `d2_reentry_max` times per run (default 1).

**Operator selection:** a brief can carry 1 or 2 question types. The operator set is the union.

| Question type | Operators |
|---|---|
| `profile` | L1 + L2 |
| `compare` | L4 + L2 (+ L1 per side) |
| `synthesize` | L3 + L2 |

**L3 always runs**, because it is the only atom that merges per-region subgraphs. For `profile` and
`compare` it merges along interfaces only.

## Budget tiers

All values are **defaults** in `config/defaults.toml`. The config file is hashed, and the hash is
recorded in every run.

| Tier | Review passes | Max regions | k (judge samples, odd) | G3 pair cap | R2 | p50 / p90 wall time target |
|---|---|---|---|---|---|---|
| quick | 1 | 3 | 1 | 50 | off | 3 / 6 min |
| standard | 2 | 6 | 3 | 200 | off | 8 / 15 min |
| deep | 3 | 12 | 5 | 600 | on | 20 / 40 min |

The wall-time targets are initial estimates, recalibrated from Phase A logs.

## Cost model (per run)

Each stage's cost is estimated as calls × k × caps:

```
intake     ≈ I1 (k) + I2 (1) + I3 (1) + G5 (k × semantic_fields) × (1 + re_expands)
gather     ≈ regions × gather_calls_per_region          # cheap tier model
lens       ≈ L1: probes × k  +  L4: squares × k  +  G3: min(pairs, pair_cap) × k
             + G4: internal_nodes × (composed + collapsed) × k
artifact   ≈ A1 (1) + A2 (1 + render)
review     ≈ passes × R1 (1) + [deep] R2 (baseline 1 + min(findings, r2_cap) × k)
```

Every token figure is taken from **harness or API usage logs**, never from the model's own report.
Where a surface exposes no usage data, the run records call counts and is labelled
`cost=estimated`.

## Core types (JSON Schema in `kernel/schemas/`)

| Type | Fields |
|---|---|
| `Brief` | `raw`, `question_types[1..2]`, `objective`, `audience`, `constraints[]`, `targets[]` (L1 subjects), `views?[2]` (L4), `medium?`, `stakes`, `attachments[]`, `assumptions[]`, `vocab_allowlist[]` |
| `Spec` | `objective`, `slots[]` (name, type, required, origin ∈ {stated, inferred, default}), `constraints[]` (+origin), `success_criteria[]` (each → a finding kind or gate), `hash` |
| `Cover` | `universe[]`, `tree` (nodes: id, question, answer_type, children[]), `regions[]` (id, element ids ⊆ universe), computed `overlaps[]` and `triple_overlaps[]`, `relation_vocab[]` (each: name, `value_type` ∈ {qualitative, numeric, relative}, unit dimension?) |
| `Claim` | `id` (kernel ULID), `subject`, `relation`, `object`, `value?`, `value_type`, `unit?` (normalized by the kernel to the canonical unit in config), `as_of?`, `source` (locator, retrieved_at, independence_key), `region`, `stage_token` (kernel-issued: run, region, agent, model, prompt_hash) |
| `RegionGraph` / `DomainGraph` | `objects[]` (id, kind, labels), `morphisms[]` (id, src, dst, relation_type, claim_ids[]), `shape_metrics` (hub_concentration, longest_chain, two_sidedness, temporal_density) |
| `Interface` | `regions[2]`, `identified_objects[]`, `identified_relations[]`, `declared_by` |
| `ConflictEdge` | `claim_a`, `claim_b`, `block_key`, `judgment_ids[]`, `votes` (per verdict), `verdict` |
| `G3Report` | per key/component: `h0`, `discord` (per edge, per region), `harmonic_share` + `harmonic_null`, `balance` (balanced/frustrated + section), `near_zero_ambiguous`, `judged_pairs`, `candidate_pairs` |
| `Finding` | `id`, `kind` ∈ {verified, single_source, exception, tension, obstruction, profile_trait, coarse_model, uncertain}, `sides?[2]` (each: label, claim_ids[]), `evidence[]` (claim ids), `region?`, `confidence` ∈ {settled, single_source, contested, uncertain}, `plain_text` |
| `Judgment` | `id`, `task` (judge or answer), `prompt_hash`, `model`, `temperature`, `producer_model`, `output`, `blind` (aggregates expose vote counts, not only the majority) |
| `DriftReport` | `fields[]` (name, match, judgment_ids[]), `drift`, `verdict` |
| `DeltaReport` | `baseline_hash`, `delta_finding_ids[]`, `count`, `verdict` |
| `ArtifactPlan` | `form`, `layout`, `host_type?`, `rule_id`, `rationale` |
| `SectionMap` | `sections[]` (id, tag ∈ {main, contested, uncertain, single_source, how_checked}, finding_ids[]) |
| `EditScript` | `edits[]` (target_stage, target_id, op, rationale, priority P1–P3), `reviewer_model`, `blind` |
| `RunRecord` | `run_id`, `brief_hash`, `tier`, `surface`, `kernel` ∈ {present, absent}, `model_map`, `prompt_hashes{}`, `config_hash`, `schema_version`, `playbook_hash`, `gate_results[]`, `usage` (source ∈ {api, harness, estimated}), `wall_ms`, `seed` |
| `PlaybookEntry` | `id`, `rule`, `trigger`, `evidence_runs[]`, `elevation` ∈ {1_observation, 2_field_note} (v2.0) |

## Kernel CLI (`meta`) → atoms

| Command | Atom |
|---|---|
| `meta version` | surface probe (stage 0) |
| `meta validate <file>` | G1 |
| `meta claim add` (stdin JSON; returns ULID or a schema error) | D2, G1 |
| `meta stage begin <stage> [--region r]` → stage token (env) | O1 |
| `meta sample issue \| record \| aggregate` | J1 |
| `meta gate G2\|G3\|G4\|G5\|G6\|G7 <run_dir>` | gates |
| `meta graph build \| shape <run_dir>` | D3 |
| `meta interface declare \| merge <run_dir>` | L3 (refuses a merge without an interface) |
| `meta delta <run_dir>` | R2 aggregation |
| `meta ledger append \| show <run_dir>` | S1 |
| `meta playbook propose \| list` | S2 |
| `meta loop next <run_dir>` | O1 |
| `meta sheaf analyze <in.json> --out <report.json>` | the standalone gluing check (D-048): the same engine as G3 layers 1–3, usable on any JSON sheaf without running the pipeline |

**Concurrency:**
- Each gather agent appends to its own shard, `claims/<agent>.jsonl`.
- The kernel merges the shards at the stage boundary.
- The ledger is append-only under a file lock.

## Surfaces and capabilities

The capability matrix is resolved in the M0 spike:

| Capability | Claude Code | Claude app (uploaded skill) | Cowork | non-Claude hosts |
|---|---|---|---|---|
| exec binary | **verified** in the cloud sandbox (M0 spike), including the noexec fallback | **to verify** with `meta-probe` | expected | no |
| subagents (blind contexts) | yes | to verify | yes | varies |
| network for gather | yes | yes | yes | varies |
| headless browser (G7) | **verified**: Chrome for Testing + Playwright in the cloud sandbox | to verify (`meta-probe` step 6) | to verify | no |
| cross-family judge | via eval harness only | no | no | — |
| eval harness | `claude --bare -p … --output-format json`. Tokens are taken from `modelUsage` (top-level `usage` excludes subagents); dollars are estimates unless reconciled | — | — | — |

**Degraded mode:**
- **Binary missing or not executable:** the run is labelled `kernel=absent`. Every gate is reported
  `unchecked` in the artifact and the ledger. **The model never simulates a gate verdict.**
- **No subagents:** the blind attestation is recorded as `blind=false`.
- **Platform builds:** v2 ships static binaries for linux-amd64, linux-arm64, darwin-arm64,
  darwin-amd64 and windows-amd64. M0 measured them at about 3.2 MB each, about 6.7 MB zipped for all
  five.
- **Invocation is always `sh bin/meta …`**, never `./bin/meta`. Python's zip extraction strips the
  exec bit, and the launcher restores it.
- **Noexec fallback.** If the skill folder is mounted noexec, the launcher copies the binary to
  `$TMPDIR` and runs it from there. This was verified on a noexec tmpfs.
- **Windows without a POSIX shell** gets a `meta.cmd` launcher.
- **Deterministic output.** Every kernel command sorts its output (for example, errors by instance
  path, then keyword) so the same input gives byte-identical output.
- **Distribution:** v2 is excluded from the non-Claude bundles in `dist/` until v2.x.

## Repository layout (v2 additions)

```
kernel/                       Go module: cmd/meta, internal/{schema,gates,ledger,loop,sample,graph,sheaf}
  schemas/*.json              core types (versioned)
  config/defaults.toml        every threshold, tier value, rule table and canonical-unit table (hashed per run)
  internal/sheaf/             Go port of the sheaf-core graph subset, MIT (D-038, D-047); CI parity vs golden reports
  testdata/sheaf/             golden fixtures + oracle reports (`frontier-sheaf/v1`)
skills/meta-lens/
  SKILL.md                    L0 trailhead, ≤500 words; step 0 = kernel probe
  index.md                    L1: question types → operators → files
  ops/                        L2: one operator spec per file (program + JSON contract, XML-rendered)
  reference/                  L3: plain-language explainers, artifact patterns
  bin/                        built `meta` binaries (5 platforms)
evals/
  dev-briefs/                 development only; never used in eval
  sealed/                     SHA-256 of eval briefs, rubrics, held-out domain list + draw procedure
  report/
spec/v2/                      this spec
```
