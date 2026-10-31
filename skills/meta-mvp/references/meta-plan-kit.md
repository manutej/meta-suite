# Meta-plan kit (self-contained planning layer)

Embeds the planning scaffold so meta-mvp plans deeply without depending on /meta-planning.

## The one idea
A plan is a prompt whose audience is executors. Structure in the domain → structure in the plan:
entities → typed work products; relationships → interface contracts; independence → parallel lanes;
ordering invariants → sequential gates.

## Produce `04_PLAN.md` with these typed parts

### 1. Domain map
Entities, invariants, stakeholders (incl. the named human-of-record), known unknowns (each marked
SPIKE / ASSUME / ASK-{stakeholder}). A non-trivial domain with zero unknowns is hiding something.

### 2. Work-unit graph (the table)
| # | Work unit | Output type | Depends on | Gate (done when…) | Owner |
Every unit has a **named output type** (if you can't name it, split or spike it) and a
**done-criterion written before implementation**. Mark the critical path.

### 3. Parallel lanes
Two units run in parallel ONLY if they share no write target and neither needs the other's output.
A false parallel usually costs more than a false sequence. Name fan-ins and contended resources.

### 4. Human gates
Anything irreversible / externally-visible / stakeholder-facing is a checkpoint bound to a named
human. Reversibility unknown or contested → defaults to a human gate.

### 5. Kaizen loop
After each build batch, append to an edit history answering: who/what fell through? which connection
did we miss? what is the work asking for that we're not naming? Replan triggers vs stop-the-line
triggers listed explicitly.

## Sub-agent spawning protocol (this is how depth scales)
For independent lanes — especially the per-module SPEC stage — spawn one sub-agent per unit, in a
single message so they run concurrently. Each sub-agent prompt is a **context-complete packet**:
- The shared world facts (so it doesn't hallucinate the domain).
- The exact spec/template shape it must follow (point it at the relevant reference file).
- The exact output path to write to.
- An instruction to return only the path + a one-line summary.
Then the orchestrator verifies each output against its done-criterion before integrating.

## Anti-slop
Slop is unverified plausibility. Every claim in the plan ("this parallelizes", "this library does X")
is unverified until a gate checks it. Plans assert; gates know. Keep a non-empty unknowns ledger.

## Quality bar
Self-sufficient (a fresh agent could execute a unit from its packet) · everything typed · honest
unknowns · safely parallel · evaluation-first (done-criteria before implementation) · kaizen-structured.
