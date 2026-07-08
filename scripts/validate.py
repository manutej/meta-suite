#!/usr/bin/env python3
"""
meta-suite repo integrity check — the suite dogfoods its own bar.

Runs in CI (.github/workflows/validate.yml) and locally:  python scripts/validate.py

Checks, with zero third-party dependencies (stdlib only):
  1. .claude-plugin/plugin.json and marketplace.json are valid JSON with the required keys.
  2. Exactly 9 skills/*/SKILL.md, each opening with YAML frontmatter carrying name + description,
     and each frontmatter `name:` equal to its directory name.
  3. The bundled palette tool bin/dict_color.py byte-compiles (py_compile).
  4. Every committed .svg parses as well-formed XML and contains no <script>.
  5. Every relative Markdown link / image path resolves on disk (anchors, URLs, and
     brace/wildcard globs are skipped by design).

Exit code is non-zero if any check fails, so it can gate a merge.
"""
from __future__ import annotations
import json, re, sys, py_compile
from pathlib import Path
from xml.dom import minidom

ROOT = Path(__file__).resolve().parent.parent
errors: list[str] = []
checks = 0


def ok(cond: bool, msg: str) -> None:
    global checks
    checks += 1
    if not cond:
        errors.append(msg)


def ignored(p: Path) -> bool:
    parts = p.relative_to(ROOT).parts
    return any(seg in ("node_modules", ".git", ".vercel", "__pycache__", "_archive") for seg in parts)


# 1 — manifests
for name, required in (
    (".claude-plugin/plugin.json", ("name", "version", "description")),
    (".claude-plugin/marketplace.json", ("name", "owner", "plugins")),
):
    f = ROOT / name
    try:
        data = json.loads(f.read_text())
        for key in required:
            ok(key in data, f"{name} missing required key: {key}")
    except Exception as e:  # surface any manifest failure clearly
        ok(False, f"{name} is not valid JSON: {e}")

# 2 — skills + frontmatter
skills = sorted(ROOT.glob("skills/*/SKILL.md"))
ok(len(skills) == 9, f"expected 9 skills, found {len(skills)}")
for skill in skills:
    rel = skill.relative_to(ROOT)
    text = skill.read_text()
    if not text.startswith("---"):
        ok(False, f"{rel}: missing YAML frontmatter")
        continue
    fm = text.split("---", 2)[1]
    m = re.search(r"^name:\s*(.+)$", fm, re.M)
    ok(bool(m), f"{rel}: frontmatter missing `name:`")
    ok("description:" in fm, f"{rel}: frontmatter missing `description:`")
    if m:
        ok(m.group(1).strip() == skill.parent.name,
           f"{rel}: name '{m.group(1).strip()}' != dir '{skill.parent.name}'")

# 3 — bundled tool compiles
tool = ROOT / "bin" / "dict_color.py"
ok(tool.exists(), "bin/dict_color.py missing")
if tool.exists():
    try:
        py_compile.compile(str(tool), doraise=True)
    except py_compile.PyCompileError as e:
        ok(False, f"bin/dict_color.py does not compile: {e}")

# 4 — SVGs well-formed, no scripts
for svg in ROOT.rglob("*.svg"):
    if ignored(svg):
        continue
    rel = svg.relative_to(ROOT)
    try:
        minidom.parseString(svg.read_bytes())
    except Exception as e:
        ok(False, f"{rel}: not well-formed XML: {e}")
    ok("<script" not in svg.read_text(errors="ignore").lower(), f"{rel}: contains <script>")

# 5 — relative markdown links resolve
link_re = re.compile(r"(?<!\!)\[[^\]]+\]\(([^)]+)\)|!\[[^\]]*\]\(([^)]+)\)")
for md in ROOT.rglob("*.md"):
    if ignored(md):
        continue
    rel = md.relative_to(ROOT)
    for a, b in link_re.findall(md.read_text()):
        target = (a or b).split("#")[0].strip()
        if not target or target.startswith(("http://", "https://", "mailto:")):
            continue
        if any(c in target for c in "{}*"):  # globs by design
            continue
        ok((md.parent / target).resolve().exists(), f"{rel}: dead link -> {target}")

print(f"validate.py: {checks} checks, {len(errors)} error(s)")
for e in errors:
    print(f"  ✗ {e}")
sys.exit(1 if errors else 0)
