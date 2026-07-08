---
name: meta-design
description: Derive a typed, medium-specific DESIGN SCHEMA from any design brief, then execute it into an audited, self-contained artifact plus the design system that generates it. Use whenever the user wants to design — but design rigorously and reusably rather than by imitation — teaching material, slides, animations, web pages, dashboards, or adjacent media. Triggers include "meta-design", "/meta-design", "derive a design system for X", "design schema", "build the system then all the instances", "make a design generator not a one-off", "audited design with WCAG predicates", "tokens + component library + composed artifact", or any request to map a brief into a structure-preserving design meta-prompt. This skill operates one level up: it emits the medium-specific schema that designs the whole category the brief belongs to, executes it through a render-and-verify loop, and refines via critique edit-scripts until convergence. Do NOT use for a single ad-hoc visual with no reuse and no audit requirement.
---

> **Meta-Suite doctrine (overrides anything below).** Never ask the user clarifying questions —
> infer what you need from the request, choose sensible defaults, and state your assumptions in one
> line before proceeding. Never surface category theory, functors, monads, morphisms, or other
> typed-formalism vocabulary to the user; keep that machinery internal and reveal it only if the
> user explicitly asks to review the internal work. Keep all user-facing output in plain language.

# Meta-Design

The design-discipline sibling of `/meta-prompting`, `/meta-planning`, and `/meta-mvp`. Where meta-prompting emits the meta-prompt that solves a task category and meta-planning emits the meta-plan that an agent fleet executes, **meta-design emits the medium-specific design *schema*** — a typed, example-agnostic generator — then fills it into an audited artifact plus the system that produces all of its instances.

Treat design as a formal discipline. Briefs are objects in a category of tasks; design schemas are objects in a category of prompts; your work is the structure-preserving map between them. **Never reason by imitating reference artifacts.** Derive, execute, and refine typed schemas, and trust nothing that has not been rendered and audited.

The complete operating contract — roles, phases, typed acceptance predicates, and the output syntax — lives in `references/meta-meta-design-prompt.xml`. **Read that file in full before producing anything.** This SKILL.md is the on-ramp; the XML is the law.

## On invocation — do this

1. **INTAKE.** Collect the three input slots: `brief` (Medium × Audience × Purpose × Constraints × Success criteria), `content` (the material to be designed — lesson copy, data, text), and optional `brand_inputs` (existing palettes, tokens, assets). Do not ask the user anything; infer the slots from the brief and record genuine unknowns honestly as "unknown," stating your assumptions in one line.
2. **LOAD THE CONTRACT.** Read `references/meta-meta-design-prompt.xml`. It defines the Proposer → Executor → Auditor role sequence, the three phases, the seven typed acceptance predicates, and the `<task_schema>` output syntax. Follow it literally.
3. **RUN THE THREE ROLES IN SEQUENCE, IN ONE RESPONSE**, keeping their boundaries visible in the trace.

## The three roles

- **Proposer (Phase 1 — schema derivation).** Classify the brief as an object `Medium × Audience × Purpose × Constraints`. Instantiate the base schema and adapt it by *named, schema-preserving edits (morphisms)* — e.g. `AdaptMedium(web→slides)`. If the brief composes simpler briefs (`dashboard = data-viz ∘ layout ∘ tokens`), derive its schema as the composition of their schemas, in order — never derive monolithically what can be composed. Emit the derived schema explicitly, example-agnostic, **before** executing it.
- **Executor (Phase 2 — execution).** Resolve every component through `[Decision] → [Sketch] → [Render] → [Audit] → [Resolved]`. A render is the design analogue of code execution: no claim about a component is trusted until rendered and audited. Resolve the **system once** (tokens + component library + composition rules), then generate all instances by instantiation in a single pass. Resolved tokens are immutable downstream; failing components are excluded from the inheritance chain.
- **Auditor (Phase 3 — refinement).** Each critique pass is an edit script `σ` over schema and tokens, with stated rationale, recorded in the refinement log. Normalize conflicting edits by precedence: **accessibility predicates ⊐ brief constraints ⊐ aesthetic preference**. Converged ⇔ all predicates pass ∧ the latest pass is the identity edit. Cap at 3 passes; if unconverged, ship the best artifact with open issues stated.

## Typed acceptance predicates (compute, never estimate)

> **Non-visual artifacts:** if the brief's medium is prose, process, or policy rather than a rendered visual, replace the visual predicates (contrast, target size, breakpoints, motion) with their domain analogues — legibility, accessibility of understanding, and fitness for use — and audit those instead. The render-and-verify discipline is the point; the specific predicates adapt to the medium.

Contrast (WCAG ≥ 4.5 / ≥ 3.0 large), interactive target size (≥ 44px preferred, never < 24px), motion (`prefers-reduced-motion` fallback, ≤ 3 flashes/sec, easings from MotionCurve tokens), hierarchy (strictly nested headings, reading order ≡ source order), rhythm (spacing ∈ SpacingRhythm, body measure 45–75ch), responsive (no overflow at any breakpoint), and — for teaching media — pedagogy (one concept per view, ≤ 7 new elements, progressive disclosure). The full statements are in the XML; report each as a computed pass/fail in the audit table.

## Optionally compose with design skills (when available)

This skill is self-contained and produces its full deliverable on its own. When these companion design skills happen to be installed, it can optionally orchestrate them for higher-fidelity tokens — but it never depends on them and degrades gracefully without them:

- **Color / palette:** `dict-color`, `owala-design:owala-colors`, `owala-design:dark-mode-pairing`
- **Tokens / rhythm / motion:** `owala-design:design-tokens`, `owala-design:layout-rhythm`, `owala-design:motion-curves`
- **Contrast audit:** `owala-design:wcag-contrast`, `design:accessibility-review`
- **Data encoding:** `bertin` (pick visual variables from the data's structure)
- **Medium build skills:** `frontend-design`, `shadcn-concept-page`, `editorial-dashboard` / `milton`, `ceti-explainer` (animation), `motion-stack` / `motion-media-l2`, `pptx` (slides), `docx` / `pdf` (print), `generative-ui`
- **Critique:** run the medium's critique loop (e.g. the milton agent) as part of Phase 3.

## Output contract

Deliver exactly what `<task_schema>` specifies: **(a)** the design system — full typed token set (`Palette: list[OKLCH]`, `TypeScale: list[rem]`, `SpacingRhythm: list[rem]`, Radius/Shadow/Z, `MotionCurve: cubic-bezier`, `Breakpoints: list[px]`) plus the component library; **(b)** the code that generates ALL instances from the system; **(c)** the self-contained rendered artifact in the medium-appropriate format (HTML/React/SVG/PDF/.pptx); **(d)** the audit report as a predicate table with computed values and pass/fail. Keep the derived schema, execution trace, and refinement log visible so the derivation is auditable.
