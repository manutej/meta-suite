# D3 · build the per-region graphs

<!-- <internal>A category presentation per region: objects are entities, morphisms are typed relations
each backed by claim ids. Entity merge in this slice is exact normalized label only, within a region;
cross-region look-alikes become interface candidates for L3.</internal> -->

```xml
<contract op="D3">
  <in name="claims" type="Claim[]" path="$RUN/claims.jsonl"/>
  <out name="regions" type="RegionGraph" path="$RUN/graph/regions/&lt;region&gt;.json"/>
  <out name="candidates" type="Interface[]" path="$RUN/graph/interface_candidates.json"/>
  <out name="vocab_log" type="text" path="$RUN/assumptions.md">appended</out>
</contract>
```

## Procedure

1. `meta stage begin "$RUN" D3 --agent main --model <id> --prompt-file "$RUN/prompts/D3.md"`.
2. **Label hygiene (before building).** Read the subagents' replies for relation gaps and the
   claim labels in `claims.jsonl`. For each gap, decide: map it onto an existing relation, or accept
   it as new. Log each decision as one line in `$RUN/assumptions.md`. Do not edit `claims.jsonl`; if
   two labels in one region plainly name the same thing but differ, note it there (exact-label merge
   only in this build).
3. `meta graph build "$RUN"`.
4. `meta validate "$RUN/graph/regions/<r>.json" --type RegionGraph` for each region file.
5. Read `graph/interface_candidates.json`. Its size is the first sign of whether regions talk about
   the same things. Zero candidates with two or more regions means the regions do not connect:
   add the "regions share no entities" note to `$RUN/assumptions.md`; the page must say the parts were
   researched separately.

## Kernel errors

| Command | Error | Do |
|---|---|---|
| `graph build` | claim refused (bad relation or missing source) | The kernel names the claim. Leave it out of the graph by re-gathering that region only if it is the region's main evidence; otherwise accept the drop and note it. |
| `graph build` | other | Retry once. Still failing: continue `kernel=absent`; group claims by region and subject yourself for A2, and mark "evidence structure" unchecked. |
| `validate` | field path | A kernel-written file failing its own schema is a kernel fault: save the error to `$RUN/gates/D3.error.json`, mark "graph format check" unchecked, continue. |
