# Dashboard playbook — how to hit WOW, repeatably

Default architecture: **Motion editorial for the daily tool, Aceternity skin for the showcase.**
One shared data model + one shared SVG icon set feed both; only the motion layer swaps.

## The shared spine (build once, reuse)

1. **One typed data model** (stage 5) — both dashboards bind to the same `seed_data.json`.
2. **A custom inline-SVG icon set** — never stock icons. Draw simple, consistent line icons in the
   brand palette. Reuse across both builds.
3. **A section grammar** — hero → glance/stats → map or relationship view → segments/explorer →
   the load-bearing panel (weave / approval queue / whatever the domain needs) → interactive viewer.
   A new dashboard = new data binding against this grammar.

## Palette (dict-color — bundled, stdlib-only)

`scripts/dict_color.py` ships inside this skill; no install needed. Run it as `<skill-dir>/scripts/dict_color.py`,
where `<skill-dir>` is the folder holding this skill's SKILL.md. Prefer a richer external dict-color tool if present;
degrade gracefully (derive/audit by hand against the ≤12-core-colors rule) if the script can't be found.

- Generate a palette: `python3 "<skill-dir>/scripts/dict_color.py" generate --mood <warm|earthy|dark|vivid|cool|mono|editorial> --size 4`.
- Name the colors when you present (names carry intent).
- If a real brand palette exists, use it; if you can't extract exact hex, derive + flag honestly.
- Audit coherence: `python3 "<skill-dir>/scripts/dict_color.py" audit <file>.html --threshold 12` → expect PASS.

## Encoding (/bertin)

- Map data type → visual variable: ordered status (warmth, stage) → an ordered hue/value ramp;
  quantity (ties, hours) → size/length; nominal (region, segment) → position/hue.
- Don't decorate; encode. Every visual element should answer a question the stakeholder has.

## Build A — Aceternity (showcase skin)
Dark, jewel-toned, effect-rich. Signature moves: aurora/spotlight background, bento grid,
gradient-border glow on hover, animated beams between nodes, sparkles, number tickers, a slide-in
drawer as the interactive viewer. Use for launch/showcase/front-door surfaces. Compose
/editorial-dashboard + /shadcn-concept-page + /frontend-design for polish.

## Build B — Motion (daily tool)
Warm/light editorial. Signature moves: layout-animated segmented control (sliding indicator),
scroll-reveal stagger, draw-on SVG (stroke-dashoffset), eased number tickers, a spring modal as the
interactive viewer. Motion is load-bearing, not decorative; respects `prefers-reduced-motion`.

## The top animation (/ceti-explainer, inline)
For a single deployable file, implement the ceti **method** inline rather than the separate build
pipeline: one deterministic clock (`render(t)` pure function of time), a persistent **anchor** +
cross-faded **working-zone** scenes (one lit at a time) + a **detail band**, exactly **8 beats**,
~35–45s, with play/pause, a scrub slider, chapter dots, captions, and reduced-motion → final frame.
Use role-token colors so it re-skins with the brand. Apply /motion-media-l2 motion laws (purposeful,
one focal motion per beat, ≥30% whitespace, no AI-slop).

## Mobile-ready (always)
- `<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">`
- Responsive breakpoints (e.g. 1040 / 760 / 600); multi-column grids collapse to one column.
- SVG uses `viewBox` + `width:100%` so it scales; the ceti player scales; controls wrap.
- Tap targets ≥ ~38px; no hover-only affordances for primary actions.

## Anti-slop
No stock icons, no perpetual ambient motion that competes with data on a daily tool, no hype words,
no unlabeled abstract shapes, no raw category-theory vocabulary in client outputs. If it looks like a
default LLM page, redo it.
