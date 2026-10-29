# G7 · gate.render

**Purpose:** Every artifact must pass the accessibility and layout checks before it is shown (P11).
The checks are carried over from v1 meta-design and made computable.

**In:** the artifact (HTML, SVG, or a host artifact type's data)
**Out:** a predicate table with the computed values

**Predicates**
- **Contrast:** at least 4.5 for every text pair, or at least 3.0 for large text. Computed from the
  resolved colours, not estimated.
- **Headings:** nested strictly, and reading order matches source order.
- **Interactive targets:** at least 24px, with 44px preferred.
- **Motion:** a `prefers-reduced-motion` fallback exists, and nothing flashes more than 3 times per
  second.
- **Layout:** no horizontal overflow at 360px and at the declared breakpoints.

**Contract**
- Static predicates (contrast, headings, targets declared in CSS) are computed in Go by parsing the
  HTML and CSS.
- Overflow checks need a real render. Use a headless browser if one is present in the environment.
  If not, mark the predicate `unchecked` rather than `pass`.
- For host artifact types (docs, slides and similar), use the host's own checks where they exist,
  and record which checks were not available.

**Check:** every predicate is `pass` or explicitly `unchecked`. None is `fail`.

**Open:** whether a headless browser is available in each surface's sandbox. This is resolved in
the M0 spike.
