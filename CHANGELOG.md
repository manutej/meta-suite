# Changelog

## 1.2.0 — 2026-10-21 · cross-tool packaging

**Fixes**
- `/plugin install meta-suite@meta-suite` failed on current Claude Code (`author: expected object`).
  The marketplace entry's `author` is now an object.
- Seven of the nine `SKILL.md` files had frontmatter that strict YAML parsers reject (unquoted `: `).
  Claude Code tolerated it, but Codex, Cursor, claude.ai upload and `skills-ref validate` did not.
  All descriptions are now `>-` blocks.
- `meta-mvp`'s description was 1,270 characters, over the 1,024 Agent Skills limit; trimmed to 954.
- `meta-mvp`'s palette tool moved from `bin/` (only reachable inside a Claude Code plugin through
  `$CLAUDE_PLUGIN_ROOT`) to `skills/meta-mvp/scripts/dict_color.py`, so it travels with the skill
  into every tool.

**New**
- [`INSTALL.md`](INSTALL.md): exact steps for Claude Code, the Claude app, the Claude API, Codex,
  Cursor, Grok, ChatGPT, and raw APIs.
- `dist/`: per-skill upload zips (≤200-char descriptions for claude.ai), an all-skills zip,
  single-file prompts, and a ChatGPT/Grok Project router. Built by `scripts/package.py`.
- Every skill carries `metadata.short-description` for upload surfaces.
- `scripts/validate.py` now parses frontmatter strictly and checks spec limits, upload-safe keys,
  manifest types, and version parity. CI also checks that `dist/` is fresh and runs
  `claude plugin validate --strict`.

## 1.1.0 — 2026-07-08
- Nine task-agnostic meta-skills.
