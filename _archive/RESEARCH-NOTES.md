# Research notes (synthesized, internal)

Five parallel research agents informed this build:
1. **CRAFT (manutej/craft):** not an acronym; the "craft of senior engineering." 13-rule constitution
   around ETC; rule #8 is literally "Clear beats clever." Three-tier skills (SKILL.md reflex →
   references → examples). Earned palette: ink/cream + gold #E2B65A + teal #43C9BC, **no SaaS purple/gradient**;
   Hanken Grotesk + JetBrains Mono. README: shields badge row, why→how→architecture(mermaid)→skills
   (details accordions)→rules table→install→sources. Single-source-of-truth (RULES.md) → generated ports.
2. **compound-engineering-plugin (EveryInc):** Plan→Work→Review→Compound→Repeat loop; 11 pillars
   (compounding axiom; invest in the system 50/50; taste in systems; teach once; safety nets not
   gatekeeping; agent-native parity; parallelize; plans are the new code; capture learnings;
   discoverability; refresh against rot). Artifacts live in the user's repo.
3. **compound-knowledge-plugin (EveryInc):** knowledge-work port of compounding. Capture-so-it-accrues;
   retrieval closes the loop; plain greppable markdown + YAML atoms; anti-rot stale-checker; quality
   over quantity; honesty about uncertainty; read/write separation; search-before-generate;
   auto-classify don't ask; pipeline mode.
4. **Plugin/README best practice:** only `.claude-plugin/plugin.json` required; components at root;
   marketplace.json (name, owner, plugins[]); description = routing rule with trigger phrases; lean
   SKILL.md + progressive disclosure; `claude plugin validate`; install via /plugin marketplace add + /plugin install.
5. **Citation:** primary = Zhang, Yuan & Yao, "Meta Prompting for AI Systems" (arXiv:2311.11482) —
   example-agnostic meta-prompts via type+category theory. Lineage: Suzgun & Kalai, "Meta-Prompting:
   Enhancing LMs with Task-Agnostic Scaffolding" (arXiv:2401.12954) for the "task-agnostic scaffolding"
   phrase; de Wynter et al., "On Meta-Prompting" (arXiv:2312.06562) for the categorical "prompting-to-
   obtain-prompts" analysis.
