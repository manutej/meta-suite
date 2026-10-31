# Surface: a strategy, ICP, report or research-digest page

**Grammar** (adapt it, don't pad it):
1. **Hero:** the verdict in one sentence, a confidence chip, and 2–3 key numbers. Fits on a phone's first screen.
2. **Glance:** one figure that carries the whole argument, e.g. a positioning map, a ranked bar chart, a range chart.
3. **The body modules** from BRIEF.md. Each opens with its claim as the heading and its figure, then the evidence.
4. **The decision:** the recommendation, the runner-up, what would flip it, and the next 90 days.
5. **Evidence:** a sources list with A/B/C grades and dates, plus open questions. Collapsed by default when long.

**Figures** are custom inline SVG with a viewBox, drawn from seed_data:
- Every axis and value is labelled.
- Text colours come from theme tokens.
- No decorative blobs.
- Icons don't count as figures.

**Interaction**, only where it answers a question: a weighting slider that re-ranks, a filter, a tab between segments.

**Tags:** a small Fact or Inference chip on every claim. Estimates are labelled "≈" with their basis.

**Mobile:**
- one viewport meta;
- grids collapse to one column below about 640 px;
- tap targets of about 40 px or more;
- no hover-only actions;
- tables scroll inside their own container.

**Palette:** 4–6 tokens plus a dark twin. Run `dict_color.py audit` and expect PASS.

**Length:** say it once. If a section only repeats sources, collapse it.
