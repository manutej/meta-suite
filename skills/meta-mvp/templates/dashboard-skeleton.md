# Dashboard skeletons (copy-paste) — Motion tool + Aceternity showcase

Build from these, don't start blank. Both share one palette + one custom-SVG icon set + the section
grammar. Fill data from `seed_data.json`. Always: custom inline SVG (never stock icons), mobile-ready.

## Shared head (both dashboards)
```html
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
<style>
:root{
  /* paste the /dict-color palette here as tokens; example warm-editorial set: */
  --paper:#f4ecdd;--card:#fffdf6;--ink:#231c14;--soft:#5e5644;--faint:#9b9077;--line:#e7dcc6;
  --accent:#b1502b;--accent2:#4f7355;--warm:#df9038;--cool:#2f6f8f;
  --serif:Georgia,serif;--sans:-apple-system,system-ui,sans-serif;
}
*{box-sizing:border-box;margin:0;padding:0}
/* responsive: collapse multi-col grids → 1 col; SVG scales via viewBox+width:100% */
@media(max-width:760px){ .grid3{grid-template-columns:1fr} .grid4{grid-template-columns:1fr 1fr} }
@media(max-width:460px){ .grid4{grid-template-columns:1fr} }
@media(prefers-reduced-motion:reduce){*{animation:none!important;transition:none!important}}
</style>
```

## Section grammar (both — fill per domain)
1. Hero (eyebrow · serif headline w/ one italic accent · lede)
2. Glance/stat strip (4 numbers, encode quantity → size/length per /bertin)
3. Map or relationship view (custom SVG; node size = magnitude, hue = nominal, ramp = ordered)
4. Segments / explorer (interactive filter)
5. The load-bearing panel (the domain's core: approval queue / weave / threads)
6. Interactive viewer (drawer or modal showing one entity's relationships)

## Build B — Motion (the DAILY TOOL; warm/editorial; default for the product)
Signature, load-bearing motion (CSS + small JS):
- **Layout-animated segmented control** (sliding pill): `pill.style.transform = translateX(activeBtn.offsetLeft)`.
- **Scroll-reveal stagger**: IntersectionObserver adds `.in`; `.rev{opacity:0;transform:translateY(24px);transition:.7s}` → `.in{opacity:1;transform:none}`.
- **Draw-on SVG**: `path.style.strokeDasharray=len; strokeDashoffset=len → 0` on reveal.
- **Spring modal** viewer: `transform:scale(.85)→scale(1)` with `cubic-bezier(.34,1.56,.64,1)`.
- **Eased number tickers** (requestAnimationFrame, cubic ease-out).
Robustness: also render content statically so it reads even if scripts are sandboxed.

## Build A — Aceternity (the SHOWCASE skin; dark/jewel; for stakeholder/front-door)
Signature effects (mostly CSS):
- **Aurora background**: 2–3 blurred radial-gradient blobs drifting via `@keyframes` (opacity ~.3).
- **Gradient-border glow on hover**: `::before` with `background:linear-gradient(...)` + mask-composite.
- **Animated beams** between nodes: a moving highlight via `@keyframes` on a gradient strip.
- **Bento grid** of stats; **sparkles**; **number tickers**.
- Put the **ceti engine** (templates/ceti-engine.md) at the very top as the explainer.

## Custom SVG icon rule
Draw simple, consistent 1.6px-stroke line icons in `currentColor` so they inherit the palette. Reuse
the same set across both builds. Never import a stock icon font/CDN.

## Mobile checklist (must pass)
viewport meta · ≥2 breakpoints · grids → 1 col · SVG scales (viewBox+width:100%) · tap targets ≥38px
· no primary action behind hover-only · ceti player + controls wrap.

## Verify (run before ship — see references/quality-gates.md)
JS parses · ceti contract (8 beats / 35–45s / pure render) · palette coherence · mobile markers ·
no stock icons · no hype words · no category-theory vocabulary in client output.
```
