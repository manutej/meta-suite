# Anti-Pattern Catalog — Preventing AI-Generated Slop

Slop is **unverified plausibility**: output that pattern-matches quality without having earned it through a check. This catalog exists to be *sliced into context packets* — each work unit's `anti_pattern_watchlist` carries only the entries relevant to its level. The entries are constructive: each names the failure, why it happens, and the move that prevents it. They are guardrails, not straitjackets — when a watchlist entry conflicts with the domain's actual needs, say so in the edit script rather than silently complying.

## Planning level

**The 90% plan.** Early units are detailed, late units say "then integrate and polish." Integration is where projects die; the last units need the strongest contracts. *Move:* write the fan-in units first, with full contracts, then work backward.

**Optimistic parallelism.** Units declared parallel because they sound unrelated. *Move:* run the five criteria (no shared mutable state, no ordering invariant, independently verifiable, named fan-in, no contended finite resource — rate limit, quota, budget, single human's attention) on every parallel pair; record the verdicts, and any shared resource by name, in the graph. (Full text: generator step 4.)

**Context dumping.** Every agent gets the full history "to be safe." Agents reason worse with bloated context — relevant signals drown. *Move:* curate the domain slice; ask "does this change the agent's decisions?" for every inclusion.

**Phantom capability.** The plan assumes tools, access, or skills the environment lacks. *Move:* the capability inventory is verified, not assumed; "unknown" entries get spikes.

**Plan-as-scripture.** No replan triggers; reality diverges and agents keep executing the original. *Move:* explicit drift checks each cycle; the plan is a versioned artifact with edit history.

**Padding decomposition.** Work units invented to make the plan look thorough (a "research best practices" unit with no consumer for its output). *Move:* every unit's output must be an input to another unit or to the deliverable; orphan outputs mean the unit is decorative.

## Architecture level

**Speculative generality.** Abstractions, plugin systems, and config layers for futures nobody asked for. AI generators love these because they pattern-match "good engineering." *Move:* every abstraction must have ≥2 concrete consumers in the current plan; otherwise inline it.

**Framework reflex.** Reaching for a heavy framework/queue/microservice split when a function would do. *Move:* state the simplest architecture that satisfies the invariants; require justification (in the plan, auditable) for each step up in complexity.

**Distributed monolith.** Parallel-built components that secretly share state or deploy order. The architecture mirrors a bad dependency graph. *Move:* interface contracts between lanes are written and frozen *before* lanes start; changes to a contract are replan triggers, not quiet edits.

**Resume-driven novelty.** Choosing tech because it's current rather than because the domain map demands it. *Move:* every technology choice cites the domain entity/invariant/constraint that motivates it.

**Invented interfaces.** Code written against APIs as the generator wishes they were (hallucinated methods, parameters, response shapes). *Move:* probe the real interface (call it, read its schema) before any unit builds against it; record findings in the knowns ledger.

## Design level

**Generic-ness.** The default-gradient, default-font, lorem-shaped output that says nothing about this domain. *Move:* design decisions trace to the brief's audience and purpose; if a decision would fit any project, it hasn't been made yet.

**Decoration over hierarchy.** Visual effects added where information structure is missing. *Move:* establish reading order, grouping, and emphasis first; anything decorative must survive the question "what does the user understand faster because of this?"

**Unverified accessibility.** "Accessible" asserted, never measured. *Move:* contrast, target sizes, reduced-motion, and keyboard paths are machine checks in the gate, not adjectives in the plan.

**Sample-data theater.** Mockups beautiful with three short perfect rows, broken by real data (long names, empty states, 10k rows). *Move:* gates test with adversarial data: empty, overflow, unicode, error states.

## Code level

**Plausible-but-untested.** Code that reads correctly and was never run. In our experience the most consequential slop category (named assumption, not a measurement). *Move:* a unit's `done` is never "code written" — it is "check passed"; every code unit names its executable check.

**Happy-path-only.** No failure handling because the prompt described success. *Move:* failure semantics (idempotency, retry, rollback, blast radius) are part of the work unit type — the slot forces the thought.

**Comment inflation & narration.** Comments restating the line below, docstrings padding for thoroughness. *Move:* comments explain *why*, never *what*; gates may include a comment-density review for generated code.

**Duplicated near-twins.** The generator re-derives a helper slightly differently in each unit because no unit could see the others. *Move:* shared utilities are named artifacts in the plan with one owner unit; context packets reference them by contract.

**Dead scaffolding.** Unused imports, vestigial parameters, half-wired options left from the generation process. *Move:* a cleanup check (lint, dead-code detection) sits in the implementation gate, mechanical and non-negotiable.

**Silent swallowing.** `try/except: pass` and defaulted fallbacks that convert errors into wrong-but-quiet behavior. *Move:* errors either propagate with context or are handled with a decision recorded in the unit's notes; "handled" without a decision is swallowing.

## Process / content level

For work units that are not software — research, reports, operations workflows, approvals, communications. (Examples B and C in `examples.md` live mostly here.)

**Approval theater.** A "review" gate where the reviewer rubber-stamps because the rubric is vague or the reviewer lacks standing. *Move:* every HumanCheck names its human and its rubric; a gate nobody could fail is not a gate.

**Citation laundering.** A claim cites a source that cites a source that turns out to be speculation, or quotes the one supportive source and stops. *Move:* claims carry provenance grades; contradiction sweeps are a planned fan-in unit, not a virtue.

**Metric gaming.** The done-check measures what's easy (word count, section presence) rather than what the stakeholder needs, so units optimize the proxy. *Move:* done-criteria trace to the decision or stakeholder need the deliverable serves.

**Unverified sign-off.** "Stakeholders aligned" asserted because nobody objected in the window. *Move:* consent status is a recorded field per stakeholder, not an inference from silence.

**Cargo-cult cadence.** Ceremonies (standups, retros, status docs) preserved as form after their information content died. *Move:* each kaizen-loop element names what decision it feeds; feed nothing, cut it.

## Using this catalog

When assembling a context packet, copy in only the relevant level(s) — a documentation unit doesn't need the code list. Watchlists **may and should include domain-derived entries this catalog lacks**: derive them from the domain map's invariants ("what would violating this invariant look like as a plausible-seeming output?"). When auditing (see `audit-protocol.md`), OPUS walks this catalog explicitly. When an executing agent hits a watchlist conflict, the licensed behavior is defined in its output contract (`execution-handoff.md`): flag it in the `<flags>` section — silently obeying and silently ignoring are both slop with extra steps.
