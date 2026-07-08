# Review lens catalog (internal reference)

A lens is an independent evaluator that looks for one *class* of failure. Running several in
parallel and merging their findings beats a single holistic pass, because each lens is blind to the
others' concerns. Select the subset that fits the artifact; always run Adversarial and Clarity.

- **Correctness** — Does it produce the right result on representative and boundary inputs? Trace
  the main path and at least one edge. For claims/analysis: do the numbers and logic hold?
- **Robustness** — Behavior under malformed input, partial failure, high scale, concurrency, or an
  adversary. What degrades, what breaks, what is unrecoverable?
- **Clarity** — Can the next person understand and safely change this? Naming, structure, hidden
  cleverness, missing "why." Clear beats clever.
- **Safety / security** — Injection, leakage, privilege, destructive actions, irreversibility,
  harm to users or third parties. For non-code: reputational, legal, ethical exposure.
- **Scope & feasibility** — Right-sized for the intent? Over- or under-built? Actually achievable
  with the stated resources and constraints?
- **Coherence** — Do all parts agree with each other and with the stated goal? Contradictions,
  drift between spec and implementation, mismatched assumptions.
- **Evidence** — Is each load-bearing claim backed by a runnable proof, a citation, or a rendered
  result — or just asserted? Flag assertion-without-evidence.
- **Adversarial / devil's advocate** — Construct the strongest argument that the work is wrong,
  unnecessary, or will fail in practice. This lens is mandatory.

Domain add-ons: data-integrity (analysis/ETL), accessibility & contrast (UI), factual-accuracy &
source-credibility (research/writing), performance (systems), migration-safety (schema/state).

## Prioritization precedence
Correctness & Safety ⊐ Scope & Feasibility ⊐ Coherence ⊐ Clarity & Polish. When two findings
conflict (e.g. "add a guard" vs "keep it minimal"), the higher band wins. This yields a single,
unambiguous ordered fix list.
