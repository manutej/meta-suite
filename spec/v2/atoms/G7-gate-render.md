# G7 · gate.render

**Purpose:** Every artifact passes computed accessibility and layout checks before it is shown (P11).

**Internal construct:** —

**In:** the artifact (HTML, SVG, or a host artifact type's data)
**Out:** a predicate table with computed values, and a count of `unchecked` predicates

**Contract**
- **Contrast** is computed over the **declared colour-token pairs** (foreground/background tokens and
  the components that use them). The default thresholds are 4.5, or 3.0 for large text
  (config: `gates.g7.*`).
  - On a hand-built page every text colour must come from a token, so contrast can **never** be
    `unchecked`.
  - On host artifact types, the host's own checks are used where they exist, and the gaps are
    recorded.
- **Static predicates** computed in Go:
  - heading nesting
  - reading order matching source order
  - declared target sizes (default minimum 24px, 44px preferred)
  - a `prefers-reduced-motion` fallback
  - no flashing above 3 per second in declared animations
- **Render predicates** (horizontal overflow at 360px and the declared breakpoints) run only when a
  headless browser is present. Otherwise they are `unchecked`, never `pass`.
- **Visibility.** The number of `unchecked` predicates appears in "how this was checked" and in the
  ledger.

**Check:** no predicate is `fail`. Contrast is never `unchecked` on a hand-built page.

**Guards against:** inaccessible deliverables, and "looks fine to me."

**Evidence:** —
