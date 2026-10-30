package findings

import (
	"testing"

	"github.com/manutej/meta-suite/next/kernel/internal/model"
)

// SYNTHETIC claims (placeholder entities and publishers) exercising the L2 rules.

var p = Params{MinIndependent: 2, NumericRelTol: 0.01, StripPrefixes: []string{"www."}}

func c(id, s, rel, o string, v any, unit, key, region string) model.Claim {
	return model.Claim{ID: id, Subject: s, Relation: rel, Object: o, Value: v, Unit: unit, Region: region,
		Source: model.Source{Locator: "https://" + key + "/" + id, IndependenceKey: key}}
}

func TestVerifiedNeedsDistinctKeys(t *testing.T) {
	fs, conf := Verify([]model.Claim{
		c("01", "Entity A", "part_of", "Entity B", nil, "", "pub-a.example", "r1"),
		c("02", "entity a", "part_of", "entity b.", nil, "", "www.pub-b.example", "r2"),
		c("03", "Entity A", "located_in", "Entity C", nil, "", "pub-a.example", "r1"),
		c("04", "Entity A", "located_in", "Entity C", nil, "", "www.pub-a.example", "r1"), // same publisher twice
	}, p)
	if len(fs) != 2 || len(conf.Conflicts) != 0 {
		t.Fatalf("findings %+v conflicts %+v", fs, conf)
	}
	kinds := map[string]string{}
	for _, f := range fs {
		kinds[f.Key.Relation] = f.Kind
	}
	if kinds["part_of"] != "verified" || kinds["located_in"] != "single_source" {
		t.Fatalf("kinds %v", kinds)
	}
}

func TestNumericToleranceAndConflict(t *testing.T) {
	fs, conf := Verify([]model.Claim{
		c("01", "E", "length", "E", 100.0, "m", "a.example", "r1"),
		c("02", "E", "length", "E", 100.5, "m", "b.example", "r2"), // within 1%
		c("03", "E", "length", "E", 120.0, "m", "c.example", "r3"), // disagrees
	}, p)
	if len(fs) != 2 || len(conf.Conflicts) != 1 {
		t.Fatalf("findings %+v conflicts %+v", fs, conf)
	}
	for _, f := range fs {
		if f.Confidence != "contested" {
			t.Fatalf("finding on a contested key must be contested: %+v", f)
		}
	}
	sides := conf.Conflicts[0].Sides
	if len(sides) != 2 || sides[0].Value.(float64) != 100 || sides[1].Value.(float64) != 120 {
		t.Fatalf("sides %+v", sides)
	}
	// the agreeing pair is verified (2 keys), the outlier single-source
	var v, s int
	for _, f := range fs {
		if f.Kind == "verified" {
			v++
		} else {
			s++
		}
	}
	if v != 1 || s != 1 {
		t.Fatalf("v=%d s=%d", v, s)
	}
}

func TestTextValuesAndUnitsDiffer(t *testing.T) {
	_, conf := Verify([]model.Claim{
		c("01", "E", "status", "E", "Active", "", "a.example", "r1"),
		c("02", "E", "status", "E", "active.", "", "b.example", "r1"),
		c("03", "E", "size", "E", 3.0, "acre", "a.example", "r1"),
		c("04", "E", "size", "E", 3.0, "hectare", "b.example", "r1"),
	}, p)
	if len(conf.Conflicts) != 1 || conf.Conflicts[0].Reason != "sources give values in different units that the kernel cannot convert" {
		t.Fatalf("%+v", conf)
	}
}

func TestDeterministicIDs(t *testing.T) {
	in := []model.Claim{c("02", "X", "r", "Y", nil, "", "b.example", "r1"), c("01", "X", "r", "Y", nil, "", "a.example", "r1")}
	a, _ := Verify(in, p)
	b, _ := Verify([]model.Claim{in[1], in[0]}, p)
	if a[0].ID != b[0].ID || a[0].Evidence[0] != "01" {
		t.Fatalf("%+v vs %+v", a, b)
	}
}

