# L3 · lens.glue

**Purpose:** Merge the per-region subgraphs into one synthesis without silently conflating things.
This is the **only** merge in the pipeline.

**Internal construct:** a pushout (colimit in general), taken along declared interfaces.

**In:** `RegionGraph[]`, `Interface[]` candidates from D3, `Cover.overlaps`
**Out:**
- the merged `DomainGraph`
- the accepted `Interface[]`
- conflict candidates inside each interface, passed to G3
- a synthesis skeleton for A2

**Contract**
- **No interface, no merge.** `meta interface merge` refuses any merge without a declared interface.
- **Declaring an interface.** For each candidate from D3 the model either:
  - accepts it, or
  - edits it, with a reason.

  Identifying two objects requires a J1 `judge` verdict.
- **Outside the interface** everything stays separate and keeps its sources.
- **Conflicts inside the interface** go to G3. They are never averaged.
- **Always runs.** For `profile` and `compare` it merges along interfaces only. For `synthesize` the
  skeleton follows the merged structure, not a stock report template.

**Check:** the kernel refuses undeclared merges. G3 gets the conflicts.

**Guards against:** drift through conflation, and "both sides" mush.

**Evidence:** E3 (colimit synthesizer), E11, E12.
