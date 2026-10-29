# 11 · v1 inventory: where every v1 element goes (VR-00)

VR-00 asked for **all** the elements of the v1 skills to be rewritten for Opus 5.5. v2.0 builds the
kernel plus `meta-lens` first. This table makes sure no v1 element is lost: each one is absorbed in
v2.0, scheduled for a port onto the kernel, or explicitly dropped.

| v1 skill | Element | v2 home | When |
|---|---|---|---|
| meta-research | decompose, fan out, credibility triage, adversarial cross-verify, cited brief | D1, D2, L2, G2, G3, R1, A2 | **v2.0: absorbed** (meta-lens replaces it) |
| meta-research | reusable knowledge atoms | S2 (recorded); domain residue | v2.0 records / v2.1 loads |
| meta-prompting | task → category → typed slots → procedure → bind → render | I1, I2 (spec expansion is meta-prompt generation), the ops/ contract format | **v2.0: inside meta-lens** |
| meta-prompting | recursive refinement (RMP) and edit-script history | O1 + R1 edit scripts, stage-boundary adaptation (D-057), S1 | **v2.0** |
| meta-prompting | standalone "emit a reusable meta-prompt for category X" plus a sibling-instance test | port onto I1–I3/G5 + J1 | v2.1 (port 2) |
| meta-review | lens selection, adversarial review, P1–P3, blind spots | R1 (inside the pipeline) | **v2.0: internal** |
| meta-review | standalone review of the user's own deliverable; lesson atom | port R1 + J1 for any artifact; S2 | v2.1 (port 1) |
| meta-design | typed acceptance predicates (WCAG etc.), computed | G7 | **v2.0** |
| meta-design | medium-specific design schema, tokens, render-audit loop, design-system output | A1/A2 + host design types; full schema derivation ported | v2.2 (port 4) |
| meta-mvp | one-liner → grounded research → artifact | Frame → Map → See → Make | **v2.0** |
| meta-mvp | WOW standard, mobile-ready, show what the code does | M8 artifact bar (D-056), P12, P14 | **v2.0** |
| meta-mvp | vision, execution plan, data model, prototype/API, animated explainer, palette tool | bundle builders on top of meta-lens output | v2.1 (port 3) |
| meta-planning | domain mapping before planning | D1 cover + tree | **v2.0** |
| meta-planning | meta-plan for agent fleets, work packets, audits, anti-slop | port onto the kernel (O1 + J1 + G gates) | v2.2 |
| meta-craft | craft constitution, Definition of Done | principles P1–P14 inside v2; standalone port | v2.2 |
| meta-compound | compounding map, capture/retrieve/refresh | S2 + domain residue; standalone port | v2.1 (residue) / v2.2 |
| meta-context | context strategy, substrate scaffold, handoff | L0–L4 disclosure, S1 ledger, residue; standalone port | v2.2 |

**Dropped:** nothing. Anything dropped later needs a decision entry and the owner's sign-off.

**Proposed port order** (cheapest first, since they reuse the most atoms):

| Port | Skill | Version |
|---|---|---|
| 1 | meta-review | v2.1 |
| 2 | meta-prompting | v2.1 |
| 3 | meta-mvp | v2.1 |
| 4 | meta-design | v2.2 |
| 5 | planning, craft, compound and context | v2.2 |

**Plan B if Phase A is a no-go.** v1 stays the default install. The kernel's value does not depend
on the lens: provenance, gates, the ledger, judgment sampling and plug-and-play all hold without it.
So the v1 skills are ported onto the kernel **without** the lens operators, and the lens is
re-attempted in v2.1 with the Phase A findings.
