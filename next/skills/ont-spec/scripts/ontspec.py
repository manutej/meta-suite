#!/usr/bin/env python3
"""ontspec — check, compose-verify and render an ontological specification.

The spec is ONE file (YAML or JSON) of typed atoms arranged in a refinement tree of depth <= 4.
Everything else (contracts, element tables, the approval document) is generated from it.

Commands
  check    SPEC                 ontology laws O1-O11 -> certificate (exit 0 GREEN, 1 AMBER, 2 RED, 3 invalid)
  operad   SPEC [--stitch PATH] project the tree onto the operadic-stitch engine and run it: typed
                                plugging, closed ports per refinement, one value per key across depth
                                seams, bounded recursion, budgets (same exit codes)
  delta    OLD NEW              structural change between two passes; says whether the loop converged
  contracts SPEC --out DIR      one typed contract per leaf Operation/Invariant (+ contracts.json)
  elements  SPEC                the architecture forces the spec implies (Markdown)
  render    SPEC                the human approval document (Markdown)

Findings use the format: FINDING @ where · law · atoms · the one question that resolves it.
Decision logic is pure; file and process I/O stay in main() and run_stitch().
"""

from __future__ import annotations

import argparse
import json
import os
import re
import subprocess
import sys
import tempfile
from collections import Counter, defaultdict
from dataclasses import dataclass, field
from pathlib import Path

KINDS = {
    "System", "Component", "Entity", "Value", "State", "Relation", "Event",
    "Operation", "Invariant", "Boundary", "Obligation", "Open",
}
DEFINES = {"Entity", "Value", "State", "Relation", "Event", "Boundary"}  # things others refer to
STATUSES = {"stated", "inferred", "open"}
EFFECTS = {"pure", "deterministic", "effectful", "llm", "human"}
METHODS = {"proof", "model-check", "property-test", "differential", "review"}
MAX_DEPTH = 4
DEPTH_BUDGETS = {2: 8000, 3: 4000, 4: 2000}  # tokens per refinement brief; must strictly decrease
ID_RE = re.compile(r"^[a-z0-9][a-z0-9-]{0,62}$")
KEY_REF = re.compile(r"\{([A-Za-z0-9._-]+)\}")
COMMON_FIELDS = {"id", "kind", "depth", "parent", "status", "source", "statement", "notes"}
KIND_FIELDS = {
    "System": set(), "Component": set(),
    "Entity": {"type", "durable", "external_input", "shared_writers"},
    "Value": {"key", "value", "type"},
    "State": {"type", "durable", "shared_writers", "external_input"},
    "Relation": {"type", "uses"}, "Event": {"type", "external_input"}, "Boundary": {"type"},
    "Operation": {"inputs", "outputs", "reads", "writes", "effect", "pre", "post", "errors"},
    "Invariant": {"uses", "formal"},
    "Obligation": {"target", "method", "evidence"},
    "Open": {"blocks", "resolve_by"},
}
EFFECT_RANK = {"pure": 0, "deterministic": 1, "effectful": 2, "llm": 2, "human": 2}
CONTRACT_FIELDS = ("effect", "pre", "post", "formal", "errors", "status", "method", "target",
                   "shared_writers", "durable", "blocks", "resolve_by", "evidence")
EXIT = {"GREEN": 0, "AMBER": 1, "RED": 2}


@dataclass(frozen=True)
class Finding:
    law: str
    level: str  # RED | AMBER | NOTE
    where: str
    atoms: tuple[str, ...]
    question: str

    def line(self) -> str:
        return f"{self.level} @ {self.where} · {self.law} · {', '.join(self.atoms) or '-'} · {self.question}"


@dataclass
class Spec:
    meta: dict
    atoms: dict[str, dict]
    order: list[str]
    children: dict[str, list[str]] = field(default_factory=lambda: defaultdict(list))

    def ancestors(self, aid: str) -> list[str]:
        out, cur = [], self.atoms[aid].get("parent")
        while cur:
            out.append(cur)
            cur = self.atoms[cur].get("parent")
        return out

    def subtree(self, aid: str) -> list[str]:
        out, stack = [], [aid]
        while stack:
            cur = stack.pop()
            out.append(cur)
            stack.extend(reversed(self.children.get(cur, [])))
        return out

    def leaves(self) -> list[str]:
        return [a for a in self.order if not self.children.get(a)]


# ── loading (boundary) ────────────────────────────────────────────────────────


class SpecError(ValueError):
    """The file is not a spec at all; distinct from a spec that breaks a law."""


