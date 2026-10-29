# D2 · gather evidence, one subagent per region (core op)

<!-- <internal>Local sections: each claim holds only for its region and its source. No synthesis here;
the only merge is L3.</internal> -->

```xml
<contract op="D2">
  <in name="cover" type="Cover" path="$RUN/cover.json"/>
  <in name="probes" type="json" path="$RUN/probes.json"/>
  <out name="shards" type="Claim[]" path="$RUN/claims/&lt;agent&gt;.jsonl">written only by `meta claim add`</out>
  <out name="claims" type="Claim[]" path="$RUN/claims.jsonl">kernel merge</out>
  <out name="g2" type="gate report" path="$RUN/gates/G2.json"/>
</contract>
```

## Procedure (main thread)

1. For each region `r` in `cover.json`, write the subagent brief below to `$RUN/prompts/D2-<r>.md`,
   filled with that region's element questions, its probes, the relation vocabulary, `$ROOT` and `$RUN`.
2. For each region, get its token:
   `meta stage begin "$RUN" D2 --region <r> --agent gather-<r> --model <gather model id> --prompt-file "$RUN/prompts/D2-<r>.md"`.
   Use the strongest available model for regions marked `contested` in `probes.json`, a fast model
   for the rest. Record which in `$RUN/assumptions.md`.
3. Launch **all region subagents in parallel**, each with web search and its token. If the host has
   no subagents, run the regions one after another in the main thread and add "gathering was not
   isolated" to `$RUN/assumptions.md`.
4. When all return: `meta claims merge "$RUN"`, then `meta gate G2 "$RUN" > "$RUN/gates/G2.json"`.
5. **Thin regions.** A region with fewer than 3 claims, or whose claims share one `independence_key`,
   gets one targeted re-gather (counts against `loop.d2_reentry_max`), then re-run step 4.

**`kernel=absent`:** no tokens. Each subagent appends its claim lines to
`$RUN/claims/gather-<r>.jsonl` directly (brief step 3 replaced). These claims carry no kernel ids, so
nothing downstream may be called verified.

## Subagent brief (write this, filled, to the prompt file)

> You research one region of a larger question. Gather evidence only; do not summarise or conclude.
> Region: {region id and element questions}. Probes: {probes}. Relations you may use: {vocab}.
>
> 1. Search the web. Read primary and practitioner sources first (official documents, filings,
>    standards, datasets, trade and specialist press, named experts) and follow the probes. Aim for
>    at least three different publishers.
> 2. For each fact you can source, write one claim, one fact per claim:
>    `{"subject":…, "relation":<from the list>, "object":…, "value":…?, "value_type":…, "unit":…?,
>    "as_of":…?, "region":"{r}", "source":{"locator":<URL or document reference>, "retrieved_at":<UTC
>    ISO time>, "independence_key":<publisher's registrable domain, e.g. the site's own domain>}}`.
>    Name subjects and objects with the plainest standard label, the same way every time.
> 3. Record it: `printf '%s' '<claim JSON>' | sh "{ROOT}/bin/meta" claim add "{RUN}" --token "{TOKEN}"`.
>    The kernel stamps id, run, region, agent and model; never type those yourself.
> 4. On a schema error, fix the named field and resend once. If it fails again, drop that claim.
> 5. Never write a claim without a real locator you opened. No long quotes. If a fact needs a
>    relation that is not on the list, use the closest one and note the gap in your reply.
> 6. Reply with: claims added, distinct publishers, the probes you could not answer, relation gaps.
>    Nothing else.

## Kernel errors

| Command | Error | Do |
|---|---|---|
| `stage begin` | any | That region cannot write claims. Retry once; if it still fails, skip the region and list it as "not gathered". |
| `claim add` | schema | Subagent fixes and resends once (brief step 4). |
| `claim add` | lock or I/O | Subagent waits a moment and resends once; then reports the loss. |
| `claims merge` | any | Retry once. Still failing: kernel results stop here; continue `kernel=absent` using the shards as written. |
| `gate G2` | non-zero with findings | Not a crash: the gate failed. Open the report, re-gather the regions it names (within the re-gather cap), re-run. Remaining failures are shown on the page. |
