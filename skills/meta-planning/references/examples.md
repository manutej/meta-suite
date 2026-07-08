# Worked Instantiations — Four Dissimilar Inputs, One Procedure

These show the *shape* the procedure produces, not templates to copy. Each is deliberately compressed — a real meta-plan fills every slot in the output schema. Notice what stays constant across all four (typed units, five-criteria parallelism, fan-in units, evaluation-first gates, kaizen loop) and what changes (everything domain-shaped). That constancy is the structure-preserving move; the variation is why a generic plan template would fail all of them. Example D's input is itself a meta-prompt — the pathway the skill is named for.

---

## Example A — Multi-agent software build
**Category:** "Build a tool with a defined external interface against a third-party API" (instance: an MCP server for a ticketing system).

**Domain map highlights.** Entities: the API (its real endpoints), the MCP protocol contract, auth flow, rate limits, the agent-consumers of the tool. Critical relationship: tool ergonomics ← consumer needs, not ← API shape (the classic failure is mirroring the API 1:1). Invariant: every exposed tool must be safe to retry. Unknowns ledger: actual API pagination behavior, token refresh semantics → both become **spikes** (call the real API, record findings in knowns ledger) before any code unit starts. This removes the main source of the *invented interfaces* anti-pattern — for the surfaces the spikes actually probed; executors can still hallucinate unprobed endpoints, which is why the watchlist entry stays in their packets.

**Work breakdown / graph.**
```
S1 spike: probe API auth+pagination ──┐
S2 spike: enumerate consumer tasks ───┤→ U1 design tool surface (contract doc)
                                      │     ├─∥─ U2 build read tools ──┐
                                      │     ├─∥─ U3 build write tools ─┤→ U5 fan-in:
                                      │     └─∥─ U4 build eval suite ──┘  integration
                                      └──────────────────────────────────  + e2e gate
```
U2/U3/U4 pass criteria (a)–(d): separate artifacts, contract frozen by U1 so no ordering, each has its own check, U5 is the named fan-in. Criterion (e) — contended resources — *fails naively*: U2 and U3 both hit the same rate-limited API the domain map names. The fix is in the graph, not in hope: the rate limit is named as a shared resource, and U2/U3 get a throttle budget each (or, if the limit is tight, are sequenced despite passing a–d). This is the criterion most "sounds independent" plans miss. U4 (evals) running parallel to implementation is the evaluation-first principle made physical — the eval suite is built from U1's contract, not from the code.

**Context packet sample (U3, write tools).** Meta-plan ref; domain slice = auth findings from S1 + the write-relevant endpoints only; interface contract = U1's tool-surface doc (immutable; change = replan trigger); done = U4's write-tool evals pass + idempotency check; watchlist = code level + *invented interfaces*, *happy-path-only*, *silent swallowing*.

**Kaizen.** Cycle = one tool group; each cycle ends integrated and eval-green. Replan triggers: spike contradicts contract; eval suite fails the same tool twice.

---

## Example B — Research-to-report pipeline
**Category:** "Produce a decision-grade, multi-source research deliverable" (instance: market entry analysis as a Word report).

**Domain map highlights.** Entities: claims, sources, the decision the report serves, the audience's prior beliefs. Key relationship: every claim → ≥1 source with provenance grade; the *decision* (not the topic) determines what's in scope. Invariant: no claim in the deliverable that isn't in the verified claims ledger. Unknowns: source availability per sub-question → resolved by a breadth-scan spike before committing the outline.

