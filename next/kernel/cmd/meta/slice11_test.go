package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Slice 1.1 seam fixes (D-087): one test per defect found by the first real run. Inputs are SYNTHETIC.

// K2: run init creates every folder an op writes into, so `> "$RUN/gates/G2.json"` cannot fail in the shell.
func TestRunInitCreatesWorkDirs(t *testing.T) {
	dir := newRun(t, "fast", false)
	for _, d := range []string{"gates", "prompts", "review", "artifact", "artifact/shots", "claims", "tokens", "graph/regions"} {
		if st, err := os.Stat(filepath.Join(dir, d)); err != nil || !st.IsDir() {
			t.Errorf("run init did not create %s/", d)
		}
	}
	// Empty folders change no verdict: G2 on a fresh run is still "unchecked" (no claims), not a pass.
	if r := meta(t, "", "gate G2", dir); r.code != 1 || r.json(t)["verdict"] != "unchecked" {
		t.Fatalf("%+v", r)
	}
}

// coverWith writes a cover with extra relations and an optional labels list, and validates it.
func coverWith(t *testing.T, dir string, labels []string) {
	t.Helper()
	var cov map[string]any
	json.Unmarshal(readFixture(t, "cover.json"), &cov)
	cov["relation_vocab"] = append(cov["relation_vocab"].([]any),
		map[string]any{"name": "share", "value_type": "numeric", "unit_dimension": "ratio"},
		map[string]any{"name": "max_fine", "value_type": "numeric", "unit_dimension": "currency"})
	if labels != nil {
		cov["labels"] = labels
	}
	b, _ := json.Marshal(cov)
	os.WriteFile(filepath.Join(dir, "cover.json"), b, 0o644)
	must(t, "", "validate", filepath.Join(dir, "cover.json"))
}

func numClaim(rel string, v float64, unit string) string {
	return `{"subject":"Org Alpha","relation":"` + rel + `","object":"Org Alpha","value":` + strings.TrimSuffix(strings.TrimSuffix(jsonNum(v), ".0"), ".") +
		`,"unit":"` + unit + `","value_type":"numeric","source":{"locator":"https://a.example/n","retrieved_at":"2026-10-29T12:00:00Z","independence_key":"a.example"}}`
}

func jsonNum(v float64) string { b, _ := json.Marshal(v); return string(b) }

// K4: a canonical unit ("%") is accepted without a false note; currency codes are a non-convertible dimension.
func TestUnitsCanonicalAndCurrency(t *testing.T) {
	dir := newRun(t, "fast", false)
	coverWith(t, dir, nil)
	tok := token(t, dir, "D2", "region-a", "gather-a")
	for _, c := range []struct{ rel, unit string }{{"share", "%"}, {"max_fine", "EUR"}, {"max_fine", "USD"}} {
		out := must(t, numClaim(c.rel, 35, c.unit), "claim add", dir, "--token", tok).json(t)
		if out["notes"] != nil || out["normalized"] != nil {
			t.Errorf("%s %s: want no note and no conversion, got %v", c.rel, c.unit, out)
		}
	}
	if out := must(t, numClaim("share", 50, "bp"), "claim add", dir, "--token", tok).json(t); out["normalized"] == nil {
		t.Errorf("bp must still convert to %%: %v", out)
	}
	for _, c := range []struct{ rel, unit, want string }{
		{"max_fine", "€", "is not a code of it"},   // a currency relation needs an ISO 4217 code
		{"max_fine", "ABC", "is not a code of it"}, // code-shaped, but not in ISO 4217
		{"max_fine", "%", "measures ratio"},        // a ratio unit on a currency relation
		{"share", "EUR", "measures currency"},      // a currency code on a ratio relation
		{"headcount", "EUR", "measures currency"},  // ... or on a count relation
	} {
		r := meta(t, numClaim(c.rel, 1, c.unit), "claim add", dir, "--token", tok)
		if r.code != 1 || !strings.Contains(r.stdout, c.want) {
			t.Errorf("%s %s: want rejection %q, got %+v", c.rel, c.unit, c.want, r)
		}
	}
	// D-097: no unit at all on a currency relation is refused.
	noUnit := strings.Replace(numClaim("max_fine", 5, "EUR"), `"unit":"EUR",`, "", 1)
	if r := meta(t, noUnit, "claim add", dir, "--token", tok); r.code != 1 || !strings.Contains(r.stdout, "unit is required") {
		t.Errorf("currency claim without a unit: %+v", r)
	}
	// D-097: lower-case codes are upper-cased; a code-shaped non-currency unit ("FTE") on a count relation is
	// an unknown unit (noted), not money.
	must(t, numClaim("max_fine", 7, "eur"), "claim add", dir, "--token", tok)
	if out := must(t, numClaim("headcount", 3, "FTE"), "claim add", dir, "--token", tok).json(t); out["notes"] == nil {
		t.Errorf("FTE must be noted as an unknown unit: %v", out)
	}
	stored := false
	for _, c := range readJSONL(t, filepath.Join(dir, "claims", "gather-a.jsonl")) {
		stored = stored || (c["value"] == 7.0 && c["unit"] == "EUR")
	}
	if !stored {
		t.Error("eur was not stored as EUR")
	}
	// A unit in no table is still kept as given, with a note.
	if out := must(t, numClaim("share", 3, "widgets"), "claim add", dir, "--token", tok).json(t); out["notes"] == nil {
		t.Errorf("unknown unit must be noted: %v", out)
	}
}

