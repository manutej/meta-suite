LIVE regression fixture (D-079), recorded from the first real meta-lens run (slice 1, 2026-10-29/30, run
01M651YMAFM8EPYHJ0AB13C2WQ, deep_narrow, one runner context on claude-opus-5-5). Question: which EU AI Act
obligations a mid-size SaaS company building on third-party models must meet, by when, and what is unsettled.

Kept as the run wrote them (D-097: absolute sandbox paths in ledger.jsonl, how_checked.json and gates/G6.json
were rewritten run-relative, e.g. ".../runs/<run_id>/prompts/I1.md" -> "prompts/I1.md", and the ledger's
hash chain was recomputed so `meta ledger show` still verifies it): run.json, spec.json, cover.json, claims.jsonl (201 claims, each
carrying its kernel-issued stage token), findings.jsonl (103: 73 verified, 30 single-source), conflicts.json
(3), how_checked.json, ledger.jsonl (308 events, hash chain intact, 42 deviations) and the kernel gate outputs
gates/{D3.graph,G2,G6,L2.verify,L3.merge}.json.

Trimmed: no fetched page bodies (claims hold typed facts and locators only), no rendered HTML (artifact/),
no screenshots, no prompts, no claim shards or token files (the replay rebuilds both from the stage tokens
stamped on each claim), and not gates/A2.selfcheck.json (it quotes rendered page text) or
gates/report.stdout.json (a copy of how_checked.json).

Replayed by cmd/meta/replay_test.go (TestReplayLiveEUAIAct): merge, graph build, plain interface merge,
findings verify, G2 and report are re-run on these claims and the recorded counts asserted. Where a slice-1.1
fix (D-087) legitimately changes a number, the test says which fix and why.
