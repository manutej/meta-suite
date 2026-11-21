# Ontology schema

One YAML (or JSON) file: a `spec` header and a flat `atoms` list. The tree comes from `parent`.

```yaml
spec:
  id: kebab-id
  title: Human title
  version: 0.1
  max_depth: 4
  verbal: >-
    The one-paragraph depth-0 statement.
atoms:
  - id: kebab-id          # unique, lower-case, <= 63 chars
    kind: Operation        # see table
    depth: 2               # parent depth + 1; the System root is depth 0
    parent: other-id       # omitted only on the System root
    status: stated         # stated (with source) | inferred | open
    source: "S1 §3.2 Eq.(4)"  # required when stated: a span, not just a document
    statement: One plain sentence.
    notes: Optional free text for anything the schema has no field for.
    # kind-specific fields below
```

## Fields by kind

| kind | means | required at a leaf | optional |
|---|---|---|---|
| System | the whole thing (exactly one, the root) | — | — |
| Component | a named part that groups others | must have children | — |
| Entity | a thing with identity or a data shape | `type` | `durable`, `external_input` (enters from outside without a Boundary) |
| Value | a fixed parameter | `key`, `value` (or status open), `type` with unit/range | — |
| State | something that changes over time | `type` (fields, or enum with its order) | `durable`, `shared_writers` (a reason) |
| Relation | how things connect | `type` (signature) | `uses` |
| Event | something that happens and is recorded | `type` | — |
| Boundary | what crosses into or out of the system | `type` (message type) | — |
| Operation | what changes things or computes | `inputs`, `outputs`, `effect`, `pre`, `post`; `writes` at most one | `reads`, `errors` |
| Invariant | what must always hold | `uses`, `formal` (a checkable predicate) | — |
| Obligation | how a claim is verified | `target`, `method` | `evidence` (test names, commands, prior results) |
| Open | an unknown | `blocks`, `resolve_by` | — |

`effect`: `pure` (no writes, no outside world), `deterministic` (may write owned state; replays
identically), `effectful` (outside world, needs idempotency), `llm` (model call), `human`.
`method`: `proof`, `model-check`, `property-test`, `differential` (against a reference or the
source's own numbers), `review`.

## References between atoms

| field | may point at |
|---|---|
| inputs, outputs | Entity, Value, State, Relation, Event, Boundary |
| reads | Entity, State, Value, Boundary |
| writes | Entity, State |
| uses | anything except Obligation and Open |
| target | anything except Obligation and Open |
| blocks | anything |

Definitions follow lexical scope: an atom may use definitions in its own subtree, or declared
directly under any of its ancestors. Reaching into a sibling branch is a finding (lift the
definition to a common ancestor). `target` and `blocks` may point anywhere.

## Writing good atoms

- Write constants in text as `{key}` (`ρ̄ ≥ {theta}`), never as bare numbers: O13 checks the key exists, and
  O7 then guards its single value. A bare 0.85 in a post-condition can drift silently.
- Unknown fields are reported (O12). Put anything else in `notes`.

- One sentence per statement. If it needs "and", it is two atoms.
- Types carry units and ranges: `f64 in [0,1]`, `u32 tokens/s`, `enum{a < b < c}`.
- Put a parameter where it is first needed, at the highest level that needs it. Refinements inherit it.
- A source is a span (section, page, equation, table, line), not a document.
