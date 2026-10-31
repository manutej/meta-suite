# The 8-stage pipeline (full contracts)

Each stage has a typed input → output and a gate. Don't advance past a failing gate.

## 0 · INTAKE  → `00_INTAKE.md`
- **In:** the user's minimal prompt (could be one line).
- **Do:** resolve the four unknowns (audience, outcome, aesthetic, scope) per `intake.md`.
  **Default: derive them** by scoring the plausible values against the prompt + research (never
  interrogate the user); for high-stakes ambiguity run the 2–3 agent panel; only use
  `AskUserQuestion` if the user explicitly opts in (`--ask` / "ask me"), and fall back to the
  categorical path if that tool is unavailable.
- **Gate:** you can name the audience, the single outcome, the aesthetic, and the scope tier,
  each with a confidence + deciding signal — without having blocked on the user.

## 1 · GROUND  → `01_RESEARCH.md`
- **In:** intake + prompt.
- **Do:** Research before building. Public facts via web search + fetch. The stakeholder's *real*
  world via connected MCPs (Carbon Voice transcripts, Gmail, Notion, Calendar, Drive…). Apply
  /research-synthesis (sources→themes→insights→implications) and /user-research (personas, journey).
  Brand: attempt to extract real colors (site CSS, brand guide); if blocked, derive + flag honestly.
- **Gate:** every downstream claim traces to a cited source or a flagged assumption. No invented facts.

## 2 · EXPAND  (internal; feeds everything)
- **In:** grounded findings.
- **Do:** category-theory expansion (see `category-expansion.md`): list entities (objects),
  relationships (morphisms), invariants, and the few free variables vs fixed structure. Decide the
  decomposition — it determines the shape of the meta-prompt, the plan, and the data model.
- **Gate:** every entity is typed; every relationship named; ≥1 honest unknown recorded.

## 3 · META-PROMPT  → `03_META_PROMPT.md`
- **Do:** compose /meta-prompting. Emit the reusable, example-agnostic system prompt with typed
  slots, a procedure, bound rules, and a checkable output contract.
- **Gate:** zero solved examples; serves the whole category, tested against 2–3 sibling instances.

## 3b · SPEC LAYER  → `03b_COMPONENT_SPEC_META_PROMPT.md` + `SPEC_01..NN_*.md`
- **Do:** the depth-defining stage, easy to skip and the #1 regression. Decompose the system into 3–6
  modules. Write the component-spec meta-prompt (`meta-prompt-kit.md` Part B), then **spawn one
  sub-agent per module concurrently** to write each two-layer spec (`spec-layer.md`).
- **Gate:** ≥3 module specs exist; each has both layers, typed entities, a no-dead-state machine, ≥3
  edge cases with responses, ≥1 marked spike, and links the canonical data-model entity by name.

## 4 · PLAN  → `04_PLAN.md`
- **Do:** compose /meta-planning. Typed work-unit graph: units, contracts, gates, parallel lanes,
  human checkpoints. Spawn sub-agents for independent build lanes.
- **Gate:** every unit has an output type + a done-criterion; parallel pairs are truly independent.

## 5 · MODEL  → `05_DATA_MODEL.md` + `seed_data.json`
- **Do:** build the common data model from stage-1 facts. Entities for the domain (people, orgs,
  initiatives, conversations, channels, … as the domain needs). Emit a valid JSON seed of real instances.
- **Gate:** JSON parses; sensitive facts carry `sensitivity:"personal"` and are excluded from
  customer-facing artifacts; types match what the dashboards consume.

## 6 · BUILD  → `dashboard_tool.html`, `dashboard_showcase.html` (+ any docs)
- **Do:** see `dashboard_playbook.md`. Default architecture: Motion editorial for the daily tool;
  Aceternity skin for the showcase, with a ceti-method animated explainer at the top. Custom inline
  SVG icons. Palette via /dict-color. Encoding via /bertin. Mobile-ready throughout.
- **Gate:** the quality gates in `quality-gates.md` pass.

## 7 · VERIFY
- **Do:** run all gates (JS parse, ceti contract, palette audit, mobile breakpoints, anti-slop,
  fact-check). For high-stakes work, run a verification sub-agent.
- **Gate:** all green; nothing unverified ships unflagged.

## 8 · PACKAGE & SHIP
- **Do:** present the bundle with `present_files`. Offer `APPENDIX_how-we-thought.md`. Offer a
  Vercel preview deploy **only on explicit approval**.
- **Gate:** files named in order; WOW confirmed; no side-effectful action taken without a yes.
