# D3 · build the per-region graphs

<!-- <internal>A category presentation per region: objects are entities, morphisms are typed relations
each backed by claim ids. Entity merge in this slice is exact normalized label only, within a region;
cross-region look-alikes become interface candidates for L3.</internal> -->

```xml
<contract op="D3">
  <in name="claims" type="Claim[]" path="$RUN/claims.jsonl"/>
  <out name="regions" type="RegionGraph" path="$RUN/graph/regions/&lt;region&gt;.json"/>
  <out name="candidates" type="Interface[]" path="$RUN/graph/interface_candidates.json"/>
</contract>
```

Start with `meta stage begin "$RUN" D3 …`, then `meta graph build "$RUN"` and
`meta validate` each region file.

## Judgment (before and after the build)

- **Gaps.** Read the agents' replies for relation and label gaps. For each, decide: map it onto an
  existing relation or label, or accept it. Log each decision with
  `meta ledger append "$RUN" --kind assumption --stage D3 --note "<gap>: <decision>"`. Never edit
  `claims.jsonl`.
- **Near-duplicate labels.** Two labels in one region that plainly name one thing stay two entities
  in this build (exact-label merge only): log it the same way. If the pair carries the region's main
  evidence, re-gather with the label from `cover.json` (counts against the re-gather cap).
- **Candidates.** `graph/interface_candidates.json` is the first sign of whether regions talk about the
  same things. Zero candidates with two or more regions: add "The parts of the research share no
  common subjects, so each was researched separately." to `$RUN/assumptions.md`.

## Kernel errors

| Command | Error | Do |
|---|---|---|
| `graph build` | claim refused | The kernel names it. Re-gather only if it is the region's main evidence; else accept the drop and log it. |
| `graph build` | other | Retry once. Still failing: continue `kernel=absent`; group claims by region and subject for A2; mark "evidence structure" unchecked. |
| `validate` | field path | A kernel file failing its own schema is a kernel fault: save to `$RUN/gates/D3.error.json`, mark "graph format check" unchecked. |
