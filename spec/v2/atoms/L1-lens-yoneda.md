# L1 · lens.yoneda

**Purpose:** Define an entity by everything it relates to, then prove the profile is fine-grained
enough to tell it apart from its nearest neighbours.

**Internal construct:** Yoneda (an object is determined by its relations), made practical as a probe
set plus a discriminability test.

**In:** the merged `DomainGraph`, `Brief.targets[]`, `relation_vocab`
**Out:** `Finding[kind=profile_trait]`, a discriminability score, and possibly a
`Finding[kind=coarse_model]`

**Contract**
- **Probes come only from `relation_vocab`.** Each relation type is probed in both directions: what
  the target does, and what is done to it. There are no hardcoded "standard probes per kind."
- **Empty probes.** A probe with no answer triggers a targeted D2 re-gather, within the loop caps.
- **Discriminability test.**
  - The profile must separate the target from the `lens.yoneda.neighbours` nearest objects
    (default 2) on at least 1 sourced probe.
  - Whether two answers are "the same" is a J1 `judge` task.
  - If the test fails, add probes up to `lens.yoneda.max_probes`. If it still fails, emit a
    `coarse_model` finding.

**Check:** discriminability ≥ `lens.yoneda.min_discriminability`. Every trait passes G2.

**Guards against:** generic profiles that would fit any entity of the same kind.

**Evidence:** E3 (the Yoneda probe in categorical deep-research evaluation).
