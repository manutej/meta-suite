# A2 · compose the page (core op)

<!-- <internal>A structure-preserving map from the merged graph and its findings to page sections: every
section is the image of a set of findings; nothing unanchored survives. "How this was checked" is
rendered only from kernel output (P12). G6 enforces that no internal construct reaches the reader.</internal> -->

```xml
<contract op="A2">
  <in name="plan" type="ArtifactPlan" path="$RUN/plan.json"/>
  <in name="merged" type="DomainGraph" path="$RUN/graph/merged.json"/>
  <in name="findings" type="Finding[]" path="$RUN/findings.jsonl"/>
  <in name="conflicts" type="Conflicts" path="$RUN/conflicts.json"/>
  <in name="coverage" type="Coverage" path="$RUN/coverage.json"/>
  <in name="assumptions" type="text" path="$RUN/assumptions.md"/>
  <in name="edits" type="EditScript" path="$RUN/review/edits.&lt;n&gt;.json" optional="true"/>
  <out name="page" type="html" path="$RUN/artifact/index.html"/>
  <out name="section_map" type="SectionMap" path="$RUN/section_map.json"/>
  <out name="how_checked" type="HowChecked" path="$RUN/how_checked.json">kernel-written by `meta report`</out>
</contract>
```

## Order (matters: "how this was checked" must report the page's own G6)

1. `meta stage begin "$RUN" A2 …`, then `meta gate G2 "$RUN" > "$RUN/gates/G2.json"`: it must be
   `pass` now (claims, graphs and findings all exist).
2. **Section map** in `plan.json`'s order: each section an `id`, a `tag`, the `finding_ids` it rests
   on. Tags **index** findings, they do not split them: a finding sits in its reasoning-order `main`
   section and may also be listed in a `contested`, `single_source` or `uncertain` index section.
   `contested` covers every `conflicts.json` entry and every "kept apart" label; `uncertain` the
   unsettled or unanswered leaves in `coverage.json`. A section with no finding ids is cut, except
   `how_checked` and `assumptions`. `meta validate "$RUN/section_map.json"` (checks every id).
3. **Verdict** from verified findings only: one or two sentences with the key number. Nothing
   verified → say so plainly; the page is **Provisional**.
4. **Draft** the page: with `plan.host_type` set, that host type carrying the same sections, labels
   and checks; otherwise `$RUN/artifact/index.html` to the scaffold below.
5. `meta gate G6 "$RUN/artifact/index.html" > "$RUN/gates/G6.json"`; reword every hit and re-run until
   zero (terms in `brief.vocab_allowlist` are exempt).
6. `meta report "$RUN"` → `how_checked.json`, now recording this page's G6.
7. **Final render:** fill "How this was checked" from `how_checked.json`, then G6 once more. Any hit:
   reword, then steps 5–7 again, so the page never reports a check it did not pass.
8. **Self-check before R1.** Every finding id in `section_map.json` is an element id on the page;
   every number traces to the data block; every status shows a word, not just a colour. With a
   headless browser, save screenshots at 360×740 and 1280×800 (light and dark) to
   `$RUN/artifact/shots/`; without one, R1 reviews the source and says so.

## Page scaffold (hand-built page)

Design guidance, not a template to copy. The user's or host's design system overrides token values
and font families; the structure does not bend.

**Data first.** Embed `how_checked.json`, the findings used and `section_map.json` in one
`<script type="application/json" id="run-data">`. Every number in the markup carries
`data-src="<json path>"` naming where it came from. No number is typed by hand. Counts always show
their denominator ("5 of 9"). A "number" here is a count or measurement; dates, article numbers and
the verification rule's own "two" are text.

**Tokens.** 9–19 role tokens on `:root`, and no raw colour anywhere else. Default set (each text pair
was computed at ≥ 5.7:1 in both themes):

| Role | Light | Dark |
|---|---|---|
| `--ground` | `#f9f6f2` | `#0f1216` |
| `--surface` | `#fefdfa` | `#171b20` |
| `--ink` | `#1b1f27` | `#eae7e2` |
| `--muted` | `#565b63` | `#aeaaa2` |
| `--rule` | `#d3d1cb` | `#35383d` |
| `--accent` | `#0c60a3` | `#7fb6ee` |
| `--pass` | `#206b38` | `#82cb92` |
| `--warn` | `#8d5406` | `#e7b369` |
| `--fail` | `#ac312a` | `#f08f83` |

