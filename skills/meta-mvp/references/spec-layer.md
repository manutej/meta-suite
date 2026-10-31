# Spec layer — generate per-module specs with sub-agents (the missing piece)

This stage is what made the original build feel deep. Do NOT skip it. After the system meta-prompt
and the data model, the system is decomposed into 3–6 modules, and EACH gets its own two-layer spec
(concept brief + build spec, format in `meta-prompt-kit.md` Part B), written by a dedicated sub-agent.

## Procedure

1. **Decompose.** From the data model + plan, list the system's modules (e.g. for a community
   coordination product: Invite/Onboarding, Conversation/Threads, Match/Weave, Home/Architecture).
   Aim for 3–6 — each a coherent object that another module links to by name.

2. **Write the component-spec meta-prompt to a file** (`03b_COMPONENT_SPEC_META_PROMPT.md`) so the
   spec shape is explicit and shared. Use Part B of `meta-prompt-kit.md` verbatim, adapted to the domain.

3. **Spawn one sub-agent per module, concurrently** (single message, multiple Agent calls). Each
   packet contains: the shared world facts; the path to the component-spec meta-prompt to obey; the
   canonical data-model entity names to link to; and the exact output path
   `SPEC_0N_<module>.md`. Tell each to reply only with the path + one line.

   Sub-agent prompt skeleton:
   ```
   You are writing ONE module spec for {PROJECT}.
   WORLD FACTS: {3–6 grounded bullets — never contradict them}
   1. Read the spec contract at {path to component-spec meta-prompt}.
   2. Write the spec for the {MODULE} module following it EXACTLY (Layer 1 concept brief, then
      Layer 2 build spec: typed entities, no-dead-state machine, surfaces, in-voice copy, ≥3 edge
      cases w/ responses, open questions marked SPIKE/ASSUME/ASK).
   3. You own entities {…}; link to {canonical Member/etc.} by name.
   4. Save to {abs path}/SPEC_0N_<module>.md. Reply only with the path + a one-line summary.
   ```

4. **Verify each spec** against the output contract (both layers, typed entities, no dead states,
   in-voice copy, ≥3 edge cases, ≥1 spike) before integrating. A quick programmatic grep works:
   each file has "Concept brief" + "Build spec" + links the canonical entity + contains "SPIKE".

5. **Cross-check composition:** every cross-module link must name another module's entity exactly as
   that module names it. Fix mismatches so the specs compose into one product.

## Why sub-agents (not inline)
Parallel + context-isolated: each spec gets full attention and a clean packet, and they run
concurrently so depth doesn't cost wall-clock. This is the mechanism that turns "a dashboard" into
"a fully-specified system."

## Output of this stage
`03b_COMPONENT_SPEC_META_PROMPT.md` + `SPEC_01..NN_*.md` — the spec layer the thin skill was missing.
