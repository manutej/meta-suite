# R2 · review.expert_delta

**Purpose:** Measure how much an expert would learn. Count the sourced insights a plain,
unstructured answer would not contain.

**Internal construct:** —

**In:** the artifact's findings and a plain baseline answer to the same `Brief`. The baseline uses
the same model, with no lens, no gather and the quick budget.
**Out:** `DeltaReport`

**Contract**
- **Full in the deep tier and in eval. R2-lite in the standard tier** (D-060): `review.r2.lite_max_findings` sampled findings, one baseline. Config: `tiers.<tier>.r2`. This keeps standard runs
  affordable.
- **Fresh baseline every run.** The baseline is never cached across briefs.
- **Eligible findings.** Only these kinds count:
  - `verified`
  - `exception`
  - `tension`
  - `obstruction`
  - a `profile_trait` that discriminates

  At most `review.r2.max_findings` are judged, highest-confidence first. Presentation differences
  never count.
- **Judging.** Whether a finding is absent from the baseline is a J1 `judge` task.
- **Low count.** If `count < review.r2.min_delta`, R2 emits an edit script asking for a deeper
  re-gather in the regions that produced the fewest delta findings.

**Check:** `count ≥ review.r2.min_delta`. Every delta finding passes G2.

**Guards against:** output that looks expert but says what any chatbot would. Agents already score
about 90 on presentation and about 40 on recall.

**Evidence:** E11, E12, E13.
