# Compound Engineering — the compounding invariant for meta-plans

> **Self-contained pointer.** The compounding map this plan-level invariant draws on
> ships *inside this suite* as the **`meta-compound`** skill
> (`../meta-compound/SKILL.md` + `references/pillars.md`, when the full suite is installed). Read that when the TASK is a
> recurring capability; this file only states *why meta-planning cares* and *what to lift*.
> (If the workspace also has a standalone `/compound` skill, it drives the same loop
> interactively — but nothing here depends on it being installed.)

## Why this matters here

Meta-planning's "append-only edit history / Writer monad" framing (SKILL.md, the second of
the two ideas) and its kaizen loop **are** compound engineering's algebra: each cycle must
leave the next one cheaper. `meta-compound` names the law the edit history is *for* —
**non-decreasing ease** (its compounding axiom) — and the concrete arrow that enforces it:
capture a `LEARNING`, then re-read it.

## What to lift into a meta-plan

When the TASK is for an **ongoing / recurring capability** (not a one-shot build), add the
embeddable invariant as a typed work unit, not an afterthought:

> Every generated plan must name its **compound step** — the `LEARNING` it will capture and
> *which store* it goes to (project-local `docs/solutions/` + a CONCEPTS/index file, and/or a
> cross-session memory) — and its **return arrow**: which later stage reads that learning as
> grounding. A plan with no compound step is a waterfall; it fails the non-decreasing-ease
> gate and must be revised before execution.

Map this onto the existing pipeline:
- The **compound step** is a work unit appended after the work's EVALUATION stage.
- The **return arrow** is a dependency edge from that learning store back into the next
  cycle's DOMAIN MAP / PLAN inputs.
- It joins the **KAIZEN LOOP** check (generator step 7): "did this cycle lower the next
  cycle's cost, and is the evidence captured where the next cycle will find it?"

## When NOT to lift it

One-shot, non-recurring tasks have no "next comparable unit", so the compounding invariant is
inert — don't manufacture a learning store for work that will never recur (that is the
"ceremony inflation" anti-pattern). Right-sizing governs here exactly as it governs everywhere
in the suite: an honest "this task is too small for this instrument" is a passing answer.