def load(path: Path) -> Spec:
    text = path.read_text()
    if path.suffix in (".yaml", ".yml"):
        import yaml  # only needed for YAML specs

        try:
            raw = yaml.safe_load(text)
        except yaml.YAMLError as err:
            raise SpecError(f"{path}: not valid YAML: {err}") from err
    else:
        raw = json.loads(text)
    if not isinstance(raw, dict) or not isinstance(raw.get("atoms"), list):
        raise SpecError(f"{path}: expected a mapping with an 'atoms' list")
    atoms, order = {}, []
    for i, a in enumerate(raw["atoms"]):
        if not isinstance(a, dict) or not isinstance(a.get("id"), str):
            raise SpecError(f"{path}: atom #{i} has no string id")
        if a["id"] in atoms:
            raise SpecError(f"{path}: duplicate atom id {a['id']!r}")
        atoms[a["id"]] = a
        order.append(a["id"])
    spec = Spec(meta=raw.get("spec", {}), atoms=atoms, order=order)
    for aid in order:
        parent = atoms[aid].get("parent")
        if parent in atoms:
            spec.children[parent].append(aid)
    return spec


def refs(a: dict) -> dict[str, list[str]]:
    """Every outgoing reference of an atom, by field."""
    return {f: list(a.get(f) or []) for f in ("inputs", "outputs", "reads", "writes", "uses", "target", "blocks")}


# ── ontology laws (pure) ──────────────────────────────────────────────────────


