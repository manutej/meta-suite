# Typed contracts

`ontspec.py contracts` writes `contracts.json` and one Markdown file per leaf operation and
invariant. A contract is language-neutral; it says what any implementation must satisfy.

```json
{
  "id": "...", "kind": "operation", "effect": "deterministic",
  "signature": {"inputs": [{"id": "...", "type": "..."}], "outputs": [{"id": "...", "type": "..."}]},
  "reads": [], "writes": [],
  "requires": ["precondition"], "ensures": ["postcondition"], "errors": [],
  "source": "span", "status": "stated", "verified_by": [{"id": "...", "method": "proof"}]
}
```

## Projections (used by the investigation, not chosen here)

| target | an operation contract becomes | an invariant becomes |
|---|---|---|
| A proof-carrying language (e.g. Verus for Rust) | a function with `requires` / `ensures` | a predicate proved preserved by every writer |
| A model checker | a transition | a property checked over all reachable states |
| A typed application language (e.g. TypeScript) | a typed function plus a runtime validator at the boundary | a test or assertion |
| A durable workflow engine (e.g. Temporal) | `deterministic`: workflow code; `effectful`/`llm`: an activity with retries and idempotency; `human`: a signal or wait | a check before each state commit |
| A parallel functional language (e.g. Bend) | `pure` operations only | properties over large generated inputs |

The effect class decides the projection, which is why O11 (effect honesty) is RED.

## Regeneration rule

Code is regenerated from contracts, never the reverse. When code and contract disagree, the
contract is the reference: either the code is wrong, or the ontology changes, the contract is
regenerated, and every obligation on it runs again. Approval is of the ontology and its
contracts, by a human, recorded in the hand-kept APPROVALS.md (SPEC.md is regenerated and
would lose it).

Each contract lists `blocked_by`: the `Open` atoms on the contract itself, on an enclosing
operation, or on a definition it uses. A contract with an empty `blocked_by` is final unless the
ontology changes; the rest are the review queue.
