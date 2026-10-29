# G7 · gate.render

**Purpose:** Every artifact passes computed accessibility and layout checks, and the rubric's machine checks, before it is shown (P11).

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
  - declared target sizes (config: `gates.g7.target_min_px`, default 24; `gates.g7.target_pref_px`, default 44)
  - a `prefers-reduced-motion` fallback
  - no flashing above `gates.g7.flash_max_hz` (default 3 per second) in declared animations
- **Render predicates** (horizontal overflow at `gates.g7.viewports`, default 320/360/390/768px, and the
  declared breakpoints) run only when a headless browser is present. Otherwise they are `unchecked`,
  never `pass`.
- **Rubric machine checks (D-069).** `meta rubric` runs the config table `rubric.checks`, one row per
  [rubric](../12-artifact-rubric.md) *Check* line: viewports, both themes, overflow, targets, contrast
  on resolved colours, first-viewport verdict presence, token count, reduced-motion end frame. A check
  that needs a browser that isn't there is `unchecked`, and the artifact may not claim L2 on that
  dimension.
- **Visibility.** The number of `unchecked` predicates appears in "how this was checked" and in the
  ledger.

**Check:** no predicate is `fail`. Contrast is never `unchecked` on a hand-built page.

**Guards against:** inaccessible deliverables, and "looks fine to me."

**Evidence:** —
