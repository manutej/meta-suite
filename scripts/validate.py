#!/usr/bin/env python3
"""
meta-suite repo integrity check — the suite dogfoods its own bar.

Runs in CI (.github/workflows/validate.yml) and locally:  python scripts/validate.py

Checks, with zero third-party dependencies (stdlib only):
  1. .claude-plugin/plugin.json and marketplace.json are valid JSON with the required keys.
  2. Exactly 9 skills/*/SKILL.md, each opening with YAML frontmatter that a STRICT parser accepts
     (Codex, Cursor, claude.ai upload, the Skills API — not just Claude Code's lenient loader):
     name == directory, name matches the Agent Skills spec (<=64, [a-z0-9-]), description is a
     block scalar or a safe plain/quoted scalar, <=1024 chars, metadata.short-description <=200
     (the claude.ai upload limit), and only upload-safe keys are used.
     marketplace `author`/`owner` are objects and plugin/marketplace versions agree.
  3. The bundled palette tool skills/meta-mvp/scripts/dict_color.py byte-compiles (py_compile).
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

# 2 — skills + frontmatter (strict: what the Agent Skills spec and upload surfaces accept)
UPLOAD_KEYS = {"name", "description", "license", "compatibility", "metadata", "allowed-tools"}
NAME_RE = re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*$")


def frontmatter(text: str) -> dict[str, str]:
    """Tiny strict parser for the subset we use: `key: value` plus `key: >-`/`|` blocks.
    Raises ValueError on anything a real YAML parser would reject (e.g. unquoted ': ')."""
    out: dict[str, str] = {}
    lines = text.strip("\n").split("\n")
    i = 0
    while i < len(lines):
        line = lines[i]
        m = re.match(r"^([A-Za-z][\w-]*):(?:\s+(.*))?$", line)
        if not m:
            raise ValueError(f"unparseable frontmatter line: {line[:60]!r}")
        key, val = m.group(1), (m.group(2) or "").strip()
        i += 1
        if not val:  # one-level map, e.g. `metadata:` -> `  short-description: "..."`
            while i < len(lines) and lines[i].startswith((" ", "\t")):
                sub = frontmatter(lines[i].strip())
                out.update({f"{key}.{k}": v for k, v in sub.items()})
                i += 1
            out.setdefault(key, "")
            continue
        if val[:1] in (">", "|"):
            block = []
            while i < len(lines) and (lines[i].startswith((" ", "\t")) or not lines[i].strip()):
                block.append(lines[i].strip())
                i += 1
            val = " ".join(b for b in block if b)
        elif val[:1] in ("'", '"'):
            if len(val) < 2 or val[-1] != val[0]:
                raise ValueError(f"{key}: unterminated quoted value")
            val = val[1:-1]
        elif ": " in val or " #" in val or val[:1] in "&*!%@`{[,":
            raise ValueError(f"{key}: plain value contains ': ' / ' #' or a YAML indicator — "
                             "use a `>-` block or quotes")
        out[key] = val
    return out


skills = sorted(ROOT.glob("skills/*/SKILL.md"))
ok(len(skills) == 9, f"expected 9 skills, found {len(skills)}")
for skill in skills:
    rel = skill.relative_to(ROOT)
    text = skill.read_text()
    if not text.startswith("---\n"):
        ok(False, f"{rel}: missing YAML frontmatter on line 1")
        continue
    try:
        fm = frontmatter(text.split("---", 2)[1])
    except ValueError as e:
        ok(False, f"{rel}: invalid frontmatter: {e}")
        continue
    name, desc = fm.get("name", ""), fm.get("description", "")
    ok(name == skill.parent.name, f"{rel}: name '{name}' != dir '{skill.parent.name}'")
    ok(bool(NAME_RE.match(name)) and len(name) <= 64, f"{rel}: name violates Agent Skills spec")
    ok(0 < len(desc) <= 1024, f"{rel}: description is {len(desc)} chars (must be 1-1024)")
    ok("<" not in desc and ">" not in desc, f"{rel}: description contains XML angle brackets")
    short = fm.get("metadata.short-description", "")
    ok(0 < len(short) <= 200,
       f"{rel}: metadata.short-description is {len(short)} chars (claude.ai upload needs 1-200)")
    extra = {k for k in fm if "." not in k} - UPLOAD_KEYS
    ok(not extra, f"{rel}: non-portable frontmatter keys {sorted(extra)} (rejected by claude.ai/API upload)")

try:
    plugin = json.loads((ROOT / ".claude-plugin/plugin.json").read_text())
    market = json.loads((ROOT / ".claude-plugin/marketplace.json").read_text())
    ok(isinstance(market.get("owner"), dict), "marketplace.json: owner must be an object")
    for entry in market.get("plugins", []):
        ok(isinstance(entry.get("author", {}), dict),
           f"marketplace.json: plugin '{entry.get('name')}' author must be an object, not a string")
        ok(entry.get("version") == plugin.get("version"),
           f"marketplace.json version {entry.get('version')} != plugin.json {plugin.get('version')}")
except Exception:
    pass  # JSON validity is already reported by check 1

# 3 — bundled tool compiles
tool = ROOT / "skills" / "meta-mvp" / "scripts" / "dict_color.py"
ok(tool.exists(), "skills/meta-mvp/scripts/dict_color.py missing")
if tool.exists():
    try:
        py_compile.compile(str(tool), doraise=True)
    except py_compile.PyCompileError as e:
        ok(False, f"skills/meta-mvp/scripts/dict_color.py does not compile: {e}")

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
