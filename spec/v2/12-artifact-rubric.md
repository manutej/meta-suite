# 12 · Artifact quality rubric

**Status:** draft for M0 freeze (D-056). **Scope:** every rendered artifact v2 ships, any domain or form.

## Purpose and use

L2 (minimum bar) is the owner's recent work, measured. L3 goes beyond it.

| Used by | How |
|---|---|
| **A2 build bar** | No handoff until every dimension is at L2 and every machine check passes. |
| **R1 rendered review** | Scores the render, not the source. Any L1 is a P1 edit. A judged L2 with a clear route to L3 is a P3 edit. |
| **Eval M8** | Blind judges and SMEs score D1, D2, D9 and D10 using the 60-second protocol. Machine results sit beside, never blended in. |

**Ship:** no dimension at L1, and all machine checks green.

## Dimensions

**D1 · Verdict on the first screen** (machine + judged)
- **L1:** The first screen holds only a title, a logo or an intro. The finding sits below the fold.
- **L2:** At 360×740, the first screen shows the title, a one- or two-sentence verdict that includes its number, and a status signal: a stamp row, a strip of 4–12 counts, or one decisive figure. Any erratum or "provisional" caveat is also above the fold.
- **L3:** The first screen alone lets an SME act. The verdict states its own limits, and every number in the strip links to its evidence.
- *Check:* Playwright at 360×740 and 1280×800: verdict element and ≥1 number in view. Judge confirms it is the real finding.

**D2 · Layering and progressive disclosure** (judged; structure counted by machine)
- **L1:** One flat scroll, or everything expanded at once.
- **L2:** Three levels. Verdict, then 4–8 named views (tabs or sections), then raw evidence behind a `details` element, drawer or sheet. Every level opens by tap. Nothing primary is hover-only.
- **L3:** Disclosure follows the reader's question. An "up a level" control. Linked views, so one choice updates all of them. A dismissible first-run coach that shows 3–4 moves.
- *Check:* DOM count of disclosure elements and views; a hover-only primary control fails.

**D3 · Palette and contrast** (machine)
- **L1:** Pure #FFF or #000 ground, raw hex inside components, any text pair below AA, or light theme only.
- **L2:** 9–19 role tokens on `:root`: ground, surface, ink, muted, rule, one accent, and pass/warn/fail. A dark twin under both `prefers-color-scheme` and `[data-theme]`. Every text pair passes AA (4.5:1, or 3:1 for large text) in both themes. The accent covers about 5% of pixels or less. No status is shown by color alone: each carries a word or a glyph.
- **L3:** The dark twin is derived by method (OKLCH lightness inversion with hue held). The page shows its contrast pairs, or lets the reader adjust them live.
- *Check:* Computed contrast on the rendered DOM in both themes; token count; hex outside `:root`; pixel scan for accent share.

**D4 · Typography** (machine + judged)
- **L1:** One default face for everything, body text under 15px, or lines longer than 80ch.
- **L2:** Exactly three families, each with a job: display, body, and mono for labels and numbers. Body text at 15–18px (17–20px for long reads), line height 1.5–1.62, measure 60–72ch. Headings sized with `clamp()` and set with `text-wrap: balance`. Mono eyebrows at 10.5–13px, uppercase, tracked. Numbers use `tabular-nums`.
- **L3:** Type carries meaning, not just rank. For example, a statement face paired with an italic "turn" face in a couplet headline.
- *Check:* Computed styles (families, minimum size, measure, numerals). Judge rates each face's job.

**D5 · Mobile and touch** (machine)
- **L1:** Sideways page scroll at 360px, targets under 24px, or no viewport meta.
- **L2:** No page-level overflow at 320, 360, 390 or 768px. Wide tables and figures scroll inside their own container and say so. Gutters of at least 16px plus the safe-area inset. Primary targets at least 44px. At least two breakpoints, with grids collapsing to one column.
- **L3:** Built for the thumb. Frequent controls sit in the bottom third. A bottom sheet with resting heights. Landscape and screens of 340px or less handled. Forced-colors mode supported.
- *Check:* Playwright asserts `scrollWidth ≤ innerWidth` at each width, measures target rects, and counts media queries.

