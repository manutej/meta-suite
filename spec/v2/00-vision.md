# 00 · Vision

## North star

Hand meta-suite a rough, half-formed question in any domain. Get back a research view or artifact
whose nuance makes domain experts lean in:
- the tensions practitioners argue about
- the exceptions only insiders know
- the contested facts, named as contested
- the facts that rest on only one source, marked as such

Every claim is sourced. Structural checks run as code, and semantic judgments are sampled and recorded (J1), so every verdict is auditable. The run needs less human input than a
traditional harness, and its cost is capped and measured. Second-order effects are a v2.1 goal
(composition operator, see D-031).

## Why the category-theory lens is the core

LLMs are inherently unstructured, non-deterministic, and prone to drift and hallucination. The lens
imposes a **hypothesis-and-experiment structure** on that mechanism:

- **The domain model is a hypothesis.** It is the set of objects and the typed relations between them.
- **Each construct is an experiment.** It is a check with a pass/fail result that code can compute.
- **A failed check is a finding, not an error.** A commuting square that fails is a real tension. A
  gluing obstruction is a contested fact. A probe set that cannot tell two entities apart is a
  coarse model.

The lens drives research structure and relational domain mapping. This is what made v1 work (for
example, defining a customer entirely by its relations in ICP work), and v2 strengthens it.

**The lens never reaches the reader.** Artifacts speak the domain's language. The machinery is
recorded in the run ledger and shown only on request.

## What changed since v1 (why v2 now)

Evidence is in [evidence](07-evidence.md).
- **Frontier models no longer need procedural hand-holding.** Persona and step-by-step boilerplate
  does not help reasoning models.
- **Self-critique alone does not reliably improve output.** Blind, external or deterministic
  verification does.
- **The best prompt pipelines are now measured search loops.** They have metrics, held-out sets,
  and small playbook edits instead of full rewrites.
- **Compact skills beat comprehensive ones.** Self-generated skills hurt unless validated.
- **Categorical structure is now measurable.**
  - Operadic consistency predicts accuracy without labels.
  - Pullback verifiers and colimit synthesizers have published ablations.
- **Deep-research agents already score about 90 on presentation but about 40–50 on recall and
  analysis.** The lens targets exactly that gap.

## Product horizon

| Horizon | Shape |
|---|---|
| **v2.0** | **A new suite**: `meta-lens` as the front door plus the nine skills rebuilt on one shared Go kernel and shared lanes. Every skill has a deep mode (default) and `--fast`. Each skill is measured against its v1 version. v2.0 ships the `meta` CLI as an **agent-callable harness core**: any agent can call its gates, ledger and sheaf check. A stable public API/SDK is v3. |
| **v2.x** | Cross-run learning and the domain memory bank (v2.1, D-063). **Watch and alerts** (milestone M7, D-074): re-run on a cadence, detect changed facts and gate results, notify the user. More question types and operators come in from the parking lot. |
| **v3** | A stable public API/SDK and middleware layer any agent can call: intake→spec, lens operators, gates, alerts. It composes with the Stitch decomposition/composition infrastructure as the reasoning layer above it. |

## Success, defined

v2.0 succeeds when the pre-registered evaluation in [eval](05-eval.md) shows seven things:
1. v2 beats v1 on expert-relevant recall and analysis.
2. Removing the lens measurably hurts (so the lens is proven, not assumed).
3. No model in the matrix regresses.
4. Cost stays within the pre-registered ceiling over v1's.
5. Artifacts meet the [rubric](12-artifact-rubric.md)'s minimum bar and are at least as good as v1's best artifacts.
6. `--fast` is no slower than v1.
7. Each specialist skill is at least as good as its v1 version.
