# Meta-prompt kit (self-contained — do NOT depend on other skills being installed)

This file embeds the meta-prompting procedure AND the spec-layer generator so meta-mvp produces the
deep result even when no sibling skill is loaded. Use it verbatim.

---

## Part A — the system meta-prompt (the "engine" for the project)

A meta-prompt specifies *how to think* about the whole task category, not one instance. Build it in 5 steps:

1. **Name the category, not the example.** Not "message Priya" but "personally coordinate ~100
   entrepreneurs, voice-first." State what varies (free variables) and what's fixed (invariants).
2. **Decompose into typed slots.** Every input/intermediate/output is a named typed term:
   `members: list[Member]`, `threads: list[Thread]`, `intent: str`, output a typed JSON contract.
3. **Specify the procedure.** The step-by-step "how to think" over those slots, example-free
   (e.g. Segment → Connect → Listen → Respond → Weave → Synthesize).
4. **Bind with rules + a checkable output contract.** Ordering, invariants, confidentiality, SLA;
   then an exact JSON output a human can scan in seconds.
5. **Render as an explicit schema** (XML below). Zero solved examples — only slots, rules, procedure.

### Template — fill the slots, keep the shape
```xml
<system>
  <role>
    You are the {ENGINE_NAME} for {PROJECT}. Your job: {ONE_OUTCOME}. {MEDIUM/CONTEXT}.
    Tone: {VOICE}.
  </role>
  <invariants>
    <!-- true for every instance; never violated -->
    - {INVARIANT_1 e.g. address each entity individually, by name + real need}
    - {INVARIANT_2 e.g. identity = relationships, not a static profile}
    - {INVARIANT_3 e.g. SLA: nothing waits past {N}}
    - {INVARIANT_4 e.g. confidentiality: nothing leaks across {boundary}}
  </invariants>
  <typed_inputs>
    <slot name="{x}" type="list[{Type}]" desc="…"/>
    <slot name="{intent}" type="str" desc="what the operator wants this cycle"/>
  </typed_inputs>
  <procedure>
    1. {STEP} … 6. SYNTHESIZE — one picture + the operator's 3 highest-leverage moves.
  </procedure>
  <output_contract type="json">
    { "{summary}":"…", "{at_risk}":[…], "{actions}":[…], "{top_3_moves}":[…] }
  </output_contract>
</system>
```

### Recursive refinement (RMP) — don't ship draft 1
Draft → critique against the whole category (cover instances you did NOT use as examples? leaked
specifics? untyped slot? unbound rule? checkable output?) → emit a named edit list → apply → repeat
3–5 passes until the edit list is empty. Keep the edit history.

### Quality bar
Example-agnostic · every slot typed · sub-tasks map to sub-prompts · rules bind the slots · output
contract checkable · verified against 2–3 sibling instances, not just the user's example.

---

## Part B — the component-spec generator (meta-prompt #2: how each module is specced)

This is the SPEC LAYER — the single biggest quality driver and the thing most often dropped. After
the system meta-prompt, decompose the system into modules and spec EACH one with the same fixed shape
so they read as one product. Hand one module per sub-agent (see `spec-layer.md`).

### The fixed two-layer spec shape (every module obeys it)
```
# <Module Name>
> One-sentence definition.

## LAYER 1 — Concept brief (for the stakeholder; skimmable in 60s)
- Why it matters (2–3 sentences tied to the one outcome)
- The one moment that proves it (a concrete scene of someone experiencing it)
- How it feels (3–5 adjectives + the design principle it honors)
- What success looks like (2–3 observable signals, not vanity metrics)

## LAYER 2 — Build spec (for the builder; complete enough to hand off)
- Data model: typed entities this module owns/touches; links to other modules' entities by exact name
- States & flow: the lifecycle as a state machine (states + transitions + triggers); no dead states
- Key surfaces: where it lives + what the user sees
- Copy & voice: 2–4 real example strings in the project's register
- Edge cases & failure modes: 3–5 honest "what goes wrong" + the designed response (+ the invariant)
- Open questions / spikes: each marked ASSUME / SPIKE / ASK-{stakeholder}
```

### Shared invariants across all module specs
Concept brief resolves before build spec and is immutable downstream · every build entity traces to
the brief (no orphan data) · cross-module links name the other module's entity exactly · each spec
ends with ≥1 honest open question · example-agnostic.

### Output contract for a module spec
Both layers in order · every entity typed · state machine has no dead states · copy in-voice · ≥3
edge cases each with a response · ≥1 marked spike.

---

## How meta-mvp uses this file
- Stage META-PROMPT → Part A produces `03_META_PROMPT.md`.
- Stage SPEC → Part B + `spec-layer.md` produce `SPEC_01..NN_*.md` (one per module), via sub-agents.
If `/meta-prompting` happens to be installed you may also invoke it, but this kit is the source of
truth so the result is identical with or without it.
