# Laws: what each check catches

Two checkers, one discipline: every violation is a finding with a question that resolves it.

## Ontology laws (`ontspec.py check`)

| law | level | catches | resolve by |
|---|---|---|---|
| O1 shape | RED | missing root, bad ids or kinds, depth not parent + 1, depth > 4, parent cycles | fix the tree; split the spec if it needs depth 5 |
| O2 typed references | RED | a reference to nothing, or to an atom that cannot play that role | define it, or point at the right kind |
| O3 provenance | RED | a "stated" atom with no source span; an unknown with no blocked atoms or owner | cite, downgrade to inferred, or give the unknown an owner |
| O4 ownership | AMBER | two operations write one piece of state | name the owner, or declare `shared_writers` with the reason (that reason is an architecture force) |
| O5 coverage | AMBER, RED when converging | a leaf operation or invariant nothing verifies; an obligation with no method | add an obligation |
| O6 production | RED / AMBER | an input nothing produces and that does not enter at a boundary (RED); a definition nothing uses (AMBER) | add the producer or boundary; delete the dead definition |
| O7 one value per key | RED | the same parameter with two values anywhere in the tree | decide which holds; the other is a dispute to record |
| O8 atomicity | AMBER, RED when converging | a leaf that is not yet a contract (untyped, no effect, no pre/post, writes two things) | refine one level, or complete it |
| O9 realisation | RED | an operation refined into children none of which is an operation | add the child operations |
| O10 locality | AMBER | a reference into a sibling branch's private definitions | lift the definition to the common ancestor |
| O11 effect honesty | RED | a `pure` operation that writes state | return the new state as an output, or change the effect |

## Composition laws (`ontspec.py operad`, run by the operadic-stitch engine)

The projection makes each refined atom the Goal of its own brief, nested under a frame for its
depth; a brief contains the definitions its ancestors declared plus its own children.

| engine law | in spec terms |
|---|---|
| typed plugging (exit 3) | the projection itself is malformed; a bug in the spec's kinds |
| `port.open` (RED) | an atom in a refinement uses a definition that brief cannot see: usually a sibling leak, sometimes a definition placed too low |
| `sheaf.glue` / `sheaf.seam` (RED) | a key with two values inside one refinement, or a child disagreeing with an ancestor |
| `recursion.bounded` / `wellfounded` (RED) | refinement deeper than allowed, or budgets not shrinking |
| `budget.child` (AMBER) | a refinement too large to state within its depth's budget: it is not yet atomic enough |
| `brief.no-goal` (AMBER) | a refinement with no goal of its own |

Both checkers can disagree in one useful way: O10 is AMBER (a smell) while `port.open` is RED,
because the engine sends each refinement alone. Treat the engine's verdict as the gate.

## Finding format

`LEVEL @ where · law · atoms · question`. Record every RED and AMBER in CONVERGENCE.md with the
pass that closed it. A finding is closed by changing the ontology or by turning it into an
`Open` atom with an owner, never by editing the checkers.
