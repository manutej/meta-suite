# Context policies — full detail (internal reference)

Context engineering is the discipline of deciding what an agent or collaborator should know at each
moment, and ensuring that knowledge persists and is found. Four policies plus an index cover it.

## 1. Substrate — where durable knowledge lives
Default: **plain, greppable, version-controlled markdown.** A small project context file (identity,
goal, constraints, conventions, glossary) plus a knowledge folder of tagged atoms. Plain text is
tool-agnostic, diffable, human- and machine-readable, and outlives any specific app or model. Avoid
proprietary or opaque stores unless a hard requirement forces them.

## 2. Load policy — what gets pulled in, when
Progressive disclosure. Keep an **always-loaded core** as small as possible — only what every task
needs (who/what/why, hard constraints, conventions). Everything else is a **reference** loaded on
demand when a specific task calls for it. The cost of context is real: irrelevant loaded context
crowds out the relevant. Name, for this project, what belongs in the core vs the references.

## 3. Compression policy — how knowledge stays small
Prefer summaries and atoms over raw history. One durable note per decision, learning, or pattern —
self-contained and tagged. As the store grows, **merge** duplicates and **prune** the stale. The
goal is a store whose size tracks the project's *distinct* knowledge, not its activity volume.

## 4. Persistence & handoff — what gets written back
At the close of a unit of work, write back: what changed, what was decided (and why), what's
unresolved, and the obvious next step. This is the difference between resuming warm and starting
cold. Make it a habit, not an afterthought — ideally the last step of every loop.

## 5. Discoverability index — the entry point
A single file (e.g. the project context file itself, or a short index) that points to every other
piece of context. Every new atom updates the index. This guarantees future work *finds* the
context automatically. Knowledge that exists but isn't discoverable does not compound.

## Default starter scaffold
- `CONTEXT.md` (or project context file): Identity · Goal · Hard constraints · Conventions ·
  Glossary · Pointers to knowledge & decisions.
- `knowledge/` : tagged atom files; one template `knowledge/_template.md` with frontmatter
  (type, tags, confidence, date, source).
- The context file doubles as the discoverability index until the project is large enough to need a
  dedicated one.
