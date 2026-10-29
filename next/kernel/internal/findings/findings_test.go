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
