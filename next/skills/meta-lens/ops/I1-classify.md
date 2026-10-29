# I1 · classify the request

<!-- <internal>The brief as an object: domain × audience × purpose × question type(s). Question types
select the operator set; See is off in this build, so types only steer I2, D1 and A1.</internal> -->

```xml
<contract op="I1">
  <in name="raw" type="text">the user's message, verbatim</in>
  <in name="attachments" type="file[]" optional="true"/>
  <out name="brief" type="Brief" path="$RUN/brief.json"/>
  <side name="route" type="text" optional="true">a meta-suite skill name, when this is not research</side>
</contract>
```

## Procedure

1. **Route.** If the request is really for a reusable prompt, an agent plan, a design system, a full
   bundle from one line, a review of the user's own work, a quality standard, a compounding map or a
   context setup, hand off to the matching meta-suite skill (meta-prompting, meta-planning,
   meta-design, meta-mvp, meta-review, meta-craft, meta-compound, meta-context). Tell the user in one
   sentence and stop. Otherwise continue.
2. **Mode.** `--fast` or "quick" in the request → `fast`; otherwise `deep`.
3. **Draft the Brief** to `${TMPDIR:-/tmp}/meta-lens-brief.json`:
   - `raw`: verbatim. `objective`: one sentence. `audience`: stated, else the most likely reader.
   - `question_types`: one or two of `profile` (one subject in depth), `compare` (two or more options
     side by side), `synthesize` (how a field fits together). A path, trade-off or analogy question maps
     to the nearest type; add the mapping to `assumptions[]` with the words "unsupported type".
   - `targets[]` (the subjects named), `views[]` (the two sides, compare only), `medium` (only if the
     user named one), `constraints[]` (stated limits: time window, place, budget, sources to avoid).
   - `stakes`: `high` if a wrong answer could cost money, health, legal standing or reputation;
     otherwise `normal`.
   - `vocab_allowlist[]`: plain terms of the user's own field that collide with the banned-words list,
     each with its reason. Usually empty.
   - `assumptions[]`: every guess you made, each one line, with its stakes.
4. **Deep mode only:** reread the raw message cold and classify again. If the two classifications
   differ, keep the first and add the disagreement to `assumptions[]`.
5. `meta validate "${TMPDIR:-/tmp}/meta-lens-brief.json" --type Brief`.
6. `meta run init --brief "${TMPDIR:-/tmp}/meta-lens-brief.json" --mode <mode>` → prints the run folder;
   set `RUN` to it. If `$RUN/brief.json` is missing, copy the brief there.
7. Write the classification instructions you followed to `$RUN/prompts/I1.md`, then
   `meta stage begin "$RUN" I1 --agent main --model <your model id> --prompt-file "$RUN/prompts/I1.md"`.
8. **Ask rule.** Ask the user one question only if an assumption is `high` stakes **and** costly or
   impossible to redo. Otherwise copy `assumptions[]` into `$RUN/assumptions.md` and continue.

## Kernel errors

| Command | Error | Do |
|---|---|---|
| `validate` | field path | Fix that field; retry once. Second failure: save the error to `$RUN/gates/I1.error.json` once `RUN` exists, keep the brief as written, list "request format check" as unchecked. |
| `run init` | any | Retry once. Still failing: create `runs/absent-<UTC timestamp>/` yourself, copy the brief in, and continue with `kernel=absent`. |
| `stage begin` | any | Save the error; continue. The ledger will lack this stage; say so under "how this was checked". |
