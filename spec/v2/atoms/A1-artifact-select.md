# A1 · artifact.select

**Purpose:** Choose the artifact's form and layout. The question type(s) set the form, and the shape
of the domain graph sets the layout.

**Internal construct:** —

**In:** `Brief` (question types, medium, audience), the **reasoning sequence** (the ordered chain of sub-questions from D1's tree: the order in which the answer must be built for the reader, VR-07), the merged graph's `shape_metrics`, and the
host's available artifact types
**Out:** `ArtifactPlan`

**Contract**
- **Default forms** (config: `artifact.forms`). These are defaults, not a closed list.

  | Question type | Form |
  |---|---|
  | `profile` | profile page |
  | `compare` | decision workbench, tensions first |
  | `synthesize` | atlas or report |
  | data that is monitored or explored (any type) | dashboard |
  | a mechanism or process to explain (any type) | animated explainer, through an installed explainer skill when present |

  With two question types, the form is a composite: the primary type's form, plus a section built
  for the secondary type.
- **Layout comes from the shape metrics**, using thresholds in config (`artifact.layout`):

  | Shape | Layout |
  |---|---|
  | high hub concentration | hub-first navigation |
  | long chains | sequential narrative |
  | two-sidedness | split view |
  | temporal density | timeline spine |

  Metrics can combine, for example a split view with a timeline spine.
- **The section order follows the reasoning sequence**, not the order of the graph, so the reader meets the ideas in the order they build on each other.
- **A stated medium always wins**, for example "a deck."
- **Use a host artifact type when one fits.** A hand-built page is the fallback.
- **Rules are data.** The chosen rule's id and a one-line rationale are logged. New rules are added
  in config, never in code.

**Check:** G1 on the plan.

**Guards against:** using one template for every question.

**Evidence:** E10.
