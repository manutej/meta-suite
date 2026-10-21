---
name: meta-planning
description: >-
  Generate a customized META-PLAN — a typed, example-agnostic planning scaffold — before any
  complete plan is written, for agentic tasks executed by one or many agents. Bridges
  /meta-prompting (generic tasks) to agentic planning (ultraplan-class capability). Use whenever the
  user wants to plan work that agents will execute, orchestrate subagents, decompose a complex
  build, decide what can run in parallel vs. in sequence, hand context-complete work packets to
  fresh agents, or asks for "meta-plan", "/meta-planning", "ultraplan", "plan the plan", "deep
  planning", "map the domain before planning", or "a plan another agent can execute without me".
  Also use when the input TASK is itself a meta-prompt (operating from meta-meta-prompt space), or
  when the user wants anti-slop quality gates and kaizen-style iterative execution baked into a
  plan. Do NOT use for simple single-agent to-do lists a direct plan would serve.
metadata:
  short-description: "Write a meta-plan before the plan: a typed scaffold that agents can execute, with parallel vs sequential lanes, context-complete work packets, and quality gates. For multi-agent builds."
---

> **Meta-Suite doctrine (overrides anything below).** Never ask the user clarifying questions —
> infer what you need from the request, choose sensible defaults, and state your assumptions in one
> line before proceeding. Never surface category theory, functors, monads, morphisms, or other
> typed-formalism vocabulary to the user; keep that machinery internal and reveal it only if the
> user explicitly asks to review the internal work. Keep all user-facing output in plain language.

# Meta-Planning

The bridge from **meta-prompting** (how to think about a category of tasks) to **agentic planning** (how a fleet of agents should execute a category of work). Where `/meta-prompting` emits the meta-prompt that solves a task category, `/meta-planning` emits the **meta-plan**: the typed scaffold from which a complete, high-quality, agent-executable plan is generated — *before* that plan is written.

This skill operates in **meta-meta-prompt space**: the input TASK may itself be a meta-prompt. You are not planning the task. You are constructing the *planning instrument* customized to the task's domain — the instrument any planning agent runs to produce the complete plan, and any executing agent carries as its constitution.

## On invocation — do this

