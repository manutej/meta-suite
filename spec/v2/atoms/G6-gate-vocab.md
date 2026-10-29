# G6 · gate.vocab

**Purpose:** Enforce P1. The internal formal machinery never reaches the reader.

**Internal construct:** —

**In:** all user-facing text: artifact content, captions and the chat summary
**Out:** a verdict, plus the location of every hit

**Contract**
- **Terms** (config: `gates.g6.terms`) are matched by stem or prefix, not whole word. "sheaf" also
  catches "sheaves" and "functor" catches "functorial". The list covers:
  - functor, morphism, pullback, pushout, colimit, cocone, sheaf, presheaf, stalk
  - operad, Yoneda, adjunction/adjoint, natural transformation, monad, cohomology
  - category theory, harmonic, Laplacian, nerve, cochain, H0/H1 (as tokens)
- **Phrases** (config: `gates.g6.phrases`) catch paraphrases, for example "fails to commute",
  "commuting square", "fibre product" and "gluing condition".
- **Allowlist.** Only I1's `vocab_allowlist` can exempt a term: a per-brief entry derived from the
  Spec's domain vocabulary (for example "harmonic" in a music brief), or a formal term in a brief that
  is itself about formal methods. Each exemption's reason is logged.
- **Internal-view exemption.** The other exemption is a user who explicitly asked for the internal
  view. That flag is recorded in S1.

**Check:** zero hits outside the allowlist.

**Guards against:** jargon leaking into client deliverables.

**Evidence:** —
