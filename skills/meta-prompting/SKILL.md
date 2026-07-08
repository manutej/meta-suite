---
name: meta-prompting
description: Generate a structured, example-agnostic META-PROMPT from any input task — a reusable "how to think" scaffold rather than worked examples. Use this whenever the user wants to design a prompt/system prompt for a CATEGORY of tasks, build a prompt template or reasoning harness, turn a one-off prompt into a reusable structure, create an "expert persona" system prompt with typed slots, or apply Meta Prompting / Recursive Meta Prompting (RMP). Triggers include "write a meta-prompt", "design a system prompt for X tasks", "make a reusable prompt template", "prompt that works for any category of problem", "structure the reasoning for this task", "self-improving prompt", "functor/typed-slot prompt", or any request to map an input task into a valid meta-prompt. This skill maps the meta-meta-prompt space: given a task, it emits the meta-prompt that solves the whole category.
---

> **Meta-Suite doctrine (overrides anything below).** Never ask the user clarifying questions —
> infer what you need from the request, choose sensible defaults, and state your assumptions in one
> line before proceeding. Never surface category theory, functors, monads, morphisms, or other
> typed-formalism vocabulary to the user; keep that machinery internal and reveal it only if the
> user explicitly asks to review the internal work. Keep all user-facing output in plain language.

# Meta Prompting

Based on **"Meta Prompting for AI Systems"** (Zhang, Yuan, Yao — arXiv:2311.11482, ICLR 2024 BGPT Workshop) and its `meta-prompting/meta-prompting` reference implementation.

## What this skill is for

A **meta-prompt** specifies the *formal structure* of how to solve an entire category of tasks — the syntax, the typed slots, the reasoning procedure, the output contract — **without** content-specific examples. Few-shot prompting shows the model *what* to think (concrete solved instances); meta-prompting shows it *how* to think (a reusable, example-agnostic blueprint). This is more token-efficient, more robust, and composes cleanly.

Your job here is **not** to solve the user's task. It is to operate one level up: take an input task and **emit the meta-prompt that solves the whole category that task belongs to**. You are mapping the *meta-meta-prompt space* — the space of procedures that generate valid meta-prompts.

Trigger this skill whenever the value is in a reusable structure rather than a one-shot answer: designing a system prompt, a prompt template, a reasoning harness, an agent persona, or any "make me a prompt for X-type problems."

## The one idea to hold onto

> **Meta Prompting is a functor.** It maps the *category of tasks* to the *category of prompts*, preserving structure: if a task decomposes into sub-tasks, the meta-prompt decomposes into composable sub-prompts. So your first move is always to find the task's structure, because the prompt's structure is the image of it. Get the decomposition right and the prompt writes itself.

A second idea makes it self-improving:

> **Recursive Meta Prompting (RMP) is a monad.** An LLM can generate and then refine its *own* meta-prompt. Model the refinement as a Writer monad: each pass appends an *edit script* to an accumulating history, and nested refinements collapse into one coherent edit trace. In practice: draft → critique → emit edit script → apply → repeat until stable.

Read `references/theory.md` for the functor/monad/type-theory grounding and the paper's benchmark results. Read it when the user wants the conceptual justification, asks about category theory, or you need the multi-modal typed-slot schema.

## The procedure: task → meta-prompt

Work through these five steps. They ARE the functor — each step lifts one piece of task structure into prompt structure.

### 1. Identify the task category (not the instance)
Look past the specific instance the user gave you to the **category** it represents. "Solve 3x+7=22" → the category is *single-variable algebraic equations*, or more usefully *step-by-step mathematical reasoning with a verifiable final answer*. The meta-prompt must serve the whole category, so name the category explicitly and state what varies across instances (the free variables) and what stays fixed (the invariant structure).

### 2. Decompose into typed sub-decisions
Break the category's solution into the sequence of sub-decisions any instance requires, and **assign each a type**. Treat every design decision as a typed term, never an ad-hoc value:
- inputs as typed slots — `Problem: str`, `Diagram: image/png`, `Constraints: list[str]`
- intermediate decisions as typed terms — `Plan: list[Step]`, `Subproblems: list[Problem]`
- the output as a typed contract — `FinalAnswer: LaTeX`, `Result: json{...}`