// D-087: claims in different as_of windows never verify each other, and a different value in another window
// is a change over time, not a conflict. Undated claims are their own window; undated finding ids are unchanged.
func TestAsOfWindows(t *testing.T) {
	pm := p
	pm.AsOfWindow = "month"
	dated := func(cl model.Claim, asOf string) model.Claim { cl.AsOf = asOf; return cl }
	claims := []model.Claim{
		dated(c("01", "Fines imposed", "states", "none reported", nil, "", "a.example", "r1"), "2026-08-02"),
		dated(c("02", "Fines imposed", "states", "none reported", nil, "", "b.example", "r1"), "2026-10-26"),
		dated(c("03", "Standard X", "has_status", "publication status", "pending", "", "a.example", "r1"), "2025-06-07"),
		dated(c("04", "Standard X", "has_status", "publication status", "published", "", "b.example", "r1"), "2026-04-27"),
		dated(c("05", "Org", "headcount", "Org", 10.0, "count", "a.example", "r1"), "2026-10-01"),
		dated(c("06", "Org", "headcount", "Org", 10.0, "count", "b.example", "r1"), "2026-10-30"),
		c("07", "Org", "located_in", "City", nil, "", "a.example", "r1"),
		c("08", "Org", "located_in", "City", nil, "", "b.example", "r1"),
	}
	fs, conf := Verify(claims, pm)
	// D-097: different values in disjoint windows are a change over time, recorded, never a disagreement.
	if len(conf.Conflicts) != 1 || conf.Conflicts[0].Kind != KindChanged || conf.Conflicts[0].Sides[0].Value != "pending" {
		t.Fatalf("want one changed_over_time record (pending, then published): %+v", conf)
	}
	byWindow := map[string]string{}
	for _, f := range fs {
		byWindow[f.Key.Subject+"@"+f.AsOf] = f.Kind
	}
	want := map[string]string{"fines imposed@2026-08": "single_source", "fines imposed@2026-10": "single_source",
		"standard x@2025-06": "single_source", "standard x@2026-04": "single_source", "org@2026-10": "verified", "org@": "verified"}
	for k, v := range want {
		if byWindow[k] != v {
			t.Errorf("%s = %q, want %q (all %v)", k, byWindow[k], v, byWindow)
		}
	}
	// With the window off, the old behaviour returns: the two fines claims verify each other.
	off := p
	off.AsOfWindow = "none"
	fs2, _ := Verify(claims[:2], off)
	if len(fs2) != 1 || fs2[0].Kind != "verified" || fs2[0].AsOf != "" {
		t.Fatalf("window none: %+v", fs2)
	}
	// Undated ids do not depend on the window setting.
	und1, _ := Verify(claims[6:], pm)
	und2, _ := Verify(claims[6:], off)
	if und1[0].ID != und2[0].ID {
		t.Fatal("undated finding id changed with the window setting")
	}
}

