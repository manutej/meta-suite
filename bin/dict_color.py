#!/usr/bin/env python3
"""
dict_color.py — self-contained palette generator + coherence auditor for meta-suite.

Stdlib only (no pip installs). Bundled so meta-mvp / meta-design palette gates run
with zero external dependencies. If a richer external dict-color tool is installed,
prefer it; this is the graceful-degradation fallback that always works.

Usage:
  # Audit: what fraction of a file's color usage collapses onto <= N core colors?
  python3 dict_color.py audit FILE.html --threshold 12
      -> prints "on-palette: NN%  (distinct clusters: K, core: N)"; exit 0 if coherent.

  # Generate: emit a small, harmonious palette for a mood.
  python3 dict_color.py generate --mood warm --size 4
      -> prints one hex per line.

Moods: warm, earthy, dark, vivid, cool, mono, editorial  (see MOODS below).
Coherence = share of weighted color usage covered by the top `threshold` clusters.
"""
import argparse
import colorsys
import re
import sys
from collections import Counter

# ---------- color parsing ----------

_HEX = re.compile(r"#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})\b")
_RGB = re.compile(r"rgba?\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)", re.I)


def _hex_to_rgb(h):
    h = h.lstrip("#")
    if len(h) == 3:
        h = "".join(c * 2 for c in h)
    return tuple(int(h[i:i + 2], 16) for i in (0, 2, 4))


def _rgb_to_hex(rgb):
    return "#{:02x}{:02x}{:02x}".format(*(max(0, min(255, int(round(c)))) for c in rgb))


def extract_colors(text):
    """Return Counter of rgb-tuple -> usage count found in the text."""
    counts = Counter()
    for m in _HEX.finditer(text):
        counts[_hex_to_rgb("#" + m.group(1))] += 1
    for m in _RGB.finditer(text):
        counts[(int(m.group(1)), int(m.group(2)), int(m.group(3)))] += 1
    return counts


# ---------- clustering (no numpy) ----------

def _dist2(a, b):
    # weighted RGB distance (approx perceptual: green weighted highest)
    dr, dg, db = a[0] - b[0], a[1] - b[1], a[2] - b[2]
    return 2 * dr * dr + 4 * dg * dg + 3 * db * db


def cluster(counts, merge_radius=42):
    """Greedy merge of near-identical colors. merge_radius ~ perceptual sameness."""
    thr2 = merge_radius * merge_radius * 9  # scale to weighted metric above
    clusters = []  # each: {"anchor": rgb, "weight": int}
    for rgb, w in counts.most_common():
        placed = False
        for c in clusters:
            if _dist2(rgb, c["anchor"]) <= thr2:
                c["weight"] += w
                placed = True
                break
        if not placed:
            clusters.append({"anchor": rgb, "weight": w})
    clusters.sort(key=lambda c: c["weight"], reverse=True)
    return clusters


def audit(path, threshold):
    try:
        with open(path, "r", encoding="utf-8", errors="ignore") as f:
            text = f.read()
    except OSError as e:
        print(f"dict_color: cannot read {path}: {e}", file=sys.stderr)
        return 2
    counts = extract_colors(text)
    total = sum(counts.values())
    if total == 0:
        print("on-palette: n/a  (no colors found)")
        return 0
    clusters = cluster(counts)
    core = clusters[:threshold]
    covered = sum(c["weight"] for c in core)
    pct = round(100 * covered / total)
    core_hexes = ", ".join(_rgb_to_hex(c["anchor"]) for c in core[:threshold])
    verdict = "PASS" if (len(clusters) <= threshold or pct >= 90) else "REVIEW"
    print(f"on-palette: {pct}%  (distinct clusters: {len(clusters)}, "
          f"core kept: {len(core)}/{threshold})  [{verdict}]")
    print(f"core palette: {core_hexes}")
    return 0 if verdict == "PASS" else 1


# ---------- generation ----------

# Each mood: list of (hue 0-1, sat 0-1, light 0-1) seeds; we sample `size` from a
# smooth walk across them so palettes stay harmonious at any size.
MOODS = {
    "warm":      [(0.03, 0.72, 0.55), (0.09, 0.75, 0.60), (0.13, 0.55, 0.72), (0.98, 0.45, 0.42)],
    "earthy":    [(0.08, 0.45, 0.42), (0.11, 0.38, 0.55), (0.28, 0.30, 0.45), (0.05, 0.25, 0.30)],
    "dark":      [(0.62, 0.55, 0.28), (0.72, 0.45, 0.35), (0.50, 0.60, 0.40), (0.00, 0.00, 0.10)],
    "vivid":     [(0.55, 0.85, 0.55), (0.85, 0.80, 0.58), (0.13, 0.90, 0.60), (0.33, 0.75, 0.50)],
    "cool":      [(0.55, 0.45, 0.55), (0.60, 0.40, 0.45), (0.48, 0.35, 0.62), (0.70, 0.30, 0.40)],
    "mono":      [(0.58, 0.12, 0.20), (0.58, 0.10, 0.40), (0.58, 0.08, 0.62), (0.58, 0.06, 0.85)],
    "editorial": [(0.09, 0.62, 0.55), (0.47, 0.55, 0.55), (0.00, 0.00, 0.06), (0.10, 0.20, 0.93)],
}


def generate(mood, size):
    seeds = MOODS.get(mood.lower())
    if not seeds:
        print(f"dict_color: unknown mood '{mood}'. Known: {', '.join(sorted(MOODS))}",
              file=sys.stderr)
        return 2
    out = []
    n = max(1, size)
    for i in range(n):
        # position along the seed walk
        pos = (i / max(1, n - 1)) * (len(seeds) - 1) if n > 1 else 0
        lo = int(pos)
        hi = min(lo + 1, len(seeds) - 1)
        frac = pos - lo
        h = seeds[lo][0] + (seeds[hi][0] - seeds[lo][0]) * frac
        s = seeds[lo][1] + (seeds[hi][1] - seeds[lo][1]) * frac
        li = seeds[lo][2] + (seeds[hi][2] - seeds[lo][2]) * frac
        r, g, b = colorsys.hls_to_rgb(h % 1.0, li, s)
        out.append(_rgb_to_hex((r * 255, g * 255, b * 255)))
    print("\n".join(out))
    return 0


def main(argv=None):
    p = argparse.ArgumentParser(prog="dict_color.py", description=__doc__,
                                formatter_class=argparse.RawDescriptionHelpFormatter)
    sub = p.add_subparsers(dest="cmd", required=True)

    pa = sub.add_parser("audit", help="palette coherence of a file")
    pa.add_argument("file")
    pa.add_argument("--threshold", type=int, default=12,
                    help="max core colors expected (default 12)")

    pg = sub.add_parser("generate", help="emit a harmonious palette for a mood")
    pg.add_argument("--mood", default="warm")
    pg.add_argument("--size", type=int, default=4)

    args = p.parse_args(argv)
    if args.cmd == "audit":
        return audit(args.file, args.threshold)
    if args.cmd == "generate":
        return generate(args.mood, args.size)
    return 2


if __name__ == "__main__":
    sys.exit(main())
