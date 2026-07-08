# Meta Prompting — Theory & References

Grounding for the `meta-prompting` skill. Read when you need the conceptual justification (category theory, type theory), the multi-modal schema, or the empirical results.

**Primary source:** Yifan Zhang, Yang Yuan, Andrew Chi-Chih Yao. *Meta Prompting for AI Systems.* arXiv:2311.11482 (2023); published in the ICLR 2024 BGPT Workshop. Reference implementation: github.com/meta-prompting/meta-prompting. Project page: meta-prompting.github.io. License: MIT.

## Table of contents
1. Definition and the MP vs. few-shot contrast
2. The formal framework: functor, monad, type theory
3. Multi-modal meta prompting (the typed `<task_schema>`)
4. Empirical results
5. Canonical prompt artifacts from the paper
6. Recursive Meta Prompting variants (the custom GPTs)

---

## 1. Definition and the MP vs. few-shot contrast

**Meta Prompting (MP)** is "a technique that provides a high-level, structural template for *how to think* rather than specific examples of *what to think*." Unlike few-shot prompting, which relies on content-rich examples, a meta-prompt is a **single, example-agnostic scaffold** that guides the reasoning process. It focuses on the formal procedure, syntax, and compositionality of problem-solving, and teaches the model a reusable, structured method for an entire category of tasks.

**Meta prompt — a structural template** (reusable blueprint, agnostic to the actual numbers/details):
```
Integrate step-by-step reasoning to solve mathematical problems under the following structure:
{
  "Problem": "[question to be answered]",
  "Solution": {
    "Step 1": "Begin the response with 'Let's think step by step.'",
    "Step 2": "Follow with the reasoning steps, ensuring the solution process is broken down clearly and logically.",
    "Step 3": "End the solution with the final answer encapsulated in a LaTeX-formatted box, $\\boxed{...}$"
  },
  "Final Answer": "[final answer]"
}
```

**Few-shot prompt — concrete examples** (the model must infer the pattern from solved pairs):
```
Problem: Find the domain of the expression $\frac{\sqrt{x-2}}{\sqrt{5-x}}$.
Solution: ... Therefore, the domain is $\boxed{[2,5)}$.
Final Answer: The final answer is $[2,5)$.
---------
Problem: If $\det A = 2$ and $\det B = 12$, find $\det(AB)$.
Solution: $\det(AB) = (\det A)(\det B) = 24$. ...
---------
...
```

The MP version carries no instance data; it is pure structure. That is the whole point — and why it transfers across instances and modalities, and costs far fewer tokens than stacking examples.

## 2. The formal framework

The paper's core novelty is a rigorous mathematical foundation for prompt engineering.

**The Meta Prompting Functor.** Model MP as a functor `M: Task → Prompt` mapping a category of tasks to a category of structured prompts. Functoriality means structure is preserved: if a complex task decomposes into simpler sub-tasks, the corresponding prompt is constructed by composing the corresponding sub-prompts. This is the theoretical guarantee of **modularity and compositionality** — and the practical reason the skill's procedure starts by finding the task's decomposition.

**The Recursive Monad (RMP).** Model an LLM refining its own prompts as a **monad over the category of prompts**, specifically a **Writer monad** that accumulates *edit scripts*. Let edit scripts form a monoid `(E, ·, ε)`.
- **Endofunctor `T`:** for a prompt object `P`, `T(P) = (P, E)` — the prompt paired with its edit history.
- **Unit `η`:** lifts a prompt into the refinement context with empty history `ε`.
- **Multiplication `μ`:** collapses nested levels of meta-reasoning, linearly composing recursive optimization steps into a single coherent edit trace.

This guarantees the refinement process is algebraically consistent and stable — i.e. iterated self-improvement composes instead of thrashing. Operationally: draft → critique → append edit script → flatten → repeat to a fixpoint.

**Type theory.** The framework is grounded in type theory: a meta-prompt enforces a clear, syntactic structure on inputs and outputs via **typed slots**, exactly as a type system ensures different data types are handled correctly. This is what makes the approach extend to multi-modal settings.

## 3. Multi-modal meta prompting

Typed slots create explicit "sockets" for heterogeneous data streams, letting the model synthesize across modalities. The paper's schema:

```xml
<task_schema>
    <input_slots>
        <data type="image/png" name="problem_diagram"></data>
        <data type="audio/mp3" name="verbal_instructions"></data>
        <data type="model/obj" name="object_model"></data>
    </input_slots>
    <output_schema>
        <synthesis type="text/markdown" name="analysis_summary"></synthesis>
        <result type="video/mp4" name="solution_animation"></result>
    </output_schema>
</task_schema>
```

Use this shape whenever a task mixes modalities or has several distinct inputs/outputs — give each its own typed, named slot.

## 4. Empirical results

A single, zero-shot, example-agnostic meta-prompt on a **base Qwen-72B** (no fine-tuning):

| Benchmark | Result with MP | Context |
|---|---|---|
| MATH | **46.3%** | beats GPT-4 (2023-03-14) CoT 42.5%, and Qwen-72B-MetaMathQA fine-tuned 41.7% |
| GSM8K | **83.5%** | beats MetaMath-70B 82.3% and WizardMath-70B 81.6% |
| Game of 24 | **100%** | by generating a single Python solver; ~$0.0003/sample vs. ToT 74% at $0.74/sample |

Qwen-14B (base) with MP also clears strong baselines (MATH 28.9%, GSM8K 64.8%). Headline takeaways: MP unlocks reasoning in *base* models without fine-tuning, with large **token-efficiency** gains over few-shot/CoT/ToT, and the Game-of-24 result shows the power of having the meta-prompt direct the model to emit *one program that solves every instance* rather than reasoning per-instance.

## 5. Canonical prompt artifacts

**MATH meta-prompt (used in the paper, plain-text variant):**
```
**Problem Statement**:
- **Problem**: [question to be answered]

**Solution Structure**:
1. Begin the response with "Let's think step by step."
2. Follow with the reasoning steps, ensuring the solution process is broken down clearly and logically.
3. End the solution with the final answer encapsulated in a LaTeX-formatted box, $\boxed{...}$, for clarity and emphasis.
4. Finally, state "The answer is [final answer to the problem].", with the final answer presented in LaTeX notation.
```

Both the JSON-skeleton and the plain-numbered variants are example-agnostic. Choose whichever renders the typed slots most legibly for the target model — JSON/XML when slots nest, numbered prose when the procedure is linear.

## 6. Recursive Meta Prompting variants

The reference repo ships three "CR Agent / MP" custom GPTs that operationalize RMP at different points of the loop:
- **CR Agent (Complex Reasoning)** — v0.1 and an XML v0.2: applies a fixed meta-prompt to reason through hard problems.
- **MP-PT (Prompt Revision)** — takes a prompt and revises it (one monad step, surfaced as a tool).
- **MP-ICPD (Recursive In-Context Prompt Design)** — the full recursive loop: design, critique, refine in context.

When the user wants a *self-improving* prompt, MP-ICPD is the pattern to emulate: keep the design, critique, and edit-script in the same context and iterate to a fixpoint.

## Citation
```
@article{zhang2023meta,
  title={Meta Prompting for AI Systems},
  author={Zhang, Yifan and Yuan, Yang and Yao, Andrew Chi-Chih},
  journal={arXiv preprint arXiv:2311.11482},
  year={2023},
  note={Published in ICLR 2024 BGPT Workshop}
}
```
