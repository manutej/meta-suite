# The Tri-Audit Protocol — OPUS / FABLE / MERCURIO

Three independent audit lenses applied to a draft meta-plan (or an instantiated plan) before it ships. Each lens has a distinct failure class it is best at catching; running all three is *designed* to cover truth, craft, and rigor without one perspective's blind spots dominating — a design intent, not a measured guarantee (falsifier: finding classes that recur post-ship despite passing all three lenses). Each pass emits an **edit script**; scripts accumulate in the meta-plan's single `audit_record`. (Full persona meta-prompts for FABLE and MERCURIO — hand-to-subagent verbatim — are in `auditor-meta-prompts.md`; **OPUS has no persona file: its rubric section below is its verbatim subagent prompt**, with the artifact, its declared category, and the anti-pattern catalog (`anti-patterns.md`) as inputs.)

Run lenses as **separate passes with separate mindsets** — ideally separate subagents, each given only the draft and its own lens rubric (not the other lenses' findings, until merge). Named assumption: auditors who see each other's findings anchor and converge prematurely — plausible mechanism, not a measured fact; falsifier: blinded and unblinded lens runs producing equivalent finding sets.

## Proportionality — the user picks the tier, because the user pays

The full machinery (three lenses × subagents × multiple passes) costs real tokens and time, and that cost lands on the user. Division of authority: the tier *defaults* below govern the generation-time audit; the user's informed choice — made from the gap_report's cost estimate and tier derivation — governs everything from pre-execution onward and may force a re-audit at a higher tier. Because tier-determining facts (unit count, irreversible actions, lane-sets, external stakeholders) are self-assessed by the artifact being audited, the gap_report must show the derivation, and any "unknown" among those facts escalates one tier.

- **Tier 1 — self-audit** (default for small plans — working-default boundary, not a measured threshold: ≤5 work units, no irreversible actions, single domain, no parallel lane-sets): one combined pass, all three rubrics applied by the planning agent itself in one sitting. Cheap, known-weaker (self-review shares the drafter's blind spots — disclosed, not hidden).
- **Tier 2 — sequential lenses** (default for mid-size plans): three separate passes by the same agent in fresh sittings, protocol order, separate edit scripts.
- **Tier 3 — independent subagents** (required when irreversible actions, external stakeholders, or >1 parallel lane-set exist — and whenever the user asks): one subagent per lens, blinded, merged by the orchestrator.

Silently skipping the audit because the full tier felt heavy is the worst outcome; a disclosed Tier 1 beats an undisclosed nothing.

---

## MERCURIO — truth, ethics, impact

*The question it holds: should this plan exist in this form, and is it honest?*

- **Veracity plane.** Walk every factual claim in the meta-plan: domain map entries, capability inventory, knowns ledger. Each claim has provenance (verified probe, cited source) or is moved to the unknowns ledger. Confident unsourced claims are the finding — not their truth or falsity.
- **Ethics plane.** Identify irreversible actions, data exposure, consent boundaries, and affected parties who aren't at the table. Check human checkpoints sit at every irreversibility boundary. Check the plan doesn't launder a consequential decision through "the agent decided."
- **Impact plane.** For each stakeholder the domain map names (and any it should have named): who carries the cost if a unit half-fails? Are blast radii honest? Does the failure semantics slot actually protect the people, not just the artifacts?

**MERCURIO passes when:** no unsourced claim survives outside the unknowns ledger; every irreversible edge has a human gate; no unexamined stakeholder.

## FABLE — coherence, simplicity, story

*The question it holds: can a fresh mind carry this plan, and is it the simplest plan that could work?*

- **Narrative integrity.** Read the meta-plan start to finish as a story: does each unit's *intent* follow from the domain map? Could you explain the whole plan aloud in two minutes without notes? Units that can't be motivated in one sentence are decoration or confusion.
- **Simplicity pressure.** Hunt speculative generality at the planning level: abstractions with one consumer, lanes that exist for symmetry, gates that check nothing real. Every removal proposed is an edit-script entry. The kaizen test: does the simplest version ship first, with refinement planned rather than apologized for?
- **Legibility for fresh agents.** Sample 2–3 context packets and run the fresh-agent test ruthlessly, *as* the fresh agent: attempt to restate the unit's job, inputs, done-check, and neighbors from the packet alone. Every reach for outside knowledge is a finding.

**FABLE passes when:** the two-minute retelling works; no unit lacks a one-sentence motivation; sampled packets pass the fresh-agent test cold.

## OPUS — technical depth and rigor

*The question it holds: will this plan survive contact with execution?*

- **Graph audit.** Verify acyclicity; verify every parallel pair against all five criteria with evidence, not assertion — probe contended resources (criterion e) especially, since they pass the other four silently; verify every lane-set names a fan-in with a real contract; verify the critical path is marked and plausible.
- **Contract audit.** Every artifact typed; every interface contract specific enough that two agents implementing opposite sides independently would interoperate; every done-check actually executable (run a sample machine check if possible).
- **Failure audit.** For each unit: idempotency claim tested in reasoning, rollback named, blast radius stated. Find the unit whose half-completion hurts most and check the plan survives it.
- **Anti-pattern sweep.** Walk `anti-patterns.md` level by level against the draft. Check watchlists in context packets are sliced correctly (relevant entries present, irrelevant absent).
- **Sibling-task probe.** Instantiate the meta-plan mentally against 2–3 different tasks in the category. Where it only fits the original task, the meta-plan has overfit — generalize the slot or narrow the declared category.

**OPUS passes when:** graph, contracts, and failure semantics survive adversarial reading, and the sibling probe shows category fit.

---

## Edit-script format

Each pass appends one script. Keep the full history in the meta-plan's `audit_record` — the history is how nested refinements compose instead of thrashing.

```
EDIT SCRIPT — [lens] pass [n] — [date]
VERDICT: pass | fail (n findings)
1. [slot/unit affected] — [finding] → [concrete change]
2. ...
DEFERRED: [findings acknowledged but deliberately not acted on, with reason]
```

The `DEFERRED` section matters: an audit that forces every finding into a change produces thrash; an audit with silent omissions produces drift. Name what you're not fixing.

## Sequencing, fixpoint, and the thrash rule

1. Draft → **OPUS** first (no point ethics-checking a graph that's wrong) → apply script.
2. → **FABLE** (simplify before truth-checking; fewer claims to verify) → apply script.
3. → **MERCURIO** last (audit the plan that will actually ship) → apply script. Running ethics last is a deliberate trade — it audits the real shipping artifact — and its known risk (the last pass is the one skipped under budget pressure) is mitigated by the pre-execution human gate: the checkpoint list reaches the human regardless.
4. If any pass produced findings, re-pass the *earlier* lenses on the changed slots only. Stop when a full cycle returns three empty scripts — the fixpoint — or when remaining findings are all DEFERRED with reasons the user has seen.

**Convergence is assumed, not guaranteed.** "3–5 passes" is a working default with no measured provenance. Hard rules: set a pass budget up front (default: 6 passes total); if two consecutive full cycles produce non-empty, *non-overlapping* edit scripts, the loop is thrashing — stop, ship with the open findings in DEFERRED, and tell the user convergence failed. An honest non-convergence report beats a forced fixpoint.

At Tier 3, use a stronger model for the lenses where available, and merge scripts yourself. The independence cost is the user's to accept — that's why the tier is their call.
