# Intake — resolve the four unknowns (default: derive them; opt-in: ask)

Every run needs four things settled before building: **audience**, **the single outcome**,
**aesthetic**, and **scope this run**. meta-suite doctrine is *never interrogate the user* —
so by **default you resolve these yourself** by studying the task's structure, and only ask if
the user explicitly opts in. Begin obvious research in parallel — don't gate the first web
search on resolving these.

## Default path — categorical resolution (no questions)

Treat each unknown as a **free variable** in the task's category. For each one:

1. **Enumerate the plausible values** (the fiber over that variable). e.g. audience ∈ {a
   stakeholder to win over · the team that runs it daily · the public · an investor}.
2. **Score each value against the evidence** you actually have — the prompt's wording, the
   GROUND-stage research, and anything from connected MCPs. Pick the maximum-likelihood value.
3. **Record a confidence** (high / medium / low) and the one signal that decided it.

The four variables and what each drives:

| Variable | Plausible values | What it drives |
|---|---|---|
| **Audience** | stakeholder to win · daily operators · public · investor | tone, depth, which surface leads |
| **Single outcome** | win approval · coordinate action · explain a concept · track status | the vision thesis + the meta-prompt's invariant |
| **Aesthetic** | clean/editorial · bold/energetic · sleek/premium · warm/human | palette mood + which architecture leads |
| **Scope this run** | full bundle · prototype · visuals only · strategy only | the deliverable set + the right-size tier |

## When a variable is genuinely ambiguous AND the work is high-stakes — panel review

Don't guess in the dark and don't fall back to asking. Run a small **multi-agent panel**
(2–3 independent sub-agents). Each agent independently resolves the still-ambiguous variables
from the prompt + research alone; a synthesizer reconciles:

- **Agreement** → adopt it at high confidence.
- **Divergence** → adopt the modal choice, and record the runner-up as a one-line stated
  alternative ("Built for a stakeholder pitch; if this is really a daily tool, the daily-tool
  surface leads instead"). Proceed — the split is surfaced, not blocking.

This is the categorical review the doctrine prefers over interrogation: the answer is worked
out by studying the space, with independent readings as the cross-check. Reserve the panel for
genuine high-stakes ambiguity — for an ordinary run the single-pass scoring above is enough
(don't be overkill).

## Opt-in path — actually ask the user

Only when the user **explicitly opts in** — they say "ask me", "interactive", "let me choose",
or pass `--ask` / `ask:true` in the invocation — use `AskUserQuestion` for the unresolved
high-variance variables only (2–4 questions, the recommended option first, "(Recommended)").
If the prompt already answers one, skip it; never ask more than 4.

> **Runtime guard.** `AskUserQuestion` is unavailable in some runtimes (headless / ACP / cron
> agents disable it). If an opt-in ask can't run, **fall back to the default categorical path
> automatically and note it in one line** — never hard-fail, never block the build.

## After resolution — always

Write the resolved four into `00_INTAKE.md` with each value's confidence and deciding signal,
restate in one line what you're about to build, and proceed without further gating. If the
prompt already fixes a variable, take it as given (confidence: high, signal: "stated").
