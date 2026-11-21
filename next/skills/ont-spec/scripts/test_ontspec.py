"""Behaviour tests for ontspec: each law fails on the defect it names and passes otherwise.

Run: python3 -m pytest next/skills/ont-spec/scripts -q
"""

import copy
import json
from pathlib import Path

import pytest
import yaml

import ontspec as o

FIX = Path(__file__).parent / "fixtures" / "limiter.yaml"


def spec_from(raw: dict, tmp_path: Path) -> o.Spec:
    p = tmp_path / "s.yaml"
    p.write_text(yaml.safe_dump(raw, sort_keys=False))
    return o.load(p)


@pytest.fixture
def raw():
    return yaml.safe_load(FIX.read_text())


def atom(raw, aid):
    return next(a for a in raw["atoms"] if a["id"] == aid)


def laws(findings):
    return {f.law for f in findings}


def test_fixture_is_green_and_fully_atomic(tmp_path, raw):
    s = spec_from(raw, tmp_path)
    findings = o.check(s, converging=True)
    assert o.verdict(findings) == "GREEN", [f.line() for f in findings]
    assert o.stats(s)["atomicity"] == 1.0


def test_dangling_and_ill_typed_references_are_red(tmp_path, raw):
    atom(raw, "refill")["reads"].append("nope")
    atom(raw, "take-check")["inputs"].append("ob-cap")  # an obligation cannot be an input
    f = o.check(spec_from(raw, tmp_path))
    assert {"O2.dangling", "O2.ill-typed"} <= laws(f)
    assert o.verdict(f) == "RED"


def test_stated_atom_without_source_is_red(tmp_path, raw):
    del atom(raw, "capacity")["source"]
    assert "O3.unsourced" in laws(o.check(spec_from(raw, tmp_path)))


def test_open_question_must_name_what_it_blocks_and_who_settles_it(tmp_path, raw):
    raw["atoms"].append({"id": "q-burst", "kind": "Open", "depth": 1, "parent": "limiter", "status": "open",
                         "statement": "May a client burst above capacity?"})
    assert "O3.open-unowned" in laws(o.check(spec_from(raw, tmp_path)))
    raw["atoms"][-1].update(blocks=["capacity"], resolve_by="product owner")
    assert "O3.open-unowned" not in laws(o.check(spec_from(raw, tmp_path)))


def test_two_writers_of_one_state_are_flagged_unless_declared(tmp_path, raw):
    assert "O4.multi-writer" not in laws(o.check(spec_from(raw, tmp_path)))  # declared in the fixture
    del atom(raw, "bucket")["shared_writers"]
    f = o.check(spec_from(raw, tmp_path))
    hit = next(x for x in f if x.law == "O4.multi-writer")
    assert hit.atoms == ("refill", "take-spend") and hit.level == "AMBER"


def test_unverified_leaf_is_amber_then_red_when_converging(tmp_path, raw):
    raw["atoms"] = [a for a in raw["atoms"] if a["id"] != "ob-check"]
    s = spec_from(raw, tmp_path)
    assert ("O5.unverified", "AMBER") in {(f.law, f.level) for f in o.check(s)}
    assert ("O5.unverified", "RED") in {(f.law, f.level) for f in o.check(s, converging=True)}


def test_input_nothing_produces_is_red(tmp_path, raw):
    raw["atoms"].append({"id": "quota", "kind": "Entity", "depth": 1, "parent": "limiter", "status": "inferred",
                         "type": "u32", "statement": "A per-client quota."})
    atom(raw, "take-check")["inputs"].append("quota")
    assert "O6.unproduced" in laws(o.check(spec_from(raw, tmp_path)))


def test_one_value_per_key_across_depths(tmp_path, raw):
    raw["atoms"].append({"id": "capacity-deep", "kind": "Value", "depth": 3, "parent": "take", "status": "inferred",
                         "key": "capacity", "value": 12, "type": "u32 tokens", "statement": "Capacity seen deep down."})
    f = o.check(spec_from(raw, tmp_path))
    assert "O7.sheaf" in laws(f) and o.verdict(f) == "RED"


def test_non_atomic_leaf_is_found(tmp_path, raw):
    del atom(raw, "take-check")["post"]
    f = o.check(spec_from(raw, tmp_path), converging=True)
    assert ("O8.not-atomic", "RED") in {(x.law, x.level) for x in f}


