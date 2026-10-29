# J1 · llm.sample

**Purpose:** Give every semantic judgment and every sampled answer one owner, one recording format
and one aggregation rule. This is what makes "kernel-aggregated" verdicts auditable.

**Internal construct:** —

**In:** a task spec from a requesting atom. The task is either `judge` (pick one value from a
declared verdict enum) or `answer` (produce a typed short answer, for example a G4 node answer).
**Out:** `Judgment[]` and an aggregate verdict or answer set

**Contract**
- **Batched (Opus-native, D-059).** `meta sample run` takes n tasks × k samples in one call. One subagent executes the batch in a **programmatic tool-calling loop** and records the results in bulk, so there is no chat round-trip per judgment. Stable prompt prefixes come first, so caching applies.
- **Issue.** `meta sample issue` writes k task files. Each records the prompt hash, the role (from
  the run's model map) and the producer model. k defaults to the tier value (config:
  `tiers.<tier>.k`) and is always odd.
- **Execute.** The host model, or a subagent, executes each task in a context that has no access to
  the material being judged beyond what the task contains. It records the result with
  `meta sample record`.
- **Judge must differ from producer.** The judge model must differ from the producer model where the
  surface allows it. When it cannot, the judgment is recorded with `blind=false` and
  `same_model=true`. That is allowed, but it is visible.
- **Aggregate.**
  - `judge` tasks: majority vote, **with the vote counts kept** (G3's layer 3 weights edges by vote margin). Ties are impossible because k is odd.
  - `answer` tasks: the set of answers, passed to the requesting atom.
  - The rule is set in config: `sample.aggregate`.
- Every judgment is a ledger event: model, prompt_hash, temperature and output. A verdict can be
  reproduced given the judgment files.

**Check:** G1 on every task and judgment. `aggregate` refuses to run if fewer than k judgments have
been recorded.

**Guards against:** six different hidden "the LLM decides" steps, each with its own unwritten rule.

**Evidence:** E21 (judge bias, hence recording producer vs judge).