def check(spec: Spec, converging: bool = False, max_depth: int | None = None) -> list[Finding]:
    """O1-O11. With converging=True, leaf atomicity and verification coverage become RED."""
    F: list[Finding] = []
    A = spec.atoms
    gap = "RED" if converging else "AMBER"

    # O1 shape: ids, kinds, statuses, tree, depth
    roots = [a for a in spec.order if not A[a].get("parent")]
    if len(roots) != 1 or A[roots[0]].get("kind") != "System":
        F.append(Finding("O1.root", "RED", "spec", tuple(roots), "Which single System atom is the root?"))
    for aid in spec.order:
        a = A[aid]
        if not ID_RE.fullmatch(aid):
            F.append(Finding("O1.id", "RED", aid, (aid,), "Rename to lower-case kebab-case, at most 63 characters?"))
        if a.get("kind") not in KINDS:
            F.append(Finding("O1.kind", "RED", aid, (aid,), f"Which kind is it, one of {sorted(KINDS)}?"))
        if a.get("status") not in STATUSES:
            F.append(Finding("O1.status", "RED", aid, (aid,), "Is it stated by a source, inferred by us, or open?"))
        if not str(a.get("statement") or "").strip():
            F.append(Finding("O1.statement", "RED", aid, (aid,), "What does this atom say, in one plain sentence?"))
        parent = a.get("parent")
        if parent and parent not in A:
            F.append(Finding("O1.parent", "RED", aid, (aid, parent), "Which existing atom does this refine?"))
        depth = a.get("depth")
        expected = 0 if not parent else (A[parent].get("depth", -9) + 1 if parent in A else None)
        if expected is not None and depth != expected:
            F.append(Finding("O1.depth", "RED", aid, (aid,), f"Depth should be {expected} (parent depth + 1); is the parent right?"))
        if isinstance(depth, int) and depth > spec.meta.get("max_depth", MAX_DEPTH):
            F.append(Finding("O1.max-depth", "RED", aid, (aid,), "Can this be stated at depth 4 or less, or is it a separate spec?"))

    if max_depth is not None:
        deep = [a for a in spec.order if isinstance(A[a].get("depth"), int) and A[a]["depth"] > max_depth]
        if deep:
            F.append(Finding("O1.pass-depth", "RED", "spec", tuple(deep[:8]),
                             f"This pass is limited to depth {max_depth}: move these to the next pass?"))

    # cycles in the parent relation
    for aid in spec.order:
        seen, cur = set(), aid
        while cur in A and A[cur].get("parent"):
            if cur in seen:
                F.append(Finding("O1.cycle", "RED", aid, tuple(sorted(seen)), "Which of these is really the parent?"))
                break
            seen.add(cur)
            cur = A[cur]["parent"]

    # O2 typed references: every reference exists and has a kind that can play that role
    allowed = {
        "inputs": DEFINES, "outputs": DEFINES, "reads": {"Entity", "State", "Value", "Boundary"},
        "writes": {"Entity", "State"}, "uses": KINDS - {"Obligation", "Open"},
        "target": KINDS - {"Obligation", "Open"}, "blocks": KINDS,
    }
    for aid in spec.order:
        for fld, ids in refs(A[aid]).items():
            for r in ids:
                if r not in A:
                    F.append(Finding("O2.dangling", "RED", aid, (aid, r), f"Define {r!r}, or is it a typo?"))
                elif A[r].get("kind") not in allowed[fld]:
                    F.append(Finding("O2.ill-typed", "RED", aid, (aid, r),
                                     f"{A[r].get('kind')} {r!r} cannot be used as {fld}; which atom was meant?"))

    # O3 provenance: stated atoms cite a source; open atoms say what they block and who resolves them
    for aid in spec.order:
        a = A[aid]
        if a.get("status") == "stated" and not a.get("source"):
            F.append(Finding("O3.unsourced", "RED", aid, (aid,), "Where does the source say this, or is it inferred?"))
        if a.get("kind") == "Open" and not (a.get("blocks") and a.get("resolve_by")):
            F.append(Finding("O3.open-unowned", "RED", aid, (aid,), "What does this unknown block, and who or what settles it?"))

    # O4 ownership: one writer per piece of state (shared writers must be declared with a reason)
    writers = defaultdict(list)
    for aid in spec.leaves():
        if A[aid].get("kind") == "Operation":
            for w in A[aid].get("writes") or []:
                writers[w].append(aid)
    for target, ws in writers.items():
        declared = (A.get(target) or {}).get("shared_writers")
        if len(ws) > 1 and not declared:
            F.append(Finding("O4.multi-writer", "AMBER", target, tuple(ws),
                             "Which operation owns this state, or why may several write it (declare shared_writers)?"))

    # O5 verification coverage: every leaf Operation and every Invariant has an obligation
    targeted = {t for aid in spec.order if A[aid].get("kind") == "Obligation" for t in A[aid].get("target") or []}
    for aid in spec.order:
        k = A[aid].get("kind")
        if (k == "Invariant" or (k == "Operation" and not spec.children.get(aid))) and aid not in targeted:
            F.append(Finding("O5.unverified", gap, aid, (aid,), "Which obligation (proof, model check, property test) covers this?"))
    for aid in spec.order:
        a = A[aid]
        if a.get("kind") == "Obligation" and a.get("method") not in METHODS:
            F.append(Finding("O5.method", "RED", aid, (aid,), f"How is it discharged, one of {sorted(METHODS)}?"))

    # O6 production: every operation input is produced somewhere or enters at a boundary/parameter
    produced = {o for aid in spec.order for o in (A[aid].get("outputs") or [])}
    sources = {aid for aid in spec.order if A[aid].get("kind") in ("Value", "Boundary") or A[aid].get("external_input")}
    for aid in spec.leaves():
        if A[aid].get("kind") != "Operation":
            continue
        for i in A[aid].get("inputs") or []:
            if i in A and i not in produced and i not in sources and not (A[i].get("kind") in ("Entity", "State") and A[i].get("durable")):
                F.append(Finding("O6.unproduced", "RED", aid, (aid, i),
                                 f"Nothing produces {i!r}: which operation outputs it, or is it an external input?"))
    used = {r for aid in spec.order for ids in refs(A[aid]).values() for r in ids}
    for aid in spec.order:
        if A[aid].get("kind") in DEFINES and aid not in used and not spec.children.get(aid):
            F.append(Finding("O6.unused", "AMBER", aid, (aid,), "Nothing uses this definition: delete it, or what should use it?"))

    # O7 one value per key across the whole tree (the depth seams included)
    by_key = defaultdict(list)
    for aid in spec.order:
        a = A[aid]
        if a.get("kind") == "Value" and "key" in a:
            by_key[a["key"]].append(aid)
    for key, ids in by_key.items():
        values = {json.dumps(A[i].get("value"), sort_keys=True) for i in ids if "value" in A[i]}
        if len(values) > 1:
            F.append(Finding("O7.sheaf", "RED", key, tuple(ids), f"{key} has values {sorted(values)}: which one holds?"))

    # O8 atomicity of leaves
    for aid in spec.leaves():
        problem = atomicity_gap(A[aid])
        if problem and A[aid].get("kind") not in ("System", "Open"):
            F.append(Finding("O8.not-atomic", gap, aid, (aid,), problem))

    # O9 refinement completeness: a refined Operation is realised by child operations
    for aid in spec.order:
        a = A[aid]
        kids = spec.children.get(aid, [])
        if a.get("kind") == "Operation" and kids and not any(A[k].get("kind") == "Operation" for k in kids):
            F.append(Finding("O9.unrealised", "RED", aid, (aid,), "Which child operations carry this operation out?"))

    # O10 locality: a reference must stay within the atom's own subtree or its ancestors' direct children
    for aid in spec.order:
        visible = visible_defs(spec, aid)
        for fld, ids in refs(A[aid]).items():
            if fld in ("target", "blocks"):
                continue  # obligations and unknowns may point anywhere
            for r in ids:
                if r in A and A[r].get("kind") in DEFINES and r not in visible:
                    F.append(Finding("O10.sibling-leak", "AMBER", aid, (aid, r),
                                     f"{r!r} is private to another branch: lift it to a common ancestor?"))
    # O11 effect honesty: a pure operation writes nothing, and a refined operation's declared
    # effect is at least as strong as the strongest effect among its child operations
    for aid in spec.order:
        a = A[aid]
        if a.get("kind") != "Operation":
            continue
        if a.get("effect") == "pure" and a.get("writes"):
            F.append(Finding("O11.impure", "RED", aid, (aid,),
                             "It is marked pure but writes state: return the new state as an output, or change its effect?"))
        kids = [k for k in spec.subtree(aid)[1:] if A[k].get("kind") == "Operation" and A[k].get("effect") in EFFECT_RANK]
        if a.get("effect") in EFFECT_RANK and kids:
            strongest = max(kids, key=lambda k: EFFECT_RANK[A[k]["effect"]])
            if EFFECT_RANK[A[strongest]["effect"]] > EFFECT_RANK[a["effect"]]:
                F.append(Finding("O11.understated", "RED", aid, (aid, strongest),
                                 f"It is declared {a['effect']} but contains a {A[strongest]['effect']} step: raise its effect?"))

    # O12 closed schema: a field the schema does not know is probably a typo or an unchecked idea
    for aid in spec.order:
        a = A[aid]
        extra = set(a) - COMMON_FIELDS - KIND_FIELDS.get(a.get("kind"), set())
        if extra:
            F.append(Finding("O12.unknown-field", "AMBER", aid, (aid,),
                             f"Unknown fields {sorted(extra)}: rename to a schema field, or move into notes?"))

    # O13 constants by reference: {key} in text must name a Value key, so text cannot drift from values
    keys = {A[v].get("key") for v in spec.order if A[v].get("kind") == "Value"}
    for aid in spec.order:
        a = A[aid]
        texts = [a.get("formal") or ""] + list(a.get("pre") or []) + list(a.get("post") or [])
        for k in sorted({m for t in texts for m in KEY_REF.findall(str(t))}):
            if k not in keys:
                F.append(Finding("O13.unknown-key", "RED", aid, (aid,), f"{{{k}}} names no Value key: which parameter is meant?"))

    # O5b coverage strength: a pure or deterministic leaf covered only by review could be checked mechanically
    by_target = defaultdict(set)
    for o in spec.order:
        if A[o].get("kind") == "Obligation":
            for t in A[o].get("target") or []:
                by_target[t].add(A[o].get("method"))
    for aid in spec.leaves():
        a = A[aid]
        if a.get("kind") == "Operation" and a.get("effect") in ("pure", "deterministic") and by_target.get(aid) == {"review"}:
            F.append(Finding("O5.weak", "AMBER", aid, (aid,),
                             "Only reviewed, yet it is pure or deterministic: which proof or test could check it mechanically?"))
    return F


