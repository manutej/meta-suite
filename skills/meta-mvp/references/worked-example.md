# Worked example (quality anchor) — "coordinate 100 people for my friend"

A one-line prompt → the full bundle. Use this to calibrate depth; do not copy its content, copy its *level*.

- **Prompt:** "Help my friend Kevin coordinate ~100 people to show up and act."
- **0 INTAKE:** audience = the friend (stakeholder) → outcome = coordinate action + make each person
  feel seen → aesthetic = warm/editorial → scope = full bundle.
- **1 GROUND:** web-researched the friend's org; pulled his real words from a connected voice MCP →
  discovered a 3-tier topology (core team / coaches / ecosystem), real initiatives, real tone. Cited.
- **2 EXPAND (hidden):** objects = people/orgs/initiatives; morphisms = offers/asks/intros; identity =
  relationships; colimit = one room view. Picked the decomposition into 4 modules.
- **3 META-PROMPT (`03`):** typed system prompt — slots {members, threads, invites, intent}, procedure
  Segment→Connect→Listen→Respond→Weave→Synthesize, JSON output of the operator's top-3 moves.
- **3b SPEC LAYER (`03b` + `SPEC_01..04`):** component-spec meta-prompt, then 4 sub-agents wrote
  Invite/Onboarding, Threads, Weave, Architecture — each concept brief + build spec, linking one
  canonical `Member` entity. (This is the depth the thin version lacked.)
- **4 PLAN (`04`):** typed work units + gates + parallel lanes + human gate on the SLA.
- **5 MODEL (`05` + seed.json):** grounded entities (people/orgs/initiatives/channels), sensitive
  facts gated out.
- **6 BUILD:** `dashboard_tool.html` (Motion: segmented control, scroll-reveal, spring modal viewer)
  + `dashboard_showcase.html` (Aceternity aurora/beams + the ceti 8-beat explainer on top), custom
  SVG, mobile-ready, palette from a real dictionary, 100% coherence.
- **7 VERIFY:** JS parses, ceti contract holds, palette audit ~100%, mobile markers present, no slop.
- **8 SHIP:** presented in order; offered the optional categorical appendix; offered Vercel deploy
  (did not deploy without a yes).

**The tell of a good run:** there is a SPEC layer, the meta-prompt and plan are typed and example-
agnostic, the data model is grounded in real research, and both dashboards have load-bearing motion +
a working scrubbable explainer. If any of those are missing, the run regressed — go back and add it.