1. **INTAKE.** Take the TASK from the invocation. Fill the generator's three input slots — `task_to_plan_for`, `execution_environment` (agents, tools, MCPs, human availability), `constraints` (time, budget, irreversibility boundaries, quality bar). Do not ask the user anything; infer each slot from the task and record genuine unknowns as "unknown" honestly — the generator treats "unknown" as a first-class value, and stated assumptions stand in for answers. (Only if the user explicitly opts in — `--ask` / "ask me" — may you confirm blocking unknowns; if that tool is unavailable, fall back to inference and note it.)
2. **LEVEL CHECK.** Determine whether the TASK is an instance, a category, or itself a meta-prompt. If it is a meta-prompt, consume it per the sub-procedure in `references/domain-mapping.md` (its slots become the meta-plan's free variables; its constraints become invariants; its output contract becomes the SOLUTION gate).
3. **GENERATE.** Read `references/meta-meta-planning-prompt.md` and execute its procedure — **the generator is canonical**; the seven moves below are its summary. Emit the `generated_meta_plan` and `gap_report`.
4. **VERIFY (mandatory minimum).** Instantiate one work unit and run the fresh-agent probe on it — spawn, don't simulate; simulation is screening only (`references/execution-handoff.md` §3 is the binding form and explains why). A full L1 plan is the user-optional extra.
5. **AUDIT.** Run the tri-audit per `references/audit-protocol.md`, in its order (OPUS → FABLE → MERCURIO), at the tier the user picks — protocol defaults govern generation-time; Tier 3 is a floor, not a choice, for irreversible, externally-stakeholdered, or multi-lane work. Surface the `gap_report` — checkpoint list, tier derivation, and order-of-magnitude cost estimate — to the user **before** any execution begins; their recorded approval (by a named human, never an agent) is a typed artifact every root work unit depends on.

## Level map (hold this orientation)

- **L0 — Solution**: the work itself (code, document, system).
- **L1 — Plan**: ordered work units with owners and gates (what ultraplan/plan-mode produces).
- **L2 — Meta-plan**: the typed scaffold + procedure that *generates* L1 plans for the whole task category. Domain-customized, example-agnostic.
- **L3 — This skill**: the generator that maps any TASK (even one that is itself a meta-prompt) into its L2 meta-plan.

The deliverable is always the **L2 meta-plan**. The meta-plan and any plan instantiated from it must be *self-sufficient*: handed to a fresh agent with zero conversation history, they carry everything needed to move TASK → PLAN → PROMPT → SOLUTION → IMPLEMENTATION → EVALUATION.

## The two ideas to hold onto

> **A plan is a prompt whose audience is executors.** Structure found in the domain becomes structure in the plan: objects become typed work products; relationships become interface contracts; independence in the dependency graph becomes parallel lanes; ordering invariants become sequential gates. We borrow the word *functor* for this move — loosely, as orientation: structure in, corresponding structure out. No mathematical laws are claimed or verified; the seven moves below are the actual mechanism. A good domain map strongly constrains the plan's shape; it does not uniquely determine it — when several shapes fit, the choice is a named decision, not an accident.

> **Plan refinement is an append-only edit history.** Every execution cycle appends an edit script (what we learned, what changed, why) to the plan's history — the "Writer monad" framing, again loose. Small batches, verified increments, retro after each cycle. Iterative improvement is not a phase of the plan; it is the algebra the plan runs on. This is KAIZEN made structural.

## The pipeline (every stage has a typed output)

```
TASK → PLAN → PROMPT → SOLUTION → IMPLEMENTATION → EVALUATION
        ↑__________________ kaizen loop __________________|
```

- **TASK** — the level-checked task statement (instance / category / meta-prompt, with free variables named).
- **PLAN** — the typed work-unit graph (units, contracts, gates, lanes).
- **PROMPT** — the executor prompts: each work unit's ContextPacket rendered via the ExecutorPrompt template (`references/execution-handoff.md`) into the prompt a spawned agent actually receives.
- **SOLUTION** — per unit or lane, a `SolutionSketch`: the decided approach, key decisions traced to invariants, alternatives rejected with reasons. Decided *before* building.
- **IMPLEMENTATION** — the built artifacts, machine checks green, human checkpoints signed by their named humans.
- **EVALUATION** — done-criteria (written at PLAN stage) judged; edit script appended; drift checked.

Each arrow is a gate with a typed contract; evaluation criteria are written *before* implementation so work is graded against intent, not against whatever got built.

## The procedure: TASK → META-PLAN

This is an orientation map only — **the generator is canonical**; each move names its generator step, where the full text lives. (Fixes that must land somewhere land there, once.)

- **0. LEVEL CHECK & INTAKE** — covered in "On invocation"; maps 1:1 to the generator's input block and step 1.
- **1. DOMAIN MAP** (step 2) — entities, relationships, invariants, stakeholders, capability inventory, knowns/unknowns ledgers, with depth techniques from `references/domain-mapping.md`. Gate: the closure check — a mechanical *necessary* condition that turns "deep enough" from an adverb into a checkable floor (not a depth guarantee).
- **2. TASK MAP** (step 3) — decompose into fully typed WorkUnits. If a unit's output type can't be named, it isn't understood: split it or spike it.
- **3. DEPENDENCY GRAPH** (step 4) — partial order; parallel only when all five criteria hold (the list lives in step 4); the criteria are necessary, not sufficient. A false parallel usually costs more than a false sequence, but both cost something real. Mark the critical path.
- **4. CONTEXT PACKETS** (step 5) — fill the ContextPacket type, budget and escalation rule included. Spawn the probe, don't simulate it — simulation is screening only (`references/execution-handoff.md` §3 is the binding form).
- **5. VERIFICATION GATES** (step 6) — done-criteria written before implementation guidance; HumanChecks bind to named humans; reversibility unknown or contested defaults to human; the human's approval of the checkpoint list is a typed artifact that gates every root work unit.
- **6. KAIZEN LOOP** (step 7) — small verified batches; single accumulating edit history; the retro's three questions; replan triggers vs. stop-the-line triggers (lists live in step 7). The plan must be fireable, not only revisable.
- **7. RENDER + AUDIT** (steps 8–9) — emit the schema with the ExecutorPrompt template embedded; tri-audit per `references/audit-protocol.md` at the user's tier; stop on fixpoint or the thrash rule.

## Roles (who owns what)

- **Meta-planner** (this skill's invoker): produces the L2 meta-plan.
- **Planner**: instantiates an L1 plan from it.
- **Orchestrator**: spawns executors from packets via the ExecutorPrompt template, runs gates (including the SOLUTION gate on each executor's SolutionSketch), merges fan-ins, owns the edit history, fires replan triggers, and escalates stop-the-line triggers to the human-of-record.
- **Executor**: one work unit from one packet; obeys the escalation rule in `references/execution-handoff.md`.
- **Auditors**: the tri-audit lenses.
- **Human-of-record**: a named, reachable human bound in the meta-plan's roles at render time. Owns checkpoint approval and stop-the-line responses; a halt is indefinite until they respond — an unanswered escalation is never license to resume. If human availability is none or unknown and the plan contains any irreversible unit, the generator emits a blocked `gap_report` instead of a plan.

One agent may hold several roles on small tasks; the roles still exist so every gate, merge, and trigger has an owner.

## Bridging plan mode / ultraplan

When invoked inside plan mode or an ultraplan-class flow: the L1 plan instantiated from the meta-plan **is** the plan presented for approval; verification gates map to plan checkpoints; and the meta-plan ships as an appendix of the presented plan, so approving the plan also approves the constitution executors will inherit.

## Anti-slop: the quality floor

The meta-plan carries its own immune system: `references/anti-patterns.md` catalogs failure modes at five levels (planning, architecture, design, code, process/content), sliced per-unit into watchlists — which may and should also include domain-derived entries the catalog lacks. The unifying principle: **slop is unverified plausibility**. Every claim in a plan ("this parallelizes", "this library does X") is unverified until a gate has checked it. Plans assert; gates know.

## Quality checklist before you ship a meta-plan

- **Self-sufficient**: one work unit verified by a spawned fresh-agent probe (minimum); packets curated with exclusion notes.
- **Typed**: every work unit, artifact, contract, gate, and stakeholder is a named typed slot. No "and then we integrate everything."
- **Honest about unknowns**: non-empty unknowns ledger for any non-trivial domain; each unknown spiked or assumed-with-falsifier.
- **Safely parallel**: every parallel pair passes all five criteria with evidence; fan-ins named; contended resources named.
- **Evaluation-first**: done-criteria written before implementation guidance; HumanChecks bound to named humans.
- **Kaizen-structured**: small verified batches, single accumulating edit history, replan *and* stop-the-line triggers.
- **Handoff-complete**: ExecutorPrompt template embedded; escalation rule in every packet.
- **Audited**: tri-audit run in protocol order at the right tier; checkpoint list and costs surfaced to the user before execution.

## Failure modes specific to this level

(Execution-level slop lives in `references/anti-patterns.md`; these are the L2/L3 modes.)

- **Planning the instance, not the category** — the meta-plan only fits the one task described. Test against 2–3 sibling tasks.
- **Skipping the domain map** — decomposition without the map is guesswork dressed as structure; the closure check exists to catch this.
- **Evaluation as afterthought** — done-criteria written after implementation grade the work against itself.
- **Solving instead of structuring** — if your output contains implementation, you've collapsed to L0. Stay at L2; the gates force quality at L1/L0.
- **Scaffold-worship** — running this full machinery on a task a direct plan would serve; see Limits below and the proportionality tiers in the audit protocol.

## Limits of this method (named assumptions, with falsifiers)

This skill practices the provenance discipline it preaches. Its load-bearing assumptions are currently *unmeasured*:

- **Meta-planning beats direct planning for category-shaped, multi-agent work.** Falsifier: plans instantiated from meta-plans show no better gate pass-rates or rework rates than direct plans of equal token budget.
- **The spawned fresh-agent probe predicts real executor success.** Two layered assumptions: the probe checks *legibility* (restatement), not executability; and probing one packet licenses fan-out only if packet defects are planner-systemic rather than unit-local. Falsifiers: probe-passing packets still produce blocked or divergent executors at material rates; unprobed packets block where the probed one passed.
- **Audit loops converge.** Falsifier: thrash per the stop rule above. The "3–5 passes" figure is a working default, not a measured constant.
- **Heuristics herein** (spikes cheaper than rework; integration as primary failure site; bloat-vs-starvation symmetry) are directionally defensible folk engineering, named here so they can be falsified rather than inherited as fact.

The machinery has real costs — tri-audits, spikes, spawned probes, parallel agents — paid in the **user's** tokens and time. Disclose them; let the user pick the tier. An honest "this task is too small for this instrument" is a passing answer.

## References

- `references/meta-meta-planning-prompt.md` — **the canonical generator**: role, procedure, full type system, the meta-plan output schema. The verbatim instrument.
- `references/domain-mapping.md` — depth-forcing technique: entity interrogation, second-order relations, the closure check, stakeholder mapping, and the consume-a-meta-prompt-as-TASK sub-procedure.
- `references/execution-handoff.md` — the PROMPT stage made concrete: ExecutorPrompt template, executor output contract, escalation rules, orchestrator duties.
- `references/anti-patterns.md` — the slop catalog (five levels), sliceable into watchlists.
- `references/examples.md` — four constructive instantiations, including one whose input TASK is itself a meta-prompt.
- `references/audit-protocol.md` — the OPUS → FABLE → MERCURIO tri-audit: rubrics, proportionality tiers, edit-script format, thrash rule.
- `references/auditor-meta-prompts.md` — the FABLE and MERCURIO persona meta-prompts, hand-to-subagent verbatim.
- `references/edit-history.md` — this skill's own accumulated edit scripts: the evidence that it was built by the loop it encodes.
- `references/compound-engineering.md` — the **compounding invariant** for plans whose TASK is an ongoing capability: what to lift (the *compound step* + *return arrow* as typed work units) from the sibling `meta-compound` skill. Self-contained within this plugin; read it when the work recurs.

## Composition with sibling skills

Use `/meta-prompting` upstream when the TASK first needs lifting into a meta-prompt (this skill consumes its output — see domain-mapping.md). Use `skill-creator` downstream when a meta-plan proves reusable enough to freeze into its own skill. Use domain skills (mcp-builder, generative-ui, docx/pptx/xlsx, design skills) inside work units — the capability inventory names them so executors reach for them instead of improvising.

When the TASK is for a **recurring capability**, read `references/compound-engineering.md` and embed the compounding invariant: every generated plan must name its *compound step* (the learning it captures, and to which store) and its *return arrow* (which later stage re-reads it), or it fails the non-decreasing-ease gate. The sibling `/meta-compound` skill is the map this draws on — it designs the compounding machine; this skill embeds one invariant from it into the plan.
