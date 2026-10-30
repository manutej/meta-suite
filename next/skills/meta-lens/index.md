# meta-lens · index (L1)

Lanes → ops → when each file is loaded. `RUN` is the run folder printed by `meta run init`
(`runs/<run_id>/` under the working directory; init also creates `prompts/`, `gates/`, `review/`,
`artifact/shots/`). The See lane is off in this build: Combine is a plain merge plus independent-source
verification. Every `meta <command> --help` carries that command's arguments, outputs and gotchas; the
ops below hold only contracts, order and judgment.

| Lane | Op | File | Writes in `$RUN/` | Fast | Deep |
|---|---|---|---|---|---|
| Frame | I1 classify | [ops/I1-classify.md](ops/I1-classify.md) | `brief.json` | — | at step 1 |
| Frame | I2 expand | [ops/I2-expand.md](ops/I2-expand.md) | `spec.json` | — | at step 1 |
| Map | D1 cover | [ops/D1-cover.md](ops/D1-cover.md) | `cover.json`, `probes.json` | **core** | at step 2 |
| Map | D2 gather | [ops/D2-gather.md](ops/D2-gather.md) | `claims/*.jsonl`, `claims.jsonl`, `gates/G2.json` | **core** | at step 2 |
| Map | D3 graph | [ops/D3-graph.md](ops/D3-graph.md) | `graph/regions/*.json`, `graph/interface_candidates.json` | — | at step 2 |
| Combine | L3 merge | [ops/L3-merge.md](ops/L3-merge.md) | `graph/interfaces.json`, `graph/merged.json` | — | at step 3 |
| Combine | L2 verify | [ops/L2-verify.md](ops/L2-verify.md) | `findings.jsonl`, `conflicts.json`, `coverage.json` | — | at step 3 |
| Make | A1 select | [ops/A1-select.md](ops/A1-select.md) | `plan.json` | — | at step 4 |
| Make | A2 compose | [ops/A2-compose.md](ops/A2-compose.md) | `artifact/index.html`, `section_map.json`, `how_checked.json`, `gates/G6.json` | **core** | at step 4 |
| Check | R1 review | [ops/R1-review.md](ops/R1-review.md) | `review/edits.<n>.json` | **core** | at step 5 |

"—" in Fast: the op runs from SKILL.md's one-line step; open its file only after a repeated kernel
error.

## Run-folder files

Every JSON file a model writes has a type; `meta validate <file>` infers it from the name.

| File | Type | Written by |
|---|---|---|
| `run.json` | RunRecord | kernel (`run init`; `validate` adds the spec hash) |
| `brief.json` | Brief | I1 |
| `spec.json` | Spec | I2 |
| `cover.json` | Cover | D1 |
| `probes.json` | Probes | D1 |
| `claims/<agent>.jsonl` | Claim per line | kernel (`meta claim add`) |
| `retractions.jsonl` | Retraction per line | kernel (`meta claim retract`) |
| `graph/regions/<region>.json` | RegionGraph | kernel |
| `graph/merged.json` | DomainGraph | kernel |
| `findings.jsonl`, `conflicts.json` | Finding per line, Conflicts | kernel |
| `coverage.json` | Coverage | L2 |
| `plan.json` | ArtifactPlan | A1 |
| `section_map.json` | SectionMap | A2 |
| `how_checked.json` | HowChecked | kernel (`meta report`) |
| `review/edits.<n>.json` | EditScript | R1 |
| `prompts/*.md` | — (hashed by `--prompt-file`) | each op that briefs an agent |
| `gates/*.json` | — (raw stdout of each gate) | every op that runs a gate |
| `assumptions.md` | — | I1, I2, L3; shown on the page as written |

Inside a run, `validate` also checks the ids a file names (EditScript targets, SectionMap and
Coverage finding ids, Probes regions) against the run's other files.

## Two writing rules

- **`assumptions.md` is reader text.** The page shows it verbatim. Write each line as a plain sentence
  a reader of the page needs ("The company is assumed to sell in the EU."), never internal wording:
  no mode names, model ids, file names, op names, config keys. Process notes go to the ledger:
  `meta ledger append "$RUN" --kind assumption --stage <op> --note "…"`.
- **Departures go to the ledger.** Any place you do not follow an op as written:
  `meta ledger append "$RUN" --kind deviation --stage <op> --note "what, and why"`.

## Kernel errors, common to every op

- **Launcher exit 126 or 127** (binary missing or won't start): switch the run to `kernel=absent` from
  this point. Earlier kernel results stay valid; every later check is `unchecked`.
- **JSON error with a field path** (schema): fix that field, retry once. A second failure follows the
  op's own table. Never edit kernel-written files by hand to get past an error.
- **Any other non-zero exit:** save stdout and stderr to `$RUN/gates/<op>.error.json`, then follow the
  op's table. Never write a result the kernel did not print.

## Stages and edit routing (until the loop controller ships)

An R1 edit goes to the stage that owns what must change, and everything after it re-runs:

| What must change | Stage |
|---|---|
| The question or success criteria were misread | I2 |
| A sub-question is missing, or a region is badly drawn | D1 |
| Evidence is missing, thin or misread (re-gather a region) | D2 |
| Two labels for one thing in a region | D3 |
| Two regions joined on a label that means different things | L3 |
| Page form or layout (timeline, split view, sections' order) | A1 |
| Anything on the page itself: wording, headings, labels, stamps, evidence display, layout details, how-checked text | A2 |

Page-text edits never re-run research: an edit whose target is on the page and whose fix is wording
goes to A2 even if the reviewer named an earlier stage. D2 is re-entered at most
`loop.d2_reentry_max` times per run. Review passes stop at `modes.<mode>.review_passes`.

**Caps** (region count, review passes, re-gather limit, independent-source minimum) come from
`meta config show --mode <mode>`; the config hash is in `run.json`.
<!-- <internal>Keys: modes.<mode>.max_regions, modes.<mode>.review_passes, loop.d2_reentry_max,
lens.pullback.min_independent, lens.pullback.as_of_window.</internal> --> `<mode>` is `fast`, or
`deep_narrow` / `deep_broad` per the scope I2 set with `meta run scope`. If the kernel is absent, use
the spec's fast-mode defaults (3 regions, 1 review pass, 1 re-gather, 2 independent sources) and say
so under Known limits.
