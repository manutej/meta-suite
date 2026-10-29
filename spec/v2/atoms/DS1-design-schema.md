# DS1 · design.schema

**Purpose:** Derive a medium-specific design system (tokens, components, composition rules) from a brief, then generate every instance from it.

**Internal construct:** a typed schema. Resolved tokens are immutable downstream unless a named edit changes them.

**In:** `Spec` (medium, audience, brand inputs, content)
**Out:** `DesignSystem`: typed tokens (palette, type scale, spacing rhythm, radius, shadow, motion, breakpoints) plus components and composition rules, and instances generated from it

**Contract**
- **Precedence:** accessibility > brief > taste. Brand inputs and the user's design system override the rubric's token-count and font-family numbers; the rubric's structural dimensions are non-negotiable ([rubric](../12-artifact-rubric.md), D-069). Nothing visual is hardcoded.
- **Three roles converge:** a proposer drafts the system, an executor generates instances from it, and a blind auditor (G7 + R1) returns edit scripts until no P1/P2 edits remain.
- Every instance is generated from the system, not hand-styled.
- Palettes come from `meta palette` and every token pair is contrast-checked (G7). Instances meet rubric L2 on every applicable dimension.
- **Modes:** deep by default; `--fast` loads only the core op (D-066).

**Check:** G1, G7 (`meta rubric`) on the rendered instances.

**Guards against:** one-off designs that don't generalize, and pretty designs that fail accessibility.

**Evidence:** —
