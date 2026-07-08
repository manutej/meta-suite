---
name: meta-context
description: Design the context strategy for a project or working session — what knowledge to keep, where it lives, what to load when, how to compress it, and how to hand it off — so work stays grounded and compounds across sessions instead of starting cold every time. Use whenever someone wants their setup to remember and reuse context: "set up memory for this project", "what context should the agent have", "stop re-explaining everything", "organize my project knowledge", "design a CLAUDE.md / context file", "session handoff", "context engineering", "meta-context", or "/meta-context". Task-agnostic. Produces a durable substrate plan, a load/compression policy, a persistence & handoff policy, a discoverability index, and a starter substrate scaffold — in plain language. Do NOT use for a one-off question with no project or continuity.
---

# Meta-Context

Every session that starts cold pays the same tax: re-explaining the project, re-deriving the
constraints, re-finding what was decided. **Context engineering removes that tax** by designing
*what knowledge persists, where it lives, what gets loaded when, and how it's found again.* This
skill produces that design for any project, plus a starter scaffold you can drop in today.

You operate one level up. You are not doing the project's work; you are building the memory and
context system that every future unit of work will draw on.

## Operating rules (read first)

- **Never ask the user questions.** Infer the project's shape, tools, and continuity needs from
  what they gave you; default sensibly and state assumptions in one line. Proceed.
- **Plain language out.** Deliver a usable plan and scaffold. Keep the internal typed model
  (`references/internal-model.md`) invisible unless the user explicitly asks for the internal work.
- **Less, but load-bearing.** A small, well-curated context beats a large dump. Optimize for what
  changes the work, not for completeness.
- **Pipeline-friendly.** If another skill invokes you, emit just the plan + scaffold.

## The four context policies

Full detail in `references/policies.md`. Every context strategy is these four decisions:

1. **Substrate** — *where durable knowledge lives.* Default to a plain, greppable, version-controlled
   single source-of-truth document — a project brief, a run-of-show, a vendor binder, a staffing-policy
   doc, or a context file — plus a small knowledge folder, because plain text outlives
   tools and is readable by anyone and anything.
2. **Load policy** — *what gets pulled in, when.* Use progressive disclosure: a small always-loaded
   core (the project's identity, constraints, conventions) and larger references loaded only when a
   task needs them. Don't load what you won't use.
3. **Compression policy** — *how knowledge is kept small.* Summaries and atoms over raw transcripts;
   one durable note per decision or learning; prune and merge as it grows.
4. **Persistence & handoff** — *what gets written back.* At the end of a unit of work, capture
   what changed, what was decided, and what's next, so the following session (or agent, or person)
   resumes warm instead of cold.

Plus the connective tissue:

5. **Discoverability index** — a single entry point that points to everything else, so future work
   *finds* the context automatically without being told it exists.

## What to do on invocation

1. **Classify the project silently** — its continuity needs, the tools in play, and where memory
   could live. Use `references/internal-model.md`; do not show it.
2. **Decide the four policies** — substrate, load, compression, persistence/handoff — each as a
   concrete choice for *this* project.
3. **Design the discoverability index** — the one file that routes future work to the right context.
4. **Emit a starter scaffold** — a ready-to-use skeleton (e.g. a project context file with the
   core sections filled in as far as inference allows, plus an empty tagged knowledge-atom template
   and the index). Make it droppable-in.

## Output contract

Deliver, in plain language:

- **Context plan** — the four policies (substrate · load · compression · persistence/handoff), each
  a concrete decision for this project.
- **Discoverability index** — what the entry-point file is and what it points to.
- **Starter scaffold** — copy-pasteable skeleton files (project context file + atom template +
  index), pre-filled as far as inference allows, with clearly marked blanks for the rest.
- **First move** — the one file to create now.

Offer (one line) to show the internal typed model only if asked. Never lead with it.