def visible_defs(spec: Spec, aid: str) -> set[str]:
    """Definitions an atom may use: in its own subtree, or declared at any ancestor level
    (children of an ancestor or of the root). This is lexical scoping for specs."""
    vis = set(spec.subtree(aid))
    for anc in spec.ancestors(aid) + [aid]:
        vis.update(k for k in spec.children.get(anc, []) if spec.atoms[k].get("kind") in DEFINES)
        vis.add(anc)
    return vis


def atomicity_gap(a: dict) -> str | None:
    """None if the leaf is atomic; otherwise the question that makes it atomic."""
    k = a.get("kind")
    if k == "Operation":
        if not a.get("inputs") or not a.get("outputs"):
            return "What are its typed inputs and outputs?"
        if a.get("effect") not in EFFECTS:
            return f"Which effect class is it, one of {sorted(EFFECTS)}?"
        if not a.get("pre") or not a.get("post"):
            return "What must hold before it runs, and what does it guarantee after?"
        if len(a.get("writes") or []) > 1:
            return "It writes several things: split it so each leaf writes one?"
    elif k in ("Entity", "State", "Event", "Relation"):
        if not a.get("type"):
            return "What is its type (fields, enum order, or relation signature)?"
    elif k == "Value":
        if "value" not in a and a.get("status") != "open":
            return "What is its value, or mark it open?"
        if not a.get("type"):
            return "What is its type and unit or range?"
    elif k == "Invariant":
        if not a.get("uses"):
            return "Which atoms does it constrain?"
        if not a.get("formal"):
            return "How is it written as a checkable predicate?"
    elif k == "Obligation":
        if not a.get("target") or a.get("method") not in METHODS:
            return "What does it verify, and how?"
    elif k == "Boundary":
        if not a.get("type"):
            return "What crosses this boundary (its message type)?"
    elif k == "Component":
        return "A component with no parts: what does it contain?"
    return None


def verdict(findings: list[Finding]) -> str:
    levels = {f.level for f in findings}
    return "RED" if "RED" in levels else "AMBER" if "AMBER" in levels else "GREEN"


def stats(spec: Spec) -> dict:
    A = spec.atoms
    leaves = [l for l in spec.leaves() if A[l].get("kind") not in ("System", "Open")]
    atomic = [l for l in leaves if atomicity_gap(A[l]) is None]
    obligations = [a for a in spec.order if A[a].get("kind") == "Obligation"]
    return {
        "atoms": len(spec.order),
        "by_kind": dict(sorted(Counter(A[a].get("kind") for a in spec.order).items())),
        "by_depth": dict(sorted(Counter(A[a].get("depth") for a in spec.order).items())),
        "leaves": len(leaves),
        "atomic_leaves": len(atomic),
        "atomicity": round(len(atomic) / len(leaves), 3) if leaves else 0.0,
        "open": sum(1 for a in spec.order if A[a].get("kind") == "Open"),
        "obligations_by_method": dict(Counter(A[o].get("method") for o in obligations)),
    }


