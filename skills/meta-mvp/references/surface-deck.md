# Surface: an HTML pitch deck

**One slide visible at a time.** Hide every other slide with `display:none` (or `visibility:hidden` and `opacity:0`) and let nothing bleed through. Stacked or ghosted slides are the most common deck failure, and the render audit checks for it after stepping through slides.

**Navigation:**
- arrow keys, plus Space and PageDown;
- swipe on touch;
- visible Prev and Next buttons;
- a slide counter "n / N";
- a contents jump;
- the `#n` hash updates, and `ArrowRight` works on load without a click.

**It must run in a sandboxed viewer** (app file previews, embedded frames). Wrap every `history.*`, `location` write, `localStorage`/`sessionStorage` and clipboard call in `try { … } catch (e) {}`. One unguarded `history.replaceState` crashes the whole deck in a preview, even though it works when opened as a file.

**Shape, typically 12–20 slides:**
1. Title.
2. **The ask and the verdict** (decision, budget, timeline).
3. The problem in the audience's own terms and numbers.
4. The evidence: one claim per slide, each with a figure (forest plot, dot plot, before and after, range).
5. Where the evidence is weak or mixed: say so plainly.
6. Options and costs.
7. Risks, each with its mitigation.
8. The pilot design: metrics, a decision point, kill criteria.
9. The ask again.
10. Sources, graded A/B/C.

**An interactive model**, at least one, bound to seed_data, typically on the costs or pilot slide:
- a budget or ROI calculator with editable inputs (seats, months, rates) that recomputes totals live;
- or a weighting or scenario control that changes the recommendation.

The decision-maker should be able to test the decision with their own numbers.

**Every slide:**
- a heading that states the claim (not a topic label);
- one figure or one tight list;
- a Fact or Inference tag on each claim;
- a source footnote.

**Phone:**
- slides reflow to one column;
- text 16 px or more;
- no slide content cut off;
- figures scale with their viewBox.

**Palette:** 4–6 tokens plus a dark twin. Run `dict_color.py audit` and expect PASS.
