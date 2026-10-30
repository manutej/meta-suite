package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// TestEndToEnd drives the slice-1 pipeline on the SYNTHETIC scenario in testdata/synthetic/e2e:
// init → validate spec/cover → stage begin (4 gather agents, 3 regions) → claim add (8 claims, incl. a
// conflict, a duplicate publisher, a unit conversion and a region override) → claims merge → graph build →
// interface merge --plain → findings verify → gate G2 → gate G6 → ledger append → report, and asserts the
// outputs of every step.
func TestEndToEnd(t *testing.T) {
	dir := newRun(t, "deep", true)
	var rec map[string]any
	readJSON(t, filepath.Join(dir, "run.json"), &rec)
	if rec["mode"] != "deep_narrow" || len(rec["config_hash"].(string)) != 64 || rec["kernel"] != "present" {
		t.Fatalf("run.json %v", rec)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "brief.json")); string(b) != string(readFixture(t, "brief.json")) {
		t.Fatal("brief.json is not a byte copy of the brief")
	}
	os.WriteFile(filepath.Join(dir, "spec.json"), readFixture(t, "spec.json"), 0o644)
	if h := must(t, "", "validate", filepath.Join(dir, "spec.json"), "--type", "Spec").json(t)["hash"].(string); len(h) != 64 {
		t.Fatalf("spec hash %q", h)
	}

	// --- D2: one token per (agent, region); claims only through `claim add` ---
	tokens := map[string]string{}
	lines := strings.Split(strings.TrimSpace(string(readFixture(t, "claims.jsonl"))), "\n")
	var ids []string
	for i, l := range lines {
		var in struct {
			Agent  string          `json:"agent"`
			Region string          `json:"region"`
			Claim  json.RawMessage `json:"claim"`
		}
		if err := json.Unmarshal([]byte(l), &in); err != nil {
			t.Fatal(err)
		}
		key := in.Agent + "/" + in.Region
		if tokens[key] == "" {
			os.WriteFile(filepath.Join(dir, "prompt-"+in.Agent+".md"), []byte("synthetic prompt for "+in.Agent), 0o644)
			tokens[key] = strings.TrimSpace(must(t, "", "stage begin", dir, "D2", "--region", in.Region, "--agent", in.Agent,
				"--model", "test-model", "--prompt-file", filepath.Join(dir, "prompt-"+in.Agent+".md")).stdout)
		}
		out := must(t, string(in.Claim), "claim add", dir, "--token", tokens[key]).json(t)
		ids = append(ids, out["id"].(string))
		if i == 7 { // the last claim says region-a, its token says region-c: the token wins, with a note
			if out["region"] != "region-c" || len(out["notes"].([]any)) != 1 {
				t.Fatalf("region override: %v", out)
			}
		}
		if i == 1 { // 12 thousand -> 12000 count
			n := out["normalized"].(map[string]any)
			if n["from_unit"] != "thousand" || n["from_value"].(float64) != 12 {
				t.Fatalf("normalized %v", out)
			}
		}
	}
	if !sort.StringsAreSorted(ids) {
		t.Fatalf("ULIDs not increasing: %v", ids)
	}
	if shards, _ := filepath.Glob(filepath.Join(dir, "claims", "*.jsonl")); len(shards) != 4 {
		t.Fatalf("shards %v", shards)
	}

	// --- merge ---
	m := must(t, "", "claims merge", dir).json(t)
	if m["claims"].(float64) != 8 || m["independence_keys"].(float64) != 3 {
		t.Fatalf("merge %v", m)
	}
	merged := readJSONL(t, filepath.Join(dir, "claims.jsonl"))
	for i := 1; i < len(merged); i++ {
		if merged[i-1]["id"].(string) >= merged[i]["id"].(string) {
			t.Fatal("claims.jsonl not sorted by id")
		}
	}

	// --- D3 ---
	g := must(t, "", "graph build", dir).json(t)
	if g["interface_candidates"].(float64) != 3 || len(g["regions"].([]any)) != 3 {
		t.Fatalf("graph build %v", g)
	}
	var ra map[string]any
	readJSON(t, filepath.Join(dir, "graph", "regions", "region-a.json"), &ra)
	// region-a: "Org Alpha"/"org alpha" merge, "Unit Delta"/"Unit Delta." merge -> objects org alpha, city one, unit delta
	if n := len(ra["objects"].([]any)); n != 3 {
		t.Fatalf("region-a objects %d: %v", n, ra["objects"])
	}
	must(t, "", "validate", filepath.Join(dir, "graph", "regions", "region-a.json"))

	// --- L3 plain ---
	refused := meta(t, "", "interface merge", dir) // no declared interfaces, 3 regions, not plain
	if refused.code != 1 || refused.json(t)["merged"] != false {
		t.Fatalf("non-plain merge without interfaces must be refused: %+v", refused)
	}
	im := must(t, "", "interface merge", dir, "--plain").json(t)
	if im["merged"] != true || im["identifications"].(float64) != 3 { // org alpha x3 regions (2) + city one in a and b (1)
		t.Fatalf("interface merge %v", im)
	}
	var dg struct {
		Objects []struct {
			ID      string   `json:"id"`
			Regions []string `json:"regions"`
			Members []string `json:"members"`
		} `json:"objects"`
		ShapeMetrics map[string]float64 `json:"shape_metrics"`
	}
	readJSON(t, filepath.Join(dir, "graph", "merged.json"), &dg)
	var hub string
	for _, o := range dg.Objects {
		if len(o.Regions) == 3 {
			hub = o.ID
		}
	}
	if hub != "region-a:org alpha" || dg.ShapeMetrics["hub_concentration"] <= 0.5 {
		t.Fatalf("merged objects %+v shape %v", dg.Objects, dg.ShapeMetrics)
	}
	must(t, "", "validate", filepath.Join(dir, "graph", "merged.json"), "--type", "DomainGraph")

	// --- L2 ---
	fv := must(t, "", "findings verify", dir).json(t)
	byKind := fv["by_kind"].(map[string]any)
	if fv["conflicts"].(float64) != 1 || byKind["verified"].(float64) != 2 || byKind["single_source"].(float64) != 3 {
		t.Fatalf("findings %v", fv)
	}
	findings := readJSONL(t, filepath.Join(dir, "findings.jsonl"))
	got := map[string]string{}
	for _, f := range findings {
		k := f["key"].(map[string]any)
		label := k["subject"].(string) + "|" + k["relation"].(string) + "|" + k["object"].(string)
		if v, ok := f["value"]; ok {
			label += "=" + strconv.FormatFloat(v.(float64), 'f', -1, 64)
		}
		got[label] = f["kind"].(string) + "/" + f["confidence"].(string)
	}
	want := map[string]string{
		"org alpha|located_in|city one":       "verified/settled",            // a.example + b.example
		"org alpha|has_part|unit delta":       "single_source/single_source", // a.example twice = one publisher
		"org alpha|headcount|org alpha=12000": "verified/contested",          // a.example (12 thousand) + c.example
		"org alpha|headcount|org alpha=15000": "single_source/contested",
		"org beta|has_part|unit gamma":        "single_source/single_source",
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("finding %s = %q, want %q (all: %v)", k, got[k], v, got)
		}
	}
	var conf struct {
		Conflicts []struct {
			Sides []struct {
				Value float64 `json:"value"`
			} `json:"sides"`
		} `json:"conflicts"`
	}
	readJSON(t, filepath.Join(dir, "conflicts.json"), &conf)
	if len(conf.Conflicts) != 1 || conf.Conflicts[0].Sides[0].Value != 12000 || conf.Conflicts[0].Sides[1].Value != 15000 {
		t.Fatalf("conflicts %+v", conf)
	}
	must(t, "", "validate", filepath.Join(dir, "findings.jsonl"))

	// --- G2 ---
	g2 := must(t, "", "gate G2", dir).json(t)
	if g2["verdict"] != "pass" || g2["complete"] != true || len(g2["warnings"].([]any)) != 0 {
		t.Fatalf("G2 %v", g2)
	}

	// --- G6 on a page that passes, then the report ---
	page := filepath.Join(dir, "artifact", "index.html")
	os.MkdirAll(filepath.Dir(page), 0o755)
	os.WriteFile(page, []byte(`<!doctype html><html><head><title>Org Alpha</title></head><body><p>Two independent publishers agree it is in City One.</p></body></html>`), 0o644)
	must(t, "", "gate G6", page, "--run", dir)
	must(t, "", "ledger append", dir, "--kind", "deviation", "--stage", "D2", "--note", "synthetic: two agents gathered region-b")
	hc := must(t, "", "report", dir).json(t)
	c := hc["counts"].(map[string]any)
	if c["claims"].(float64) != 8 || c["independent_sources"].(float64) != 3 || c["verified"].(float64) != 2 ||
		c["conflicts"].(float64) != 1 || c["deviations"].(float64) != 1 || c["regions_gathered"].(float64) != 3 {
		t.Fatalf("how_checked counts %v", c)
	}
	status := map[string]string{}
	for _, g := range hc["gates"].([]any) {
		gm := g.(map[string]any)
		status[gm["gate"].(string)] = gm["status"].(string)
	}
	if status["G1"] != "passed" || status["G2"] != "passed" || status["G6"] != "passed" || status["G3"] != "unchecked" || status["G7"] != "unchecked" {
		t.Fatalf("gate statuses %v", status)
	}
	if hc["gates_passed"].(float64) != 3 || hc["gates_unchecked"].(float64) != 4 {
		t.Fatalf("gate counts %v %v", hc["gates_passed"], hc["gates_unchecked"])
	}
	// The report's own words must pass G6.
	must(t, "", "gate G6", filepath.Join(dir, "how_checked.json"), "--run", dir)

	// --- ledger: every command logged, hash chain intact, every event valid ---
	led := must(t, "", "ledger show", dir).stdout
	cmds := map[string]int{}
	for _, l := range strings.Split(strings.TrimSpace(led), "\n") {
		var ev map[string]any
		json.Unmarshal([]byte(l), &ev)
		cmds[ev["command"].(string)]++
	}
	for _, c := range []string{"run init", "validate", "stage begin", "claim add", "claims merge", "graph build", "interface merge",
		"findings verify", "gate G2", "gate G6", "ledger append", "report"} {
		if cmds[c] == 0 {
			t.Errorf("no ledger event for %q (%v)", c, cmds)
		}
	}
	if cmds["claim add"] != 8 || cmds["stage begin"] != 4 { // one token per (agent, region) pair
		t.Fatalf("ledger counts %v", cmds)
	}
	must(t, "", "validate", filepath.Join(dir, "ledger.jsonl"))

	// --- staleness: editing a claim after G2 makes the report mark G2 unchecked, and G2 now fails ---
	shard := filepath.Join(dir, "claims", "gather-c.jsonl")
	b, _ := os.ReadFile(shard)
	os.WriteFile(shard, []byte(strings.Replace(string(b), "c.example/data", "c.example/edited", 1)), 0o644)
	hc2 := must(t, "", "report", dir).json(t)
	for _, g := range hc2["gates"].([]any) {
		if gm := g.(map[string]any); gm["gate"] == "G2" && gm["status"] != "unchecked" {
			t.Fatalf("stale G2 reported %v", gm)
		}
	}
	if r := meta(t, "", "gate G2", dir); r.code != 1 || r.json(t)["verdict"] != "fail" {
		t.Fatalf("G2 after edit: %+v", r)
	}
}
