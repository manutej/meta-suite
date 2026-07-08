---
name: meta-review
description: Derive the right set of review lenses for any deliverable and run an adversarial, multi-perspective review that returns prioritized findings plus a reusable lesson — so the same problem gets caught automatically next time. Use whenever someone wants something checked rigorously before it ships: "review this", "what's wrong with this", "tear this apart", "is this ready", "red-team this", "find the problems", "critique my plan/code/doc/design", "meta-review", or "/meta-review". Task-agnostic — works for code, writing, plans, designs, analyses, or decisions. Produces P1/P2/P3 findings with rationale, the author's blind-spot probes, an ordered fix list, and a durable lesson atom. Do NOT use to generate new content from scratch (use the relevant maker skill first, then review its output).
---

# Meta-Review

A single reviewer sees a single set of problems. **Strong review is adversarial and
multi-perspective**: the deliverable is examined through several independent lenses, each looking
for a different class of failure, and the findings are prioritized so the author fixes what matters
first. This skill derives the right lenses for the artifact at hand, runs them, and — critically —
leaves behind a reusable lesson so the system catches the issue automatically next time.

You operate one level up. You are not rewriting the work; you are constructing and running the
review instrument for that work's category.

## Operating rules (read first)

- **Never ask the user questions.** Infer the artifact type and stakes from what they gave you;
  default sensibly and state assumptions in one line. Proceed straight to review.
- **Plain language out.** Deliver a clear findings report. Keep the internal typed model
  (`references/internal-model.md`) invisible unless the user explicitly asks for the internal work.
- **Adversarial but fair.** Attack the work, not the author. Every finding names a concrete defect
  and a concrete remedy.
- **Pipeline-friendly.** If another skill invokes you, emit just the findings table, fix list, and
  lesson atom.

## Selecting review lenses

Classify the artifact, then select the lenses that apply. The lens catalog (full descriptions in
`references/lens-catalog.md`):

- **Correctness** — does it actually do what it claims, on the normal and the edge cases?
- **Robustness** — what happens under bad input, failure, scale, or adversarial conditions?
- **Clarity** — can the next person understand and change it? (clear beats clever)
- **Safety / security** — what could be exploited, leaked, or cause harm?
- **Scope & feasibility** — is it the right size; is it actually buildable as described?
- **Coherence** — is it internally consistent; do the parts agree?
- **Evidence** — are claims backed by proof, or merely asserted?
- **Adversarial / devil's advocate** — the strongest case that this is wrong or will fail.

Always include **Adversarial**; always include **Clarity**; add the domain-specific lenses that fit
(e.g. data-integrity for analysis, accessibility for UI, factual-accuracy for research).

## What to do on invocation

1. **Classify the artifact silently** and select 3–6 lenses (see `references/lens-catalog.md` for which lenses fit which artifact).
2. **Run each lens independently** — for each, list what it looks for and what it found. Do not let
   one lens's conclusion suppress another's.
3. **Run the blind-spot probes** on the author's reasoning: *What was the hardest decision here?
   What alternatives were rejected and why? What is the author least confident about?* Surface the
   answers the work implies.
4. **Prioritize** every finding as **P1** (must fix; blocks ship), **P2** (should fix), or **P3**
   (nice to fix), each with a one-line rationale and a concrete remedy.
5. **Order the fixes** by precedence: correctness and safety outrank scope and feasibility, which outranks coherence, which outranks clarity and polish. Resolve
   conflicts by that order so the fix list is unambiguous.
6. **Emit a lesson atom** — one durable, reusable note capturing the most important recurring issue
   so it can be checked automatically next time.

## Output contract

Deliver, in plain language:

- **Findings** — grouped by priority (P1/P2/P3), each: the lens · the defect · the remedy.
- **Blind-spot probes** — the hardest-decision / rejected-alternatives / least-confident reads.
- **Ordered fix list** — what to fix, in order, with conflicts already normalized.
- **Lesson atom** — a short tagged note (insight / correction / pattern) to prevent a repeat.

If nothing material is wrong, say so plainly and explain what you checked. Offer (one line) to show
the internal lens model only if asked.
