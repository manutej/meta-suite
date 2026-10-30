# A1 · choose the page's form and layout

<!-- <internal>Form from question type(s); layout from the merged graph's shape metrics (hub
concentration, longest chain, two-sidedness, temporal density). Rules are config data
(`artifact.forms`, `artifact.layout`), never code.</internal> -->

```xml
<contract op="A1">
  <in name="brief" type="Brief" path="$RUN/brief.json"/>
  <in name="sequence" type="text[]">D1 tree order: the reasoning sequence</in>
  <in name="shape" type="json" path="$RUN/graph/merged.json#shape_metrics"/>
  <in name="host_types" type="text[]">artifact types the host offers for this run's output, if any</in>
  <out name="plan" type="ArtifactPlan" path="$RUN/plan.json"/>
</contract>
```

Start with `meta stage begin "$RUN" A1 …`. Rules come from `meta config show` (`artifact.forms`,
`artifact.layout`); the one that fired goes in `rule_id`.

1. **Medium.** A medium the user named wins outright (`rule_id: medium.stated`).
2. **Form** by the brief's first question type: `profile` → profile page; `compare` → decision
   workbench, disagreements first; `synthesize` → atlas or report; data to watch → dashboard; a process
   to explain → animated explainer (only through an installed explainer skill; else a report with a
   stepped figure). A second question type adds one section built for it.
3. **Layout** from `shape_metrics` against `artifact.layout`: one dominant hub → hub-first navigation;
   long chains → sequential narrative; two-sidedness → split view; dated relations (dates in values or
   as_of) → timeline spine. Combine when two are over threshold. No metrics → sequential narrative;
   say why in `rationale`.
4. **Section order** = the reasoning sequence from D1's tree, never graph order.
5. **Host type.** Set `host_type` only when the run's output will be published as that host type (a
   document, a deck, a dashboard). When the output is a local page, omit `host_type`.
6. `meta validate "$RUN/plan.json"`.

## Kernel errors

| Error | Do |
|---|---|
| `validate` field path | Fix; retry once. Second failure: keep the plan, mark "page plan format" unchecked. |
| config unreadable | Use the defaults above; `rule_id` = `default.<form>`; log a deviation. |