# ── operadic projection (pure) + engine run (I/O) ─────────────────────────────


def _safe(text: str) -> str:
    """stitch fact values may not hold newlines or angle brackets."""
    return re.sub(r"\s+", " ", str(text)).replace("<", " lt ").replace(">", " gt ").strip()


def project(spec: Spec) -> tuple[dict, str]:
    """Map the refinement tree onto a stitch snapshot and one term.

    Definitions become Fact atoms (`def.<id>` = kind:type, and Values also carry their key).
    Operations, Invariants and Obligations become Method, Constraint and Gate atoms that
    require the facts of everything they reference. Each refined atom becomes a nested brief
    under the frame for its children's depth; a brief inherits its ancestors' definitions.
    So: a reference to a sibling branch's private definition opens a port (RED), a key with two
    values across a depth seam is sheaf.seam (RED), and a brief over its depth budget is AMBER.
    """
    A = spec.atoms
    kind_map = {"Operation": "Method", "Invariant": "Constraint", "Obligation": "Gate",
                "Component": "Goal", "System": "Goal", "Open": "Evidence"}
    atoms = [
        {"id": "role.refiner", "kind": "Role", "scope": ["global"],
         "body": "You refine one part of a specification and nothing else."},
        {"id": "contract.refine", "kind": "Contract", "scope": ["global"],
         "body": "The children below must jointly realise the parent goal, using only definitions in this brief."},
    ]
    for d, budget in DEPTH_BUDGETS.items():
        atoms.append({"id": f"frame.depth{d}", "kind": "Frame", "scope": ["global"],
                      "requires": ["Contract", "Role"],
                      "facts": {"max_depth": 3, "budget": budget, "returns": "Evidence"},
                      "body": f"Refinement brief at depth {d}."})
    for aid in spec.order:
        a = A[aid]
        k = a.get("kind")
        sid = f"a.{aid}"
        body = _safe(a.get("statement", ""))
        if k in DEFINES and spec.children.get(aid):
            # a refined definition is still a definition (Fact) but also the Goal of its own brief
            atoms.append({"id": f"g.{aid}", "kind": "Goal", "scope": ["global"], "requires": ["Contract"],
                          "body": _safe(f"Refine the definition: {a.get('statement', aid)}")})
        if k in DEFINES:
            facts = {f"def.{aid}": _safe(f"{k}:{a.get('type', '?')}")}
            if k == "Value" and "key" in a and "value" in a:
                facts[_safe_key(a["key"])] = _safe(json.dumps(a["value"]))
            atoms.append({"id": sid, "kind": "Fact", "scope": ["global"], "facts": facts, "body": body})
        else:
            need = sorted({f"key:def.{r}" for fld, ids in refs(a).items() if fld not in ("target", "blocks")
                           for r in ids if r in A and A[r].get("kind") in DEFINES})
            # A refined atom is the Goal of its own brief; a leaf keeps its section colour.
            refined = bool(spec.children.get(aid))
            atom = {"id": sid, "kind": "Goal" if refined else kind_map.get(k, "Constraint"),
                    "scope": ["global"], "body": body or aid}
            if refined:
                atom["requires"] = ["Contract"]
            if need:
                atom["requires"] = atom.get("requires", []) + need
            atoms.append(atom)
    snapshot = {"_comment": f"ontspec projection of {spec.meta.get('id', 'spec')}",
                "scopes": {"global": None}, "overlaps": [], "confidential": [], "invariants": [],
                "atoms": atoms, "operations": []}
    root = next(a for a in spec.order if not A[a].get("parent"))
    return snapshot, _term(spec, root, inherited=[])


def _safe_key(key: str) -> str:
    return re.sub(r"[^A-Za-z0-9._-]", "_", key)


def _term(spec: Spec, aid: str, inherited: list[str]) -> str:
    A = spec.atoms
    kids = spec.children.get(aid, [])
    defs_here = [k for k in kids if A[k].get("kind") in DEFINES]
    context = inherited + defs_here
    own = [f"@g.{aid}", f"@a.{aid}"] if A[aid].get("kind") in DEFINES else [f"@a.{aid}"]
    parts = ["@role.refiner", "@contract.refine"] + own + [f"@a.{d}" for d in dict.fromkeys(context) if d != aid]
    for k in kids:
        if A[k].get("kind") in DEFINES and not spec.children.get(k):
            continue  # already in context
        if spec.children.get(k):
            depth = A[k].get("depth", 1) + 1
            parts.append(f"(nest @frame.depth{min(depth, max(DEPTH_BUDGETS))} {_term(spec, k, context)})")
        else:
            parts.append(f"@a.{k}")
    return "(glue " + " ".join(parts) + ")"


