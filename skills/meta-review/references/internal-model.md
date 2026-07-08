# Internal typed model (do NOT surface to the user unless they ask to review internal work)

Let A be the artifact under review. A **lens** ℓ is a functor from the artifact to a set of typed
findings: `ℓ : A → Findings`. The review is the product of selected lenses
`R = ⟨ℓ₁, …, ℓₙ⟩ : A → ∏ Findings_i`, run independently so no lens's image collapses another's.

Findings carry a priority in the totally-ordered set {P1 > P2 > P3} and a band in the precedence
order Correctness/Safety ⊐ Scope ⊐ Coherence ⊐ Clarity. Ordering the merged findings by (band, then
priority) is a sort into a unique normal form — the fix list — so conflicting remedies resolve
deterministically.

The **lesson atom** is the capture functor applied to the highest-recurrence finding: it maps a
one-time defect into a durable check that pre-composes with future reviews of the same artifact
class, closing the loop so the defect is caught automatically next time.

This framing is a private reasoning aid. Never present it, name the lenses-as-functors structure,
or ask the user about it unless they explicitly request the internal work.