def test_refined_operation_must_have_operation_children(tmp_path, raw):
    raw["atoms"] = [a for a in raw["atoms"] if a.get("parent") != "take" and a["id"] not in ("ob-check",)]
    raw["atoms"].append({"id": "take-note", "kind": "Value", "depth": 3, "parent": "take", "status": "inferred",
                         "key": "note", "value": 1, "type": "u8", "statement": "A note, not an operation."})
    assert "O9.unrealised" in laws(o.check(spec_from(raw, tmp_path)))


def test_reference_into_a_sibling_branch_is_a_leak(tmp_path, raw):
    # refill-rate is defined inside `admit`; a top-level invariant reaching into it leaks.
    atom(raw, "never-negative")["uses"].append("refill-rate")
    assert "O10.sibling-leak" in laws(o.check(spec_from(raw, tmp_path)))


def test_depth_must_follow_the_parent_and_stay_within_four(tmp_path, raw):
    atom(raw, "take-check")["depth"] = 2
    assert "O1.depth" in laws(o.check(spec_from(raw, tmp_path)))


def test_pure_operation_that_writes_state_is_red(tmp_path, raw):
    atom(raw, "refill")["effect"] = "pure"
    assert "O11.impure" in laws(o.check(spec_from(raw, tmp_path)))


def test_delta_detects_convergence(tmp_path, raw):
    a = spec_from(raw, tmp_path)
    b_raw = copy.deepcopy(raw)
    atom(b_raw, "take-check")["statement"] = "Reworded only."
    p = tmp_path / "b.yaml"; p.write_text(yaml.safe_dump(b_raw)); b = o.load(p)
    assert o.delta(a, b)["converged"]
    atom(b_raw, "take-check")["reads"] = []
    p.write_text(yaml.safe_dump(b_raw)); b = o.load(p)
    assert o.delta(a, b)["structurally_changed"] == ["take-check"]


def test_contracts_cover_every_leaf_operation_and_invariant(tmp_path, raw):
    cs = o.contracts(spec_from(raw, tmp_path))
    assert {c["id"] for c in cs} == {"refill", "take-check", "take-spend", "never-negative"}
    refill = next(c for c in cs if c["id"] == "refill")
    assert refill["ensures"] == ["bucket.tokens <= capacity"] and refill["verified_by"][0]["method"] == "proof"


def test_elements_report_effect_classes_and_owners(tmp_path, raw):
    md = o.elements(spec_from(raw, tmp_path))
    assert "| effectful |" in md and "| bucket | yes | refill, take-spend |" in md


# ── operadic projection, run through the real engine when it is installed ────

STITCH = o.find_stitch(None)
needs_stitch = pytest.mark.skipif(STITCH is None, reason="operadic-stitch engine not installed")


@needs_stitch
def test_operad_green_on_the_fixture(tmp_path, raw):
    snap, term = o.project(spec_from(raw, tmp_path))
    cert = o.run_stitch(STITCH, snap, term)
    assert cert["verdict"] in ("GREEN", "NOTE"), json.dumps(cert, indent=1)[:2000]
    assert cert["briefs"] == 2  # admit (depth 2) and take (depth 3)


@needs_stitch
def test_operad_catches_a_value_that_changes_across_a_depth_seam(tmp_path, raw):
    raw["atoms"].append({"id": "capacity-deep", "kind": "Value", "depth": 2, "parent": "admit", "status": "inferred",
                         "key": "capacity", "value": 12, "type": "u32 tokens", "statement": "Capacity, redefined."})
    snap, term = o.project(spec_from(raw, tmp_path))
    cert = o.run_stitch(STITCH, snap, term)
    assert cert["verdict"] == "RED"
    assert any(str(f).find("sheaf") >= 0 for f in cert["findings"])


@needs_stitch
def test_operad_catches_a_reference_into_a_sibling_branch(tmp_path, raw):
    # take-check (inside take) uses refill-rate, which is defined inside admit's brief and inherited,
    # so that is fine; but a top-level invariant using it is not in the root brief: port.open.
    atom(raw, "never-negative")["uses"].append("refill-rate")
    snap, term = o.project(spec_from(raw, tmp_path))
    cert = o.run_stitch(STITCH, snap, term)
    assert cert["verdict"] == "RED"
    assert any("port" in str(f) for f in cert["findings"])
