# Installing meta-suite — every tool, one page

Meta-suite is nine [Agent Skills](https://agentskills.io) (`skills/<name>/SKILL.md`). The same
skills ship in three shapes, so pick the row for your tool and follow only that section.

| You use… | Shape | Section |
|---|---|---|
| **Claude Code** (CLI, desktop, web, IDE) | plugin — installs from this repo | [Claude Code](#claude-code-recommended) |
| **Claude app** (claude.ai, Claude Desktop) | one zip per skill | [Claude app](#claude-app-claudeai--claude-desktop) |
| **Claude API** | one zip per skill | [Claude API](#claude-api) |
| **OpenAI Codex**, **Cursor**, **Grok Build**, GitHub Copilot, Gemini CLI, OpenCode, Windsurf, … | skill folders | [Coding agents](#coding-agents-codex-cursor-grok-build-copilot-gemini-cli-) |
| **ChatGPT** | zip per skill, or instructions + one knowledge file | [ChatGPT](#chatgpt) |
| **Grok** (grok.com / apps) | zip per skill, or instructions + one knowledge file | [Grok](#grok-grokcom-and-apps) |
| **OpenAI API**, **xAI API**, any other model | single-file prompt | [Any other model](#any-other-model-or-api) |

Every non-plugin file lives in [`dist/`](dist/) and is rebuilt by `python scripts/package.py`
(CI fails if it goes stale), so what you download always matches `skills/`.

---

## Claude Code (recommended)

Inside a Claude Code session:

```
/plugin marketplace add manutej/meta-suite
/plugin install meta-suite@meta-suite
```

Or from a shell, non-interactively:

```bash
claude plugin marketplace add manutej/meta-suite
claude plugin install meta-suite@meta-suite
```

Start a new session (or run `/reload-plugins`) and the nine skills are live. They auto-trigger from
what you ask (*"review this plan"*, *"make this project compound"*), or invoke one directly — plugin
skills are namespaced by the plugin name:

```
/meta-suite:meta-review   /meta-suite:meta-planning   /meta-suite:meta-mvp   …
```

**Check it worked:** `claude plugin details meta-suite@meta-suite` lists `Skills (9)`.

**Update later:** `claude plugin marketplace update meta-suite && claude plugin update meta-suite@meta-suite`,
then start a new session. Run `claude plugin list` to see the version you have (current: 1.2.0).

**Try without installing** (from a clone): `claude --plugin-dir ./meta-suite`.

**Prefer bare `/meta-review` names?** Install the skills as personal skills instead of a plugin
(same as [coding agents](#coding-agents-codex-cursor-grok-build-copilot-gemini-cli-) below, with
`-a claude-code`). Don't do both, or every skill appears twice.

## Claude app (claude.ai / Claude Desktop)

1. Make sure **code execution** is enabled in Settings (skills need it).
2. Download the zip for each skill you want from [`dist/upload/`](dist/upload/)
   — e.g. [`meta-review.zip`](dist/upload/meta-review.zip). Don't unzip it.
3. **Customize → Skills → Upload skill**, pick the zip, and toggle the skill on. Repeat per skill.

Each zip holds one folder named after the skill with `SKILL.md` inside, and uses a short
(≤200-character) description, as the Claude app requires. Uploaded skills are per-user and don't sync
to Claude Code — install the plugin there separately.

## Claude API

Upload each zip from [`dist/upload/`](dist/upload/) to `POST /v1/skills`, then attach it to a request
with `container.skills` (code execution tool required). See the
[Skills guide](https://platform.claude.com/docs/en/build-with-claude/skills-guide).

## Coding agents (Codex, Cursor, Grok Build, Copilot, Gemini CLI, …)

One command installs all nine skills for the agents you name, using the open
[`skills` CLI](https://github.com/vercel-labs/skills) (Node 22+):

```bash
npx skills add manutej/meta-suite --skill '*' -g -a codex -a cursor -a grok
```

Drop `-g` to install into the current project instead of your home directory. Other agent ids
include `claude-code`, `github-copilot`, `gemini-cli`, `opencode`, and `windsurf`; run
`npx skills add manutej/meta-suite --list` to preview.

**Manual alternative:** unzip [`dist/meta-suite-skills.zip`](dist/meta-suite-skills.zip) (nine
folders) into the agent's skills directory:

| Agent | Personal (all projects) | Per project | Invoke |
|---|---|---|---|
| OpenAI Codex | `~/.agents/skills/` | `.agents/skills/` | `$meta-review`, or `/skills` |
| Cursor | `~/.cursor/skills/` or `~/.agents/skills/` | `.cursor/skills/` or `.agents/skills/` | `/meta-review` |
| Grok Build | `~/.grok/skills/` or `~/.agents/skills/` | `.grok/skills/` | `/meta-review` |
| Claude Code (non-plugin) | `~/.claude/skills/` | `.claude/skills/` | `/meta-review` |

All of them also auto-trigger a skill when your request matches its description. Restart the agent
after installing.

## ChatGPT

**If your plan has Skills:** upload each zip from [`dist/upload/`](dist/upload/) under
**Skills → Create → Upload**, one skill per zip.

**Otherwise (any plan), as a Project or Custom GPT:**

1. Paste [`dist/instructions.md`](dist/instructions.md) into the Project's (or GPT's) instructions —
   it fits the 8,000-character limit.
2. Upload [`dist/meta-suite-prompt.md`](dist/meta-suite-prompt.md) as a file / knowledge file. It
   holds all nine skills with every reference inlined.
3. Ask normally (*"meta-review this deck"*). The instructions route each request to the right
   skill in the file.

## Grok (grok.com and apps)

**With Grok Skills:** upload a zip from [`dist/upload/`](dist/upload/) per skill.
**Otherwise:** create a Project, paste [`dist/instructions.md`](dist/instructions.md) into its
instructions, and add [`dist/meta-suite-prompt.md`](dist/meta-suite-prompt.md) as a file.

## Any other model or API

For an API with no skills feature (xAI API, OpenAI Chat/Responses without the shell tool, a local
model), use [`dist/prompts/<skill>.md`](dist/prompts/) as the **system prompt**. Each file is one
complete skill with its references inlined, so one skill per conversation costs about 2–22k tokens rather
than the ~60k for the whole suite. The OpenAI Responses API *does* support skills natively
(shell tool + `POST /v1/skills`); upload the zips from `dist/upload/` there.

> **What changes outside Claude Code.** Every skill works with nothing but a model. A few steps are
> better with tools, such as meta-mvp's palette check (`scripts/dict_color.py`, plain Python) or
> fanning out sub-agents in meta-planning and meta-review. Where a tool is missing, the skill does the
> step by hand or says in one line that it was skipped. It never stalls.

---

## Troubleshooting

| Symptom | Fix |
|---|---|
| `/meta-review` is "unknown" after plugin install | Plugin skills are namespaced: use `/meta-suite:meta-review`, or just describe the task. |
| `/plugin install` errors with `author: expected object` | You're on a pre-1.2.0 copy of the repo. Run `claude plugin marketplace update meta-suite`. |
| Skills don't appear after install | Start a new session (or `/reload-plugins`); in other agents, restart the app. |
| A skill appears twice | It's installed both as a plugin and as personal skills; remove one. |
| Old behaviour after updating | Run `claude plugin list` to check the version, then `claude plugin update meta-suite@meta-suite` and start a new session. |
| claude.ai rejects the upload | Upload the zip from `dist/upload/` as-is, without unzipping and rezipping. The folder must be the zip's root. |

## For maintainers

```bash
python scripts/validate.py          # manifests, strict frontmatter, links, SVGs, bundled tool
python scripts/package.py           # rebuild dist/ after editing anything in skills/
claude plugin validate . --strict   # Claude Code's own manifest + skill check
```

Bump `version` in **both** `.claude-plugin/plugin.json` and `.claude-plugin/marketplace.json`
whenever skills change. Installed users stay on the version they have until it changes.
