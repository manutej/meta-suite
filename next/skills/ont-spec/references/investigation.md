# The investigation brief and the verification plan

The spec ends by launching work it does not do itself. Two documents.

## ARCH-INVESTIGATION.md

```markdown
# Architecture investigation: <system>

## The question
Find the simplest architecture that satisfies every element in ELEMENTS.md and lets every
obligation in VERIFICATION-PLAN.md be discharged.

## Forces (from ELEMENTS.md, with counts)
- effect classes: <n pure, n deterministic, n effectful, n llm, n human>
- state owners and durability: <...>
- boundaries: <...>
- shared writers and their reasons: <...>

## Questions an architecture must answer
One per force that admits more than one answer, e.g. where durable state lives, how a
deterministic step stays replayable, how model calls are recorded and evaluated, how a human
approval waits, how contracts are enforced at each boundary.

## Candidates to evaluate
List the requester's preferred stack first, as a candidate, then one or two simpler alternatives.
For each: which forces it satisfies natively, which need custom work, which obligations its
tooling can discharge.

## Spikes that would decide
Small, time-boxed experiments, each tied to the question it settles and the evidence that would
settle it.

## Open questions carried from the spec
Every `Open` atom: what it blocks and who settles it.

## Not decided here
State plainly that this brief does not choose the architecture.
```

## VERIFICATION-PLAN.md

```markdown
# Verification plan: <system>

Budget: about 80% verification, 20% build. Build only what an obligation needs to run against.

| order | obligation | targets | method | discharged when | effort |
|---|---|---|---|---|---|

Order: proofs and model checks on the pure core first (cheapest to change now), then
differential tests against the source's own numbers, then property tests at boundaries, then
review of what cannot be mechanised.

## Evidence ledger
Each discharged obligation records the command, its output digest and the contract version.
```

## The 80:20 rule in practice

For every build task, list the obligations it unblocks; a build task that unblocks none is cut.
Track effort per obligation and per build task, and report the ratio at each milestone.
