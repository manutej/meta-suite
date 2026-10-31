# eval/: the habit check

A cheap, automatic floor that every v1.3 change must clear **before** it is spent on a blind test (spec/v1.3/SPEC.md §5, step 0).

```
node eval/habit_check.js <page.html> --case 4A            # check against the frozen v1 winner
node eval/habit_check.js <page.html>                      # measure only
node eval/habit_check.js <page.html> --case 4B --shots out/   # also save 360/390/1280 screenshots
node eval/habit_check.js --make-baseline <case> <frozen_v1_final.html>
```

Needs Node and Playwright (Chromium). Exit codes: 0 = pass, 1 = below the floor, 2 = could not run.

## What it measures

These are the habits the owner rewarded in blind tests, plus a render audit:

| Check | Floor (relative to the frozen v1 page) |
|---|---|
| Data figures: an SVG with drawn shapes or a canvas. Icons and sprite sheets are not counted | ≥ 75% of v1 |
| Fact/inference/evidence tags | ≥ 75% of v1 |
| Interactive controls | ≥ 75% of v1 |
| Distinct source links; source-quality grades | ≥ 75% of v1 |
| Decision stated on the first screen (or the first slide after the title) | required if v1 had it |
| A runner-up named; what would flip the ranking | required if v1 had it |
| Length: the dossier guard | ≤ 2× v1's words |
| Render at 360, 390 and 1280 px: horizontal overflow, script errors, overlapping text (also after two slide steps) | must be clean |

## Calibration (2026-10-30, lanes 4A and 4B, 8 frozen pages)

The check passes both v1 pages and fails every page the owner ranked lower:

| Lane | v1 | plain Opus | v2 deep | v2 fast |
|---|---|---|---|---|
| 4A ICP | PASS | FAIL: grades, runner-up | FAIL: figures, tags, runner-up, length | FAIL: figures, tags, grades, runner-up, length |
| 4B deck | PASS | FAIL: figures, controls | FAIL: figures, verdict, length, **overlapping slides at all three widths** | FAIL: figures, tags, verdict, flip, overlap at 1280 |

It catches the stacked-slide failure the owner flagged in v2 deep's deck. No other check did.

**Fixed 2026-10-31:** the first-screen detector ignored headings that wrap a styled word, which gave a false "no verdict" on two pages. It now reads every visible text node. The re-run kept 8/8 agreement.

**Limits, stated plainly:**
- It was calibrated on these same 8 pages, so its 8/8 agreement is in-sample.
- It is a floor, not a ranker: it does not order plain Opus against v2.
- The decision and runner-up detectors are keyword heuristics. A page can game them, or phrase things in ways they miss.
- The owner's blind judgment stays the final gate. Re-check agreement every time a new judged pair exists.

The frozen pages it was calibrated on live outside the repo (they are generated test outputs). `baselines.json` keeps only the numbers.
