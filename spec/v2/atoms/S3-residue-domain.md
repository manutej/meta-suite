# S3 · residue.domain

**Purpose:** Save the domain graph a run built, so the next run in an overlapping domain starts warm
instead of from nothing.

**In:** the final `DomainGraph`, `relation_vocab`, findings, and timestamps
**Out:** `residue/<domain_key>/` in the user's workspace

**Contract**
- A later run loads matching residue as a **prior**. Claims older than the freshness window (set in
  config) are re-verified before use; they are never trusted as-is.
- A residue claim used as evidence must still count as a fresh, independent source. It does not
  bypass L2.
- Residue is never committed to the repo, and it is never shared across users or tenants (P13).

**Check:** G1 on load. Freshness is enforced by the kernel.

**Guards against:** re-paying for the same research, and stale facts presented as current.

**Evidence:** E16, E25.