Typing is what guarantees compositionality and lets the structure transfer across modalities (text, image, audio, code). See the multi-modal `<task_schema>` pattern in `references/theory.md`.

### 3. Specify the reasoning procedure
State the *process* the model should follow over those slots — the "how to think," kept abstract and example-free. E.g. "Begin with 'Let's think step by step.' Decompose into sub-problems. Solve each. Compose. Verify against the original constraints." Procedure, not worked examples. If the category is better served by emitting a *program* than by reasoning in prose (as Game-of-24 is solved by generating one Python solver), say so — meta-prompts can direct the model to write code that handles every instance at once.

### 4. Bind the slots with rules and an output contract
Write the rules that connect the slots: ordering, what must reference what, invariants ("every claim cites a search result"), and the exact output format. A resolved slot is **immutable downstream** — later parts inherit it unchanged unless an explicit named refinement (a morphism) transforms it. Make the output contract unambiguous so the result is machine-checkable.

### 5. Render the meta-prompt in a clean structural format
Emit the meta-prompt as an explicit schema. **XML tags or JSON skeletons** work best because they make the typed slots and their nesting legible to the model. Keep it example-agnostic: the finished meta-prompt should contain *zero* solved instances — only slots, rules, and procedure. The two canonical shapes:

**JSON-skeleton style** (good for reasoning tasks with a fixed output shape):
```
Integrate step-by-step reasoning under the following structure:
{
  "Problem": "[question to be answered]",
  "Solution": {
    "Step 1": "Begin with 'Let's think step by step.'",
    "Step 2": "Break the reasoning down clearly and logically.",
    "Step 3": "End with the final answer in $\\boxed{...}$."
  },
  "Final Answer": "[final answer]"
}
```

**XML-schema style** (good for personas, multi-slot or multi-modal tasks): a `<system>` block (role + typed instructions) plus a `<syntax>` block defining the slot tree the model fills. The user's design-artifact meta-prompt in `references/meta-meta-prompt.md` is a full worked example of this shape.

## Recursive refinement (RMP) — when to loop

For high-stakes or fuzzy tasks, don't ship the first draft. Run the monad loop, typically **3–5 passes**:

1. **Draft** the meta-prompt via the five steps.
2. **Critique** it against the category: Does it cover every instance, not just the user's example? Any leaked specifics? Any untyped slot? Any unbound rule? Is the output contract checkable?
3. **Emit an edit script** — a concrete, named list of changes (this is the Writer-monad log; keep the running history so refinements compose rather than thrash).
4. **Apply** and repeat until the edit script comes back empty (a fixpoint) or you stop improving.

The synthesis example in `references/meta-meta-prompt.md` is explicitly designed to be looped this way — its render-and-verify structure is the monad made concrete for design artifacts.

## The generator template (the meta-meta-prompt)

`references/meta-meta-prompt.md` contains a ready-to-use **meta-meta-prompt**: a prompt you can hand any task to and get back a valid meta-prompt, plus the full design-artifact synthesis example that composes this skill with the `generative-ui` skill. Read it when the user wants the reusable generator itself, or wants the design/UI synthesis.

## Quality checklist before you ship a meta-prompt

- **Example-agnostic:** no solved instances anywhere — only slots, rules, procedure.
- **Typed:** every input, intermediate, and output is a named typed slot.
- **Compositional:** sub-tasks map to sub-prompts; structure mirrors the task.
- **Bound:** rules connect the slots; the output contract is unambiguous and checkable.
- **Category-complete:** it serves the whole category, verified against instances *other* than the one the user gave you.
- **Lean:** every line earns its place; prefer explaining *why* a constraint matters over piling on rigid MUSTs.

## Common failure modes

- **Overfitting to the example instance** — the single most common error. You write a "meta"-prompt that only really works for the one problem the user showed you. Always test it mentally against 2–3 *different* instances of the category.
- **Smuggling in few-shot examples** — if you find yourself writing a fully solved example, you've dropped back to few-shot. Convert it into a slot + rule instead.
- **Untyped, ad-hoc slots** — values like "put the answer here" with no type. Name the type.
- **Unbound slots** — slots with no rule saying how they relate or what order they resolve in.
- **Solving the task instead of structuring it** — if your output answers the user's specific question, you've collapsed back to level zero. Stay one level up.
