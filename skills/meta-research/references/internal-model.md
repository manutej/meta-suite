# Internal typed model (do NOT surface to the user unless they ask to review internal work)

Model the inquiry as a question object Q decomposed into a diagram of sub-questions {qᵢ}, with
load-bearing qᵢ marked as the limit-determining objects (the answer is the limit of the diagram;
its shape is fixed by exactly those objects).

- **Fan-out** assigns each qᵢ a set of source-morphisms `sᵢⱼ : Source → Claim`.
- **Triage** is a weighting on sources; **independence** is the condition that two source-morphisms
  do not factor through a common origin.
- **Cross-verification** accepts a claim only if it is the (independent) coproduct of ≥2
  independent source-morphisms; single-sourced claims are tagged as such.
- **Synthesis** is the limit: the cone over verified claims, annotated with a confidence functor
  into {high > medium > low}.
- **Capture** maps the highest-value verified claims into durable atoms in a persistent knowledge
  sub-category, retrievable by future inquiries.

This framing is a private reasoning aid. Never present it, name the structure, or ask the user about
it unless they explicitly request the internal work. The user sees a plain, sourced, confidence-
annotated brief.
