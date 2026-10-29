# L3 · combine the regions (plain merge)

<!-- <internal>Plain merge along declared interfaces: the v2−lens path (D-058, D-068). The full version
is a pushout/colimit with J1 same/distinct verdicts and union-find coherence; here `--plain` declares
every candidate. The kernel still refuses any merge with no declared interface.</internal> -->

```xml
<contract op="L3">
  <in name="regions" type="RegionGraph[]" path="$RUN/graph/regions/"/>
  <in name="candidates" type="Interface[]" path="$RUN/graph/interface_candidates.json"/>
  <out name="interfaces" type="Interface[]" path="$RUN/graph/interfaces.json"/>
  <out name="merged" type="DomainGraph" path="$RUN/graph/merged.json">with shape_metrics</out>
</contract>
```

## Procedure

1. `meta stage begin "$RUN" L3 --agent main --model <id> --prompt-file "$RUN/prompts/L3.md"`.
2. **Review the candidates** before merging. For each candidate, ask: do both regions mean the same
   thing by this label? If a pair plainly does not (same word, different meaning in each region),
   write it to `$RUN/assumptions.md` as "kept apart: <label>, <region a> vs <region b>, <reason>".
   This build's plain merge cannot keep it apart, so A2 must label any finding resting on it
   **Contested**.
3. `meta interface merge "$RUN" --plain`.
4. `meta validate "$RUN/graph/merged.json" --type DomainGraph`.
5. Read `shape_metrics` from `merged.json` for A1.

## Kernel errors

| Error | Do |
|---|---|
| refused: no declared interface | Expected when no candidates exist. Do not hand-write `interfaces.json` to force it. Continue: L2 works on claims, and A2 presents regions side by side with the note "researched separately; no shared entities were found". |
| other non-zero | Retry once. Still failing: mark "combining the regions" unchecked, continue to L2. |
| `validate` fails on `merged.json` | Kernel fault: save the error, mark "combined graph format" unchecked, continue. |
