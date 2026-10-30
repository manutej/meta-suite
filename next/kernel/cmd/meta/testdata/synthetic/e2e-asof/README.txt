SYNTHETIC, hand-written: the ORIGINAL slice-1 e2e claims (before D-087 aligned their as_of values), kept as a
second scenario for the D-097 time rules. It uses ../e2e's brief.json and cover.json.
Outcome under D-097, asserted by TestAsOfScenario:
  - Org Alpha located_in City One: as of 2026, a.example with undated b.example supporting -> verified;
    b.example's own undated finding -> single-source;
  - headcount 12000 in 2026-01: a.example (12 thousand) with undated c.example supporting -> verified;
    c.example's own undated 12000 -> single-source, contested by 15000 (b.example, 2025): one disagreement;
    12000 in 2026-01 vs 15000 in 2025 are disjoint periods: one change-over-time record;
  - has_part Unit Delta (a.example twice) and Org Beta has_part Unit Gamma -> single-source.
claims.jsonl lines: {"agent": <gather agent>, "region": <region>, "claim": <ClaimInput>}.
