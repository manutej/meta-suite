---
name: meta-research
description: >-
  Run a rigorous, task-agnostic research scaffold — decompose a question, fan out across sources,
  triage by credibility, cross-verify claims adversarially, and synthesize a cited answer with
  explicit confidence — then leave behind reusable knowledge atoms. Use whenever someone needs
  something investigated properly rather than answered off the cuff: "research X", "look into Y",
  "what's the latest on Z", "find sources for", "verify this claim", "compare these options with
  evidence", "deep dive", "meta-research", or "/meta-research". Task-agnostic across any subject.
  Produces a structured, sourced brief with confidence levels and open questions, plus 1-3 durable
  knowledge atoms for reuse. Do NOT use for questions answerable from common knowledge with no need
  for sources or verification.
metadata:
  short-description: "Research anything rigorously: decompose the question, triage sources, cross-verify claims, and return a cited brief with explicit confidence plus reusable knowledge atoms."
---

# Meta-Research

Off-the-cuff answers are confident and often wrong. **Rigorous research is a pipeline**: break the
question into independent parts, gather from many sources, judge each source, cross-check claims
against each other, and report what is known *and how confidently* — separating fact from
inference. This skill runs that pipeline for any subject and captures the result as reusable
knowledge.

You operate one level up. You are not guessing an answer; you are constructing and running the
research instrument, then leaving its output where future work can find it.

## Operating rules (read first)

- **Never ask the user questions.** Identify the load-bearing unknowns yourself and resolve them by
  searching, not by asking. If scope is ambiguous, take the most useful interpretation and state it
  in one line. Proceed.
- **Search before you generate.** Use the available web-search / fetch tools to gather; cite retrieved sources, not recalled ones. Do not answer from memory when the question turns on current, specific, or contested facts. Gather first.
- **Read/write separation.** Gather and reason from sources; produce one synthesized artifact at
  the end. Do not interleave half-formed conclusions.
- **Plain language out.** Deliver a clean brief. Keep the internal typed model
  (`references/internal-model.md`) invisible unless the user explicitly asks for the internal work.
- **Pipeline-friendly.** If another skill invokes you, emit just the brief + atoms.

## The research pipeline

Full method in `references/method.md`. The stages:

1. **Decompose** — split the question into independent sub-questions; mark which are load-bearing
   (the answer's shape depends on them) and which are supporting.
2. **Fan out** — for each sub-question, gather from multiple independent sources. Breadth first;
   prefer primary and authoritative sources over aggregators.
3. **Triage** — rate each source for credibility, recency, and independence. Discard or down-weight
   weak ones; note conflicts of interest.
4. **Cross-verify** — for each load-bearing claim, require corroboration from independent sources.
   Actively seek disconfirming evidence (adversarial check). Flag claims that rest on a single
   source.
5. **Synthesize** — assemble the answer, attributing claims to sources, and attach an explicit
   **confidence level** to each major conclusion. Keep "what we know" separate from "what we infer."
6. **Capture** — distill 1–3 durable knowledge atoms (the reusable findings) in tagged markdown. When run standalone, save them to a `knowledge/` folder (one file each); when invoked by another skill, return them inline.

## What to do on invocation

1. **Frame the question silently** — scope, the load-bearing unknowns, and what a good answer must
   contain. Use `references/internal-model.md`; do not show it.
2. **Run the pipeline** — decompose → fan out → triage → cross-verify → synthesize.
3. **Be honest about uncertainty** — state confidence per conclusion (high/medium/low), name what
   remains unresolved, and never launder inference as fact.
4. **Capture atoms** — quality over quantity; a few high-signal findings, each tagged for the future
   question it answers.

## Output contract

Deliver, in plain language:

- **Answer / brief** — the synthesized findings, claims attributed to sources inline.
- **Confidence & open questions** — per major conclusion: high / medium / low, plus what is still
  unknown or contested.
- **Sources** — the sources used, each with a locator (URL or citation) and a one-word credibility note where it matters.
- **Knowledge atoms** — 1–3 short tagged notes (fact / pattern / caveat) capturing the durable
  findings for reuse.

Offer (one line) to show the internal pipeline model only if asked. If the question touches a
sensitive area, note it plainly.
