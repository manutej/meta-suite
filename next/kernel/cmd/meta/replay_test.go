package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/manutej/meta-suite/next/kernel/internal/claims"
	"github.com/manutej/meta-suite/next/kernel/internal/config"
	"github.com/manutej/meta-suite/next/kernel/internal/findings"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
	runs "github.com/manutej/meta-suite/next/kernel/internal/rundir"
)

// TestReplayLiveEUAIAct replays the first real run (D-079): testdata/live/eu-ai-act, LIVE data recorded from
// meta-lens on Opus 5.5 (see its README.txt). It rebuilds the claim shards and token files from the stage
// tokens stamped on each recorded claim, re-runs merge → graph build → interface merge --plain → findings
// verify → G2 → report, and asserts the recorded counts. Expectations differ from the recording only where a
// slice-1.1 fix (D-087) changes them on purpose; each such place says which fix and why.
func TestReplayLiveEUAIAct(t *testing.T) {
	const live = "../../testdata/live/eu-ai-act"
	dir := filepath.Join(t.TempDir(), "runs", "01M651YMAFM8EPYHJ0AB13C2WQ")
	for _, f := range []string{"run.json", "spec.json", "cover.json", "ledger.jsonl", "claims.jsonl"} {
		b, err := os.ReadFile(filepath.Join(live, f))
		if err != nil {
			t.Fatal(err)
		}
		os.MkdirAll(dir, 0o755)
		os.WriteFile(filepath.Join(dir, f), b, 0o644)
	}

	// --- the recording itself: what the run reported ---
	recFindings := readJSONL(t, filepath.Join(live, "findings.jsonl"))
	recKinds := map[string]int{}
	for _, f := range recFindings {
		recKinds[f["kind"].(string)]++
	}
	var recConf struct {
		Conflicts []map[string]any `json:"conflicts"`
	}
	readJSON(t, filepath.Join(live, "conflicts.json"), &recConf)
	if len(recFindings) != 103 || recKinds["verified"] != 73 || recKinds["single_source"] != 30 || len(recConf.Conflicts) != 3 {
		t.Fatalf("recording changed: %d findings %v, %d conflicts", len(recFindings), recKinds, len(recConf.Conflicts))
	}

	// --- rebuild shards and token files from the recorded claims (byte-identical lines) ---
	raw, _ := os.ReadFile(filepath.Join(live, "claims.jsonl"))
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	byAgent := map[string][]string{}
	var all []model.Claim
	for _, l := range lines {
		var c model.Claim
		if err := json.Unmarshal([]byte(l), &c); err != nil {
			t.Fatal(err)
		}
		all = append(all, c)
		byAgent[c.StageToken.Agent] = append(byAgent[c.StageToken.Agent], l)
		b, _ := json.Marshal(c.StageToken)
		os.MkdirAll(filepath.Join(dir, "tokens"), 0o755)
		os.WriteFile(filepath.Join(dir, "tokens", c.StageToken.TID+".json"), b, 0o644)
	}
	os.MkdirAll(filepath.Join(dir, "claims"), 0o755)
	for a, ls := range byAgent {
		os.WriteFile(filepath.Join(dir, "claims", a+".jsonl"), []byte(strings.Join(ls, "\n")+"\n"), 0o644)
	}
	must(t, "", "ledger show", dir) // the recorded hash chain is intact

	// K4: every recorded unit is accepted as the run's own vocabulary declares it ("%" canonical, "EUR" in the
	// currency dimension), with no false "not in the table" note and no conversion.
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	run, err := runs.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	vocab, err := claims.LoadVocab(run)
	if err != nil {
		t.Fatal(err)
	}
	units := map[string]int{}
	for _, c := range all {
		if c.Unit == "" {
			continue
		}
		cc := c
		known, rej := claims.Normalize(&cc, cfg, vocab)
		if !known || rej != nil || cc.Normalized != nil {
			t.Errorf("claim %s unit %q: known=%v rejection=%v normalized=%v", c.ID, c.Unit, known, rej, cc.Normalized)
		}
		units[c.Unit]++
	}
	if units["%"] != 12 || units["EUR"] != 9 {
		t.Fatalf("units %v", units)
	}

	// D-097 (N1) on live data: the run's cross-month "no fines" claims (wsgrdataadvisor.com, 2026-08;
	// reedsmith.com, 2026-10) stay two findings even when an undated agreeing claim from a third publisher is
	// added: it supports each month without joining them, and no widened "2026" finding appears.
	undated := model.Claim{}
	for _, c := range all {
		if c.ID == "01M653AP2F2XQYBGWJB53TVMY6" {
			undated = c
		}
	}
	undated.ID, undated.AsOf = "01M659ZZZZZZZZZZZZZZZZZZZZ", ""
	undated.Source = model.Source{Locator: "https://third.example/fines", RetrievedAt: "2026-10-30T00:00:00Z", IndependenceKey: "third.example"}
	fs, _ := findings.Verify(append(append([]model.Claim(nil), all...), undated), findings.Params{MinIndependent: cfg.Lens.Pullback.MinIndependent,
		NumericRelTol: cfg.Lens.Pullback.NumericRelTol, StripPrefixes: cfg.Independence.StripPrefixes, AsOfWindow: cfg.Lens.Pullback.AsOfWindow})
	fines := map[string]model.Finding{}
	for _, f := range fs {
		if f.Key.Subject == "ai act fines imposed" {
			fines[f.AsOf] = f
		}
	}
	if len(fines) != 3 || fines["2026-08"].Kind != "verified" || fines["2026-10"].Kind != "verified" || fines[""].Kind != "single_source" ||
		len(fines["2026-08"].Evidence) != 1 || len(fines["2026-10"].Evidence) != 1 {
		t.Fatalf("no-fines split with an undated claim: %+v", fines)
	}

	// K9: the spec hash lands in run.json.
	must(t, "", "validate", filepath.Join(dir, "cover.json"))
	h := must(t, "", "validate", filepath.Join(dir, "spec.json")).json(t)["hash"]

	// --- replay ---
	m := must(t, "", "claims merge", dir).json(t)
	if m["claims"].(float64) != 201 || m["independence_keys"].(float64) != 49 {
		t.Fatalf("merge %v", m)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "claims.jsonl")); string(b) != string(raw) {
		t.Fatal("re-merged claims.jsonl differs from the recorded one")
	}
	g := must(t, "", "graph build", dir).json(t)
	if len(g["regions"].([]any)) != 6 {
		t.Fatalf("graph %v", g)
	}
	// K5: the run's labels with brackets ("transparency obligations (Article 50)") stay balanced.
	for _, f := range append(g["files"].([]any), "graph/interface_candidates.json") {
		b, _ := os.ReadFile(filepath.Join(dir, f.(string)))
		if strings.Contains(string(b), "(article 50\"") || strings.Contains(string(b), "(article 6\"") {
			t.Errorf("%s still holds an unbalanced label", f)
		}
	}
	im := must(t, "", "interface merge", dir, "--plain").json(t)
	// K8: dates given as values (applies_from, published_on) now count as dated content. The run measured
	// 0.08 on its as_of-only definition (deviation seq 254). It is now 0.29: still under the timeline
	// threshold (artifact.layout.timeline_min_temporal_density = 0.4), because most of this run's relations
	// (states, obliges, position_on) carry no date. The threshold is not tuned to this one run; it is an open
	// calibration item.
	if td := im["shape_metrics"].(map[string]any)["temporal_density"].(float64); td != 0.29 {
		t.Errorf("temporal_density %v, want 0.29", td)
	}

	fv := must(t, "", "findings verify", dir).json(t)
	byKind := fv["by_kind"].(map[string]any)
	// K6 (as_of window = month): two recorded "verified" findings rested on claims from different months, and
	// now split into single-source findings per window:
	//   - "AI Act fines imposed / no fines reported": 2026-08 (wsgrdataadvisor.com) vs 2026-10 (reedsmith.com),
	//     the runner's own deviation (seq 205/206);
	//   - "harmonised standards / publication status: pending": 2025-06 vs 2026-04.
	// So 103 → 105 findings, 73 → 71 verified, 30 → 34 single-source. Conflicts are unchanged at 3.
	// D-097 (containment, undated-as-any-period, UTC) changes nothing here: neither split key has an undated or
	// a coarser claim, and the run has no same-key value change across periods (changed_over_time = 0).
	if fv["findings"].(float64) != 105 || byKind["verified"].(float64) != 71 || byKind["single_source"].(float64) != 34 ||
		fv["conflicts"].(float64) != 3 || fv["changed_over_time"].(float64) != 0 {
		t.Fatalf("findings %v", fv)
	}
	// Every recorded finding not touched by K5 (a bracketed label) or K6 (a dated claim in its evidence) comes
	// back with the same id, kind, confidence and evidence.
	dated := map[string]bool{}
	for _, c := range all {
		if c.AsOf != "" {
			dated[c.ID] = true
		}
	}
	now := map[string]map[string]any{}
	for _, f := range readJSONL(t, filepath.Join(dir, "findings.jsonl")) {
		now[f["id"].(string)] = f
	}
	evidence := map[string]map[string]any{} // evidence set -> new finding
	for _, n := range now {
		evidence[jsonStr(n["evidence"])] = n
	}
	same, bracketed, windowed := 0, 0, 0
	for _, f := range recFindings {
		k := f["key"].(map[string]any)
		isDated := false
		for _, ev := range f["evidence"].([]any) {
			isDated = isDated || dated[ev.(string)]
		}
		if isDated {
			windowed++
			continue
		}
		if strings.Contains(k["subject"].(string)+k["object"].(string), "(") {
			// K5: the id changes with the balanced label, but the same claims still form one finding of the
			// same kind and confidence.
			bracketed++
			n := evidence[jsonStr(f["evidence"])]
			if n == nil || n["kind"] != f["kind"] || n["confidence"] != f["confidence"] {
				t.Errorf("bracketed finding %s lost its evidence group: %v", f["id"], n)
			}
			continue
		}
		n, ok := now[f["id"].(string)]
		if !ok || n["kind"] != f["kind"] || n["confidence"] != f["confidence"] || jsonStr(n["evidence"]) != jsonStr(f["evidence"]) {
			t.Errorf("recorded finding %s changed: %v -> %v", f["id"], f, n)
			continue
		}
		same++
	}
	if same != 77 || bracketed+windowed != 26 {
		t.Fatalf("unchanged %d, bracketed %d, dated %d", same, bracketed, windowed)
	}

	g2 := must(t, "", "gate G2", dir).json(t)
	if g2["verdict"] != "pass" || g2["complete"] != true {
		t.Fatalf("G2 %v", g2["verdict"])
	}
	hc := must(t, "", "report", dir).json(t)
	c := hc["counts"].(map[string]any)
	if c["claims"].(float64) != 201 || c["sources"].(float64) != 65 || c["independent_sources"].(float64) != 49 ||
		c["regions_gathered"].(float64) != 6 || c["verified"].(float64) != 71 || c["single_source"].(float64) != 34 ||
		c["conflicts"].(float64) != 3 || c["deviations"].(float64) != 42 || c["shared_entities"].(float64) != 19 {
		t.Fatalf("how_checked counts %v", c)
	}
	// K8: dated_claims counts dated content: the recording's 13 as_of claims plus 53 claims with an ISO date
	// as their value (applies_from, published_on) = 66.
	if c["dated_claims"].(float64) != 66 || c["changed_over_time"].(float64) != 0 {
		t.Fatalf("dated_claims %v, changed_over_time %v", c["dated_claims"], c["changed_over_time"])
	}
	if hc["spec_hash"] != h {
		t.Fatalf("spec hash %v vs %v", hc["spec_hash"], h)
	}
	must(t, "", "gate G6", filepath.Join(dir, "how_checked.json"), "--run", dir)

	// K7 on live data: retract the runner's mis-encoded claim (deviation seq 206: the 2 August source was
	// recorded as "no fines as of late October"). Its window's single-source finding disappears.
	must(t, "", "claim retract", dir, "--id", "01M653AP3AR8PABW94M9E5J6KM", "--reason",
		"The source is dated 2 August 2026 and does not support a late-October statement.", "--stage", "L2")
	must(t, "", "claims merge", dir)
	must(t, "", "graph build", dir)
	must(t, "", "interface merge", dir, "--plain")
	fv2 := must(t, "", "findings verify", dir).json(t)
	if fv2["claims"].(float64) != 200 || fv2["findings"].(float64) != 104 || fv2["by_kind"].(map[string]any)["single_source"].(float64) != 33 {
		t.Fatalf("after retraction %v", fv2)
	}
	must(t, "", "gate G2", dir)
	if hc2 := must(t, "", "report", dir).json(t); hc2["counts"].(map[string]any)["retracted_claims"].(float64) != 1 {
		t.Fatalf("report after retraction %v", hc2["counts"])
	}

	// K11 on live data: an edit naming a real finding passes; the reviewer's mis-cited style of id does not.
	fid := readJSONL(t, filepath.Join(dir, "findings.jsonl"))[0]["id"].(string)
	es := filepath.Join(dir, "review", "edits.1.json")
	os.MkdirAll(filepath.Dir(es), 0o755)
	os.WriteFile(es, []byte(`{"edits":[{"target_stage":"A2","target_id":"`+fid+`","op":"rewrite","rationale":"wording","priority":"P2"}],"reviewer_model":"m","blind":true}`), 0o644)
	must(t, "", "validate", es)
	os.WriteFile(es, []byte(`{"edits":[{"target_stage":"D3","target_id":"F-000000000000","op":"rewrite","rationale":"wording","priority":"P2"}],"reviewer_model":"m","blind":true}`), 0o644)
	if r := meta(t, "", "validate", es); r.code != 1 {
		t.Fatalf("unknown target id accepted: %+v", r)
	}
}

func jsonStr(v any) string { b, _ := json.Marshal(v); return string(b) }
