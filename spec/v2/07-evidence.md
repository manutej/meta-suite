# 07 · Evidence base

Each item notes its peer-review status and what it justifies in this spec. Gains the authors report
in their own papers are treated as hypotheses until the eval reproduces them.

| ID | Source | Status | Justifies |
|---|---|---|---|
| E1 | Zhang, Yuan, Yao — *Meta Prompting for AI Systems*, [arXiv:2311.11482](https://arxiv.org/abs/2311.11482) (latest version 2026; recursive meta prompting = proposer/validator/executor; authors note the CT structure does not guarantee correctness and that convergence is not guaranteed) | preprint | typed slots, edit logs, O1 roles; measure, don't assume |
| E2 | de Wynter et al. — *On Meta-Prompting*, [arXiv:2312.06562](https://arxiv.org/abs/2312.06562) | preprint | prompts that produce prompts, with declared types |
| E3 | Liu et al. — *From Intent to Evidence: a categorical approach to evaluating deep research agents*, [arXiv:2603.25342](https://arxiv.org/abs/2603.25342) | preprint | Yoneda probe, pullback verifier, colimit synthesizer; modest ablation gains (+1.8 to +2.8) |
| E4 | Bottman, Liu, Richardson — *Operadic consistency*, [arXiv:2606.13649](https://arxiv.org/abs/2606.13649) | preprint | G4: label-free, r = 0.86–0.94 with accuracy across 12 LLMs |
| E5 | Bottman, Richardson — *Operads for compositional reasoning in LLMs*, [arXiv:2606.13634](https://arxiv.org/abs/2606.13634) | theory | D1 decomposition tree |
| E6 | Kim — *HOLOGRAPH*, [arXiv:2512.24478](https://arxiv.org/abs/2512.24478) | preprint | G3 per-axiom reporting (locality fails separately from gluing) |
| E7 | He et al. — *Global consistency repair for NL claim sets (QXR)*, [arXiv:2601.13600](https://arxiv.org/abs/2601.13600) | preprint | G3 minimal conflicting sets |
| E8 | Olivieri, Hernández — *Sheaf-theoretic transport and obstruction*, [arXiv:2605.14033](https://arxiv.org/abs/2605.14033) | preprint | an obstruction score that breaks down into parts |
| E9 | Zhou et al. — *Ontology-to-tools compilation*, [arXiv:2602.03439](https://arxiv.org/abs/2602.03439) | preprint | D2/G1 typed writes (+27 F1 vs post-hoc checking) |
| E10 | Xiang et al. — *When to use graphs in RAG*, [arXiv:2506.05690](https://arxiv.org/abs/2506.05690) | ICLR 2026 | the lens helps synthesis and hurts simple lookup, so I1 gates which operators run |
| E11 | Li et al. — *DeepResearch Bench II*, [arXiv:2601.08536](https://arxiv.org/abs/2601.08536) | preprint | agents score about 90 on presentation and about 40–50 on recall/analysis; M1 rubric design |
| E12 | Han et al. — *DEER*, [arXiv:2512.17776](https://arxiv.org/abs/2512.17776) | ICML 2026 | expert-derived binary rubric items |
| E13 | Si, Yang, Hashimoto, [arXiv:2409.04109](https://arxiv.org/abs/2409.04109); *Ideation–Execution Gap*, [arXiv:2506.20803](https://arxiv.org/abs/2506.20803) | preprint | experts impressed at first glance, less so after execution; M3 SME panel |
| E14 | Agrawal et al. — *GEPA*, [arXiv:2507.19457](https://arxiv.org/abs/2507.19457) | ICLR 2026 | improvement as a measured search loop with held-out sets |
| E15 | Zhang et al. — *ACE: Agentic Context Engineering*, [arXiv:2510.04618](https://arxiv.org/abs/2510.04618) | ICLR 2026 | S2 delta edits; avoids brevity bias and context collapse |
| E16 | Wang et al. — *Agent Workflow Memory*, [arXiv:2409.07429](https://arxiv.org/abs/2409.07429) | ICML 2025 | S2/S3 reusable workflows and residue |
| E17 | Li et al. — *SkillsBench*, [arXiv:2602.12670](https://arxiv.org/abs/2602.12670) | NeurIPS 2026 | compact skills, 2–3 modules; self-generated skills −1.3 points without validation |
| E18 | Meincke, Mollick et al. — *Prompting Science Report 2*, [arXiv:2506.07142](https://arxiv.org/abs/2506.07142) | report | P7: no chain-of-thought boilerplate for reasoning models |
| E19 | Basil, Mollick et al. — *Prompting Science Report 4*, [arXiv:2512.05858](https://arxiv.org/abs/2512.05858) | report | P7: no accuracy personas |
| E20 | Kamoi et al. — self-correction survey, [arXiv:2406.01297](https://arxiv.org/abs/2406.01297) | TACL 2024 | R1/I3: verification needs an external signal |
| E21 | Wataoka et al., [arXiv:2410.21819](https://arxiv.org/abs/2410.21819); Roytburg et al., [arXiv:2601.22548](https://arxiv.org/abs/2601.22548) | preprint | self-preference bias exists, though its size is contested; cross-model judges |
| E22 | Anthropic — [Effective context engineering for AI agents](https://anthropic.com/engineering/effective-context-engineering-for-ai-agents) (2025) | guidance | right altitude, sub-agents, minimal high-signal tokens |
| E23 | Anthropic — [Equipping agents with Agent Skills](https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills) (2025) | guidance | progressive disclosure; descriptions drive triggering |
| E24 | Opsahl-Ong et al. — *MIPROv2*, [arXiv:2406.11695](https://arxiv.org/abs/2406.11695) | EMNLP 2024 | optimize instructions and examples jointly (v2.x) |
| E25 | Suzgun et al. — *Dynamic Cheatsheet*, [arXiv:2504.07952](https://arxiv.org/abs/2504.07952) | preprint | S3 persistent memory; risk that wrong entries persist, hence S2 elevation |

**Searched and not found:** the gluing construction attributed to Grant & Wang (CWRU), which
`meta-score` is planned to build on. It could not be verified (see Q3).
