# D1 · split the question into regions (core op)

<!-- <internal>A decomposition tree of sub-questions plus a cover of a declared universe by regions.
Overlaps and triple overlaps are computed by the kernel, never declared; they become the cell complex
for the gluing gate in a later slice.</internal> -->

```xml
<contract op="D1">
  <in name="spec" type="Spec" path="$RUN/spec.json"/>
  <in name="max_regions" type="int">modes.&lt;mode&gt;.max_regions, read from config (see index.md)</in>
  <out name="cover" type="Cover" path="$RUN/cover.json"/>
  <out name="probes" type="json" path="$RUN/probes.json">{region_id: [question, …]}</out>
</contract>
```

## Procedure

1. `meta stage begin "$RUN" D1 --agent main --model <id> --prompt-file "$RUN/prompts/D1.md"`.
2. **Tree.** Break the objective into sub-questions, root first. Each node: `id`, `question`,
   `answer_type` (fact, number, list, judgment), `children[]`. Order children in the order a reader
   needs them answered: this order is the reasoning sequence A1 uses.
3. **Universe.** `universe[]` = every leaf sub-question id, plus the kinds of entity the answer will
   name (as `kind:<name>` ids).
4. **Regions.** Group the universe into at most `max_regions` regions, each a set of element ids.
   - Every element sits in at least one region.
   - Neighbouring regions must share at least one element, so every region connects to the rest
     through shared elements. If one stands alone, add a bridging element to it and a neighbour.
   - Name each region for its angle on the question (a source community, a time span, a
     sub-system), never for the answer you expect.
5. **Relation vocabulary.** A short controlled list (aim for 6–15) of relations claims may use. Each
   has `name` (snake_case), `value_type` (`qualitative`, `numeric` or `relative`) and, for numeric or
   relative, a unit dimension.
6. **Insider probes.** For each region write 3–5 questions to `$RUN/probes.json`: what practitioners
   argue about, where the official rule does not hold in practice, how it is actually enforced or
   measured, and which practitioner sources (trade press, standards bodies, filings, forums of
   record) to read beyond top search results. Questions only, never answers. Mark a region `contested`
   in the probes file if its probes point at a known dispute.
7. `meta validate "$RUN/cover.json" --type Cover`. The kernel rejects a cover that misses an element
   or whose regions do not connect.

## Kernel errors

| Error | Do |
|---|---|
| uncovered element | Add the element to the closest region; retry. |
| regions not connected | Add a bridging element shared by the isolated region and its nearest neighbour; retry. |
| too many regions | Merge the two regions with the largest overlap; retry. |
| other field path | Fix; retry once. Second failure: keep the cover, mark "split covers the question" unchecked. |
