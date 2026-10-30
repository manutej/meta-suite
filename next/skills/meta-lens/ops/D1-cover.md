# D1 · split the question into regions (core op)

<!-- <internal>A decomposition tree of sub-questions plus a cover of a declared universe by regions.
Overlaps and triple overlaps are computed by the kernel, never declared; they become the cell complex
for the gluing gate in a later slice.</internal> -->

```xml
<contract op="D1">
  <in name="spec" type="Spec" path="$RUN/spec.json"/>
  <in name="max_regions" type="int">`meta config show --mode &lt;run mode&gt;`</in>
  <out name="cover" type="Cover" path="$RUN/cover.json">tree, universe, regions (with contested), relation_vocab, labels</out>
  <out name="probes" type="Probes" path="$RUN/probes.json">{region_id: [question, …]}, questions only</out>
</contract>
```

Start with `meta stage begin "$RUN" D1 … --prompt-file "$RUN/prompts/D1.md"`. Then write the cover
in this order, because each part names the one before it.

1. **Tree.** Sub-questions, root first; each node `id`, `question`, `answer_type` (fact, number,
   list, judgment), `children[]`. Children in the order a reader needs them answered: A1 uses this
   order as the page's reasoning sequence.
2. **Universe.** Every leaf id, plus the kinds of entity the answer will name (`kind:<name>`).
3. **Regions** (at most `max_regions`), each a set of universe elements, named for its angle on the
   question (a source community, a time span, a sub-system), never for the expected answer. Every
   element sits in some region; regions connect through shared elements (add a bridging element when
   one stands alone). Set `"contested": true` on a region whose probes point at a known dispute: D2
   gives it the strongest model.
4. **Relation vocabulary** (aim for 6–15), see below.
5. **Labels.** `labels[]`: the canonical name for every entity you expect more than one region to
   mention (the law, each actor role, each obligation set, each institution), written exactly as
   gatherers must write them. Claims verify each other only when subject and object labels match
   after normalization, so two regions that write "GPAI model provider" and "provider of a
   general-purpose AI model" never agree. `claim add` notes any label outside this list.
6. **Probes** to `$RUN/probes.json`: 3–5 questions per region — what practitioners argue about, where
   the official rule fails in practice, how it is enforced or measured, which practitioner sources
   to read beyond top search results. Questions only.
7. `meta validate "$RUN/cover.json"` and `meta validate "$RUN/probes.json"`.

## Relation vocabulary: gotchas

Each relation: `name` (snake_case), `value_type`, `unit_dimension` for numeric or relative, and a
`description` that says what subject, object and value are. Say in the description whether the
relation **carries a value** or not; the kernel does not enforce it, gatherers follow it.

- **Statuses:** one relation (e.g. `has_status`) with the allowed values listed in its description
  (`proposed | adopted | in force | pending | withdrawn`). A fact with no fitting status gets its own
  relation, not the nearest status.
- **Counts:** a numeric relation with `unit_dimension: "count"` ("9 of 27 member states": value 9,
  the denominator in the object or a second claim). Never encode a count as a derived percentage.
- **Shares and rates:** `unit_dimension: "ratio"`, unit `%`.
- **Money:** `unit_dimension: "currency"`; unit is an ISO code (`EUR`, `USD`), value in base units.
  Codes are never converted.
- **Dates:** a fact *about* a date (applies from, published on) puts the ISO date (`YYYY-MM-DD`) in
  the value. A fact that *holds as of* a date (a status, a count, "none so far") carries `as_of`:
  claims verify each other only inside one as_of month, so the same words dated August and October
  stay two findings.
- **Positions:** `position_on` with values `supports | opposes | mixed`, for stakeholder views; a
  two-sided question also needs the role or option each side picks as a value-carrying relation.
- Keep one fallback (`states`, object = a short factual sentence, no value) and expect it to be
  single-source: free text rarely matches across publishers.

## Kernel errors

| Error | Do |
|---|---|
| uncovered element | Add it to the closest region; retry. |
| regions not connected | Add a bridging element to the isolated region and its nearest neighbour; retry. |
| too many regions | Merge the two regions with the largest overlap; retry. |
| other field path | Fix; retry once. Second failure: keep the cover, mark "split covers the question" unchecked. |
