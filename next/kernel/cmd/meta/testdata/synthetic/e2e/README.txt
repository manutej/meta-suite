SYNTHETIC end-to-end scenario, hand-written for the kernel's tests. Not recorded from a run, and not domain
data: entities ("Org Alpha", "City One"), publishers (*.example, reserved by RFC 2606) and numbers are
placeholders chosen to exercise each rule: cross-region agreement (verified), one publisher cited twice
(still single-source), a unit conversion (12 thousand -> 12000 count), a same-key numeric disagreement
(conflict), two agents gathering one region, and a single-source fact.
D-087: findings now agree only inside one as_of window (config lens.pullback.as_of_window = month), so the
claims meant to agree carry matching as_of values (year "2026" for located_in, January 2026 for headcount).
The fixture recorded from the first real meta-lens run joins this one: ../../../../../testdata/live/eu-ai-act
(D-079).
claims.jsonl lines: {"agent": <gather agent>, "region": <region>, "claim": <ClaimInput>}.
