# Meta-Plan — building & finalizing meta-suite
*(Produced with /meta-planning. Internal. Not shipped to users.)*

## Intake (typed slots)
- **task_to_plan_for:** Ship `meta-suite`, a 9-skill, standalone Claude plugin of task-agnostic
  meta-* scaffolds, with a CRAFT-grade README + landing page + diagrams, clean repo, prunable archive.
- **execution_environment:** Cowork + parallel subagents (research + verification fleet, ≥10),
  web research, file tools, packaging via zip.
- **constraints:** Skills must be (1) standalone — no references to other plugins; (2) never ask the
  user questions; (3) never surface category theory unless the user opts into internal work;
  (4) clear > clever. Repo must stay clean; internal logic quarantined in `_archive/`.

## Level check
The TASK is itself meta (build the meta-meta space of "scaffolds that solve categories"). Treated
as a category: objects = the lifecycle phases of any project; the plugin must cover the phase-space.

## Domain map — the meta-meta space (why these 9, and what would be missing)
A project's life cycle, each phase answered by one scaffold:

| Phase | Question it answers | Skill |
|---|---|---|
| Think | How do I think about this whole category? | meta-prompting |
| Plan | How will agents execute this category? | meta-planning |
| Design | What is the rigorous design schema? | meta-design |
| Build | One-liner → shippable bundle? | meta-mvp |
| Compound | How does each unit make the next easier? | meta-compound |
| Craft | What is the senior standard + done-ness? | meta-craft |
| Review | How is it verified, adversarially? | meta-review |
| Research | How is truth gathered and verified? | meta-research |
| Remember | How does context persist & compound? | meta-context |

Coverage argument: think→plan→design→build is the *forward* axis (make the thing); compound→craft→
review→research→context is the *quality & accrual* axis (make the thing good and make it last). The
two axes together close the loop — output of build feeds review feeds compound feeds context feeds
the next think. No essential phase is left uncovered; nothing is duplicated. Candidate extras
considered and deferred as non-essential / coverable by composition: meta-debug (⊂ meta-review),
meta-ship/meta-deploy (domain-specific, not task-agnostic), meta-decide (⊂ meta-research + meta-review).

## Invariants (from constraints) → enforced in every new skill
- I1 Standalone: no skill body names another plugin/skill-suite or external product as a dependency.
- I2 No questions: each skill carries an explicit "Never ask the user questions" operating rule.
- I3 No category theory surfaced: typed/categorical model lives only in `references/internal-model.md`,
  gated behind "do NOT surface unless the user asks to review internal work."
- I4 Clear > clever: plain-language output contracts; short tailored lists over exhaustive dumps.

## Work packets (parallelizable)
- WP-R research fleet ×5 (CRAFT, compound-eng, compound-knowledge, plugin best-practice, citation) — DONE.
- WP-A author 5 skills (compound, craft, review, research, context) — sequential for house-style coherence.
- WP-P presentation: README, landing page, excalidraw diagrams.
- WP-V verification fleet ×≥5 (one per new skill + README/manifest + landing + integration).
- WP-C cleanup: quarantine internals to `_archive/`; package `.plugin`.

## SOLUTION gate (done when)
All 9 skills validate (frontmatter + standalone + no-questions + no-cat-theory-surfaced), README +
landing + diagrams reference the real 9 skills, archive prunable, `.plugin` packages clean, ≥10
agents used across research+finalization.
