# meta-suite v2 — specification

**Status:** DRAFT v0.3 (after review rounds 1–2) · branch `v2` · v1 baseline = commit `8cdbf73`
(v1.2.0, untouched on `main`)

**Skill:** `meta-lens`. Five lanes run in order: **Frame → Map → See → Make → Check**.

**In one line:** v2 turns an unfinished question into an expert-grade, verified research view or
artifact. A category-theory lens is the executable engine: mandatory inside the machine, invisible
to the reader.

## How to read this spec (progressive disclosure)

| Tier | Read when | Files |
|---|---|---|
| L0 | always | this page |
| L1 | deciding *what* and *why* | [vision](00-vision.md) · [scope](01-scope.md) · [principles](02-principles.md) |
| L2 | deciding *how* | [architecture](03-architecture.md) · [atom index](04-atoms.md) |
| L3 | building one piece | [`atoms/`](atoms/): one file per atom, each a typed contract |
| L4 | proving it works, or auditing why | [eval](05-eval.md) · [roadmap](06-roadmap.md) · [evidence](07-evidence.md) · [open questions](08-open-questions.md) · [decisions](09-decisions.md) |

## Scope fence (anti-creep)

- **v2.0 is exactly the 25 atoms** in the [atom index](04-atoms.md). From M1, CI checks that this
  count matches across files.
- **A new atom enters only through a [decision](09-decisions.md)** that names what it displaces.
- **Everything else is in the [parking lot](01-scope.md#parking-lot)**, with a target version.
- **The spec freezes at M0.** After that, it changes only through the decision log.

## Refinement log

| Version | What happened |
|---|---|
| v0.1 | Drafted from the design conversation of 2026-10-29. |
| v0.2 | Blind review round 1: two independent critics, one on scope and coherence, one on feasibility and eval validity. All P1 findings were applied, recorded as D-017 to D-035. |
| v0.3 | Round 2: an independent audit against the working sheaf-core repo (70/70 tests, toy proof, Go port parity 11/11). G3 gets three real sheaf layers. L3 coherence, the D1 covering condition and Claim units were added (D-038 to D-046). |
| next | Fold in the M0 spike results, then freeze as **spec 1.0** at M0 exit. |
