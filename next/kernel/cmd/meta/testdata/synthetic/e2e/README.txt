SYNTHETIC end-to-end scenario, hand-written for the kernel's tests. Not recorded from a run, and not domain
data: entities ("Org Alpha", "City One"), publishers (*.example, reserved by RFC 2606) and numbers are
placeholders chosen to exercise each rule: cross-region agreement (verified), one publisher cited twice
(still single-source), a unit conversion (12 thousand -> 12000 count), a same-key numeric disagreement
(conflict), two agents gathering one region, and a single-source fact.
A fixture recorded from the first real meta-lens run replaces or joins this one (D-079).
claims.jsonl lines: {"agent": <gather agent>, "region": <region>, "claim": <ClaimInput>}.
