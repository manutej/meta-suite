---
name: meta-mvp
description: >-
  Turn a minimal prompt, even one paragraph or one line, into one client-ready, research-grounded
  artifact (an ICP or strategy page, an HTML pitch deck, a dashboard, a research digest, an explainer)
  built from a typed brief and typed data, then audited on phone and desktop. Use when the user says
  meta-mvp, "build an MVP / vision / prototype for X", "turn this use case into something I can show",
  "make me the pitch", or "spin up a project for a client". Resolves its own scoping and never
  interrogates the user. Do not use for pure question answering.
metadata:
  short-description: "One paragraph in, one client-ready artifact out: researched, typed, figure-led, phone-checked."
  version: "1.3.0-lean (candidate, not yet shipped)"
---

> **Doctrine (overrides anything below).** Never ask clarifying questions. Infer what you need,
> choose sensible defaults, and state your assumptions in one line. Never surface category-theory or
> other formal vocabulary to the user; it stays internal. All user-facing output is plain language.

# meta-mvp (lean)

## The five habits that win (measured in blind tests; do all five)
1. **Encode, don't decorate.** Put the key numbers into custom inline-SVG figures (bars, dot plots,
   ranges, timelines, matrices) that answer a question the reader has. One glance tells the story.
2. **Lead with the verdict.** The first screen (or first slide after the title) states the decision
   or answer and how sure you are.
3. **Tag every claim** as Fact (with its source) or Inference. Estimates are never shown as measured.
4. **Name the runner-up and what would flip it**, plus what would change your mind.
5. **Check real phone renders** and fix what you see, in one batched pass (step 5).

## The core move (do it, then hide it)
Expand the prompt into a typed domain: its entities, how they relate, and what must stay true.
Write that as a typed meta-prompt: slots, invariants, procedure and output contract. Then fill
typed data and render the artifact from that data. Structure first, prose second. The page
never shows the machinery.

## Right-size
Build **one** artifact by default: the one that carries the outcome. Add a second only if the request
clearly needs it. "This needs less" is a passing answer; say what you left out.

## Protocol
Work in one folder. Keep only these files: `BRIEF.md`, `seed_data.json`, `template.html`, `build.py` (or `.js`), `final.html`.

0. **Intake (silent).** Resolve audience, single outcome, aesthetic and scope from the prompt. Write
   them as 4 lines at the top of BRIEF.md, each with a confidence level. Start research at once.
1. **Ground.** Research before building: web search and fetch, plus any connected sources the
   request points to. Grade each source A (primary or research), B (trade press or reported) or
   C (vendor, forum or review site). No invented facts; flag what you couldn't verify. Use real
   brand colors where they exist, or derive them and say so.
2. **Brief.** Write `BRIEF.md` from `references/brief-kit.md`. It holds the thesis, the typed
   meta-prompt (the load-bearing part), the data model, 3–6 modules with their edge cases, and the
   verification checklist. Refine the meta-prompt against 2 sibling cases you name, and keep the
   real edit list. Never claim passes you didn't run.
3. **Model.** Write `seed_data.json`, the single source of truth. Every claim object carries
   `tag: "F" | "I"` and source ids that exist. Keep personal or sensitive facts out.
4. **Build.** Read the **one** surface reference for your artifact:
   - strategy, ICP, report or digest page → `references/surface-page.md`
   - pitch deck or slides → `references/surface-deck.md`
   - dashboard → `references/dashboard-playbook.md` + `templates/dashboard-skeleton.md`
   - animated explainer → `templates/ceti-engine.md`
   - a domain whose structure isn't obvious → also `references/category-expansion.md`

   Render `final.html` from `seed_data.json` with the build step; never hand-copy numbers into the page.
   Pick the palette with `scripts/dict_color.py generate`, or use the brand's.
5. **Audit (at most 2 rounds).** Run `node scripts/render_audit.js final.html --out audit/`. It
   checks 360, 390 and 1280 px for sideways overflow, script errors and overlapping text (decks
   are checked after stepping through slides), and writes one contact sheet per width. Look at the
   contact sheets **once**, list every problem, fix them all in **one** edit, and re-run. Stop
   after round 2 and report anything still open. Also run:
   `python3 scripts/dict_color.py audit final.html --threshold 12` (expect PASS), and grep the page
   for hype words and formal vocabulary (`functor|morphism|colimit|yoneda|adjunction`; must be empty).
   If Node, Playwright or Python is missing, check by hand and say which check was skipped.
6. **Ship.** Present `final.html` with its assumptions in 3–5 lines. Never publish, send or deploy
   without explicit approval.

## Helpers
Work inline by default. Use sub-agents only when a module is large and genuinely independent (it
needs no context from the others), or for an independent fact re-check on high-stakes work. If no
sub-agent tool exists, do the same work inline.

## Done means
- [ ] BRIEF.md has the intake, the thesis, the typed meta-prompt with its real edit list, the data model, the modules with their edge cases, and the checklist
- [ ] `seed_data.json` parses; every rendered claim is tagged and traces to a real source id
- [ ] All five habits are visible in `final.html`
- [ ] The render audit is clean at 360, 390 and 1280 (or any open issue is reported); palette PASS; no hype words or formal vocabulary
- [ ] One artifact, right-sized; assumptions stated
