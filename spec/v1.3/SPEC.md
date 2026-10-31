# meta-suite v1.3, Lean

**Status:** draft for the owner's review (2026-10-30). This replaces the v2 build plan as the active direction. The v2 spec and the `next/` preview stay in the repo as a record; nothing is deleted.

**In one line:** keep v1 and keep what wins, cut what costs, and ship only what beats both v1 and plain Opus in a blind test, at the same or lower cost.

## 1. Why the direction changed

Four blind comparisons, with the owner judging first where he took part:

| Test | Task | Result |
|---|---|---|
| 1 | Regulatory research brief | v1 preferred, 90% (23 vs 18 of 24) |
| 2 | Competitive landscape of a small consultancy (private) | v1 primary for the owner and in 4 of 5 other verdicts; only the design judge preferred v2 |
| 4A | ICP artifact for a public brand | v1 first ("no contest"), then plain Opus, then v2 deep, then v2 fast |
| 4B | HTML pitch deck, 4 entrants (semi-blind) | v1 first, then plain Opus, then v2 fast, then v2 deep (broken slides) |

**Cost per run, measured as billed cost from the logs:**

| Entrant | Cost per run |
|---|---|
| Plain Opus | about $2.1 |
| v1 | $3.5–3.8 |
| v2 fast | $5.3–6.8 |
| v2 deep | $9.3–12.2, plus 1–4 extra Claude sessions it starts that don't appear in its own logs |

v2 lost every time, and in the artifact lanes it finished below plain Opus. v1 beat plain Opus in both artifact lanes. That is one run per lane, so treat it as a direction, not proof (section 6, T3).

## 2. What actually makes v1 win (measured, independent research)

These are habits in the output, not the pipeline's paperwork:

1. **It encodes data as custom figures.** In 4B, v1 had 8 data charts against 0 for plain Opus. In 4A the gap was small: 2 data figures against 1. The earlier "26 vs 1" count included v1's custom icon set. (Corrected by the habit check, eval/README.md.) The owner rewarded "has figures" and "client-ready", and penalised "research dossier".
2. **It looks at real phone renders and fixes what it sees.**
3. **It writes output rules before it builds.** Each claim is tagged fact or inference, estimates are never shown as measured, the page leads with the verdict, and it names the runner-up and what would change the ranking.
4. **It right-sizes at intake.** One artifact, not a bundle.
5. **It keeps content as typed data, rendered by a build step.** This keeps numbers consistent and allows interactive weighting.

**On the mathematical framing.** Its working form carries weight, and it lives in one document: the **typed meta-prompt** (`03_META_PROMPT.md`). It has typed input slots, invariants ("every claim tagged", "never present a share as measured", "name the runner-up and what flips it"), a step-by-step procedure, and a typed output contract. That document is where habits 3 and 5 come from, so it is protected. Build from typed data (`05_DATA_MODEL.md` and `seed_data.json`). But no v1 run ever opened `category-expansion.md`, and no category-theory vocabulary appears in any intermediate file. So the wording is not yet shown to be the mechanism. We keep the framing as the method and the brand, and test the wording on its own (T4).

## 3. The original vision, carried forward

| The owner asked for | v1.3 position |
|---|---|
| The category-theory lens as the secret sauce | Kept as the working method (entities, relations and invariants, then typed data). Whether the wording adds lift is tested (T4), not assumed. |
| Flexible, task- and domain-agnostic artifacts that impress subject-matter experts | The core goal. Coverage grows lane by lane: 4C (dashboard from meeting notes) and 4D (map a knowledge space) are frozen and next. |
| Faster, modern pipeline; Go/Rust; programmatic tools | Dropped from the product path. The kernel took 20–33% of v2's run time and is a hurdle for non-technical users. Small scripts are allowed where they save turns. |
| Self-refactoring prompts during a run | Kept in its working form: one short output-rules block written before the build. |
| Progressive disclosure, nested folders | Kept as a principle: one small core file, plus the one reference the artifact type needs, loaded on demand. Not v2's ten-file split. |
| Real tests across versions and models | Kept: the blind-test protocol is v2's best legacy. Effort level and model tier are tested (T5). |
| Plug and play for non-technical people | Required: no binary, no setup. Optional checks degrade gracefully when a browser isn't available. |
| Deep by default, `--fast` as a flag | `--fast` is dropped as a separate pipeline. Effort level is the dial, if T5 supports it. |
| Harness with alerts and monitoring | Out of scope here. It belongs to a separate product. |
| A clear spec, no scope creep | This page. A change enters only through the gate in section 5. |

## 4. What changes, in order

Scope: **meta-mvp first.** It was the entry skill in every v1 run and is the owner's main workflow. The other eight skills stay as they are until a lane tests them.

