# R1 · blind review (core op)

<!-- <internal>An adversarial critic that never saw the work being made. Output is an EditScript only;
the critic cannot edit. Blind attestation is recorded; on a host without subagents, blind=false.</internal> -->

```xml
<contract op="R1">
  <in name="page" type="html" path="$RUN/artifact/index.html"/>
  <in name="shots" type="png[]" path="$RUN/artifact/shots/" optional="true"/>
  <in name="spec" type="Spec" path="$RUN/spec.json"/>
  <in name="findings" type="Finding[]" path="$RUN/findings.jsonl"/>
  <in name="how_checked" type="json" path="$RUN/how_checked.json"/>
  <out name="edits" type="EditScript" path="$RUN/review/edits.&lt;n&gt;.json"/>
</contract>
```

## Procedure (main thread)

1. `n` = this pass number, from 1. Write the reviewer brief below to `$RUN/prompts/R1-<n>.md`.
2. `meta stage begin "$RUN" R1 --agent reviewer-<n> --model <reviewer id> --prompt-file "$RUN/prompts/R1-<n>.md"`.
3. Launch **one fresh-context subagent** with the best available model. Give it only the brief and
   the paths above: no chat, no drafts, no claims, no notes. No subagents on this host → do the review
   yourself from the files alone and set `blind: false`.
4. Save its JSON to `$RUN/review/edits.<n>.json`; `meta validate "$RUN/review/edits.<n>.json" --type EditScript`.
5. **Apply.** Sort edits P1 → P3. P1 and P2: route each to the earliest stage it names (see
   index.md) and re-run from there. P3: deep mode only, on the final pass. Then review again
   (`n+1`) until no P1/P2 edits remain or `modes.<mode>.review_passes` is reached (fast: one pass).
6. Edits left open at the cap go under Known limits on the page.

## Reviewer brief (write this to the prompt file)

> You are reviewing a finished research page you did not make. Read only: {page}, {screenshots if
> any}, {spec}, {findings}, {how_checked}. You cannot edit; you return edits.
>
> 1. Open the page as a reader at phone size. In 10 seconds, can you state its claim? Would an
>    expert in this field forward it? Did you learn something in 60 seconds?
> 2. Check, citing element ids: (a) faithful to the spec's objective and criteria; (b) any statement
>    stronger than its finding's kind (single-source stated as settled, contested made one-sided);
>    (c) claims on the page with no finding behind them; (d) expert concerns missing; (e) text that
>    would fit any other subject if the nouns were swapped.
> 3. Score each rubric dimension D1–D10 (verdict first, layering, palette and contrast, type, mobile,
>    data honesty, provenance, interaction, motion, voice) 1–3 from the render, or from the source if
>    no screenshots exist (say which). Any 1 is a P1 edit. A 2 with a clear route to 3 is a P3 edit.
> 4. Return only JSON: `{"edits":[{"target_stage":"D2|D3|L3|L2|A1|A2","target_id":"<element,
>    section or finding id>","op":"add|remove|rewrite|regather|relabel","rationale":"…","priority":"P1|P2|P3"}],
>    "reviewer_model":"<your model id>","blind":true}`.

## Kernel errors

| Error | Do |
|---|---|
| `validate` field path | Ask the reviewer once to fix the named field (same isolation). Still invalid: apply only edits that parse, mark "review format" unchecked. |
| `stage begin` | Save the error; run the review anyway; the ledger lacks this pass. |
