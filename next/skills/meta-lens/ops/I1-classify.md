# I1 · classify the request

<!-- <internal>The brief as an object: domain × audience × purpose × question type(s). Question types
select the operator set; See is off in this build, so types only steer I2, D1 and A1.</internal> -->

```xml
<contract op="I1">
  <in name="raw" type="text">the user's message, verbatim</in>
  <in name="attachments" type="file[]" optional="true"/>
  <out name="brief" type="Brief" path="$RUN/brief.json">drafted outside the run, copied in by `meta run init`</out>
  <side name="route" type="text" optional="true">a meta-suite skill name, when this is not research</side>
</contract>
```

## Order (matters: no run folder exists until `run init`)

1. **Route.** If the request is really for a reusable prompt, an agent plan, a design system, a full
   bundle from one line, a review of the user's own work, a quality standard, a compounding map or a
   context setup, hand off to the matching meta-suite skill (meta-prompting, meta-planning,
   meta-design, meta-mvp, meta-review, meta-craft, meta-compound, meta-context). Tell the user in one
   sentence and stop.
2. **Mode.** `--fast` or "quick" → `fast`; otherwise `deep`.
3. **Draft the Brief** to `./meta-lens-brief.json` (the working directory; any writable path works):
   - `raw`: verbatim. `objective`: one sentence. `audience`: stated, else the most likely reader.
   - `question_types`: one or two of `profile` (one subject in depth), `compare` (options side by
     side), `synthesize` (how a field fits together). A path, trade-off or analogy question maps to
     the nearest type; record the mapping in `assumptions[]`.
   - `targets[]`, `views[]` (compare only), `medium` (only if named), `constraints[]` (stated limits).
   - `stakes`: `high` if a wrong answer could cost money, health, legal standing or reputation.
   - `vocab_allowlist[]`: the user's own field terms that collide with the banned-words list, each with
     its reason. Usually empty.
   - `assumptions[]`: every guess, one line each, with its stakes.
   - Deep mode: reread the raw message cold and classify again; if the two differ, keep the first and
     add the disagreement to `assumptions[]`.
4. `meta validate <brief> --type Brief`, then `meta run init --brief <brief> --mode <mode>` → `RUN`.
5. Write these instructions to `$RUN/prompts/I1.md`, then `meta stage begin "$RUN" I1 --agent main
   --model <your model id> --prompt-file "$RUN/prompts/I1.md"`. This token writes nothing: it records
   which model classified the request and the hash of the instructions it followed.
6. **Ask rule.** Ask one question only if an assumption is `high` stakes **and** costly or impossible
   to redo. Otherwise write the reader-relevant `assumptions[]` to `$RUN/assumptions.md` as plain
   sentences (index.md, "Two writing rules") and continue.

## Kernel errors

| Command | Error | Do |
|---|---|---|
| `validate` | field path | Fix; retry once. Second failure: keep the brief, list "request format check" as unchecked. |
| `run init` | any | Retry once. Still failing: create `runs/absent-<UTC timestamp>/` yourself, copy the brief in, continue with `kernel=absent`. |
| `stage begin` | any | Save the error; continue. Say under "how this was checked" that this stage is missing from the record. |
