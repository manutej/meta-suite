# Category expansion — minimal prompt → typed domain (the hidden engine)

This is how a one-line prompt becomes enough structure to generate a whole bundle. It runs
**internally**. None of this vocabulary appears in customer-facing deliverables.

## The map

Treat the domain as a category and read structure off it:

- **Objects** = the core entities the prompt implies (people, companies, initiatives, artifacts…).
- **Morphisms** = the relationships between them (offers, asks, intros, reports-to, funds, runs).
- **Identity = relationships (Yoneda).** Model an entity by its arrows in/out, not a static bio.
  This is why the data model stores connections and the dashboards render relationship graphs.
- **Composition** = chains of relationships (A introduces B who funds C) — the paths worth surfacing.
- **Invariants** = what must always hold (voice-first, personal-not-broadcast, SLA, confidentiality).
- **Colimit** = how many local pieces glue into one coherent whole (100 members → one room view;
  many cohorts → one ecosystem). This is the "one screen tells you everything" view.
- **Functor (loosely) = the build itself.** Structure in the domain becomes structure in the output:
  objects→data-model types; morphisms→interface contracts + intros; independent subgraphs→parallel
  build lanes; ordering invariants→sequential gates in the plan.

## The procedure

1. From the grounded research, list objects and give each a type.
2. List the morphisms; mark which are the load-bearing relationships (the ones the stakeholder
   cares about — e.g. "who can I connect to whom").
3. State invariants (the rules that never break).
4. Name the free variables (what changes per instance) vs the fixed structure (what's constant).
5. Decompose into independent sub-domains → these become the parallel build lanes and the
   data-model modules. Where several decompositions fit, pick one and name it a decision.

## The discipline: remove the abstraction from outputs

The expansion gives you a strong skeleton. Then **strip the vocabulary**:
- A "morphism" becomes "an introduction" or "who reports to whom."
- "Yoneda identity" becomes "we show each person by who they help and who they ask."
- "Colimit" becomes "one screen that rolls everyone up."
- Never print "object", "functor", "morphism", "colimit", "adjunction" in a client artifact.

The only place the categorical language is allowed is the optional `APPENDIX_how-we-thought.md`,
and only if the user asks for it. The math earns the structure; the structure ships; the math hides.
