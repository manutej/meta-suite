# Execution Handoff — the PROMPT Stage Made Concrete

The pipeline's PROMPT stage converts each `ContextPacket` into the prompt a spawned agent actually receives, with a defined output contract and escalation behavior. This file is the mechanism behind "the meta-plan can be provided to EACH new agent with specified context." The ExecutorPrompt template below is also **embedded (instantiated) inside every generated meta-plan** — the meta-plan alone must suffice to spawn agents.

## 1. The ExecutorPrompt template

```xml
<executor_prompt unit="[WorkUnit.id]">
<role>
You are the executor for one work unit of a larger plan. Your job is exactly
[WorkUnit.intent], positioned in the graph as: [upstream units you consume],
[downstream units that consume you], [your lane and fan-in, if parallel].
You build only what your packet specifies. Quality bar: your output is unverified
until your done-check passes — plans assert; gates know.
</role>

<constitution>
[Compact meta-plan: the pipeline contract, domain invariants, relevant stakeholders,
 kaizen rules, and — if the category was defined by a meta-prompt — that meta-prompt
 verbatim. This is what you reconstruct intent from, within the escalation limits.]
</constitution>

<packet>
<domain_slice>[the unit's slice of the domain map]</domain_slice>
<inputs>[Artifacts, with locations and schemas]</inputs>
<interface_contracts>[contracts this unit must honor; these are FROZEN — see escalation]</interface_contracts>
<done_check>[the executable MachineCheck and/or the named-human HumanCheck + rubric]</done_check>
<watchlist>[anti-pattern entries sliced for this unit, incl. domain-derived ones]</watchlist>
<budget>[ContextBudget, including the exclusion note — what was deliberately not given to you]</budget>
</packet>

<escalation_rule>
If you hit ambiguity touching an interface contract, a domain invariant, or anything
irreversible: HALT and emit blocked-status with your specific question. Reversibility
unknown or contested = irreversible — the plan-time default-to-human rule applies at
runtime too. Otherwise: take the smallest reversible assumption and record it in your
flags. When local instructions run out entirely, the licensed move is stop-and-escalate
(a replan trigger), NOT reconstruct-intent-and-proceed; reconstruction from the
constitution is permitted only for reversible, in-packet-scope decisions.
</escalation_rule>

<output_contract>
<solution_sketch>[FIRST deliverable, before building: SolutionSketch — approach, key
        decisions traced to invariants, alternatives rejected. The orchestrator gates
        this against the SOLUTION-stage contract before the unit proceeds to build.]</solution_sketch>
<outputs>[each declared Artifact, with done-check evidence attached]</outputs>
<flags>[assumptions taken; watchlist conflicts encountered (flag, don't silently obey
        or ignore); candidate edit-script entries for the plan's history]</flags>
<unknowns_encountered>[anything that belongs in the unknowns ledger]</unknowns_encountered>
<status>done | blocked{question} | failed{what the gate caught}</status>
</output_contract>
</executor_prompt>
```

Render rules: fill every slot from the meta-plan and packet — an empty slot is a finding against the *planner*, not the executor; keep the constitution compact (it is a reference for intent, not a transcript — context dumping in the constitution slot is still context dumping).

## 2. Orchestrator duties

The orchestrator is the owner the kaizen machinery acts through. It: spawns executors by rendering this template (one unit, one packet, one agent); gates each executor's `solution_sketch` against the SOLUTION-stage contract before the unit builds; runs gates on returned outputs — a `done` status without done-check evidence is treated as `failed`; merges fan-ins (the fan-in unit is a real unit with a real contract, not a courtesy); appends executor flags and retro answers to the meta-plan's `audit_record` (the single history); fires **replan triggers** (route to planner) and **stop-the-line triggers** (halt all lanes, escalate to the **human-of-record** with the state of every in-flight unit — the halt is indefinite until that human responds; an unanswered escalation is never license to resume). The first spawn is structurally gated: the human's recorded approval of the `gap_report` checkpoint list is a typed Artifact (contract: `recorded_by: named human, never an agent`) declared as an input to **every root work unit** — each unit with no upstream dependency — so no entry point of the graph, including parallel root spikes, can spawn before approval exists.

## 3. The fresh-agent probe (binding form — single home for the rationale)

The PROMPT-stage gate requires at least one unit's rendered ExecutorPrompt to be given to an actually context-isolated agent, which must restate: its job, its inputs, its done-check, and its neighbors — from the prompt alone. Reaches for outside knowledge are findings against the packet.

Why spawned, not simulated: a mind with full project context simulating freshness systematically tends to miss its own tribal knowledge (the curse-of-knowledge pattern) — simulation is a useful screening pass and an unreliable gate. Honest scope: the probe checks *legibility* (can the packet be restated), not *executability* (whether the work succeeds) — the done-check covers that later. And probing one packet licenses fan-out only under a named assumption: packet defects are planner-systemic (one bad habit poisons all packets) rather than unit-local. Both assumptions carry falsifiers in SKILL.md's Limits section. Scale probes with tier: Tier 1–2, probe the unit whose half-failure hurts most (consult `blast_radius`); Tier 3, probe one unit per parallel lane-set.

## 4. Degenerate cases (so the machinery scales down honestly)

- **Single-agent execution**: one agent holds orchestrator + executor roles serially. The template still applies — it becomes the per-unit working brief, and gates still run between units. What you lose is parallelism, not discipline.
- **No-subagent environments**: the fresh-agent probe degrades to the screening simulation plus a flag in `limits_and_assumptions` that the binding probe was unavailable — disclosed, not skipped silently.
- **Human executors**: a packet rendered for a human is the same packet; the done-check and escalation rule don't change. If the plan can't brief a human, it couldn't brief an agent either.