def find_stitch(explicit: str | None) -> Path | None:
    candidates = [explicit, os.environ.get("ONT_SPEC_STITCH")]
    for base in (Path.home() / ".claude", Path("/root/.claude")):
        if base.exists():
            candidates += [str(p) for p in base.glob("**/operadic-stitch/scripts/stitch.py")]
    for c in candidates:
        if c and Path(c).is_file():
            return Path(c)
    return None


def run_stitch(stitch: Path, snapshot: dict, term: str) -> dict:
    with tempfile.TemporaryDirectory() as tmp:
        snap = Path(tmp) / "snapshot.json"
        snap.write_text(json.dumps(snapshot))
        lib = subprocess.run([sys.executable, str(stitch), "check", str(snap)], capture_output=True, text=True)
        out = subprocess.run([sys.executable, str(stitch), "compose", str(snap), term, "--json"],
                             capture_output=True, text=True)
    if out.returncode >= 3:
        return {"verdict": "INVALID", "exit": out.returncode, "error": (out.stderr or out.stdout)[-2000:],
                "library_check": lib.returncode}
    cert = json.loads(out.stdout)["certificate"]
    return {"verdict": cert.get("verdict"), "exit": out.returncode, "library_check": lib.returncode,
            "findings": cert.get("findings", []), "tokens": cert.get("tokens"),
            "briefs": len(cert.get("briefs") or {}), "term_hash": cert.get("term_hash")}


# ── convergence (pure) ────────────────────────────────────────────────────────


def delta(old: Spec, new: Spec) -> dict:
    o, n = set(old.order), set(new.order)
    structural_fields = ("kind", "parent", "inputs", "outputs", "reads", "writes", "uses", "type", "value", "key")
    changed = sorted(a for a in o & n
                     if any(old.atoms[a].get(f) != new.atoms[a].get(f) for f in structural_fields))
    # Contract-level fields change what an implementation must satisfy, so they also block convergence.
    contract_changed = sorted(a for a in o & n if a not in changed
                              and any(old.atoms[a].get(f) != new.atoms[a].get(f) for f in CONTRACT_FIELDS))
    added, removed = sorted(n - o), sorted(o - n)
    return {"added": added, "removed": removed, "structurally_changed": changed,
            "contract_changed": contract_changed,
            "converged": not (added or removed or changed or contract_changed)}


# ── generated views (pure) ────────────────────────────────────────────────────


def contracts(spec: Spec) -> list[dict]:
    """One typed contract per leaf Operation and per Invariant; regenerable into any language."""
    A = spec.atoms
    out = []
    obligations, blockers = defaultdict(list), defaultdict(list)
    for aid in spec.order:
        if A[aid].get("kind") == "Obligation":
            for t in A[aid].get("target") or []:
                obligations[t].append({"id": aid, "method": A[aid].get("method"), "statement": A[aid].get("statement"),
                                       "evidence": A[aid].get("evidence")})
        if A[aid].get("kind") == "Open":
            for b in A[aid].get("blocks") or []:
                blockers[b].append(aid)

    def blocked_by(cid: str) -> list[str]:
        """Unknowns on the contract itself, on an enclosing Operation, or on a definition it uses.
        An unknown about a whole System or Component is a project risk, not a contract gap."""
        a = A[cid]
        related = [cid] + [x for x in spec.ancestors(cid) if A[x].get("kind") == "Operation"]
        related += [r for f in ("inputs", "outputs", "reads", "writes", "uses") for r in a.get(f) or []]
        return sorted({u for r in related for u in blockers.get(r, [])})
    for aid in spec.order:
        a = A[aid]
        if a.get("kind") == "Operation" and not spec.children.get(aid):
            out.append({
                "id": aid, "kind": "operation", "statement": a.get("statement"), "effect": a.get("effect"),
                "signature": {"inputs": [{"id": i, "type": A.get(i, {}).get("type")} for i in a.get("inputs") or []],
                              "outputs": [{"id": o, "type": A.get(o, {}).get("type")} for o in a.get("outputs") or []]},
                "reads": a.get("reads") or [], "writes": a.get("writes") or [],
                "requires": a.get("pre") or [], "ensures": a.get("post") or [], "errors": a.get("errors") or [],
                "source": a.get("source"), "status": a.get("status"), "verified_by": obligations.get(aid, []),
                "blocked_by": blocked_by(aid),
            })
        elif a.get("kind") == "Invariant":
            out.append({"id": aid, "kind": "invariant", "statement": a.get("statement"), "formal": a.get("formal"),
                        "constrains": a.get("uses") or [], "source": a.get("source"), "status": a.get("status"),
                        "verified_by": obligations.get(aid, []), "blocked_by": blocked_by(aid)})
    return out


