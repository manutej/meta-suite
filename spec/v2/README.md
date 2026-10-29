# meta-suite v2 — specification

**Status:** DRAFT v0.1 · branch `v2` · v1 baseline = commit `8cdbf73` (v1.2.0, untouched on `main`)

**In one line:** v2 turns an unfinished question into an expert-grade, verified research view or
artifact. It runs a category-theory lens as an executable engine: mandatory inside the machine,
invisible to the reader.

## How to read this spec (progressive disclosure)

| Tier | Read when | Files |
|---|---|---|
| L0 | always | this page |
| L1 | deciding *what* and *why* | [vision](00-vision.md) · [scope](01-scope.md) · [principles](02-principles.md) |
| L2 | deciding *how* | [architecture](03-architecture.md) · [atom index](04-atoms.md) |
| L3 | building one piece | [`atoms/`](atoms/) — one file per atom, each a typed contract |
| L4 | proving it works / auditing why | [eval](05-eval.md) · [roadmap](06-roadmap.md) · [evidence](07-evidence.md) · [open questions](08-open-questions.md) · [decisions](09-decisions.md) |

## Scope fence (the anti-creep rule)

- v2.0 is exactly the **25 atoms** in the [atom index](04-atoms.md). Nothing else ships in v2.0.
- A new atom enters only through a [decision](09-decisions.md) that names what it displaces or defers.
- Everything else lives in the [parking lot](01-scope.md#parking-lot) with a target version.
- The spec is DRAFT until milestone M0 exits. After that it changes only through the decision log.

## Refinement process for this spec

1. v0.1: drafted from the design conversation of 2026-10-29.
2. Up to two blind adversarial review rounds. Fresh-context critics return edit scripts, and every
   P1 finding is applied or rejected with a reason in the decision log.
3. Freeze as **spec 1.0** at M0 exit.