| # | Change | Expected effect | Test |
|---|---|---|---|
| C1 | **Render audit.** One small script checks for overflow, overlapping elements, stacked slides and console errors at 360/390/1280 px, plus one tiled phone contact sheet. All fixes go in one batched edit, with at most 2 rounds. | −$0.5–0.8 and −1.5–4 min per run (about half of v1's extra cost). Catches the failure that sank v2's deck. | T1 |
| C2 | **Merge, don't delete: one BRIEF.md instead of ten planning files** (about 3,000 words per run). Every job those files do is kept: the intake, the thesis (02), the research notes (01), **the typed meta-prompt with its invariants and output contract (03), which is load-bearing**, the data model (05), each module's edge cases (SPEC_01–03, condensed), and the verification checklist (from 04). Only duplicate scaffolding is dropped: the component-spec restatement (03b), the unit table where every owner is "orchestrator" (04), and the unverifiable "refined 3–5 times" claims. `seed_data.json`, `template.html` and `build.py` are untouched. | −$0.2–0.4 and −1–2 min | T1 + habit check |
| C3 | **Core skill file under about 150 lines**, plus one reference loaded per artifact type (page, deck, dashboard). "Read the kits" and the off-topic worked example are removed. | About −12–15k tokens carried per run | T1 |
| C4 | **The five habits (section 2) stated plainly at the top of the skill.** | About 0 cost; probably the main quality driver | T2 |
| C5 | **Helper agents optional, inline work by default.** | Probably the largest real-use saving; unmeasured | T7 |
| C6 | **Optional evidence appendix, collapsed:** the 3–5 headline numbers marked "two sources" or "one source". | Trust for skeptical readers, without the dossier effect | One run on the test-2 prompt |
| C7 | **Default effort for artifact runs:** low, only if T5 shows parity. | −10–20% | T5 |

## 5. How a change ships

0. **Habit check (automatic, before any blind test).** A script measures the winning features of the frozen v1 outputs:
   - the number of data figures;
   - the share of claims tagged fact or inference;
   - verdict in the first screen;
   - a named runner-up and what would flip the ranking;
   - interactive weighting where the frozen output had it;
   - graded sources;
   - a clean render-audit pass at 360/390/1280 px.
   
   A lean output that falls below the floor set from v1 on any of these (counts at 75% of v1, required flags, at most 2× v1's length, clean render) is rejected without spending a blind test. Built and calibrated: `eval/habit_check.js`; it agrees with the owner on 8 of 8 frozen pages, in-sample. Changes land one at a time: the safe cuts first (C1 batching, C3 loading, 03b), then the merge (C2), each re-checked.
1. **Fixed comparators.** Every frozen prompt keeps its judged v1 output as the fixed comparison: tests 1, 2, 4A and 4B, plus 4C and 4D once run.
2. **Pass bar.** A change ships only if, judged blind, it is **preferred or tied in at least 3 of 4 pairs** against frozen v1. It must also beat or tie the plain Opus baseline and cost **no more** in billed terms.
3. **Judging.** The owner judges first. AI judges are a pre-screen only after calibration (T6): the design judge disagreed with the owner in test 2.
4. **Hygiene for every run:**
   - billed cost taken from the logs;
   - a check for hidden `claude -p` sessions;
   - a skill-use audit;
   - one installed copy of each skill;
   - no page features in progress messages;
   - evidence appendices stripped before blind review.

## 6. Tests queued (each lane runs at Opus medium effort unless the test varies it)

| # | Test | Runs | About |
|---|---|---|---|
| T1 | v1-lean (C1–C3) vs frozen v1, on 4A and 4B, twice each | 4 | $11 |
| T2 | Plain Opus plus a 150-word card of the five habits vs frozen v1. A tie means the pipeline is dead weight. | 2 | $5 |
| T3 | Replicate unchanged v1 vs plain Opus on 4A (is "v1 beats Opus" noise?) | 2 | $6 |
| T4 | Math-wording ablation | 2 | $5 |
| T5 | Low vs medium effort | 2 | $4 |
| T6 | Calibrate an AI judge against the owner's recorded picks | 0 | $2 |
| T7 | One main-session run with helper agents, full cost logging | 1 | $6 |

T1–T3 (about $22) decide the top three changes. Run them first.

## 7. What carries over from v2

**Kept:**
- the blind-test protocol (frozen inputs, sealed keys, owner first);
- the skill-use audit;
- billed-cost accounting;
- a render audit (C1);
- an optional evidence appendix (C6);
- one optional fresh-reviewer pass for high-stakes work. It caught real overclaims, but it has not yet changed a ranking.

**Dropped from the product path:**
- the Go kernel, ledger and schemas;
- claim-level bookkeeping and findings graphs;
- sheaf gates (they stay as research in sheaf-core);
- the ten-file ops layout;
- `--fast` mode.

## 8. Unknowns to watch

- One run per lane and a single judge.
- Three copies of meta-mvp are installed (the v1.2.0 plugin, the repo, and an older account-level copy), so it's not known which one fires in plain use.
- Whether the owner's own bar is met: v1's 4B win was "best of the lot, not my favourite".
- Blinding leaks through house style.
- Hung tool calls distorting wall-time comparisons.
- 8 of 9 skills are untested.
- Published evidence that comprehensive skills can underperform focused ones (SkillsBench): add structure only where it wins.

## 9. Decisions needed from the owner

1. Approve this direction, with meta-mvp first.
2. Run T1–T3 (about $22) next.
3. Remove or update the older account-level meta-mvp copy, so one version fires.
