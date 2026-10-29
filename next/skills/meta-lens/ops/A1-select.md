# A1 · choose the page's form and layout

<!-- <internal>Form from question type(s); layout from the merged graph's shape metrics (hub
concentration, longest chain, two-sidedness, temporal density). Rules are config data
(`artifact.forms`, `artifact.layout`), never code.</internal> -->

```xml
<contract op="A1">
  <in name="brief" type="Brief" path="$RUN/brief.json"/>
  <in name="sequence" type="text[]">D1 tree order: the reasoning sequence</in>
  <in name="shape" type="json" path="$RUN/graph/merged.json#shape_metrics"/>
  <in name="host_types" type="text[]">artifact types the host offers, if any</in>
  <out name="plan" type="ArtifactPlan" path="$RUN/plan.json"/>
</contract>
```

## Procedure

1. `meta stage begin "$RUN" A1 --agent main --model <id> --prompt-file "$RUN/prompts/A1.md"`.
2. **Medium first.** A medium the user named wins outright (`rule_id: medium.stated`).
3. **Form** from config `artifact.forms` by the brief's first question type. Defaults if config is
   unreadable: `profile` → profile page; `compare` → decision workbench, disagreements first;
   `synthesize` → atlas or report; data meant to be watched → dashboard; a process to explain →
   animated explainer (only through an installed explainer skill; otherwise a report with a stepped
   figure). A second question type adds one section built for it.
4. **Layout** from `shape_metrics` against the thresholds in config `artifact.layout`: one dominant
   hub → hub-first navigation; long chains → sequential narrative; two-sidedness → split view; many
   dated claims → timeline spine. Combine when two are over threshold. No metrics (merge refused or
   `kernel=absent`) → sequential narrative, and say why in `rationale`.
5. **Section order** = the reasoning sequence from D1's tree, never graph order.
6. **Host type.** If the host offers an artifact type that fits the form (a document, a deck, a
   dashboard), set `host_type`. Otherwise leave it empty: A2 hand-builds a page.
7. Write `$RUN/plan.json` with `form`, `layout`, `host_type`, `rule_id` (the config rule that fired),
   `rationale` (one line). `meta validate "$RUN/plan.json" --type ArtifactPlan`.

## Kernel errors

| Error | Do |
|---|---|
| `validate` field path | Fix; retry once. Second failure: keep the plan, mark "page plan format" unchecked. |
| config unreadable | Use the defaults above; `rule_id` = `default.<form>`; note it in `$RUN/assumptions.md`. |