// S3 support: with a shared label list in cover.json, claim add notes labels outside it (normalized match).
func TestCoverLabelsNote(t *testing.T) {
	dir := newRun(t, "fast", false)
	coverWith(t, dir, []string{"Org Alpha", "City One"})
	tok := token(t, dir, "D2", "region-a", "gather-a")
	if out := must(t, claimA, "claim add", dir, "--token", tok).json(t); out["notes"] != nil {
		t.Fatalf("listed labels must not be noted: %v", out)
	}
	c := strings.Replace(claimA, "City One", "City 1", 1)
	out := must(t, c, "claim add", dir, "--token", tok).json(t)
	if n, _ := out["notes"].([]any); len(n) != 1 || !strings.Contains(n[0].(string), `object "City 1"`) {
		t.Fatalf("want one label note: %v", out)
	}
}

// K7: retraction is an append-only tombstone that every later command honours.
func TestClaimRetract(t *testing.T) {
	dir := pipeline(t)
	must(t, "", "interface merge", dir, "--plain")
	must(t, "", "findings verify", dir)
	must(t, "", "gate G2", dir)
	merged := readJSONL(t, filepath.Join(dir, "claims.jsonl"))
	id := merged[0]["id"].(string)

	if k := meta(t, "", "claim retract", dir, "--id", id).errKind(t); k != "usage" { // reason required
		t.Fatal(k)
	}
	if k := meta(t, "", "claim retract", dir, "--id", "01M653AP3AR8PABW94M9E5J6KM", "--reason", "x").errKind(t); k != "precondition" {
		t.Fatal(k)
	}
	out := must(t, "", "claim retract", dir, "--id", id, "--reason", "The page supports only an earlier date.", "--stage", "L2").json(t)
	if out["retractions"].(float64) != 1 {
		t.Fatalf("%v", out)
	}
	if k := meta(t, "", "claim retract", dir, "--id", id, "--reason", "again").errKind(t); k != "exists" {
		t.Fatal(k)
	}
	must(t, "", "validate", filepath.Join(dir, "retractions.jsonl"))
	// Before re-merging, G2 fails on the stale merge (and the old G2 result is stale in the report).
	if r := meta(t, "", "gate G2", dir); r.code != 1 || !strings.Contains(r.stdout, "retracted claim is still merged") {
		t.Fatalf("stale merge must fail G2: %+v", r)
	}
	m := must(t, "", "claims merge", dir).json(t)
	if m["claims"].(float64) != 7 || m["retracted"].(float64) != 1 {
		t.Fatalf("merge after retraction: %v", m)
	}
	// The shard still holds the claim (append-only).
	shards, _ := filepath.Glob(filepath.Join(dir, "claims", "*.jsonl"))
	total := 0
	for _, s := range shards {
		total += len(readJSONL(t, s))
	}
	if total != 8 {
		t.Fatalf("shards must keep every line: %d", total)
	}
	must(t, "", "graph build", dir)
	must(t, "", "interface merge", dir, "--plain")
	after := must(t, "", "findings verify", dir).json(t)
	if after["claims"].(float64) != 7 {
		t.Fatalf("findings after retraction: %v", after)
	}
	for _, f := range readJSONL(t, filepath.Join(dir, "findings.jsonl")) {
		for _, ev := range f["evidence"].([]any) {
			if ev == id {
				t.Fatal("retracted claim still used as evidence")
			}
		}
	}
	if r := must(t, "", "gate G2", dir).json(t); r["verdict"] != "pass" || r["counts"].(map[string]any)["claims"].(float64) != 7 {
		t.Fatalf("G2 after retraction: %v", r)
	}
	hc := must(t, "", "report", dir).json(t)
	if hc["counts"].(map[string]any)["retracted_claims"].(float64) != 1 || len(hc["retractions"].([]any)) != 1 {
		t.Fatalf("report: %v", hc["counts"])
	}
	led := must(t, "", "ledger show", dir).stdout
	if !strings.Contains(led, `"command":"claim retract"`) {
		t.Fatal("no ledger event for claim retract")
	}
	// A retraction naming a claim in no shard (hand-written) fails G2.
	f, _ := os.OpenFile(filepath.Join(dir, "retractions.jsonl"), os.O_APPEND|os.O_WRONLY, 0o644)
	f.WriteString(`{"id":"01M653AP3AR8PABW94M9E5J6KM","reason":"forged","at":"2026-10-30T01:00:00.000Z"}` + "\n")
	f.Close()
	if r := meta(t, "", "gate G2", dir); r.code != 1 || !strings.Contains(r.stdout, "retracts a claim that is in no shard") {
		t.Fatalf("%+v", r)
	}
}

