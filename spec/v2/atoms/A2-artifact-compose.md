# A2 · artifact.compose

**Purpose:** Build the artifact so its structure mirrors the domain's structure, then render it.

**Internal construct:** a structure-preserving map from the merged graph and its findings to the
artifact's sections.

**In:** `ArtifactPlan`, the merged `DomainGraph`, `Finding[]`, the gate reports, and the synthesis
skeleton from L3
**Out:** the artifact, plus a `SectionMap`

**Contract**
- **Every section is anchored to findings.** A section with no anchor is cut.
- **Section tags carry the confidence:**

  | Findings | Section tag |
  |---|---|
  | obstruction | `contested` |
  | uncertain | `uncertain` |
  | single_source | `single_source` (marked inline) |
  | exception, tension | visible sections of their own |

- **Assumptions are shown,** including unsupported-type mappings and residual drift.
- **"How this was checked"** is a plain-language section (P12) built only from the actual gate
  reports. It includes `unchecked` predicates and `kernel=absent` when they apply.
- **Rendering** uses the host artifact type or the existing design skills, with the user's or
  host's design system. Nothing visual is hardcoded.

**Check:** G6 and G7. Every finding the artifact relies on is covered by `section_map`, and the
stage-5 coverage checks of G3 and G4 pass.

**Guards against:** decorative structure, and claims the findings don't support.

**Evidence:** E11.
