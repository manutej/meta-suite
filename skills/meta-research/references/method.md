# Research method — full pipeline (internal reference)

A task-agnostic generalization of disciplined investigation. The pipeline is fixed; only its
content changes with the subject.

## 1. Decompose
Break the question into independent sub-questions. Classify each:
- **Load-bearing**: the overall answer changes shape depending on it. These get the most rigor.
- **Supporting**: fills detail but doesn't change the conclusion's structure.
Identify the load-bearing unknowns *yourself* — never ask the user. The default scope is the most
decision-useful reading of the question; state it in one line and proceed.

## 2. Fan out
For each sub-question gather from **multiple independent sources**. Prefer primary sources
(original documents, datasets, first-hand reports) and recognized authorities over aggregators and
secondary summaries. Breadth before depth: get several angles before committing.

## 3. Triage
Rate each source on credibility (who, expertise, track record), recency (is it current?), and
independence (does it merely echo another source?). Note conflicts of interest. Down-weight or
discard weak sources; do not let volume substitute for quality.

## 4. Cross-verify (adversarial)
Every load-bearing claim must be corroborated by independent sources. Actively search for
*disconfirming* evidence — try to break the claim, not just confirm it. A claim supported by only
one source, or by sources that all trace to the same origin, is flagged as single-sourced. Note
genuine disagreements rather than papering over them.

## 5. Synthesize
Assemble the answer with claims attributed to their sources. Attach a confidence level to each
major conclusion:
- **High** — multiple strong, independent, current sources agree.
- **Medium** — supported but with thin, dated, or partly conflicting evidence.
- **Low** — plausible but weakly sourced or contested; mark clearly.
Keep **knowledge** (sourced fact) separate from **inference** (your reasoning on top). Never present
inference as fact.

## 6. Capture
Distill 1–3 durable **knowledge atoms**: small, self-contained, tagged findings phrased around the
future question they answer. Filter hard — a few high-signal atoms beat many low ones. Default
format: plain markdown with lightweight frontmatter (type, tags, confidence, date, source).
