# A1 · artifact.select

**Purpose:** Choose the artifact form. The question decides what kind of artifact it is; the shape
of the domain graph decides how it is laid out.

**In:** `Brief` (question type, medium if stated, audience), `DomainGraph` shape metrics, the
artifact types available on the host
**Out:** `ArtifactPlan`: form, layout and the host artifact type to use

**Contract**

| Question type | Default form |
|---|---|
| `profile` | Profile page: traits, what sets it apart, evidence |
| `compare` | Decision workbench: side-by-side, tensions first |
| `synthesize` | Atlas or report: structure follows the merged graph |

| Graph shape | Layout |
|---|---|
| hub-heavy | navigate from the hubs outward |
| long chain | sequential narrative |
| two-sided | split view |
| time-dense | timeline spine |

- A medium stated by the user always wins, for example "a deck."
- Use a host artifact type when one fits (Docs, Slides, Dashboard, Design). A hand-built page is
  the fallback.
- The rules live in config, so new rows are added without code changes.

**Check:** the plan is schema-valid (G1). The chosen rule is logged.

**Guards against:** one template for every question.
