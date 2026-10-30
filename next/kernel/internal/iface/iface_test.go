package iface

import (
	"testing"

	"github.com/manutej/meta-suite/next/kernel/internal/graph"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
)

// SYNTHETIC region graphs (placeholder labels) exercising plain merge, disjoint union, transitive closure
// across three regions and the refusal of an incoherent closure.

func cl(id, region, s, rel, o string) model.Claim {
	return model.Claim{ID: id, Subject: s, Relation: rel, Object: o, Region: region, ValueType: "qualitative"}
}

func three() []model.RegionGraph {
	return []model.RegionGraph{
		graph.BuildRegion("a", []model.Claim{cl("1", "a", "X", "rel", "Y")}),
		graph.BuildRegion("b", []model.Claim{cl("2", "b", "X", "rel", "Z")}),
		graph.BuildRegion("c", []model.Claim{cl("3", "c", "X", "other", "W")}),
	}
}

func TestPlainMergeAlongCandidates(t *testing.T) {
	gs := three()
	cands := graph.Candidates(gs)
	if len(cands) != 3 {
		t.Fatalf("candidates %d", len(cands))
	}
	res := Merge(gs, cands, nil, true, nil)
	if res.Graph == nil {
		t.Fatalf("refused: %+v %+v", res.Problems, res.Collisions)
	}
	g := res.Graph
	// X is identified across a, b, c -> one object with 3 members; Y, Z, W stay apart.
	if len(g.Objects) != 4 || g.Merge.Identifications != 2 {
		t.Fatalf("objects %+v merge %+v", g.Objects, g.Merge)
	}
	if g.Objects[0].ID != "a:x" || len(g.Objects[0].Members) != 3 || len(g.Objects[0].Regions) != 3 {
		t.Fatalf("hub %+v", g.Objects[0])
	}
	if len(g.Morphisms) != 3 || g.ShapeMetrics.HubConcentration != 1 {
		t.Fatalf("morphisms %+v shape %+v", g.Morphisms, g.ShapeMetrics)
	}
}

func TestDisjointUnionWithoutInterfaces(t *testing.T) {
	res := Merge(three(), nil, nil, true, nil)
	if res.Graph == nil || len(res.Graph.Objects) != 6 || res.Graph.Merge.Identifications != 0 {
		t.Fatalf("%+v", res)
	}
}

func TestIncoherentClosureRefusedWithChain(t *testing.T) {
	gs := three()
	cands := graph.Candidates(gs)
	// a:x ~ b:x and b:x ~ c:x imply a:x ~ c:x; a J1 verdict says a:x and c:x are distinct.
	var keep []model.Interface
	for _, c := range cands {
		if !(c.Regions[0] == "a" && c.Regions[1] == "c") {
			keep = append(keep, c)
		}
	}
	res := Merge(gs, keep, []Distinct{{Objects: []string{"a:x", "c:x"}, Reason: "judged distinct"}}, false, nil)
	if res.Graph != nil || len(res.Collisions) != 1 {
		t.Fatalf("want refusal, got %+v", res)
	}
	ch := res.Collisions[0].Chain
	if len(ch) != 2 || ch[0].From != "a:x" || ch[0].To != "b:x" || ch[1].To != "c:x" {
		t.Fatalf("chain %+v", ch)
	}
}

func TestBadInterfaceReported(t *testing.T) {
	bad := []model.Interface{{Regions: []string{"a", "b"}, IdentifiedObjects: []model.IdentifiedObject{{Label: "x", Objects: []string{"b:x", "a:x"}}}}}
	res := Merge(three(), bad, nil, false, nil)
	if res.Graph != nil || len(res.Problems) != 2 {
		t.Fatalf("%+v", res.Problems)
	}
}

func TestParseDistinct(t *testing.T) {
	if _, err := ParseDistinct([]byte(`[{"objects":["a:x","a:x"]}]`)); err == nil {
		t.Fatal("self-distinct accepted")
	}
	if _, err := ParseDistinct([]byte(`[{"objects":["a:x","b:x"],"extra":1}]`)); err == nil {
		t.Fatal("unknown field accepted")
	}
}
