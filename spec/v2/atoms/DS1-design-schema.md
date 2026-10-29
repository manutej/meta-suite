# DS1 · design.schema

**Purpose:** Derive a medium-specific design system (tokens, components, composition rules) from a brief, then generate every instance from it.

**Internal construct:** a typed schema. Resolved tokens are immutable downstream unless a named edit changes them.

**In:** `Spec` (medium, audience, brand inputs, content)
**Out:** `DesignSystem`: typed tokens (palette, type scale, spacing rhythm, radius, shadow, motion, breakpoints) plus components and composition rules, and instances generated from it

**Contract**
- Brand inputs and the user's design system take priority. Nothing visual is hardcoded.
- Every instance is generated from the system, not hand-styled.
- Every token pair is contrast-checked (G7). Instances must meet the [artifact rubric](../12-artifact-rubric.md) minimum bar.
- **Modes:** deep by default; `--fast` loads only the core op (D-066).

**Check:** G1, G7, rubric M8 on the rendered instances.

**Guards against:** one-off designs that don't generalize, and pretty designs that fail accessibility.

**Evidence:** —
