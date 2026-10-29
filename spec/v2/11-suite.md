# 11 · The v2 suite (VR-00, VR-27)

**v2.0 is a new suite**, not one skill (D-065). Every v1 skill is rebuilt for Opus 5.5 on the
shared Go kernel and the shared lanes. `meta-lens` is the front door: "just ask," and it routes you
to the right specialist or specialists.

## How the suite is built

**One kernel.** Schemas, gates, J1 sampling, the ledger, the loop controller and the sheaf check.

**One set of lanes:** Frame → Map → See → Make → Check (the 25 core atoms). Every skill composes the
lanes it needs; only meta-lens runs all five on every question.

**Specialist skills are thin.** Each is:
- one SKILL.md trailhead (≤ 500 words)
- a recipe naming which shared atoms it runs
- **at most 2 atoms of its own**

**Atom budget for the whole suite: 34** = 25 core + 9 specialist atoms. Anything more goes through
a decision that displaces something.

## The suite

| v2 skill | What the user gets | Shared atoms it composes | Own atoms (new) | Rebuilt from |
|---|---|---|---|---|
| **meta-lens** (front door) | Ask anything rough; get a verified research view or artifact. Routes to specialists when the question is really a plan, a prompt, a design, a review… | all lanes; I1 also routes | — | new, plus meta-research |
| **meta-research** | The research lane with a brief as output: cited, contested points named, plus reusable takeaways. A rebuilt **alias of meta-lens in brief form** (D-072) | Frame, Map, See, Make (brief form), Check | — | meta-research |
| **meta-prompting** | A reusable, example-agnostic meta-prompt for a whole category, **proven on sibling instances** | I1–I3, G5, J1, R1, S1 | [P1 prompt.generate](atoms/P1-prompt-generate.md), [P2 prompt.sibling_test](atoms/P2-prompt-sibling-test.md) | meta-prompting |
| **meta-planning** | A plan agents can execute: context-complete work packets, a parallel/sequence graph, audit gates | Frame, D1 (domain mapping), J1, R1, O1 | [PL1 plan.packets](atoms/PL1-plan-packets.md), [PL2 gate.plan_dag](atoms/PL2-gate-plan-dag.md) | meta-planning |
| **meta-design** | A design system plus artifacts generated from it, audited | Frame, A1, A2, G7, R1, rubric (12) | [DS1 design.schema](atoms/DS1-design-schema.md) | meta-design |
| **meta-mvp** | One line → the full bundle: research, vision, plan, data model, meta-prompt, artifacts, explainer | meta-lens output + meta-planning + meta-design recipes, P1, I1 + J1 | [MV1 bundle.compose](atoms/MV1-bundle-compose.md) | meta-mvp |
| **meta-review** | Adversarial multi-lens review of **your** deliverable: P1–P3 findings, blind-spot probes, tensions between lenses, a lesson | R1 (on a user artifact), L4, J1, G6/G7 where relevant, S2 | — (lens catalog is config: `review.lenses`) | meta-review |
| **meta-craft** | A craft constitution and a checkable Definition of Done for the project | Frame, L1 (probe set), J1, R1 | [CR1 craft.dod](atoms/CR1-craft-dod.md) | meta-craft |
| **meta-compound** | A compounding map: what to capture, reuse and refresh so the next run is cheaper | Frame, S1, S2 | [CP1 compound.map](atoms/CP1-compound-map.md) | meta-compound |
| **meta-context** | A context strategy plus a droppable L0–L4 substrate for the project | Frame, G3 layer 3, S1 | [CX1 context.scaffold](atoms/CX1-context-scaffold.md) | meta-context |

**The lens in the specialists (D-072).** meta-review treats each review lens as a view of the same
artifact; where two lenses disagree, that tension (the L4 move) is shown to the user. CR1, CP1 and CX1
name their own constructs in their atom files.

Every skill keeps its v1 name, so existing users' commands still work. During the preview the suite
installs as `meta-suite-next` beside v1 (D-071).

## Two modes, for every skill (D-066)

| Mode | How to get it | What runs | Loaded |
|---|---|---|---|
| **deep** | default | Full lanes, every gate, R2. Scales with the scope of the question (`deep_narrow` or `deep_broad`), up to the caps. Deep research and deep-scope questions go here. | L0 → L3 as needed |
| **fast** | `--fast` (or "quick") | Lighter path with fewer samples, regions and review passes; no G4, no R2. **Lighter skill weight:** only L0 plus the core op is loaded, and deeper layers stay unread. | L0 + core op |

The numbers for each mode live in one place: the [modes table](03-architecture.md#modes-deep-default-and---fast-d-066).

The two modes are like Sonnet and Opus for the same skill: fast is good and cheap, deep is the full
depth. Speed is gated on fast (`--fast` p50 ≤ v1); quality is gated on deep.

## v1 elements: nothing is dropped (D-072)

| v1 element | v2 home |
|---|---|
| research pipeline, cross-verify, cited brief | meta-research / meta-lens lanes |
| credibility triage | D2 (`independence_key`, practitioner sources) and L2 (independent agreement) |
| knowledge atoms, lessons | returned to the **user** as A2's "reusable takeaways" section (meta-research, meta-review); S2 also records them in v2.0 and loads them in v2.1 (D-063) |
| meta-prompt generator, RMP loop, edit history | P1, P2, O1 stage-boundary adaptation, S1 |
| design schema, tokens, render-audit, WCAG predicates | DS1, A2, G7 (`meta rubric`), rubric |
| three-role convergence (proposer, executor, auditor); precedence accessibility > brief > taste | inside DS1 |
| MVP bundle, WOW standard, data model, meta-prompt | MV1 (meta-prompt via P1) + rubric (M8) |
| high-stakes-ambiguity agent panel | inside MV1, as an I1 classification sampled k times through J1 |
| palette tool | `meta palette` in the kernel, called by DS1, with no Python (D-004) |
| dashboard playbook, explainer | A1 forms "dashboard" and "animated explainer" (explainers through installed explainer skills when present) |
| meta-plan, work packets, audits, anti-slop | PL1, PL2, R1 |
| review lenses | config (`review.lenses`) for R1 and meta-review |
| P1–P3 findings, blind-spot probes | meta-review on R1; probes are a user-facing output section |
| craft constitution, DoD | CR1 |
| compounding pillars, capture/refresh policy | CP1, S2 |
| context policies, substrate scaffold, handoff | CX1, S1 |

## Plan B, per skill

Each specialist must beat its own v1 version in a blind comparison (see [eval](05-eval.md)). If one
doesn't, **that one skill ships as its v1 inside the v2 suite** and is fixed in v2.1. The suite
release is not held back by its weakest member.