The dark set is the light set with lightness inverted and hue held. Declare it twice:
`@media (prefers-color-scheme: dark) { :root:not([data-theme="light"]) {…} }` and
`:root[data-theme="dark"] {…}`. The accent covers about 5% of the page or less: links, focus, one
mark per figure. If you change any value, recompute every text pair (4.5:1 body, 3:1 large text).

**Three type roles.** Display (headings: a serif stack such as `"Iowan Old Style", Charter, Georgia,
serif`), body (`system-ui, -apple-system, "Segoe UI", sans-serif`), mono (`ui-monospace,
"SF Mono", Menlo, Consolas, monospace`) for labels, stamps and numbers. Body 16–18px, line-height
1.55, measure `max-width: 68ch`. Headings `font-size: clamp(…)` with `text-wrap: balance`. Eyebrows
mono 11–12px, uppercase, `letter-spacing: .08em`. Numbers `font-variant-numeric: tabular-nums`.

**First screen (360×740), in this order:** an eyebrow (date and "Provisional" when it applies), a
title that states the finding, the verdict with its number, and a status strip of 4–8 counts from
`how_checked.json` (sources, independent publishers, verified, single-source, contested, checks
passed out of run). Each count links to its section or to "How this was checked". Any correction or
`kernel=absent` notice sits here, above the fold.

**Layers.** Level 1: verdict and strip. Level 2: 4–8 sections in reasoning order, each `<section
id="s-<section id>" data-tag="<tag>">` with a heading that states its finding and a short contents
nav linking to every section. Level 3: evidence inside `<details>` per finding, one plain sentence
per source: what it says, publisher, linked document title, date ("The European Commission's
implementation timeline (read 30 Oct 2026) gives 2 August 2026."). Never show raw claim triples,
relation names (`applies_from`) or ids. Every level opens by tap; nothing is hover-only.

**Findings and labels.** Each finding is an element with `id="f-<finding id>"` and a stamp in mono:
**Verified** (with "n independent sources"), **Single source**, **Contested** (both values, each with
its sources, side by side, no winner picked), **Not settled**, or **Unverified** (`kernel=absent`).
A finding with `as_of` says "as of <month>".
Every stamp is a word plus a glyph (✓ · ! · ⇄ · ?), with colour as a third cue only.

**Mobile first.** Base styles at 320px; breakpoints at about 40rem and 64rem, grids collapse to one
column below the first. Gutters `max(16px, env(safe-area-inset-left))` (and right). No page-level
sideways scroll at 320, 360, 390 or 768px; wide tables sit in `overflow-x: auto` wrappers captioned
"scrolls sideways". Tap targets at least 44px. `<meta name="viewport" content="width=device-width,
initial-scale=1, viewport-fit=cover">`.

**Interaction.** A theme toggle (`aria-pressed`, sets `data-theme`), styled `:focus-visible`, every
control changes visible state. No external scripts; if one is unavoidable, pin its exact version. No
console errors. Motion is optional: transitions 0.2–0.9s with named easing, and under
`prefers-reduced-motion: reduce` none at all.

**How this was checked** (`id="how-checked"`, near the end, linked from the strip). Rendered only
from `how_checked.json`, in plain words: the date and mode of the run; claims gathered, sources,
independent publishers; verified (with the rule), single-source, contested; withdrawn statements and
why; each check with Passed, Failed or **Unchecked**, and why unchecked. Then **Known limits**:
regions not gathered, unanswered sub-questions, regions researched separately, review edits left
open, page checks not in this build. Then **Assumptions** from `assumptions.md` as written. Then a
dated **Sources** list (publisher, title, link, date read) and a stamp with the run id and "locked".

**Voice.** Short declarative sentences. Numbers and dates inline. Define a term at first use. No
hype words, no exclamation marks. Headings say what was found, not what the section is.

## Kernel errors

| Command | Error | Do |
|---|---|---|
| `report` | any | Build "How this was checked" from the gate files you have; any figure you cannot read from them is **Unchecked**, never estimated. |
| `gate G2` | not `pass` | `fail`: follow D2's table. `incomplete`: a step before A2 did not run; run it. |
| `validate` (SectionMap) | field path | Fix; retry once. Second failure: mark "section coverage" unchecked. |
| `gate G6` | hits | Reword; re-run. After three rounds, list remaining hits under Known limits. |
| `gate G6` | crash | Mark "plain-language check" unchecked; reread the page yourself for the words in G6's list. |
