# L2 · verify by independent agreement

<!-- <internal>Pullback: the fibre product of claims over a shared key (subject + relation + object) and
as_of window. A key is verified when ≥ lens.pullback.min_independent distinct independence_keys agree.
Same-key value disagreements are conflict candidates for G3 (not in this slice). L2 is harness, not
lens: it stays on in the v2−lens arm. Keys: lens.pullback.min_independent, lens.pullback.as_of_window.</internal> -->

```xml
<contract op="L2">
  <in name="claims" type="Claim[]" path="$RUN/claims.jsonl"/>
  <out name="findings" type="Finding[]" path="$RUN/findings.jsonl">kernel: verified | single_source</out>
  <out name="conflicts" type="Conflicts" path="$RUN/conflicts.json">kernel: same key and window, different value</out>
  <out name="coverage" type="Coverage" path="$RUN/coverage.json">yours: {leaf_id: [finding ids]}</out>
</contract>
```

Start with `meta stage begin "$RUN" L2 …`, then `meta findings verify "$RUN"` (`--help` states the
agreement rule, including the as_of window).

## Judgment

- **Read, don't rewrite.** The kernel decides `verified` and `single_source`. Never change a
  finding's kind, confidence or evidence; write your own prose for the page in A2.
- **Conflicts.** Each `conflicts.json` record has a `kind`. `disagreement`: a **Contested** item on
  the page with both values and their sources; never pick a winner or average. `changed_over_time`:
  the page shows the change, each value with its date, oldest first; it is not a dispute, but it is
  never left out and never summed up as "sources agree".
- **Dated findings.** A finding's `as_of` is the period its sources cover; the page says "as of …".
  An undated source counts for any period, so it can confirm a dated one or disagree with it.
- **Bad claims.** A finding that rests on a claim its source does not support: retract the claim
  (`meta claim retract "$RUN" --id … --reason … --stage L2`), then re-run merge, graph build,
  interface merge and findings verify. The page never quietly drops a finding.
- **Coverage.** For each leaf sub-question in `cover.json`'s tree, list the findings that answer it
  in `$RUN/coverage.json`; `meta validate` checks every leaf and finding id. A leaf with no verified
  finding is "not settled" on the page; a leaf with none is "not answered". If those exceed a third of
  the leaves and the re-gather cap allows, return to D2 for their regions.

## Kernel errors

| Error | Do |
|---|---|
| `findings verify` non-zero | Retry once. Still failing: no finding may be called verified. Present claims grouped by subject as **Unverified**; mark "independent-source agreement" unchecked. |
| `validate` coverage: unknown id | Fix the id (findings ids change when claims change: rebuild coverage after any re-run). |
