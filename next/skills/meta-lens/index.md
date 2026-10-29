# meta-lens · index (L1)

Lanes → ops → when each file is loaded. `RUN` is the run folder printed by `meta run init`
(`runs/<run_id>/` under the working directory). The See lane is off in this build: Combine is a plain
merge plus independent-source verification.

| Lane | Op | File | Writes in `$RUN/` | Fast | Deep |
|---|---|---|---|---|---|
| Frame | I1 classify | [ops/I1-classify.md](ops/I1-classify.md) | `brief.json` | — | at step 1 |
| Frame | I2 expand | [ops/I2-expand.md](ops/I2-expand.md) | `spec.json` | — | at step 1 |
| Map | D1 cover | [ops/D1-cover.md](ops/D1-cover.md) | `cover.json`, `probes.json` | **core** | at step 2 |
| Map | D2 gather | [ops/D2-gather.md](ops/D2-gather.md) | `claims/*.jsonl`, `claims.jsonl`, `gates/G2.json` | **core** | at step 2 |
| Map | D3 graph | [ops/D3-graph.md](ops/D3-graph.md) | `graph/regions/*.json`, `graph/interface_candidates.json` | — | at step 2 |
| Combine | L3 merge | [ops/L3-merge.md](ops/L3-merge.md) | `graph/interfaces.json`, `graph/merged.json` | — | at step 3 |
| Combine | L2 verify | [ops/L2-verify.md](ops/L2-verify.md) | `findings.jsonl`, `conflicts.json` | — | at step 3 |
| Make | A1 select | [ops/A1-select.md](ops/A1-select.md) | `plan.json` | — | at step 4 |
| Make | A2 compose | [ops/A2-compose.md](ops/A2-compose.md) | `artifact/index.html`, `section_map.json`, `how_checked.json`, `gates/G6.json` | **core** | at step 4 |
| Check | R1 review | [ops/R1-review.md](ops/R1-review.md) | `review/edits.<n>.json` | **core** | at step 5 |

"—" in Fast: the op runs from SKILL.md's one-line step; open its file only after a repeated kernel
error.

## Run-folder files

| File | Type (`meta validate --type`) | Written by |
|---|---|---|
| `run.json` | RunRecord | kernel (`meta run init`) |
| `brief.json` | Brief | I1 |
| `spec.json` | Spec | I2 |
| `cover.json` | Cover | D1 |
| `claims/<agent>.jsonl` | Claim per line | kernel (`meta claim add`) |
| `graph/regions/<region>.json` | RegionGraph | kernel |
| `graph/merged.json` | DomainGraph | kernel |
| `findings.jsonl` | Finding per line | kernel |
| `plan.json` | ArtifactPlan | A1 |
| `section_map.json` | SectionMap | A2 |
| `review/edits.<n>.json` | EditScript | R1 |
| `prompts/*.md` | — (hashed by `--prompt-file`) | each op that briefs a subagent |
| `gates/*.json` | — (raw stdout of each gate) | every op that runs a gate |
| `assumptions.md` | — | any op; shown on the page |

## Kernel errors, common to every op

- **Launcher exit 126 or 127** (binary missing or won't start): switch the run to `kernel=absent` from
  this point. Earlier kernel results stay valid; every later check is `unchecked`.
- **JSON error with a field path** (schema): fix that field, retry once. A second failure follows the
  op's own table. Never edit kernel-written files by hand to get past an error.
- **Any other non-zero exit:** save stdout and stderr to `$RUN/gates/<op>.error.json`, then follow the
  op's table. Never write a result the kernel did not print.

## Edit routing (until the loop controller ships)

An R1 edit goes to the earliest op its `target_stage` names, and everything after it re-runs. D2 is
re-entered at most `loop.d2_reentry_max` times per run. Review passes stop at
`modes.<mode>.review_passes`.

**Caps** (region count, review passes, re-gather limit, independent-source minimum) are read from
`$ROOT/kernel/config/defaults.toml`, whose hash is in `run.json`.
<!-- <internal>Keys: modes.<mode>.max_regions, modes.<mode>.review_passes, loop.d2_reentry_max,
lens.pullback.min_independent.</internal> --> `<mode>` is `fast`, or
`deep_narrow` / `deep_broad` per the scope I2 recorded. If the file can't be read (for example a
single-skill upload without `kernel/`), use the spec's fast-mode defaults (3 regions, 1 review pass,
1 re-gather, 2 independent sources) and add "caps taken from spec defaults, config unreadable" to
`$RUN/assumptions.md`.
