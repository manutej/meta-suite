# meta-suite v2 — specification

**Status:** DRAFT v0.3 (after review rounds 1–2) · branch `v2` · v1 baseline = commit `8cdbf73`
(v1.2.0, untouched on `main`)

**Suite:** a new v2 suite of 10 skills on one kernel. `meta-lens` is the front door, and the other nine keep their v1 names ([suite map](11-suite.md)). Every skill runs the five lanes **Frame → Map → See → Make → Check**, in **deep** mode by default or with `--fast`.

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
| L4 | proving it works, or auditing why | [eval](05-eval.md) · [roadmap](06-roadmap.md) · [evidence](07-evidence.md) · [open questions](08-open-questions.md) · [decisions](09-decisions.md) · [**vision ledger**](10-vision-ledger.md) · [suite map](11-suite.md) · [reviews](reviews/) |

## Scope fence (anti-creep)

- **v2.0 is a new suite** of 10 skills ([suite map](11-suite.md)) built on exactly **34 atoms** (25 core + 9 specialist) in the [atom index](04-atoms.md). From M1, `meta check` verifies that this
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
| v0.4 | M0 spike, MIT sheaf subset, plug and play (D-047 to D-054). The **vision ledger** was added, and the first **cross-tier fidelity review** (Opus, Sonnet, Haiku) all said "drifting". D-055 to D-061 apply the fixes: artifact measured and gated, speed criterion, in-run adaptation, walking skeleton, v1 inventory. |
| v0.4+ | Owner decisions D-062 to D-064 (port order, learning in v2.1, depth default). |
| next | Fidelity review round 2, then freeze as **spec 1.0** at M0 exit. |