def elements(spec: Spec) -> str:
    """The architecture forces implied by the spec. Forces, not decisions."""
    A = spec.atoms
    ops = [a for a in spec.leaves() if A[a].get("kind") == "Operation"]
    rows = defaultdict(list)
    for o in ops:
        rows[A[o].get("effect", "?")].append(o)
    lines = ["# Elements: what any architecture must provide", "",
             "Derived from the converged ontology. These are forces an architecture must satisfy; "
             "choosing the architecture is a separate investigation.", ""]
    lines += ["## Operations by effect class", "", "| effect | why it matters | operations |", "|---|---|---|"]
    why = {"pure": "can run anywhere, be proved, replayed, parallelised",
           "deterministic": "must replay identically (workflow-safe)",
           "effectful": "touches the outside world; needs idempotency and retries",
           "llm": "non-deterministic model call; needs budgets, recording, evaluation",
           "human": "waits on a person; needs durable waiting and audit"}
    for eff in sorted(rows):
        lines.append(f"| {eff} | {why.get(eff, '?')} | {', '.join(rows[eff])} |")
    lines += ["", "## State and its owners", "", "| state | durable | written by | read by |", "|---|---|---|---|"]
    for s in spec.order:
        if A[s].get("kind") == "State" or (A[s].get("kind") == "Entity" and A[s].get("durable")):
            w = [o for o in ops if s in (A[o].get("writes") or [])]
            r = [o for o in ops if s in (A[o].get("reads") or []) or s in (A[o].get("inputs") or [])]
            lines.append(f"| {s} | {'yes' if A[s].get('durable') else 'no'} | {', '.join(w) or '-'} | {', '.join(r) or '-'} |")
    shared = [(s, A[s]["shared_writers"]) for s in spec.order if A[s].get("shared_writers")]
    if shared:
        lines += ["", "## Shared writers and why", "", "| state | reason |", "|---|---|"]
        lines += [f"| {s} | {r} |" for s, r in shared]
    lines += ["", "## Boundaries (what crosses into or out of the system)", "", "| boundary | carries | statement |", "|---|---|---|"]
    for b in spec.order:
        if A[b].get("kind") == "Boundary":
            lines.append(f"| {b} | {A[b].get('type', '?')} | {A[b].get('statement')} |")
    lines += ["", "## Verification obligations by method", "", "| method | count | obligations |", "|---|---|---|"]
    by_m = defaultdict(list)
    for o in spec.order:
        if A[o].get("kind") == "Obligation":
            by_m[A[o].get("method")].append(o)
    for m in sorted(by_m):
        lines.append(f"| {m} | {len(by_m[m])} | {', '.join(by_m[m])} |")
    lines += ["", "## Open questions that block elements", "", "| open | blocks | settled by |", "|---|---|---|"]
    for u in spec.order:
        if A[u].get("kind") == "Open":
            lines.append(f"| {u} | {', '.join(A[u].get('blocks') or [])} | {A[u].get('resolve_by', '?')} |")
    return "\n".join(lines) + "\n"


def render(spec: Spec, findings: list[Finding]) -> str:
    A = spec.atoms
    st = stats(spec)
    m = spec.meta
    lines = [f"# {m.get('title', m.get('id', 'Specification'))}", "",
             f"Version {m.get('version', '?')} · verdict **{verdict(findings)}** · {st['atoms']} atoms · "
             f"atomicity {st['atomicity']:.0%} · {st['open']} open questions", "",
             "Generated from the ontology file. Approve or change the ontology, never this document.", ""]
    if m.get("verbal"):
        lines += ["## The specification in one paragraph", "", str(m["verbal"]).strip(), ""]
    lines += ["## Refinement tree", ""]
    root = next(a for a in spec.order if not A[a].get("parent"))
    for aid in spec.subtree(root):
        a = A[aid]
        pad = "  " * a.get("depth", 0)
        tag = f"`{a.get('kind')}` " + ("" if a.get("status") == "stated" else f"*{a.get('status')}* ")
        lines.append(f"{pad}- **{aid}** {tag}— {a.get('statement')}")
    opens = [u for u in spec.order if A[u].get("kind") == "Open"]
    if opens:
        lines += ["", "## Open questions", "", "| open | blocks | settled by | question |", "|---|---|---|---|"]
        lines += [f"| {u} | {', '.join(A[u].get('blocks') or [])} | {A[u].get('resolve_by', '?')} | {A[u].get('statement')} |" for u in opens]
    obs = [x for x in spec.order if A[x].get("kind") == "Obligation"]
    if obs:
        lines += ["", "## Verification obligations", "", "| obligation | method | verifies | evidence |", "|---|---|---|---|"]
        lines += [f"| {x} | {A[x].get('method')} | {', '.join(A[x].get('target') or [])} | {A[x].get('evidence') or '-'} |" for x in obs]
    lines += ["", "## Findings", ""]
    lines += [f"- {f.line()}" for f in findings] or ["- none"]
    return "\n".join(lines) + "\n"


# ── CLI (I/O edge) ────────────────────────────────────────────────────────────


