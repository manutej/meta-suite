# Deep Domain Mapping — the Depth-Forcing Techniques

"Map the domain deeply" is an adverb until something pushes back on shallowness. This reference supplies the pushback: interrogation protocols *designed to resist* shallow answers (nothing prevents a confident shallow answer — the backstops for a gamed interrogation are the spawned probe and the sibling-task probe), a second-order pass that finds the relationships between relationships, a checkable floor (the closure check), and the sub-procedure for consuming a TASK that is itself a meta-prompt. Consult while executing generator step 2; the generator file remains canonical for *what* the domain map contains — this file is *how to get there honestly*.

## 1. Entity interrogation protocol

For every entity in the map, answer three questions. An entity whose answers are unknown goes to the unknowns ledger — that is the protocol working, not failing.

- **Lifecycle**: who creates it, who mutates it, who destroys it, and when? (An entity nobody mutates is probably a constant — demote it. An entity two units mutate is shared mutable state — it just changed your dependency graph.)
- **Failure mode**: what breaks, and for whom, if this entity is wrong, stale, or missing? (The "for whom" feeds the stakeholder map; an entity whose wrongness harms nobody is decoration.)
- **Owner**: who can authoritatively answer questions about it during execution? (No owner = every executor will guess = divergence. Name a source or spike it.)

For every relationship: **what breaks if it inverts or goes stale?** A relationship that survives inversion ("A constrains B" vs "B constrains A" — no difference) isn't understood yet; a relationship that can silently go stale needs a drift check in the kaizen loop.

## 2. The second-order pass

After mapping entity-to-entity relationships, look for relationships *between relationships*: where satisfying one constraint changes how another applies ("the auth flow constrains how rate-limit handling can retry"; "the audience's expertise level changes which invariants are negotiable"). Second-order relations are where plans break non-locally — two units, each respecting its own constraints, jointly violating a third. Record each as an invariant or as an interface contract between the lanes it couples. Most domains need only a handful; finding zero in a multi-lane plan is a flag, not a clean bill.

## 3. The closure check (the domain map's gate)

When the task map (generator step 3) exists, run: **every input and output of every work unit must be expressible in mapped entities.** Any artifact that appears in the work breakdown but not in the domain map means the map has a hole exactly where work will happen — go back, map it, re-run. The check is cheap and mechanical, and it is a *necessary condition, not a depth guarantee*: it converts "deep enough" from an adverb into a checkable floor. Known gaming mode: entity coarsening (one mega-entity makes closure trivial) — the interrogation protocol resists this (a mega-entity fails the lifecycle question coherently), and the spawned probe is the backstop. Falsifier for the floor itself: closure-passing maps whose plans still hit unmapped artifacts mid-flight. Depth beyond closure-plus-interrogation is hobby cartography; stop there.

## 4. Consuming a meta-prompt as the TASK (the meta-meta-prompt pathway)

When `task_to_plan_for` is itself a meta-prompt (e.g., the output of `/meta-prompting`), do not flatten it into a task description — its structure is load-bearing. Map its parts onto the meta-plan systematically:

| Meta-prompt part | Becomes in the meta-plan |
|---|---|
| Typed input slots | The category's **free variables** (what varies per instance the plans will serve) |
| Persona / role block | **Capability inventory** entries + quality-bar constraints |
| Procedure / rules | **Domain invariants** (what every generated plan must preserve) |
| Output contract | The **SOLUTION-stage gate** (a solution is acceptable iff it satisfies the meta-prompt's contract) |
| The meta-prompt itself | Embedded in every **ContextPacket** — it is the category definition executors reconstruct intent from (within the escalation limits of execution-handoff.md) |

Worked shape (compressed): given a meta-prompt for "step-by-step mathematical reasoning with a verifiable boxed answer," the meta-plan's free variables are the meta-prompt's `Problem` slot; its invariant is "every solution renders through the meta-prompt's solution structure"; its SOLUTION gate is the boxed-answer contract; its work units are category-level (build the solver harness, build the verifier, build the eval set — note these are *units about serving the category*, not about solving any instance). Example D in `examples.md` runs this pathway end to end.

The level-check intuition: a meta-prompt input means your meta-plan will generate plans for *operating a capability*, not for *doing a thing once*. If the work units you're writing solve an instance, you've collapsed a level.

## 5. Stakeholder mapping

Extend the declared stakeholders with the structurally absent: end users of whatever gets built, data subjects of any pipeline, future maintainers, and the user paying for the machinery itself. For each: what do they bear if the plan half-fails at its worst point, and did anyone ask them? Stakeholders feed `blast_radius` (typed as affected parties, not prose) and decide where human checkpoints are non-negotiable. A domain map whose stakeholder list contains only the people in the room fails MERCURIO before the audit starts.

## 6. Ledger discipline

- **Knowns ledger**: facts with provenance — *verified* (you probed: ran the call, read the schema, checked the file) or *sourced* (citation). "Everyone knows" is not a provenance.
- **Unknowns ledger**: each entry resolves to a **Spike** (time-boxed investigation unit whose output is a ledger update) or a **named Assumption** with an owner and a falsifier. An empty unknowns ledger on a non-trivial domain is a red flag — it usually means unknowns were smuggled into the knowns as confident prose.
- Spike-vs-assume heuristic: spike when the answer changes the dependency graph or an interface contract; assume (with falsifier) when it only tunes a unit internally.