// D-097: the reviewer's repro cases (as_of_window = month). None may hide a disagreement.
func TestAsOfReviewCases(t *testing.T) {
	pm := p
	pm.AsOfWindow = "month"
	hc := func(id string, v float64, key, asOf string) model.Claim {
		cl := c(id, "Org", "headcount", "Org", v, "count", key, "r1")
		cl.AsOf = asOf
		return cl
	}
	kinds := func(conf Conflicts) map[string]int {
		m := map[string]int{}
		for _, x := range conf.Conflicts {
			m[x.Kind]++
		}
		return m
	}
	// 1. a dated value vs a different undated value: the undated claim overlaps every window -> disagreement.
	fs, conf := Verify([]model.Claim{hc("01", 100, "a.example", "2026-05"), hc("02", 900, "b.example", "")}, pm)
	if k := kinds(conf); k[KindDisagreement] != 1 || len(fs) != 2 || fs[0].Confidence != "contested" || fs[1].Confidence != "contested" {
		t.Errorf("case 1: %v %+v", k, fs)
	}
	// 2. different values on either side of a month boundary -> a recorded change over time.
	fs, conf = Verify([]model.Claim{hc("01", 100, "a.example", "2026-05-31"), hc("02", 900, "b.example", "2026-06-01")}, pm)
	if k := kinds(conf); k[KindChanged] != 1 || k[KindDisagreement] != 0 || len(conf.Conflicts[0].Sides) != 2 || conf.Conflicts[0].Sides[0].Value != 100.0 {
		t.Errorf("case 2: %+v", conf)
	}
	// 3. the same fact as of 2026 and as of 2026-03: the year claim supports the March slice (verified, as
	// of 2026-03) and keeps its own single-source finding at its own window; no widened window.
	fs, conf = Verify([]model.Claim{hc("01", 100, "a.example", "2026"), hc("02", 100, "b.example", "2026-03")}, pm)
	byAsOf := map[string]model.Finding{}
	for _, f := range fs {
		byAsOf[f.AsOf] = f
	}
	if len(fs) != 2 || byAsOf["2026-03"].Kind != "verified" || len(byAsOf["2026-03"].Supporting) != 1 ||
		byAsOf["2026"].Kind != "single_source" || len(conf.Conflicts) != 0 {
		t.Errorf("case 3: %+v %+v", fs, conf)
	}
	// 4. the same instant written with two offsets -> one window -> one verified finding.
	fs, conf = Verify([]model.Claim{hc("01", 100, "a.example", "2026-01-31T23:00-05:00"), hc("02", 100, "b.example", "2026-02-01T04:00Z")}, pm)
	if len(fs) != 1 || fs[0].Kind != "verified" || fs[0].AsOf != "2026-02" {
		t.Errorf("case 4: %+v", fs)
	}
	// 4b. the same instant with different values is a disagreement, not two months.
	_, conf = Verify([]model.Claim{hc("01", 100, "a.example", "2026-01-31T23:00-05:00"), hc("02", 900, "b.example", "2026-02-01T04:00Z")}, pm)
	if k := kinds(conf); k[KindDisagreement] != 1 {
		t.Errorf("case 4b: %v", k)
	}
	// 5. an undated claim supports a dated one with the same value (and keeps its own undated finding).
	fs, _ = Verify([]model.Claim{hc("01", 100, "a.example", "2026-05"), hc("02", 100, "b.example", "")}, pm)
	if len(fs) != 2 || fs[1].Kind != "verified" || fs[1].AsOf != "2026-05" || fs[1].Supporting[0] != "02" || fs[0].AsOf != "" || fs[0].Kind != "single_source" {
		t.Errorf("case 5: %+v", fs)
	}
	// Every claim is primary evidence in exactly one finding.
	all := []model.Claim{hc("01", 100, "a.example", "2026"), hc("02", 100, "b.example", "2026-03"), hc("03", 100, "c.example", "2026-05"),
		hc("04", 900, "d.example", ""), hc("05", 100, "e.example", "2027-01")}
	fs, _ = Verify(all, pm)
	seen := map[string]int{}
	for _, f := range fs {
		for _, e := range f.Evidence {
			seen[e]++
		}
	}
	for _, cl := range all {
		if seen[cl.ID] != 1 {
			t.Errorf("claim %s is in %d findings", cl.ID, seen[cl.ID])
		}
	}
}

// D-097 (N1): an undated claim supports each month it fits without chaining the months together.
func TestUndatedDoesNotChainMonths(t *testing.T) {
	pm := p
	pm.AsOfWindow = "month"
	hc := func(id, key, asOf string) model.Claim {
		cl := c(id, "Org", "headcount", "Org", 100.0, "count", key, "r1")
		cl.AsOf = asOf
		return cl
	}
	fs, conf := Verify([]model.Claim{hc("01", "a.example", "2026-05"), hc("02", "b.example", "2026-06"), hc("03", "c.example", "")}, pm)
	byAsOf := map[string]model.Finding{}
	for _, f := range fs {
		byAsOf[f.AsOf] = f
	}
	may, june, und := byAsOf["2026-05"], byAsOf["2026-06"], byAsOf[""]
	if len(fs) != 3 || len(conf.Conflicts) != 0 {
		t.Fatalf("want three findings (May, June, undated), no records: %+v %+v", fs, conf)
	}
	if _, wide := byAsOf["2026"]; wide {
		t.Fatal("a widened 2026 finding was made")
	}
	if may.Kind != "verified" || join(may.Evidence) != "01" || join(may.Supporting) != "03" {
		t.Errorf("May: %+v", may)
	}
	if june.Kind != "verified" || join(june.Evidence) != "02" || join(june.Supporting) != "03" {
		t.Errorf("June: %+v", june)
	}
	if und.Kind != "single_source" || join(und.Evidence) != "03" || len(und.Supporting) != 0 {
		t.Errorf("undated: %+v", und)
	}
	// The two months never verify each other: without the undated claim each is single-source.
	fs, _ = Verify([]model.Claim{hc("01", "a.example", "2026-05"), hc("02", "b.example", "2026-06")}, pm)
	if len(fs) != 2 || fs[0].Kind != "single_source" || fs[1].Kind != "single_source" {
		t.Errorf("two months alone: %+v", fs)
	}
}

func join(s []string) string {
	out := ""
	for i, x := range s {
		if i > 0 {
			out += ","
		}
		out += x
	}
	return out
}
