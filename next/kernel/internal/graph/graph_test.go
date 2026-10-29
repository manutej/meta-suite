package graph

import (
	"testing"

	"github.com/manutej/meta-suite/next/kernel/internal/model"
)

// SYNTHETIC claims (placeholder labels) exercising entity merge, morphism grouping, candidates and shape.

func cl(id, region, s, rel, o, asOf string) model.Claim {
	return model.Claim{ID: id, Subject: s, Relation: rel, Object: o, Region: region, AsOf: asOf, ValueType: "qualitative"}
}

func TestBuildRegionMergesExactNormalizedLabels(t *testing.T) {
	g := BuildRegion("r1", []model.Claim{
		cl("01A", "r1", "Entity A", "part_of", "Entity B", ""),
		cl("01B", "r1", "entity  a.", "part_of", "entity b", "2024"),
		cl("01C", "r1", "Entity B", "located_in", "Entity C", ""),
	})
	if len(g.Objects) != 3 {
		t.Fatalf("objects %+v", g.Objects)
	}
	if g.Objects[0].ID != "r1:entity a" || len(g.Objects[0].Labels) != 2 {
		t.Fatalf("object 0 %+v", g.Objects[0])
	}
	if len(g.Morphisms) != 2 {
		t.Fatalf("morphisms %+v", g.Morphisms)
	}
	var partOf model.Morphism
	for _, m := range g.Morphisms {
		if m.RelationType == "part_of" {
			partOf = m
		}
	}
	if len(partOf.ClaimIDs) != 2 || partOf.ClaimIDs[0] != "01A" {
		t.Fatalf("part_of %+v", partOf)
	}
	sm := g.ShapeMetrics
	if sm.LongestChain != 2 || sm.HubConcentration != 1 || sm.TemporalDensity != 0.5 {
		t.Fatalf("shape %+v", sm)
	}
}

func TestCandidates(t *testing.T) {
	a := BuildRegion("a", []model.Claim{cl("1", "a", "X", "rel", "Y", ""), cl("2", "a", "Z", "other", "Y", "")})
	b := BuildRegion("b", []model.Claim{cl("3", "b", "x", "rel", "W", "")})
	c := BuildRegion("c", []model.Claim{cl("4", "c", "Q", "rel", "R", "")})
	got := Candidates([]model.RegionGraph{c, b, a})
	if len(got) != 1 || got[0].Regions[0] != "a" || got[0].Regions[1] != "b" {
		t.Fatalf("candidates %+v", got)
	}
	if len(got[0].IdentifiedObjects) != 1 || got[0].IdentifiedObjects[0].Objects[0] != "a:x" || got[0].IdentifiedRelations[0] != "rel" {
		t.Fatalf("candidate %+v", got[0])
	}
}

func TestShapeTwoHubs(t *testing.T) {
	// two balanced hubs h1, h2 with two spokes each: two_sidedness = (2/2) * 4/4 = 1, hub = 2/4.
	ms := []model.Morphism{{ID: "1", Src: "h1", Dst: "a"}, {ID: "2", Src: "h1", Dst: "b"}, {ID: "3", Src: "h2", Dst: "c"}, {ID: "4", Src: "h2", Dst: "d"}}
	s := Shape(nil, ms, []string{"1"})
	if s.TwoSidedness != 1 || s.HubConcentration != 0.5 || s.LongestChain != 1 || s.TemporalDensity != 0.25 {
		t.Fatalf("%+v", s)
	}
}

func TestLongestChainWithCycle(t *testing.T) {
	// a->b->c->a is one cycle (collapsed), then c->d->e: chain length 2 from the collapsed node.
	ms := []model.Morphism{{Src: "a", Dst: "b"}, {Src: "b", Dst: "c"}, {Src: "c", Dst: "a"}, {Src: "c", Dst: "d"}, {Src: "d", Dst: "e"}, {Src: "e", Dst: "e"}}
	if got := longestChain(nil, ms); got != 2 {
		t.Fatalf("longest chain %d", got)
	}
}

func TestEmptyShape(t *testing.T) {
	if s := Shape(nil, nil, nil); s != (model.ShapeMetrics{}) {
		t.Fatalf("%+v", s)
	}
}