**D6 · Data honesty and encoding** (machine + judged)
- **L1:** Numbers hand-typed into markup, illustrative data passed off as real, unlabeled axes, or a multiplier with no baseline.
- **L2:** Every figure comes from embedded data or a stated computation. Counts carry their denominators. Uncertainty is shown where it exists. Illustrative figures are labeled "Illustrative". Chart form follows the shape of the data. Every chart has a text equivalent.
- **L3:** Pre-registered results sit next to post-hoc results instead of being replaced by them. Losses get the same visual weight as wins. A reader can re-derive the headline number from the page itself.
- *Check:* Numeric literals absent from the data block; `aria-label` or table on every svg/canvas. Judge looks for overclaiming.

**D7 · Provenance: how this was checked** (machine + judged)
- **L1:** No sources, or sources with no date or method.
- **L2:** A dated sources block. Per-claim status stamps, each with a text label (Proven, Single-source, Contested; Computed or Declared). A Known limits section. A version or lock stamp, with earlier versions kept. Any erratum placed on top.
- **L3:** Receipts: each claim shows its probe, a control that stays quiet, and the matched output, one tap away. What could not be checked is stated as plainly as what was checked.
- *Check:* Assert a date, a sources list, a limits section, and labeled stamps. The judge checks that the stated checks match the claims.

**D8 · Interaction and state** (machine)
- **L1:** Controls that do nothing, state shown by color only, or no visible focus.
- **L2:** `aria-pressed` or `aria-selected` on toggles. Styled `:focus-visible`. Every control changes visible state. No console errors. Script versions pinned.
- **L3:** The reader can probe the model: reroll a seed, swap a palette and watch contrast update, or pin an item and trace paths to it. State that matters survives a reload.
- *Check:* Playwright clicks each control and diffs the DOM, tabs for focus rings, collects console errors.

**D9 · Motion** (machine + judged)
- **L1:** Ambient loops, motion that competes with the data, or no reduced-motion fallback.
- **L2:** Motion is optional and explains something. One focal motion at a time. Reduced motion jumps to a legible end frame. Transitions last 0.2–0.9s and use named easings. Animated explainers are a pure function of time, so a paused frame looks the same as the playing frame.
- **L3:** The motion is the argument. A scrubbable sequence where each beat adds exactly one idea, and the last frame works as a complete summary.
- *Check:* Emulate reduced motion: end frame reached, no running or infinite animations outside the player. Judge rates meaning.

**D10 · Voice and expertise** (judged; hype words checked by machine)
- **L1:** Hype words, filler that would fit any domain, or exclamation marks.
- **L2:** Short, plain, declarative sentences. Numbers and dates inline. Terms defined at first use. Headings that state the finding. Limits named. Corrections dated and shown openly.
- **L3:** Signature moves an SME would forward. A couplet headline (a statement, then its turn). Numbered laws, each tagged with how it is checked. A section on what the method learned. A decision table the reader can reuse.
- *Check:* Grep for hype words. A blind SME answers: "Swap the nouns. Could this page belong to another domain?" A yes scores L1.

## First 60 seconds: judge protocol

1. Open at **360×740**, light theme. Start a timer; don't scroll.
2. **0–10s:** Write the page's claim in one sentence. If you cannot, D1 is L1.
3. **10–30s:** Scroll once and tap one control. Did the next level answer the question the first screen raised (D2)? Did anything scroll sideways (D5)?
4. **30–50s:** Find the evidence for the headline number. Three taps or fewer counts as L2. Are a date and a method visible (D6, D7)?
5. **50–60s:** Switch to dark, then 1280px. Note contrast, layout or hierarchy breaks.
6. **Rate:** "would forward to a client" and "learned something in 60s" (yes/no); D1, D2, D9, D10 (1–3). Judges never see arm, trace or the other arm.

## Exemplars that set Level 2

Referenced by title and type only. Exemplars are read locally and never copied here (P13).

| Title | Type | Bar it sets |
|---|---|---|
| Atlas Pocket | phone-first graph explorer | D2, D5, D8 |
| CETI Seedbank v7 | visual research instrument | D3, D4, D7, D10 |
| How to write a team wiki that an AI assistant can answer from · R6 | essay with figures and a kit | D4, D6, D7 |
| meta-score v3 Repair Loop | experiment milestone record | D1, D6, D7 |
| What the Accessibility Layer actually does | claims-with-receipts page | D7 L3 |
| Accessibility Layer 0.5.0 | release milestone | D1, D7 |
| Research Atlas | research dashboard | D1, D2, D6 |
| Calm Systems | systems-visualization lookbook | D6, D9, D10 |

**Their gaps (L3 targets):** some 30–34px chips, inner-scrolling tables at 360px, few re-derivable headline numbers, receipts on only one page.
