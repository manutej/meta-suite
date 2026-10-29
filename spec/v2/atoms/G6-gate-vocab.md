# G6 · gate.vocab

**Purpose:** Enforce P1. The internal formal machinery never reaches the reader.

**In:** user-facing text: artifact content, captions, and the final chat summary
**Out:** a verdict plus the exact locations of any hits

**Contract**
- The banned-term list lives in config and can be extended. It includes functor, morphism, pullback,
  pushout, colimit, sheaf, operad, operadic, Yoneda, adjunction, adjoint, natural transformation,
  monad and cohomology.
- Matching is whole-word and case-insensitive.
- The only exemption is a run where the user explicitly asked to see the internal view. The ledger
  records that flag.
- Domain-legitimate uses are allowed only through a per-run allowlist with a logged reason. For
  example, a brief that is itself about category theory.

**Check:** zero hits outside the allowlist.

**Guards against:** jargon leaking into client deliverables.
