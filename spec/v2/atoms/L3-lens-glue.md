# L3 · lens.glue

**Purpose:** Combine local models, whether regions, literatures or frameworks, into one synthesis
without silently treating different things as the same thing.

**Internal construct:** a pushout, or a colimit in general. Merge only along a declared interface.

**In:** per-region subgraphs, `Cover` overlaps
**Out:** merged `DomainGraph`, `interface[]`, obstruction candidates passed to G3, and a synthesis
skeleton for A2

**Contract**
- The kernel refuses any merge that has no declared interface: the objects and relations treated
  as the same across regions.
- Everything outside the interface stays separate, with its source attached.
- Conflicts inside the interface become obstruction candidates. They are never averaged away.
- The synthesis skeleton follows the shape of the merged structure, not a generic report template.

**Check:** G3 (gluing, per axiom).

**Guards against:** drift by conflation, the most common failure of LLM synthesis, and
"both-sides" mush.

**Evidence:** E3 (colimit synthesizer), E7 (minimal conflicting sets), E11 and E12 (analysis depth
is the gap).
