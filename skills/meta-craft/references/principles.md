# Craft principles — full statements (internal reference)

A domain-agnostic generalization of senior building practice. Of two ways to do something, prefer the one that leaves the result easier to change (this is the unifying frame, sometimes called ETC — Easier To Change). The opposing force is the *tactical rush* — local speed bought with global fragility.

1. **Smallest change that satisfies the intent.** Do what was asked, well. Resist scope you
   invented. Smaller changes are easier to verify, review, and reverse.
2. **No abstraction until repetition is real.** Premature generality is a cost, not a saving.
   Wait for two or three concrete instances before factoring out a shared form.
3. **Search and reuse before building new.** The thing you need often already exists in the
   project, the ecosystem, or prior work. Look first.
4. **Consolidate, don't clone.** Duplicated logic, prose, or assets drift apart and rot. Prefer
   one canonical version with variation parameterized.
5. **Match the existing conventions, not your defaults.** Consistency with the surrounding work
   is worth more than your personal preference. Fit in.
6. **Validate at the boundaries.** Check and normalize inputs where they enter the system. Fail
   loudly and early. Never swallow an error to make a symptom disappear.
7. **Keep the core clean; quarantine mess at the edges.** Keep decision logic pure and testable;
   push side effects, I/O, external calls, and unavoidable hacks to the perimeter where they're
   contained and visible.
8. **Clear beats clever.** The load-bearing rule. Code, prose, and design are read far more often
   than written. Optimize for the next person's comprehension, not for elegance or brevity that
   only you can parse.
9. **Names carry intent; explanations carry the *why*.** A good name states purpose. A good
   comment/annotation states the reason a non-obvious choice was made — not a restatement of what
   the reader can already see.
10. **Prove it works.** Evidence over assertion. Show the run, the test, the rendered result, the
    checked figure. "It should work" is not done.
11. **Don't destabilize what works.** Before changing working material, capture its current state
    (a snapshot, a pinned version, a baseline) so you can compare and revert.
12. **Verify what you depend on.** Confirm that the API, library, source, statistic, or factual
    claim you're relying on actually exists and says what you think. Do not build on an assumption.
13. **What you make outlives the moment.** Durable artifacts get revisited, extended, and migrated.
    Design for the next edit: expand safely, migrate deliberately, contract only once nothing
    depends on the old shape.

The skill selects the subset that *bites* for the work at hand and instantiates each as a concrete,
checkable rule plus a Definition-of-Done item. The full list is the menu; the constitution is the
order.
