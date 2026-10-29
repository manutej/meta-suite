# I2 · expand the request into a spec

<!-- <internal>The expanding half of an adjunction (free construction). The compressing half (I3) and
the round-trip drift gate (G5) arrive in a later slice.</internal> -->

```xml
<contract op="I2">
  <in name="brief" type="Brief" path="$RUN/brief.json"/>
  <out name="spec" type="Spec" path="$RUN/spec.json"/>
</contract>
```

## Procedure

1. `meta stage begin "$RUN" I2 --agent main --model <id> --prompt-file "$RUN/prompts/I2.md"` (write
   these instructions there first).
2. Write `$RUN/spec.json`:
   - `objective`: from the brief, sharpened to what a finished answer must let the reader do.
   - `slots[]`: every input the answer depends on (`name`, `type`, `required`, `origin`). `origin` is
     `stated` (the user said it), `inferred` (implied by what they said) or `default` (your choice).
   - `constraints[]`, each with its `origin`.
   - `success_criteria[]`: each one checkable and naming what satisfies it: a finding kind
     (`verified`, `single_source`) or a gate (`G2`, `G6`). Include at least: every sub-question has
     at least one finding; every claim on the page traces to a source.
   - Deep mode: decide scope, `narrow` or `broad`, and add it to `$RUN/assumptions.md` with a one-line
     reason ("deep research", or a request spanning several fields, is `broad`). It sets which config
     row caps D1 and R1: `modes.deep_narrow` or `modes.deep_broad`. Fast mode uses `modes.fast`.
3. **No facts.** The spec holds methodology, slots and criteria. Nothing about the subject's answer
   goes here; facts come only from D2.
4. **Stay faithful.** Each `inferred` or `default` entry must trace to something in the brief. Delete
   requirements the user never implied.
5. `meta validate "$RUN/spec.json" --type Spec`. Leave `hash` out; the kernel computes it. If the
   validator requires it, record that as a contract mismatch in `$RUN/gates/I2.error.json` and continue.

## Kernel errors

| Command | Error | Do |
|---|---|---|
| `validate` | field path | Fix; retry once. Second failure: keep the file, mark "spec format check" unchecked. |
| `stage begin` | any | Save the error; continue. |
