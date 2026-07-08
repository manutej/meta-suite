# Internal typed model (do NOT surface to the user unless they ask to review internal work)

Model the project as a category **P**:
- **Objects**: states of the project (its artifacts, knowledge, and capabilities at a point in time).
- **Morphisms**: units of work that move one state to the next.

Ordinary work: the cost function `cost: Mor(P) → ℝ⁺` is flat or increasing — later morphisms cost
as much or more (debt). **Compounding** is the condition that `cost` is *decreasing along the
project's trajectory*: each captured asset lowers the cost of a whole class of future morphisms.

- A **leverage investment** is a morphism whose codomain enriches the system so that an entire
  family of subsequent morphisms factors through cheaper sub-steps (it adds a reusable component
  to the category's structure).
- **Capture** is the functor that maps a completed morphism to a durable object (an atom) in a
  persistent sub-category **K** (the knowledge substrate).
- **Retrieval** is the adjoint pull: future morphisms are pre-composed with reads from **K**.
- **Discoverability** is the natural transformation guaranteeing every new object in **K** is
  reachable from the project's index object.
- **Rot** is drift between **K** and **P**; **refresh** is the corrective endofunctor on **K**.

The compounding map handed to the user is the plain-language image of this structure. The
category-theoretic framing is a private reasoning aid: never present it, name it, or ask the user
about it unless they explicitly request the internal work.
