# R1 · blind review (core op)

<!-- <internal>An adversarial critic that never saw the work being made. Output is an EditScript only;
the critic cannot edit. Blind attestation is recorded; blind=true only for a truly separate
context.</internal> -->

```xml
<contract op="R1">
  <in name="page" type="html" path="$RUN/artifact/index.html"/>
  <in name="shots" type="png[]" path="$RUN/artifact/shots/" optional="true"/>
  <in name="spec" type="Spec" path="$RUN/spec.json"/>
  <in name="findings" type="Finding[]" path="$RUN/findings.jsonl"/>
  <in name="how_checked" type="HowChecked" path="$RUN/how_checked.json"/>
  <out name="edits" type="EditScript" path="$RUN/review/edits.&lt;n&gt;.json">ids checked against the run</out>
</contract>
```

## Order (main thread)

1. `n` = pass number from 1. Write the reviewer brief below to `$RUN/prompts/R1-<n>.md`;
   `meta stage begin "$RUN" R1 --agent reviewer-<n> --model <reviewer id> --prompt-file …`.
2. Run the reviewer in a **fresh context** with the best available model, given only the brief and
   the input files: no chat, drafts, claims or notes. Pick the first that works:
   - the host's subagent tool → `blind: true`;
   - no subagent tool: copy only the input files into an empty folder and run a headless CLI session
     there (`claude -p "$(cat <brief>)"` with that folder as its working directory) → `blind: true`,
     because it is a separate process that sees only those files;
   - neither: review it yourself from the files alone → `blind: false`.
   Log which with `meta ledger append --kind attestation --stage R1`.
3. Save its JSON to `$RUN/review/edits.<n>.json`; `meta validate` it. Inside the run this also checks
   every `target_id` and every finding or conflict id cited in a rationale.
4. **Apply.** P1 and P2: route each edit by index.md's stage table (page wording → A2, even when the
   reviewer named an earlier stage; log the re-route as a deviation) and re-run from the earliest
   stage hit. P3: deep mode, final pass only. Review again (`n+1`) until no P1/P2 remain or
   `modes.<mode>.review_passes` is reached (fast: one pass).
5. Edits left open at the cap go under Known limits on the page.

## Reviewer brief (write this to the prompt file)

> You are reviewing a finished research page you did not make. Read only: {page}, {screenshots if
> any}, {spec}, {findings}, {how_checked}. You cannot edit; you return edits.
>
> 1. Open the page as a reader at phone size. In 10 seconds, can you state its claim? Would an
>    expert in this field forward it? Did you learn something in 60 seconds?
> 2. Check, citing ids exactly as they appear in the files: (a) faithful to the spec's objective and
>    criteria; (b) any statement stronger than its finding's kind (single-source stated as settled,
>    contested made one-sided); (c) claims on the page with no finding behind them; (d) expert
>    concerns missing; (e) text that would fit any other subject if the nouns were swapped.
> 3. Score each rubric dimension D1–D10 (verdict first, layering, palette and contrast, type, mobile,
>    data honesty, provenance, interaction, motion, voice) 1–3 from the render, or from the source if
>    no screenshots exist (say which). Any 1 is a P1 edit. A 2 with a clear route to 3 is a P3 edit.
> 4. Route each edit to the stage that must change:
>    - `A2` — anything fixable on the page: wording, headings, labels, stamps, evidence display,
>      layout details, the "how this was checked" text. Most edits are A2.
>    - `A1` — the page's form or overall layout (timeline, split view, section order).
>    - `D2` — evidence missing, thin or misread: a region must be researched again (`regather`).
>    - `D3` — one thing named two ways inside a region. `L3` — two parts joined on a word that means
>      different things in each. `D1` — a missing sub-question. `I2` — the question was misread.
> 5. Return only JSON: `{"edits":[{"target_stage":"…","target_id":"<an element id on the page, a
>    section id or a finding id>","op":"add|remove|rewrite|regather|relabel","rationale":"…",
>    "priority":"P1|P2|P3"}],"reviewer_model":"<your model id>","blind":true}`.

## Kernel errors

| Error | Do |
|---|---|
| `validate`: schema or unknown id | Ask the reviewer once (same isolation) to fix the named fields. Still invalid: apply only edits that validate on their own, mark "review format" unchecked. |
| `stage begin` | Save the error; run the review anyway; the record lacks this pass. |
