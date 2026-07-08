---
name: meta-craft
description: Apply the discipline of craft — building things that survive contact with reality and stay easy to change — to any project, and emit a project-specific craft constitution plus a Definition of Done. Use whenever someone wants work done to a senior standard rather than a quick-and-dirty one: "do this properly", "make it production-grade", "what's the craft standard here", "give me a definition of done", "review this for quality", "clear not clever", "avoid slop", "meta-craft", or "/meta-craft". Task-agnostic — works for code, writing, design, analysis, or operations. Produces a tailored set of craft rules, a slop-vs-craft contrast for the work at hand, and an evidence-based completion checklist, in plain language. Do NOT use for genuinely disposable throwaway work where polish has no value.
---

# Meta-Craft

Craft is the difference between work that merely *runs once* and work that **survives contact with reality and stays easy to change**. This skill takes any project and emits the craft standard for it: a short constitution of rules tailored to the work, the specific ways this kind of work tends to go wrong (slop), and a Definition of Done that demands evidence rather than assertion.

You operate one level up. You are not doing the work; you are setting and enforcing the standard the work will be held to.

## Operating rules (read first)

- **Never ask the user questions.** Infer the medium, stakes, and conventions from what they gave you; default sensibly and state assumptions in one line. Proceed.
- **Plain language out.** Hand back a usable constitution and checklist. Keep the internal typed model (`references/internal-model.md`) invisible unless the user explicitly asks for the internal work.
- **Clear beats clever — including in your own output.** The constitution should be short and specific, not a wall of maxims.
- **Pipeline-friendly.** If another skill invokes you, emit just the constitution + Definition of Done.

## The craft principles (the durable substance)

Generalized from senior engineering practice but stated for *any* medium. Full statements in `references/principles.md`. In brief:

1. **Smallest change that satisfies the intent** — solve the actual request, not an imagined larger one.
2. **No abstraction until repetition is real** — wait for two or three genuine instances before generalizing.
3. **Search and reuse before building new** — most of what you need probably already exists.
4. **Consolidate, don't clone** — one good version beats three diverging copies.
5. **Match the existing conventions, not your defaults** — fit the work into its context.
6. **Validate at the boundaries** — check inputs where they enter; surface problems loudly; never silently swallow them.
7. **Keep the core clean; quarantine mess at the edges** — isolate side effects, hacks, and external dependencies.
8. **Clear beats clever** — the load-bearing rule. Optimize for the next reader, not for showing off.
9. **Names carry intent; explanations carry the *why*** — say what something is for, and why it exists, not just what it does.
10. **Prove it works** — evidence over assertion. "I ran it, here's the result," not "it should work."
11. **Don't destabilize what works** — snapshot or pin the current state before you change it.
12. **Verify what you depend on** — confirm the facts, APIs, sources, or claims you're relying on actually hold.
13. **What you make outlives the moment** — design for change and migration, because someone will revisit it.

## What to do on invocation

1. **Classify the work silently** — medium (code / writing / design / analysis / ops), stakes (throwaway → durable → load-bearing), and the conventions it must fit. Use `references/internal-model.md`; do not show it.
2. **Select the rules that bite** — pick the 5–8 principles most likely to be violated for *this* work; drop the irrelevant ones. A tailored short list beats the full thirteen.
3. **Name the slop** — describe the specific ways this kind of work degrades (the failure modes), and contrast each with its craft counterpart.
4. **Write the Definition of Done** — concrete, checkable completion criteria that require *evidence*, ending with "what proof shows this is actually done."
5. **Emit the constitution** (see Output).

## Output contract

Deliver, in plain language:

- **Craft constitution for this project** — the 5–8 tailored rules, each one sentence, each clearly applicable to the work.
- **Slop vs craft** — a short contrast table of the specific failure modes for this work and their craft counterparts.
- **Definition of Done** — an evidence-based checklist; the work is not done until each item can be shown, not claimed.

Offer (one line) to show the internal typed model only if asked. Never lead with it.
