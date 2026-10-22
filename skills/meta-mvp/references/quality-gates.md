# Quality gates — run before shipping (stage 7)

All must pass. Eyeballing doesn't scale; run the checks.

## 1 · JS parses (every dashboard)
```
node -e "const fs=require('fs');const h=fs.readFileSync('FILE.html','utf8');
[...h.matchAll(/<script>([\s\S]*?)<\/script>/g)].forEach((m,i)=>{try{new Function(m[1])}catch(e){console.log('block'+i,e.message)}});
console.log('parsed')"
```

## 2 · Data logic sanity
Eval the data-generation block in Node; assert counts (members, regions sum, segments) match the
narrative. Confirm `seed_data.json` parses and counts are right.

## 3 · ceti contract (the top animation)
- Exactly **8 beats**; total **35–45s**; captions ≤118 chars.
- `render(t)` is pure (no setTimeout/scroll-driven state); paused frame == playing frame.
- Scrub to any t shows a clean frame (only one working scene + one detail scene lit).
- Reduced-motion jumps to the final, legible frame.

## 4 · Palette coherence (dict-color — bundled)
```
# dict-color ships INSIDE this skill (scripts/dict_color.py), stdlib-only — no install needed.
# <skill-dir> = the folder holding this skill's SKILL.md (works in any agent, plugin or not)
python3 "<skill-dir>/scripts/dict_color.py" audit FILE.html --threshold 12   # expect PASS (~90–100%)
```
If a richer external dict-color tool is installed, prefer it. If the script can't be located,
**degrade gracefully** (like meta-design does): check by hand that the file collapses onto ≤12
core colors and any real brand hexes are present — flag unverified ones. Never block the ship on
a missing helper. If a real brand palette is used, confirm those hexes are present.

## 5 · Mobile-ready
- Exactly one `viewport` meta with `width=device-width`.
- ≥2 responsive breakpoints; grids collapse to one column; SVG scales via viewBox.
- No primary action behind hover-only.

## 6 · Anti-slop / craft
- Custom SVG icons present (grep for inline `<svg`), no stock icon CDN.
- No hype words ("unlock", "supercharge", "seamless", "revolutionize", "game-changing").
- No raw category-theory vocabulary in any customer-facing file
  (`grep -iE "functor|morphism|colimit|yoneda|adjunction"` over client artifacts → empty).

## 7 · Fact-check / grounding
- Every non-obvious claim in `01_RESEARCH.md` / `02_VISION.md` has a source or a flagged assumption.
- Sensitive facts (`sensitivity:"personal"`) appear in NO customer-facing artifact.
- For high-stakes bundles, run a verification sub-agent to re-check claims + math independently.

## 8 · Packaging
- Files numbered in order; `present_files` called.
- Optional appendix offered, not forced.
- No deploy / publish / send performed without explicit user approval.
