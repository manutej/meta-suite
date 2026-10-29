# L2 · verify by independent agreement

<!-- <internal>Pullback: the fibre product of claims over a shared key (subject + relation + object). A
key is verified when ≥ lens.pullback.min_independent distinct independence_keys agree. Same-key value
disagreements are conflict candidates for G3 (not in this slice). L2 is harness, not lens: it stays on
in the v2−lens arm.</internal> -->

```xml
<contract op="L2">
  <in name="claims" type="Claim[]" path="$RUN/claims.jsonl"/>
  <out name="findings" type="Finding[]" path="$RUN/findings.jsonl">kinds verified | single_source</out>
  <out name="conflicts" type="json" path="$RUN/conflicts.json">same key, different value</out>
</contract>
```

## Procedure

1. `meta stage begin "$RUN" L2 --agent main --model <id> --prompt-file "$RUN/prompts/L2.md"`.
2. `meta findings verify "$RUN"`.
3. Validate one line in ten (at least three) of `findings.jsonl`: write each to a temp file and
   `meta validate <file> --type Finding`.
4. **Read, don't rewrite.** The kernel decides `verified` and `single_source`. You may write a
   `plain_text` sentence for A2 in your notes, but never change a finding's kind or evidence.
5. Every entry in `conflicts.json` becomes a **Contested** item on the page, showing both values and
   their sources. Never pick a winner or average the values.
6. **Coverage check.** For each leaf sub-question in `cover.json`'s tree, judge which findings
   answer it (their evidence claims come from a region holding that leaf). Write the map to
   `$RUN/coverage.json` as `{leaf_id: [finding ids]}`. A leaf with no verified finding is "not
   settled" on the page; a leaf with no finding is "not answered". If those exceed a third of the
   leaves and the re-gather cap allows, return to D2 for their regions.

## Kernel errors

| Error | Do |
|---|---|
| any non-zero from `findings verify` | Retry once. Still failing: no finding may be called verified. Present claims grouped by subject as **Unverified**, mark "independent-source agreement" unchecked. |
| `validate` on a finding fails | Kernel fault: save the error, mark "finding format" unchecked, continue. |
