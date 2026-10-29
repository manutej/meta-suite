# L3 · lens.glue

**Purpose:** Merge the per-region subgraphs into one synthesis without silently conflating things.
This is the **only** merge in the pipeline.

**Internal construct:** a pushout (colimit in general) along declared interfaces. Because
identifications compose transitively across three or more regions, they must be checked for
coherence.

**In:** `RegionGraph[]`, `Interface[]` candidates from D3, `Cover` overlaps and triple overlaps
**Out:**
- the merged `DomainGraph`
- the accepted `Interface[]`
- conflict candidates inside each interface, passed to G3
- a synthesis skeleton for A2

**Contract**
- **No interface, no merge.** `meta interface merge` refuses any merge without a declared interface.
- **Declaring an interface.**
  - The model accepts or edits each candidate interface, with a reason.
  - "Same object" and "distinct objects" are J1 `judge` verdicts.
- **Coherent identifications.** `meta interface merge` computes the transitive closure of all
  "same" identifications with union-find. It **refuses** the merge if the closure would join two
  objects that J1 judged distinct, and reports the chain of identifications that caused the
  collision.
- **Outside the interface** everything stays separate and keeps its sources.
- **Conflicts inside the interface** go to G3 and are never averaged.
- **Always runs.** For `profile` and `compare` it merges along interfaces only. For `synthesize` the
  skeleton follows the merged structure.

**Check:** the kernel refuses merges that are undeclared or incoherent. G3 receives the conflicts.

**Guards against:**
- drift through conflation
- A=B and B=C silently implying A=C when the evidence says A≠C
- "both sides" mush

**Evidence:** E3 (colimit synthesizer), E11, E12.
