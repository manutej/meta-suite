---
name: meta-mvp
description: >-
  Turn a minimal prompt — even one line — into a complete, WOW-grade deliverable bundle right-sized
  to the domain: grounded research, a project vision, a reusable meta-prompt, an execution plan, a
  grounded data model, and the production-quality artifacts the work needs (interactive dashboards,
  a report or doc site, an API/spec, a prototype), plus an optional animated SVG explainer.
  Research-grounded and mobile-ready. Use when the user says meta-mvp, "build an MVP / vision /
  prototype for X", "go from this one-liner to something I can show", "make me the full pitch", or
  "spin up a new project for a client". Resolves its own scoping (never interrogates you; for
  high-stakes ambiguity it cross-checks with a small agent panel). Optionally composes other
  installed skills (meta-prompting, meta-planning, frontend-design, and similar) and degrades
  gracefully without them. Do not use for a single fully-specified one-off file or pure question
  answering.
metadata:
  short-description: "Turn a one-line idea into a full, research-grounded deliverable bundle: vision, meta-prompt, plan, data model, and production-quality artifacts such as dashboards, docs or a prototype."
---

> **Meta-Suite doctrine (overrides anything below).** Never ask the user clarifying questions —
> infer what you need from the request, choose sensible defaults, and state your assumptions in one
> line before proceeding. Never surface category theory, functors, monads, morphisms, or other
> typed-formalism vocabulary to the user; keep that machinery internal and reveal it only if the
> user explicitly asks to review the internal work. Keep all user-facing output in plain language.

# meta-mvp

Take a minimal prompt — even one line — and ship a complete, stakeholder-ready bundle that produces a
WOW effect from the first glance.

> **This skill is SELF-CONTAINED.** It does not assume any sibling skill is installed. The procedures,
> the spec layer, and the copy-paste templates it needs are in `references/` and `templates/` here.
> If sibling skills happen to be loaded you may also use them, but the bundled kits are the source of
> truth — the result must be identical with or without them. **Read the kits; don't wing it from the
> description.** Thin runs that skip the kits regress to generic output. A bundled palette tool ships
> at `scripts/dict_color.py` in this skill's folder (stdlib-only) so the palette gate runs with no external install.

## The core move (do it, then hide it)
A minimal prompt is an object; the bundle is its image under a structure-preserving map. Expand the
prompt into a typed domain (entities, relationships, invariants), then map that structure onto
deliverables: entities→data-model types, relationships→contracts and intros, independence→parallel
build lanes. **Category theory is the engine, not the output** — every customer-facing artifact hides
the abstraction. (Optional `APPENDIX_how-we-thought.md` is the only place it may appear.)

## What it ships — a right-sized bundle (depth is constant; surface adapts)
Two layers. The **depth spine is always produced** — it is what separates a real MVP from a generic
page. The **surface artifacts adapt** to the domain and scope tier: don't force two dashboards onto
work that isn't a dashboard, and don't ship the full bundle for a prototype. A run that skips the SPEC
layer, the meta-prompt, or the plan has regressed — but the *surface* is chosen, not fixed.

### The depth spine — always, every domain
1. `00_INTAKE.md` — the four resolved unknowns, each with confidence + deciding signal
   (`references/intake.md`). Derived by default; never a wall of questions.
2. `01_RESEARCH.md` — grounded findings + sources (web + connected MCPs). No invented facts.
3. `02_VISION.md` — vision + scope in the stakeholder's language.
4. `03_META_PROMPT.md` — the reusable system meta-prompt (`references/meta-prompt-kit.md` Part A).
5. `03b_COMPONENT_SPEC_META_PROMPT.md` + `SPEC_01..NN_*.md` — **the SPEC LAYER**: decompose into 3–6
   modules and spec each (concept brief + build spec) via sub-agents (`references/spec-layer.md`). The
   depth most often dropped — never skip it in any domain (a prototype may run 1–2 modules).
6. `04_PLAN.md` — the typed execution plan (`references/meta-plan-kit.md`).
7. `05_DATA_MODEL.md` + `seed_data.json` — the grounded common data model (when the domain has data).

### The surface artifacts — pick what fits the domain
Map the domain to its natural WOW artifact(s). Build the ONE that carries the outcome first; add a
second only if the scope tier calls for it. Every surface still earns WOW: custom inline SVG, real
palette, purposeful motion, mobile-ready.
- **Data / pitch / ops product →** the two-dashboard pattern (Motion editorial daily tool + Aceternity
  showcase) — `templates/dashboard-skeleton.md`. Default *for this domain*, not for all work.
- **Explainer / concept / teaching →** the ceti animated explainer (`templates/ceti-engine.md`) as lead.
- **Document / report / strategy →** a typeset long-form artifact or doc site; dashboards optional/dropped.
- **API / system / technical spec →** the spec bundle + a reference doc or interactive schema viewer.
- **Anything else →** derive the artifact from the data model + outcome. The rule is "one glance tells
  the story", not "there must be a dashboard".
