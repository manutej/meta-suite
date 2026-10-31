# Brief kit: one BRIEF.md replaces the old intake, vision, meta-prompt, spec, plan and data-model files

Copy the shape and fill it. Aim for 600–1,200 words. Every section is here because it changed the
output in tests. Don't add sections that nothing downstream reads.

```markdown
# BRIEF · {project}

## Intake (resolved; never asked)
- Audience: {who reads it}: {confidence H/M/L}, deciding signal: {…}
- Single outcome: {the one decision or action it should cause}: {confidence}
- Aesthetic: {editorial / bold / premium / warm}: {confidence}
- Scope: {one artifact type}; out of scope: {…}

## Thesis
{2–4 sentences in the reader's language: the answer, why, and how sure. The first screen says this.}

## Meta-prompt (typed; the load-bearing part)
<system>
  <role>You are the {ENGINE} for {category, not this example}. Job: {one outcome}. Tone: {voice}.</role>
  <invariants>
    - Every claim is tagged FACT (with a source id) or INFERENCE (with the facts it rests on).
    - Estimates are never presented as measured.
    - Quotes are verbatim and attributed.
    - The recommendation names a runner-up and the condition that would flip it.
    - {2–3 domain invariants: what must stay true for every instance of this category}
  </invariants>
  <typed_inputs>
    <slot name="sources" type="list[Source{id,url,publisher,date,grade A|B|C}]"/>
    <slot name="{domain_slot}" type="list[{Type}{fields}]"/>
  </typed_inputs>
  <procedure>
    1. GROUND … 2. ANCHOR on the 2–4 numbers that constrain the answer … n. RECOMMEND (with runner-up, flip condition, what would change our mind)
  </procedure>
  <output_contract type="json">{ … the exact shape seed_data.json follows … }</output_contract>
</system>

Refinement: tested against {sibling case 1}, {sibling case 2}. Edits made: {real list, or "none needed"}.

## Data model
Canonical entities, with names used exactly everywhere: `Source{…}`, `Claim{t, tag:F|I, s:[Source.id]}`, {domain entities with typed fields}.
`seed_data.json` is the single source of truth; the page renders from it.

## Modules (3–6)
For each module, one block:
- **{Module}**: what it shows, and the reader question it answers.
  - Figure: {which data → which visual encoding (position / length / hue ramp / size)}
  - States: {empty / loading / filtered / mobile collapse}, with no dead ends
  - Edge cases (≥3): {what goes wrong → the designed response}
  - Open: {ASSUME / SPIKE / ASK-{who}}

## Verification checklist
- Five habits visible (figures encode data · verdict first · tags · runner-up + flip · phone render checked)
- seed_data.json parses; every source id resolves; counts match the text
- render_audit clean at 360/390/1280; palette PASS; no hype or formal words
- {domain-specific checks, e.g. totals add up, dates are in order}
```

## Rules of thumb
- **Category, not example.** The meta-prompt should work unchanged for a sibling case. Name two and check them.
- **Typed, not prose.** If a slot or an entity has no type, the build will drift. Give it fields.
- **Encode on purpose.** Ordered values get a ramp. Quantities get length or size. Categories get position or hue. Every figure answers a reader question.
- **Honest edits only.** Record the refinement edits you actually made. Never claim passes you didn't run.
