# A2 · artifact.compose

**Purpose:** Build the artifact so its structure mirrors the structure of the domain, then render it.

**Internal construct:** a structure-preserving map from the domain graph and findings to the
artifact's sections.

**In:** `ArtifactPlan`, `DomainGraph`, `Finding[]`, gate reports
**Out:** the artifact, plus a `section_map` (section → finding ids)

**Contract**
- Every section is anchored to findings or graph objects. A section with no anchor is cut.
- **Contested** findings (from G3) and **tensions** (from L4) get their own visible sections.
- **Single-source** findings are marked as such.
- A plain-language **"how this was checked"** section is generated from the actual gate reports
  (P12). Examples: "Every fact shown as verified was found in at least two independent sources";
  "These 3 points remain disputed."
- Rendering uses the host artifact type or the existing design skills. The design system is the
  host's or the user's; nothing is hardcoded.
- The artifact contains no internal formal vocabulary (G6).

**Check:** G6, G7, and section-anchor coverage of 100%.

**Guards against:** decorative structure, and claims in the artifact that no finding supports.
