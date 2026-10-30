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
	if len(conf.Conflicts) != 0 {
		t.Fatalf("values in different windows are not a conflict: %+v", conf)
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