Offer an animated explainer on top when it sharpens the story.

8. `APPENDIX_how-we-thought.md` — optional internal appendix (offer; don't force).

### Right-size the run (don't be overkill)
Read the tier from the resolved `scope`:
- **Prototype** → light spine (intake + lean vision + one surface artifact); SPEC fan-out may be 1–2
  modules; skip the second dashboard. Say it's a prototype.
- **Full bundle** → the whole spine + the domain's full surface set + a verify pass.
- **Visuals only / strategy only** → just that slice, still grounded and gated.
When torn between tiers, choose the smaller and say so — "this needs less" is a passing answer.

## The protocol (run in order; gate between stages)
Full contracts in `references/pipeline.md`. The stages:

0. **INTAKE** — resolve audience, outcome, aesthetic, and scope from the prompt + research. Derive by default (never interrogate); panel-review high-stakes ambiguity; ask only on explicit opt-in. See `references/intake.md`.
1. **GROUND** — research before building: web search/fetch for public facts; connected MCPs for the
   stakeholder's real world (voice transcripts, mail, docs). Brand: pull real colors or derive + flag.
2. **EXPAND** (internal) — `references/category-expansion.md`. Decompose; name modules + invariants.
3. **META-PROMPT** — `references/meta-prompt-kit.md` Part A.
4. **SPEC LAYER** — `references/meta-prompt-kit.md` Part B + `references/spec-layer.md`. **Sub-agents,
   one per module, concurrently.** This is the depth most often dropped — never skip it.
5. **PLAN** — `references/meta-plan-kit.md` (typed units, gates, parallel lanes, sub-agent protocol).
6. **MODEL** — grounded entities; emit valid `seed_data.json`; gate sensitive facts out.
7. **BUILD** — build the domain's surface artifact(s) (see "What it ships" + `dashboard-playbook.md`):
   dashboards via `templates/dashboard-skeleton.md`, explainer via `templates/ceti-engine.md`, or the
   domain-appropriate artifact. Palette via bundled `scripts/dict_color.py`; custom SVG; mobile.
8. **VERIFY** — every gate in `references/quality-gates.md` passes. For high stakes, a verify sub-agent.
9. **SHIP** — present in order; offer the appendix; offer a Vercel deploy only on explicit approval.

## Hard completion gate (check before you say "done")
- [ ] Intake resolved (derived by default, or asked only on opt-in); vision in the stakeholder's language.
- [ ] Research grounded + cited; palette real or flagged.
- [ ] Meta-prompt typed + example-agnostic (tested vs 2–3 sibling instances).
- [ ] **SPEC LAYER present**: component-spec meta-prompt + module specs (≥3 for a full bundle, ≥1 for a
      prototype), each with both layers, typed entities, no dead states, ≥3 edge cases, ≥1 spike, linking one canonical entity.
- [ ] Plan: typed units + done-criteria + gates + parallel lanes.
- [ ] Data model grounded; `seed_data.json` valid; sensitive facts gated out (when the domain has data).
- [ ] Every surface artifact chosen for the domain: custom SVG, load-bearing motion where used,
      mobile-ready, palette-coherent (dict-color gate PASS or hand-checked), JS parses.
- [ ] If an explainer ships: 8 beats, ~35–45s, deterministic, scrubbable.
- [ ] Right-sized to the scope tier — no artifact built that the tier didn't call for.
- [ ] No category-theory vocabulary in any client artifact; appendix offered.
- [ ] Files numbered; presented; deploy offered (not executed) on approval.

Missing any box = not done. See `references/worked-example.md` for the target depth.

## Non-negotiables
WOW from glance one (custom SVG, real palette, purposeful motion) · research first, no invented facts
· abstraction hidden in outputs · mobile-ready · right-sized to the domain + scope (never overkill) ·
honest flags on anything unverified · never publish/send without approval · gate personal/sensitive
data out of client artifacts · verify before ship.

## References & templates
- `references/pipeline.md` — the 9 stage contracts.
- `references/intake.md` — the short-intake questions.
- `references/category-expansion.md` — minimal prompt → typed domain (the hidden engine).
- `references/meta-prompt-kit.md` — meta-prompt (Part A) + the component-spec generator (Part B).
- `references/spec-layer.md` — per-module specs via sub-agents (the missing piece).
- `references/meta-plan-kit.md` — the typed planning scaffold + sub-agent protocol.
- `references/dashboard-playbook.md` — Motion+Aceternity patterns, palette, encoding, mobile.
- `templates/ceti-engine.md` — the inline deterministic 8-beat animation engine (copy-paste).
- `templates/dashboard-skeleton.md` — the two dashboard skeletons (copy-paste).
- `references/quality-gates.md` — the verification gates + commands.
- `references/worked-example.md` — the quality anchor.