// K9: validating the run's spec.json records its hash in run.json and the ledger; the report carries it.
func TestSpecHashRecorded(t *testing.T) {
	dir := newRun(t, "deep", true)
	os.WriteFile(filepath.Join(dir, "spec.json"), readFixture(t, "spec.json"), 0o644)
	h := must(t, "", "validate", filepath.Join(dir, "spec.json")).json(t)["hash"].(string)
	var rec map[string]any
	readJSON(t, filepath.Join(dir, "run.json"), &rec)
	if rec["spec_hash"] != h {
		t.Fatalf("run.json spec_hash %v, want %s", rec["spec_hash"], h)
	}
	must(t, "", "validate", filepath.Join(dir, "run.json"))
	if !strings.Contains(must(t, "", "ledger show", dir).stdout, `"spec_hash_recorded":"run.json"`) {
		t.Fatal("ledger lacks the recorded spec hash")
	}
	if hc := must(t, "", "report", dir).json(t); hc["spec_hash"] != h {
		t.Fatalf("report spec_hash %v", hc["spec_hash"])
	}
	// A spec outside the run's spec.json path is hashed but not recorded.
	other := filepath.Join(dir, "prompts", "spec-draft.json")
	os.WriteFile(other, readFixture(t, "spec.json"), 0o644)
	must(t, "", "validate", other, "--type", "Spec")
}

