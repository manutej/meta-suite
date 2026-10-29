# P1 · prompt.generate

**Purpose:** Turn a task into an example-agnostic meta-prompt that serves its whole category.

**Internal construct:** a structure-preserving map from the category of tasks to the category of prompts. A task's decomposition becomes the prompt's decomposition.

**In:** `Spec` (from I1–I3) describing the task instance and its category
**Out:** `MetaPrompt`: category, free variables, invariants, typed slots, procedure, bindings, output contract (JSON canonical, XML-rendered)

**Contract**
- Name the category, its free variables and its invariants. The category is never just the instance.
- Slots are typed, rules bind the slots, and the output contract is checkable.
- **Zero solved instances.** A leaked example is a defect. Whether one leaked is a J1 `judge` task.
- Where one program would solve every instance, the procedure says so (procedure-as-program).
- **Modes:** deep by default; `--fast` loads only the core op (D-066).

**Check:** G1 on `MetaPrompt`; J1 example-leak check; then P2.

**Guards against:** few-shot prompts in disguise, and meta-prompts that only fit the example given.

**Evidence:** E1, E2
