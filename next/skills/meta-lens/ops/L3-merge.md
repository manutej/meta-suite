# L3 · combine the regions (plain merge)

<!-- <internal>Plain merge along declared interfaces: the v2−lens path (D-058, D-068). The full version
is a pushout/colimit with J1 same/distinct verdicts and union-find coherence; here `--plain` declares
every candidate. The kernel still refuses any merge with no declared interface.</internal> -->

```xml
<contract op="L3">
  <in name="regions" type="RegionGraph[]" path="$RUN/graph/regions/"/>
  <in name="candidates" type="Interface[]" path="$RUN/graph/interface_candidates.json"/>
  <out name="interfaces" type="Interface[]" path="$RUN/graph/interfaces.json"/>
  <out name="merged" type="DomainGraph" path="$RUN/graph/merged.json">with shape_metrics for A1</out>
</contract>
```

Start with `meta stage begin "$RUN" L3 …`; then `meta interface merge "$RUN" --plain` and
`meta validate "$RUN/graph/merged.json"`.

## Judgment (before the merge)

For each candidate: do both regions mean the same thing by this label? A pair that plainly does not
(same word, different meaning) cannot be kept apart by this build's plain merge, so:
log it (`meta ledger append --kind assumption --stage L3 --note "kept apart: <label>, <a> vs <b>,
<reason>"`), add a plain sentence for the reader to `$RUN/assumptions.md` ("'Provider' means the model
maker in one part of this research and the app maker in another."), and A2 labels any finding resting
on it **Contested**.

## Kernel errors

| Error | Do |
|---|---|
| refused: no declared interface | Expected when no candidates exist. Do not hand-write `interfaces.json`. Continue: L2 works on claims; A2 presents regions side by side as researched separately. |
| other non-zero | Retry once. Still failing: mark "combining the regions" unchecked, continue to L2. |
| `validate` fails on `merged.json` | Kernel fault: save the error, mark "combined graph format" unchecked. |