**Work breakdown / graph.** Spike: breadth scan → U1 question decomposition + report contract (sections, audience, the decision served, done-rubric — *written before any drafting*). Then parallel research lanes per sub-question (safe: disjoint claim ledgers, independent verification, fan-in named) → U_f fan-in: cross-lane contradiction sweep (the unit most plans omit — sources that disagree across lanes are findings, not noise) → drafting units sequenced (narrative has ordering invariants; prose lanes don't merge well) → format unit (reads the docx skill only now — capability inventory names it).

**What the domain map changed.** In Example A, parallel lanes split by *artifact*; here they split by *question*, and the fan-in is epistemic (contradiction resolution) rather than technical (integration tests). Same slots, domain-shaped content.

**Watchlists here:** planning + design levels mostly — *padding decomposition* (research units nobody consumes), *sample-data theater*'s cousin: quoting the one supportive source and stopping.

---

## Example C — Design-system rollout
**Category:** "Define a coherent standard, then apply it across many heterogeneous artifacts" (instance: brand refresh across site, decks, docs).

**Domain map highlights.** Entities: tokens, components, the artifact inventory, the humans who must adopt the standard. Invariant: WCAG AA on every pair, every artifact. Crucial relationship: token definitions → *everything*; tokens are the most upstream mutable state in the domain.

**Work breakdown / graph.** The dependency analysis forces an unusual shape: token definition is a **strictly sequential head** (it is shared mutable state for every downstream unit — criterion (a) fails for anything running beside it), followed by a **wide parallel body** (per-artifact application lanes, each independently auditable via contrast checks), then a fan-in consistency review. A naive planner parallelizes token design with artifact work "to save time" and pays for it in N reworked artifacts — this example exists to show the five criteria *producing* the sequence/parallel split rather than intuition guessing it. (Criterion (e) also bites here: the wide parallel body shares one human brand-approver — their attention is a named contended resource, so lanes batch their approvals rather than all blocking on one person at once.)

**Kaizen.** Pilot one artifact lane to completion first (smallest batch), run the audit, fold findings into the token set *as a versioned edit*, then fan out wide. The pilot-then-fan-out move generalizes to any category with a standard-then-apply shape.

---

## Example D — the TASK is itself a meta-prompt (the meta-meta pathway)

**Input:** not a task but a meta-prompt — say, the `/meta-prompting` output for the category *"build an MCP server for any third-party API"*: typed slots (`API: spec`, `ConsumerTasks: list[str]`), a persona block (tool-design principles), a procedure (probe → design surface → implement → verify), and an output contract (a server whose every tool is retry-safe and consumer-shaped).

**Level check.** This is a meta-prompt, so the meta-plan will generate plans for *operating a capability* — "stand up the MCP-server-building capability and run it per instance" — not for building any one server. Watch the collapse risk: the moment a work unit mentions a specific API, you've dropped a level.

**Consumption mapping (domain-mapping.md §4, applied).** The meta-prompt's slots (`API`, `ConsumerTasks`) become the meta-plan's **free variables** — each instantiated plan binds them. Its persona principles become **invariants** (e.g., "tools are consumer-shaped, never API-mirrored") that every instantiated plan's SOLUTION gate enforces. Its output contract becomes the category-wide **SOLUTION gate**. The meta-prompt itself embeds in every ContextPacket's constitution slot.

**Work breakdown (category-level units — note none of them solve an instance):** U1 build the instance-intake checklist (binds free variables per instance); U2 build the probe harness (generic spike runner for criterion-(e) data: rate limits, auth); U3 build the reusable eval scaffold parameterized by `ConsumerTasks`; U4 fan-in: dry-run the whole capability on one synthetic instance and audit the run itself. The instantiated L1 plan for a *particular* API then reuses Example A's shape — which is the bridge working: meta-prompt → meta-plan (this example) → per-instance plan (Example A).

**Why this example exists.** It is the one the spec leads with and the easiest to get wrong: the gravitational pull toward solving an instance is strong, and the consumption table is what resists it.

---

## Reading these well

Apply the sibling-task test in reverse: if your generated meta-plan for a new task looks structurally identical to one of these despite a different domain map, you've copied an example instead of running the derivation — the shape must be *derived from the domain map*, and these pages are only evidence of what derivations look like.
