---
name: ont-spec
description: >-
  Turn a plain-language description of a system into an ontological specification: typed atoms
  in a refinement tree four levels deep, checked for composition at every level, looped until it
  converges, then mapped to typed contracts, the architecture forces they imply, and a brief that
  launches an architecture investigation with a verification-first plan. Use when someone wants
  to specify before building: "spec this before we write code", "ontological spec", "ont-spec",
  "/ont-spec", "make a spec we could regenerate the code from", "reverse engineer this paper or
  product into a spec", "what would we need to build X", "prove the design before coding". The
  spec is the living artifact a team approves and code is regenerated from. Never decides the
  architecture. Do NOT use for a small change to existing code or a task whose design is settled.
metadata:
  short-description: "Specify before you build: typed atoms, checked composition at four depths, convergence, typed contracts, architecture forces and a verification-first investigation brief."
---

# Ont-spec

Code is cheap now; a wrong design is not. An ontological spec states **what exists, what can
change it, and what must always hold**, as typed atoms that refer to each other. Writing it that
way surfaces bugs prose hides: a value with two definitions, state with two owners, an input
nothing produces, a "pure" step that writes, a sub-part leaning on a sibling's private detail, a
claim with no way to check it. The spec, not the code, is the artifact a team approves, keeps,
and regenerates code from.

You operate one level up: you build and verify the specification; you do not write the system,
and you do not choose its architecture. You end by launching the investigation that will.

## The pipeline

Work in a spec folder (`specs/<id>/`). `ontology.yaml` is the only hand-edited file; everything
else is generated or logged. The checker is `scripts/ontspec.py` (stdlib plus PyYAML).

**S0 · Pin the sources.** List every source (paper with version, product, brief, conversation)
in `00-sources.md` with an id and, where possible, a digest. Every later atom cites one or is
marked `inferred` or `open`.

**S1 · Verbal spec (depth 0).** Write the system as one plain paragraph and a numbered list of
linear statements, each with its source (`01-verbal.md`, and `spec.verbal` in the ontology).
This is what a stakeholder would say. Do not structure it yet.

**S2 · Ontology (depth 1).** Turn each statement into typed atoms under one `System` root.
Kinds: `Component`, `Entity`, `Value`, `State`, `Relation`, `Event`, `Boundary` (things that
exist or cross in), `Operation` (what changes things), `Invariant` (what must always hold),
`Obligation` (how a claim will be verified), `Open` (an unknown that blocks something and has an
owner). Fields and examples of each shape: `references/schema.md`.

**S3 · Verify composition.** Run both checks; read every finding.
- `ontspec.py check ontology.yaml`: the ontology laws O1–O11 (shape, typed references,
  provenance, one owner per state, coverage by obligations, every input produced, one value per
  key, leaf atomicity, refinements realised, no reaching into a sibling branch, honest effects).
- `ontspec.py operad ontology.yaml`: projects the tree onto the operadic-stitch engine (each
  refinement becomes a separate brief under a depth frame, inheriting its ancestors' definitions)
  and runs it: typed plugging, closed ports per refinement, one value per key across every depth
  seam, bounded depth, shrinking budgets. If the engine is not installed, say the check did not
  run; never report it as passed.
Each law, the bug it catches and how to resolve it: `references/laws.md`. A RED is never
"fixed" by editing the checker or silencing a finding; it is resolved in the ontology or turned
into an `Open` atom with an owner.

**S4 · Deepen and loop.** Refine to depth 2 (sub-mechanisms), 3 (algorithms and data shapes),
then 4 (atomic contracts), re-running S3 after each level. Every pass appends an edit script to
`CONVERGENCE.md`: what changed, which findings it closed, which it opened. Stop when the loop
converges: no RED, every leaf atomic (`check --converging` GREEN), the operad check GREEN, and
`ontspec.py delta` between the last two passes reports no structural change. If a branch will
not converge within four levels, it is a separate system: split the spec. Depth meanings,
atomicity and the stopping rule: `references/loop.md`.

**S5 · Typed contracts.** `ontspec.py contracts ontology.yaml --out contracts/` writes one
contract per leaf operation and invariant: signature, reads and writes, requires, ensures,
errors, and the obligations that verify it. Contracts are language-neutral; how they project to
proof tools, typed languages and durable workflows is in `references/contracts.md`.

**S6 · Elements.** `ontspec.py elements ontology.yaml > ELEMENTS.md` lists the forces any
architecture must satisfy: operations by effect class (pure, deterministic, effectful, model
call, human), state and its owners, durability, boundaries, obligations by method, and the open
questions that block each. Forces, not decisions.

**S7 · Launch the investigation.** Write `ARCH-INVESTIGATION.md` and `VERIFICATION-PLAN.md`
from the templates in `references/investigation.md`: the questions an architecture must answer,
candidate stacks to evaluate (the user's preferred stack first, as a candidate to test, not a
choice), the spikes that would decide, and a plan that spends about 80% of effort on
verification and 20% on building. Render the approval document with
`ontspec.py render ontology.yaml > SPEC.md`.

## Rules that bind the stages

- One source of truth: `ontology.yaml`. SPEC.md, contracts and elements are regenerated, never edited.
- Every atom is `stated` (with a source span), `inferred` (our reading, said so), or `open` (with
  what it blocks and who settles it). Never upgrade an inference to stated.
- A value has one definition. Refinement may add detail, never change a parent's value.
- A leaf operation writes at most one thing and declares its effect honestly.
- Every leaf operation and every invariant has at least one obligation before the loop can converge.
- Unknowns are atoms, not prose: they block named atoms and have an owner.
- Report findings as `LEVEL @ where · law · atoms · question`, and never resolve one silently.
- Do not choose the architecture, language or runtime. Record preferences as candidates.

## Output contract

```
specs/<id>/
  00-sources.md          pinned sources
  01-verbal.md           depth-0 statement list
  ontology.yaml          THE spec (typed atoms, depth <= 4)
  certificates/pass-N/   check.json, operad.json, delta.json for every pass
  CONVERGENCE.md         the edit script per pass, and why the loop stopped
  contracts/             generated: contracts.json + one .md per contract
  ELEMENTS.md            generated: architecture forces
  ARCH-INVESTIGATION.md  the brief that starts the architecture search
  VERIFICATION-PLAN.md   obligations, methods, order, the 80:20 budget
  SPEC.md                generated: the document a team approves
```

Done means: the last pass is converged with both certificates GREEN (or AMBER with every
finding explained in CONVERGENCE.md); every `Open` atom is listed in ARCH-INVESTIGATION.md
with what it blocks; and every contract names at least one obligation.

## Composes with

`meta-atomize` for reconciling several sources before S2; `operadic-stitch` (the engine behind
`operad`); `meta-planning` to turn the investigation brief into agent work packets; `meta-review`
for an adversarial pass on the converged spec. The category-level meta-prompt this skill
implements, for reuse in any tool: `references/meta-prompt.md`.
