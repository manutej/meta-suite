---
name: meta-compound
description: >-
  Derive the core essential pillars that make any project COMPOUND — where each unit of work makes
  the next one easier instead of harder — and emit a concrete, task-agnostic compounding map for the
  project at hand. Use whenever someone wants work to accrue rather than decay: "make this project
  compound", "set up a system that gets faster over time", "where should I invest so future work is
  cheaper", "stop redoing the same work", "build a flywheel / leverage map", "turn lessons into
  reusable assets", "meta-compound", or "/meta-compound". Works for engineering, writing, design,
  research, ops, or any repeated effort. Produces a compounding loop, durable anchors, the
  highest-leverage investments, and a capture/retrieval/refresh policy — in plain language. Works
  even for a single high-stakes deliverable: compound the method, not just the artifact — leave
  behind a reusable template, checklist, or playbook for the next one like it. Do NOT use for a
  trivial one-off where no future work will ever resemble it.
metadata:
  short-description: "Make any project compound: map the loop, durable anchors and highest-leverage investments so each unit of work makes the next one easier. Use for flywheels, reuse, \"stop redoing work\"."
---

# Meta-Compound

Most work *decays*: every task makes the next one a little harder (debt accumulates, context is lost, the same problems get re-solved). **Compounding work inverts this** — each unit of work makes the next one easier. This skill takes any project and emits the concrete system that makes it compound.

You operate one level up. You are not doing the project's work; you are designing the **compounding machine** for that project's category of work, then handing back a map the user (or another agent) can run immediately.

## Operating rules (read first)

- **Never ask the user questions.** Infer the project's shape from what they gave you. If something is unknown, pick the sensible default and state the assumption in one line. Proceed.
- **Plain language out.** Deliver a usable map. Keep all internal machinery (the typed model in `references/internal-model.md`) invisible unless the user explicitly asks to see the internal work.
- **Concrete, not abstract.** Every pillar must land as a specific artifact, file, or routine *for this project* — not a slogan.
- **Pipeline-friendly.** If invoked by another skill, skip preamble and emit just the structured map.

## The compounding pillars (the durable substance)

These are the load-bearing principles. Read `references/pillars.md` for the full statements. In brief:

1. **The compounding axiom** — each unit of work must leave the system more capable than it found it.
2. **The loop is the unit** — Frame → Do → Review → **Capture** → Repeat. Capture is the step that separates compounding work from ordinary work; never skip it.
3. **Invest in the system, not just the instance** — spend a meaningful share of effort building reusable assets (templates, checks, generators, playbooks). An hour spent making a repeatable asset routinely saves many hours later.
4. **Encode judgment into durable artifacts** — taste and hard-won decisions belong in reusable files, not in your head or in one-off review.
5. **Teach once, reuse forever** — context handed to the system pays compounding dividends; re-deriving it each time does not.
6. **Build safety nets, not gatekeeping** — trust comes from verification infrastructure (automated checks, review routines, monitoring, easy reversal — whatever fits the medium), not from manual re-inspection.
7. **Parallelize independent work** — separate the work that can run concurrently from the work that must be sequential.
8. **The durable artifact is the source of truth** — a written plan/spec/anchor outlives memory and steers everything downstream.
9. **Capture learnings as retrievable atoms** — small, tagged, greppable notes that future work can find.
10. **Close the discoverability loop** — every capture must update an index so future work actually finds it; uncaptured-but-undiscoverable knowledge does not compound.
11. **Maintain against rot** — learnings decay; keep, update, replace, or archive them on a cadence.

## What to do on invocation

1. **Classify the project silently** along: unit-of-work (what repeats?), substrate (where could durable memory live?), cadence (how often does the loop run?), and leverage surface (where is effort most repeated/wasted?). Use `references/internal-model.md` for the typed classification; do not show it.
2. **Instantiate the loop for this project** — name the project's concrete Frame/Do/Review/Capture/Repeat steps.
3. **Locate the durable anchors** — the 1–3 source-of-truth artifacts this project should keep (e.g. a strategy/spec doc, a decisions log, a knowledge folder).
4. **Pick the top 3 leverage investments** — the highest-return reusable assets to build first, each with the future cost it removes.
5. **Define the capture → retrieval → refresh policy** — what gets captured, in what format (default: plain markdown atoms with lightweight tags), how it's found later, and when it's reviewed for rot.
6. **Emit the compounding map** (see Output).

## Output contract

Deliver, in plain language:

- **Compounding map** — the project's loop, its durable anchors, and its substrate, in a form they can set up today.
- **Top 3 leverage investments** — each as: the asset to build · the future cost it removes · roughly when it pays back.
- **Capture/retrieval/refresh policy** — concrete formats and cadence.
- **First move** — the single highest-leverage thing to do right now.

Offer (one line) to show the internal typed model only if they want it. Never lead with it.
