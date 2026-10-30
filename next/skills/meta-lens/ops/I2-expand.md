# I2 · expand the request into a spec

<!-- <internal>The expanding half of an adjunction (free construction). The compressing half (I3) and
the round-trip drift gate (G5) arrive in a later slice.</internal> -->

```xml
<contract op="I2">
  <in name="brief" type="Brief" path="$RUN/brief.json"/>
  <out name="spec" type="Spec" path="$RUN/spec.json">kernel records its hash in run.json</out>
  <side name="scope" type="narrow|broad">deep mode only: `meta run scope` and spec.scope</side>
</contract>
```

## Writing the spec

Start with `meta stage begin "$RUN" I2 … --prompt-file "$RUN/prompts/I2.md"` (these instructions).

- `objective`: the brief's objective, sharpened to what a finished answer must let the reader do.
- `slots[]`: every input the answer depends on, each with `origin`: `stated` (the user said it),
  `inferred` (implied by what they said) or `default` (your choice).
- `constraints[]`, each with its `origin`.
- `success_criteria[]`: each checkable, naming what satisfies it. A finding kind: `verified`,
  `single_source`, `uncertain` (also for "disagreements shown with both sides": the page's contested
  items are findings with confidence contested). Or a gate: `G1`, `G2`, `G6`. Include at least:
  every sub-question has a finding; every claim on the page traces to a source.
- **Scope (deep mode).** Decide `narrow` or `broad` ("deep research", or a request spanning several
  fields, is `broad`). Set `spec.scope`, then `meta run scope "$RUN" narrow|broad`: it switches the
  config row (`modes.deep_narrow` / `modes.deep_broad`) that caps D1's regions and R1's passes. Fast
  mode has no scope.

## Gotchas

- **No facts.** The spec holds methodology, slots and criteria. Facts come only from D2.
- **Stay faithful.** Each `inferred` or `default` entry traces to something in the brief; delete
  requirements the user never implied.
- `meta validate "$RUN/spec.json"`. Leave `hash` out: the kernel computes it and records it in
  `run.json` and the ledger. Re-validate after any change to the spec.
- An assumption the reader needs (an inferred slot such as "sells in the EU") goes to
  `$RUN/assumptions.md` as a plain sentence; the scope reason goes to the ledger, not the page.

## Kernel errors

| Command | Error | Do |
|---|---|---|
| `validate` | field path | Fix; retry once. Second failure: keep the file, mark "spec format check" unchecked. |
| `run scope` | any | Save the error; keep `spec.scope`; caps stay at `deep_narrow`. |
| `stage begin` | any | Save the error; continue. |