def main(argv: list[str] | None = None) -> int:
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    sub = ap.add_subparsers(dest="cmd", required=True)
    c = sub.add_parser("check"); c.add_argument("spec"); c.add_argument("--converging", action="store_true",
                                                                     help="treat non-atomic or unverified leaves as RED")
    c.add_argument("--max-depth", type=int, help="this pass may not go deeper than N")
    v = sub.add_parser("converged", help="the stopping rule: converging GREEN, operad GREEN, no change since PREV")
    v.add_argument("prev"); v.add_argument("spec"); v.add_argument("--stitch")
    o = sub.add_parser("operad"); o.add_argument("spec"); o.add_argument("--stitch"); o.add_argument("--emit", help="write snapshot + term to this directory")
    d = sub.add_parser("delta"); d.add_argument("old"); d.add_argument("new")
    k = sub.add_parser("contracts"); k.add_argument("spec"); k.add_argument("--out", required=True)
    sub.add_parser("elements").add_argument("spec")
    sub.add_parser("render").add_argument("spec")
    args = ap.parse_args(argv)
    try:
        if args.cmd == "delta":
            print(json.dumps(delta(load(Path(args.old)), load(Path(args.new))), indent=2))
            return 0
        if args.cmd == "converged":
            prev, cur = load(Path(args.prev)), load(Path(args.spec))
            conv = verdict(check(cur, converging=True))
            stitch = find_stitch(args.stitch)
            op = run_stitch(stitch, *project(cur))["verdict"] if stitch else "UNAVAILABLE"
            d = delta(prev, cur)
            ok = conv == "GREEN" and op in ("GREEN", "NOTE") and d["converged"]
            print(json.dumps({"converged": ok, "check_converging": conv, "operad": op, "delta": d}, indent=2))
            return 0 if ok else 1
        spec = load(Path(args.spec))
    except (SpecError, OSError, json.JSONDecodeError) as err:
        print(f"invalid: {err}", file=sys.stderr)
        return 3

    if args.cmd == "check":
        findings = check(spec, converging=args.converging, max_depth=args.max_depth)
        v = verdict(findings)
        print(json.dumps({"verdict": v, "stats": stats(spec), "findings": [f.line() for f in findings]}, indent=2))
        return EXIT[v]
    if args.cmd == "operad":
        snapshot, term = project(spec)
        if args.emit:
            out = Path(args.emit); out.mkdir(parents=True, exist_ok=True)
            (out / "snapshot.json").write_text(json.dumps(snapshot, indent=1))
            (out / "term.txt").write_text(term + "\n")
        stitch = find_stitch(args.stitch)
        if stitch is None:
            print("operad check unavailable: operadic-stitch engine not found (pass --stitch or set ONT_SPEC_STITCH)",
                  file=sys.stderr)
            return 3
        cert = run_stitch(stitch, snapshot, term)
        print(json.dumps(cert, indent=2))
        return EXIT.get(cert["verdict"], 3)
    if args.cmd == "contracts":
        out = Path(args.out); out.mkdir(parents=True, exist_ok=True)
        cs = contracts(spec)
        (out / "contracts.json").write_text(json.dumps(cs, indent=2) + "\n")
        for cdef in cs:
            (out / f"{cdef['id']}.md").write_text(render_contract(cdef))
        print(f"wrote {len(cs)} contracts to {out}")
        return 0
    if args.cmd == "elements":
        print(elements(spec), end="")
        return 0
    if args.cmd == "render":
        print(render(spec, check(spec)), end="")
        return 0
    return 3


def render_contract(c: dict) -> str:
    lines = [f"# Contract: {c['id']}", "", c.get("statement") or "", "",
             f"- kind: {c['kind']} · status: {c.get('status')} · source: {c.get('source') or 'inferred'}"]
    if c["kind"] == "operation":
        lines += [f"- effect: {c.get('effect')}", "", "## Signature", ""]
        lines += [f"- in  `{i['id']}`: {i.get('type')}" for i in c["signature"]["inputs"]]
        lines += [f"- out `{o['id']}`: {o.get('type')}" for o in c["signature"]["outputs"]]
        lines += ["", f"reads {', '.join(c['reads']) or '-'} · writes {', '.join(c['writes']) or '-'}", "",
                  "## Requires", ""] + [f"- {p}" for p in c["requires"]]
        lines += ["", "## Ensures", ""] + [f"- {p}" for p in c["ensures"]]
        if c["errors"]:
            lines += ["", "## Errors", ""] + [f"- {e}" for e in c["errors"]]
    else:
        lines += ["", "## Predicate", "", f"`{c.get('formal')}`", "", f"constrains: {', '.join(c['constrains'])}"]
    lines += ["", "## Verified by", ""] + ([f"- {v['id']} ({v['method']}): {v['statement']}" + (f" Evidence: {v['evidence']}" if v.get("evidence") else "")
                                            for v in c["verified_by"]] or ["- nothing yet"])
    lines += ["", "## Blocked by", ""] + ([f"- {b}" for b in c.get("blocked_by") or []] or ["- nothing: final unless the ontology changes"])
    return "\n".join(lines) + "\n"


if __name__ == "__main__":
    sys.exit(main())
