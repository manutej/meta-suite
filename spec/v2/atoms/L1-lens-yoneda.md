# L1 · lens.yoneda

**Purpose:** Know an entity by everything it relates to: what it buys, sells, depends on, measures,
fears and displaces, and what depends on it. Then prove the profile is fine-grained enough to tell
this entity apart from its nearest neighbours.

**Internal construct:** Yoneda. An object is determined by its relations. In practice this means a
chosen probe set and a discriminability test.

**In:** `DomainGraph`, target object(s)
**Out:** `Finding[kind=profile_trait]`, `discriminability` score

**Contract**
- The probe set comes from `relation_vocab`, plus standard probes for the object's kind. Probes are
  questions, never answers.
- An empty probe triggers a targeted D2 re-gather, up to the tier cap.
- **Discriminability test:** the profile must separate the target from at least 2 nearest-neighbour
  objects on at least 1 sourced probe. If it fails, add probes (up to the cap). If it still fails,
  report "coarse model" as a finding.

**Check:** discriminability ≥ the threshold in config. G2 applies to every trait.

**Guards against:** generic profiles that would fit any entity in the category, which are the
opposite of what experts recognize.

**Evidence:** E3 (Yoneda probe in categorical deep-research evaluation).