// K10 + K11: run-aware id checks for EditScript, SectionMap, Coverage and Probes.
func TestValidateCrossRefs(t *testing.T) {
	dir := pipeline(t)
	must(t, "", "interface merge", dir, "--plain")
	must(t, "", "findings verify", dir)
	fs := readJSONL(t, filepath.Join(dir, "findings.jsonl"))
	fid := fs[0]["id"].(string)
	write := func(name, body string) string {
		p := filepath.Join(dir, name)
		os.MkdirAll(filepath.Dir(p), 0o755)
		os.WriteFile(p, []byte(body), 0o644)
		return p
	}
	os.WriteFile(filepath.Join(dir, "artifact", "index.html"), []byte(`<html><body><section id="s-verdict"><p id="f-`+fid+`">x</p></section><div id="how-checked"></div></body></html>`), 0o644)
	write("section_map.json", `{"sections":[{"id":"verdict","tag":"main","finding_ids":["`+fid+`"]},{"id":"how_checked","tag":"how_checked","finding_ids":[]}]}`)
	must(t, "", "validate", filepath.Join(dir, "section_map.json"))

	good := write("review/edits.1.json", `{"edits":[
		{"target_stage":"A2","target_id":"f-`+fid+`","op":"rewrite","rationale":"Finding `+fid+` reads stronger than it is.","priority":"P2"},
		{"target_stage":"A2","target_id":"how-checked","op":"rewrite","rationale":"x","priority":"P3"},
		{"target_stage":"D2","target_id":"region-b","op":"regather","rationale":"thin","priority":"P1"},
		{"target_stage":"A2","target_id":"verdict","op":"rewrite","rationale":"x","priority":"P3"},
		{"target_stage":"D1","target_id":"q-size","op":"add","rationale":"x","priority":"P3"}],
		"reviewer_model":"m","blind":true}`)
	v := must(t, "", "validate", good).json(t)
	if v["xref"] == nil {
		t.Fatalf("no xref info: %v", v)
	}
	bad := write("review/edits.2.json", `{"edits":[
		{"target_stage":"D3","target_id":"sec-made-up","op":"rewrite","rationale":"See F-000000000000.","priority":"P2"}],
		"reviewer_model":"m","blind":true}`)
	r := meta(t, "", "validate", bad)
	if r.code != 1 || !strings.Contains(r.stdout, "/edits/0/target_id") || !strings.Contains(r.stdout, "cites F-000000000000") {
		t.Fatalf("bad edit ids must be refused: %+v", r)
	}
	// --run: the same file outside the run, checked against the run.
	outside := filepath.Join(t.TempDir(), "edits.9.json")
	b, _ := os.ReadFile(bad)
	os.WriteFile(outside, b, 0o644)
	if r := meta(t, "", "validate", outside, "--run", dir); r.code != 1 {
		t.Fatalf("--run cross-check: %+v", r)
	}
	if r := meta(t, "", "validate", outside); r.code != 0 { // no run: schema only
		t.Fatalf("outside a run: %+v", r)
	}

	write("coverage.json", `{"q-where":["`+fid+`"],"q-size":[]}`)
	must(t, "", "validate", filepath.Join(dir, "coverage.json"))
	write("coverage.json", `{"root":["`+fid+`"],"q-parts":["F-000000000000"]}`)
	if r := meta(t, "", "validate", filepath.Join(dir, "coverage.json")); r.code != 1 || !strings.Contains(r.stdout, "not a leaf") || !strings.Contains(r.stdout, "not in findings.jsonl") {
		t.Fatalf("coverage: %+v", r)
	}
	write("probes.json", `{"region-a":["Q?"],"region-z":["Q?"]}`)
	if r := meta(t, "", "validate", filepath.Join(dir, "probes.json")); r.code != 1 || !strings.Contains(r.stdout, "region-z") {
		t.Fatalf("probes: %+v", r)
	}
	write("section_map.json", `{"sections":[{"id":"verdict","tag":"main","finding_ids":["F-000000000000"]}]}`)
	if r := meta(t, "", "validate", filepath.Join(dir, "section_map.json")); r.code != 1 {
		t.Fatalf("section map with an unknown finding: %+v", r)
	}
	must(t, "", "validate", filepath.Join(dir, "conflicts.json"))
}

// ---- D-097: review fixes ----

func findingsByLabel(t *testing.T, dir string) map[string]string {
	t.Helper()
	got := map[string]string{}
	for _, f := range readJSONL(t, filepath.Join(dir, "findings.jsonl")) {
		k := f["key"].(map[string]any)
		label := k["subject"].(string) + "|" + k["relation"].(string) + "|" + k["object"].(string)
		if v, ok := f["value"].(float64); ok {
			label += "=" + jsonNum(v)
		}
		got[label] = f["kind"].(string) + "/" + f["confidence"].(string)
	}
	return got
}

