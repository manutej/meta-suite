SYNTHETIC, hand-written: the ORIGINAL slice-1 e2e claims (before D-087 aligned their as_of values), kept as a
second scenario for the D-097 time rules. It uses ../e2e's brief.json and cover.json.
Intended outcome, asserted by TestAsOfScenario:
  - Org Alpha located_in City One: a.example as_of 2026 + b.example undated -> verified (an undated statement
    counts for any period);
  - headcount 12000: a.example (12 thousand, as_of 2026-01) + c.example (undated) -> verified, contested;
    headcount 15000: b.example as_of 2025 -> single-source, contested; one disagreement (the undated 12000
    overlaps 2025), no change-over-time record;
  - has_part Unit Delta (a.example twice) and Org Beta has_part Unit Gamma -> single-source.
claims.jsonl lines: {"agent": <gather agent>, "region": <region>, "claim": <ClaimInput>}.
