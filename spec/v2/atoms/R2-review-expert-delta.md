# R2 · review.expert_delta

**Purpose:** Measure the "SME lean-in." These are the sourced insights a plain, unstructured answer
would not contain.

**In:** the artifact's findings, and a plain baseline answer to the same `Brief`. The baseline uses
the same model with no lens and no gather, at a quick budget.
**Out:** `delta[]` (findings absent from the baseline, each sourced) and their count

**Contract**
- The baseline is generated fresh in every run. It is never cached across briefs.
- The model judges whether each finding is "absent from the baseline," sampled k times; the kernel
  aggregates.
- Only `verified`, `tension`, `obstruction` and discriminating `profile_trait` findings count.
  Presentation differences do not.
- If `count` is below the threshold in config, R2 emits an edit script that asks for a deeper
  re-gather in the regions with the fewest delta findings.

**Check:** `count ≥ N_min`, and every delta finding passes G2.

**Guards against:** artifacts that look expert but say what any chatbot would say. Agents already
score about 90 on presentation but about 40 on recall.

**Evidence:** E11, E12, E13.