// P2-3: the original e2e scenario (2025 vs 2026 headcounts, undated claims corroborating dated ones) gives its
// intended outcome under the D-097 rules.
func TestAsOfScenario(t *testing.T) {
	b, err := os.ReadFile("testdata/synthetic/e2e-asof/claims.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	dir := pipelineFrom(t, b)
	must(t, "", "interface merge", dir, "--plain")
	fv := must(t, "", "findings verify", dir).json(t)
	if fv["conflicts"].(float64) != 1 || fv["changed_over_time"].(float64) != 0 {
		t.Fatalf("findings verify %v", fv)
	}
	want := map[string]string{
		"org alpha|located_in|city one":       "verified/settled",
		"org alpha|has_part|unit delta":       "single_source/single_source",
		"org alpha|headcount|org alpha=12000": "verified/contested",
		"org alpha|headcount|org alpha=15000": "single_source/contested",
		"org beta|has_part|unit gamma":        "single_source/single_source",
	}
	got := findingsByLabel(t, dir)
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q (all %v)", k, got[k], v, got)
		}
	}
	if len(got) != len(want) {
		t.Errorf("findings %v", got)
	}
	must(t, "", "gate G2", dir)
	must(t, "", "validate", filepath.Join(dir, "conflicts.json"))
}

// P1-1: a change over time is recorded, counted and said; how_checked never claims no disagreement then, and
// its rule names the window.
func TestChangeOverTimeReported(t *testing.T) {
	dir := newRun(t, "deep", true)
	tok := token(t, dir, "D2", "region-a", "gather-a")
	for _, c := range []struct {
		v         float64
		pub, asOf string
	}{{100, "a.example", "2026-05-31"}, {900, "b.example", "2026-06-01"}} {
		in := strings.Replace(numClaim("headcount", c.v, "count"), `"value_type"`, `"as_of":"`+c.asOf+`","value_type"`, 1)
		in = strings.ReplaceAll(in, "a.example", c.pub)
		must(t, in, "claim add", dir, "--token", tok)
	}
	must(t, "", "claims merge", dir)
	must(t, "", "graph build", dir)
	must(t, "", "interface merge", dir, "--plain")
	fv := must(t, "", "findings verify", dir).json(t)
	if fv["conflicts"].(float64) != 0 || fv["changed_over_time"].(float64) != 1 || !strings.Contains(fv["rule"].(string), "calendar month") {
		t.Fatalf("%v", fv)
	}
	var cf struct {
		Conflicts []struct {
			Kind  string `json:"kind"`
			Sides []struct {
				Value float64 `json:"value"`
			} `json:"sides"`
		} `json:"conflicts"`
	}
	readJSON(t, filepath.Join(dir, "conflicts.json"), &cf)
	if len(cf.Conflicts) != 1 || cf.Conflicts[0].Kind != "changed_over_time" || cf.Conflicts[0].Sides[0].Value != 100 || cf.Conflicts[0].Sides[1].Value != 900 {
		t.Fatalf("conflicts.json %+v", cf)
	}
	must(t, "", "gate G2", dir)
	hc := must(t, "", "report", dir).json(t)
	text := jsonStr(hc["summary"])
	if hc["counts"].(map[string]any)["changed_over_time"].(float64) != 1 || strings.Contains(text, "No sources disagreed") ||
		!strings.Contains(text, "changed over time") || !strings.Contains(hc["rule"].(string), "calendar month") || hc["rule"] != fv["rule"] {
		t.Fatalf("how_checked %v / %v", hc["summary"], hc["rule"])
	}
}

// P1-2: G2 refuses findings that predate the claims.
func TestG2StaleFindings(t *testing.T) {
	dir := pipeline(t)
	must(t, "", "interface merge", dir, "--plain")
	must(t, "", "findings verify", dir)
	must(t, "", "gate G2", dir)
	tok := token(t, dir, "D2", "region-c", "gather-late")
	must(t, strings.Replace(claimA, "Org Alpha", "Org Late", 1), "claim add", dir, "--token", tok)
	must(t, "", "claims merge", dir)
	must(t, "", "graph build", dir)
	must(t, "", "interface merge", dir, "--plain")
	r := meta(t, "", "gate G2", dir) // findings verify skipped
	if r.code != 1 || r.json(t)["verdict"] != "fail" || !strings.Contains(r.stdout, "findings are stale") {
		t.Fatalf("stale findings must fail G2: %+v", r)
	}
	must(t, "", "findings verify", dir)
	must(t, "", "gate G2", dir)
}

