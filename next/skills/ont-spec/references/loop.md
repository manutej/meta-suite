# Depth, atomicity and convergence

## What each depth holds

| depth | holds | a reviewer asks |
|---|---|---|
| 0 | the system, as one paragraph | is this the thing we mean? |
| 1 | its components, the definitions they share, its invariants | are these the right parts, with the right shared language? |
| 2 | each component's mechanisms | does each part work on its own terms? |
| 3 | algorithms and data shapes | could two engineers implement this identically? |
| 4 | atomic contracts | could this be proved or tested exactly as written? |

Not every branch needs all four. A branch stops refining when its leaves are atomic.

## Atomic leaf

A leaf is atomic when it can become one contract with nothing left to decide:
- an Operation has typed inputs and outputs, an honest effect, pre and post conditions, and
  writes at most one thing;
- a definition has a type with units or range (and a Value has its value, or is open);
- an Invariant names what it constrains and states a checkable predicate;
- an Obligation names its targets and its method.

`ontspec.py check --converging` turns every non-atomic or unverified leaf into a RED.

## The pass

1. Edit `ontology.yaml` (add a depth, or resolve findings).
2. Run `check`, `check --converging` (once depth 4 is reached), `operad`, and `delta` against the previous
   pass. Save the outputs in `certificates/pass-N/`.
3. Append to CONVERGENCE.md:
   ```
   ## Pass N: <what this pass did>
   - edit: <atom ids added, removed, restructured, and why>
   - closed: <findings closed>
   - opened: <findings opened, and the decision on each>
   - verdicts: check <V> · converging <V> · operad <V> · delta <added/removed/changed counts>
   ```
4. Repeat.

The log is cumulative: each pass appends; earlier passes are never rewritten. That is what lets a
reviewer see why every atom exists.

## Stopping

Stop when, in the same pass:
- `check --converging` is GREEN (or AMBER only with each finding explained);
- `operad` is GREEN (NOTEs are fine);
- `delta` against the previous pass shows no atoms added, removed or structurally changed.

The last condition needs one quiet pass after the last edit; that pass is the convergence proof.
If a pass keeps reopening the same branch, the branch is a separate system: give it its own spec
and reference it at a Boundary. If an unknown blocks convergence, it stays an `Open` atom and
moves to the investigation brief; convergence means the spec is stable, not that the world is known.
