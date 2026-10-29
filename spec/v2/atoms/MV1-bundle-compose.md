# MV1 · bundle.compose

**Purpose:** Turn a one-liner into a right-sized bundle by chaining suite recipes: research view, vision, plan, data model, artifacts, optional explainer.

**Internal construct:** an operadic composition of skill outputs. Each part's output types must match the next part's input types.

**In:** a rough one-liner
**Out:** `Bundle`: an ordered set of artifacts with a manifest (what each is, how it was checked)

**Contract**
- Size the bundle to the domain. The default is meta-lens view + vision + plan + one primary artifact.
- Every part comes from a suite recipe (meta-lens, meta-planning, meta-design). The kernel checks that the output types match the next part's inputs.
- The manifest includes "how this was checked" for every part.
- **Modes:** deep by default; `--fast` loads only the core op (D-066).

**Check:** G1 on the manifest; each part's own gates; rubric M8 on the primary artifact.

**Guards against:** a bundle of disconnected pieces, and a padded bundle.

**Evidence:** —
