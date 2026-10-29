---
name: meta-lens
description: >-
  Ask a rough question; get back a checked, sourced answer as a clear page. meta-lens plans the
  question, researches it from several angles in parallel, counts a point as verified only when
  independent sources agree, labels what rests on one source or is disputed, and builds a page that
  leads with the answer and shows how it was checked. Use for "look into", "research", "compare", "give
  me the full picture", "is it true that", "meta-lens" or "/meta-lens". Deep by default; add --fast or
  say "quick" for a lighter pass. No setup. Hands off to the matching meta-suite skill when the request
  is really a prompt, plan, design, review or bundle. Not for common-knowledge questions.
metadata:
  short-description: "Ask a rough question; get a sourced, cross-checked answer page that leads with the verdict and shows how it was checked."
---

# meta-lens

Ask anything, however rough. You get a page that leads with the answer, labels how sure each point
is, and shows how it was checked. Deep by default; `--fast` or "quick" runs a lighter pass.

## Rules

- Never invent a source, number or check result. A check that did not run is reported as unchecked.
- Don't quiz the user. State assumptions on the page and proceed. Ask only when a wrong assumption is
  high-stakes **and** costly or impossible to undo.
- The user sees plain language only: no internal names, commands or file paths unless asked.

## Step 0 · probe

`ROOT` = two folders up from this one. `meta` means `sh "$ROOT/bin/meta"`.

Run `meta version`. Exit 0 → `kernel=present`. Otherwise `kernel=absent`: run the same steps without
kernel commands, report every check **unchecked**, mark the page **Provisional**, call nothing
verified.

## Program

**Fast** loads this file plus the core ops (D1, D2, A2, R1). **Deep** opens each op file as its step
starts. Map: [index.md](index.md).

1. **Frame.** [I1](ops/I1-classify.md): brief, mode, `meta run init` → `RUN`.
   [I2](ops/I2-expand.md) → `$RUN/spec.json`.
2. **Map.** [D1](ops/D1-cover.md) → `$RUN/cover.json`. [D2](ops/D2-gather.md): one subagent per
   region in parallel, then `meta claims merge`, `meta gate G2`. [D3](ops/D3-graph.md):
   `meta graph build`.
3. **Combine.** [L3](ops/L3-merge.md): `meta interface merge --plain`. [L2](ops/L2-verify.md):
   `meta findings verify`.
4. **Make.** [A1](ops/A1-select.md) → `$RUN/plan.json`. [A2](ops/A2-compose.md): `meta report`,
   the page, `$RUN/section_map.json`, `meta gate G6` to zero hits.
5. **Check.** [R1](ops/R1-review.md): a fresh-context reviewer returns edits. Re-run from the
   earliest step a P1/P2 edit touches, up to the pass cap.
6. **Record.** `meta ledger show "$RUN" > "$RUN/ledger.txt"`.

Each step starts with `meta stage begin` and validates what it writes with
`meta validate <file> --type <T>`. On a kernel error, fix the named field and retry once, then follow
the op's error table.

## Reply

Two or three plain sentences: the verdict; how many points are verified, single-source or disputed;
anything unchecked. Then the page: a host artifact when one fits, else `$RUN/artifact/index.html`.
Pass `meta gate G6` on the reply first.