// P2-4: retraction is idempotent for the same reason, and a hand-written tombstone fails G2.
func TestRetractIdempotentAndLedgerBound(t *testing.T) {
	dir := pipeline(t)
	ids := readJSONL(t, filepath.Join(dir, "claims.jsonl"))
	id, other := ids[0]["id"].(string), ids[1]["id"].(string)
	must(t, "", "claim retract", dir, "--id", id, "--reason", "Misread date.")
	again := must(t, "", "claim retract", dir, "--id", id, "--reason", "Misread date.").json(t)
	if again["note"] == nil || again["retractions"].(float64) != 1 || len(readJSONL(t, filepath.Join(dir, "retractions.jsonl"))) != 1 {
		t.Fatalf("idempotent retract: %v", again)
	}
	must(t, "", "claims merge", dir)
	must(t, "", "graph build", dir)
	must(t, "", "interface merge", dir, "--plain")
	must(t, "", "findings verify", dir)
	must(t, "", "gate G2", dir)
	// A tombstone written by hand for a real claim has no ledger event.
	f, _ := os.OpenFile(filepath.Join(dir, "retractions.jsonl"), os.O_APPEND|os.O_WRONLY, 0o644)
	f.WriteString(`{"id":"` + other + `","reason":"by hand","at":"2026-10-30T01:00:00.000Z"}` + "\n")
	f.Close()
	must(t, "", "claims merge", dir)
	must(t, "", "graph build", dir)
	must(t, "", "interface merge", dir, "--plain")
	must(t, "", "findings verify", dir)
	if r := meta(t, "", "gate G2", dir); r.code != 1 || !strings.Contains(r.stdout, "no matching ledger event") {
		t.Fatalf("hand-written tombstone: %+v", r)
	}
}

// P2-6: an invalid spec.json clears the recorded hash; an edited spec drops it from the report.
func TestSpecHashCleared(t *testing.T) {
	dir := newRun(t, "deep", true)
	sp := filepath.Join(dir, "spec.json")
	os.WriteFile(sp, readFixture(t, "spec.json"), 0o644)
	must(t, "", "validate", sp)
	var spec map[string]any
	json.Unmarshal(readFixture(t, "spec.json"), &spec)
	spec["objective"] = "An edited objective."
	b, _ := json.Marshal(spec)
	os.WriteFile(sp, b, 0o644)
	if hc := must(t, "", "report", dir).json(t); hc["spec_hash"] != nil {
		t.Fatalf("edited spec: report keeps a stale hash %v", hc["spec_hash"])
	}
	os.WriteFile(sp, []byte(`{"objective":""}`), 0o644)
	if r := meta(t, "", "validate", sp); r.code != 1 {
		t.Fatalf("%+v", r)
	}
	var rec map[string]any
	readJSON(t, filepath.Join(dir, "run.json"), &rec)
	if _, ok := rec["spec_hash"]; ok {
		t.Fatalf("invalid spec: run.json keeps spec_hash %v", rec["spec_hash"])
	}
}

// P3-11: an edit's target_stage must fit what its target is.
func TestEditStageFitsTarget(t *testing.T) {
	dir := pipeline(t)
	must(t, "", "interface merge", dir, "--plain")
	must(t, "", "findings verify", dir)
	fid := readJSONL(t, filepath.Join(dir, "findings.jsonl"))[0]["id"].(string)
	os.WriteFile(filepath.Join(dir, "artifact", "index.html"), []byte(`<p id="lede">x</p>`), 0o644)
	es := filepath.Join(dir, "review", "edits.1.json")
	write := func(stage, target string) res {
		os.WriteFile(es, []byte(`{"edits":[{"target_stage":"`+stage+`","target_id":"`+target+`","op":"rewrite","rationale":"x","priority":"P2"}],"reviewer_model":"m","blind":true}`), 0o644)
		return meta(t, "", "validate", es)
	}
	for _, ok := range [][2]string{{"A2", "lede"}, {"A2", fid}, {"D2", fid}, {"D1", "region-a"}, {"I2", "q-size"}} {
		if r := write(ok[0], ok[1]); r.code != 0 {
			t.Errorf("%v must pass: %s", ok, r.stdout)
		}
	}
	for _, bad := range [][2]string{{"D3", "lede"}, {"I1", fid}, {"A2", "region-a"}, {"R1", "q-size"}} {
		if r := write(bad[0], bad[1]); r.code != 1 || !strings.Contains(r.stdout, "/edits/0/target_stage") {
			t.Errorf("%v must be refused: %+v", bad, r)
		}
	}
}
