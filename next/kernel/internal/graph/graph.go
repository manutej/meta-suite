// Package graph implements D3's `meta graph build`: one RegionGraph per region from the merged claims, plus
// interface candidates for L3, and the shape metrics A1 reads.
//
// Objects are the subjects and objects of claims. Within a region, two labels are the same object when their
// normalized labels are equal (norm.Label); relation-profile merging (graph.merge_threshold < 1) arrives with
// D3's full version, so this build refuses any other threshold rather than silently ignoring it. Across
// regions nothing is merged here: a normalized label present in two regions becomes an interface candidate.
// Morphisms group claims by (src, relation, dst) and carry their claim ids.
package graph

import (
	"fmt"
	"sort"
	"strings"

	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
	"github.com/manutej/meta-suite/next/kernel/internal/norm"
)

// ObjectID is the id of a region's object with normalized label n.
func ObjectID(region, n string) string { return region + ":" + n }

// MorphismID is a content-derived morphism id, stable under claim order.
func MorphismID(scope, src, rel, dst string) string {
	return "m-" + fsx.SHA256([]byte(scope + "\x00" + src + "\x00" + rel + "\x00" + dst))[:16]
}

// BuildRegion builds one region's graph from its claims.
func BuildRegion(region string, claims []model.Claim) model.RegionGraph {
	labels := map[string]map[string]bool{} // normalized -> raw labels
	type mk struct{ src, rel, dst string }
	morphs := map[mk][]string{}
	dated := map[mk]bool{}
	addLabel := func(raw string) string {
		n := norm.Label(raw)
		if labels[n] == nil {
			labels[n] = map[string]bool{}
		}
		labels[n][strings.TrimSpace(raw)] = true
		return n
	}
	for _, c := range claims {
		s, o := addLabel(c.Subject), addLabel(c.Object)
		k := mk{ObjectID(region, s), c.Relation, ObjectID(region, o)}
		morphs[k] = append(morphs[k], c.ID)
		if Dated(c) {
			dated[k] = true
		}
	}
	g := model.RegionGraph{Region: region, Objects: []model.Object{}, Morphisms: []model.Morphism{}}
	for n, raw := range labels {
		g.Objects = append(g.Objects, model.Object{ID: ObjectID(region, n), Kind: "entity", Labels: sortedKeys(raw), Regions: []string{region}})
	}
	sort.Slice(g.Objects, func(i, j int) bool { return g.Objects[i].ID < g.Objects[j].ID })
	var datedIDs []string
	for k, ids := range morphs {
		sort.Strings(ids)
		id := MorphismID(region, k.src, k.rel, k.dst)
		g.Morphisms = append(g.Morphisms, model.Morphism{ID: id, Src: k.src, Dst: k.dst, RelationType: k.rel, ClaimIDs: dedupe(ids)})
		if dated[k] {
			datedIDs = append(datedIDs, id)
		}
	}
	sort.Slice(g.Morphisms, func(i, j int) bool { return g.Morphisms[i].ID < g.Morphisms[j].ID })
	g.ShapeMetrics = Shape(g.Objects, g.Morphisms, datedIDs)
	return g
}

// Dated reports whether a claim carries dated content (D-087): an as_of, or an ISO calendar date (at least
// year and month, norm.IsDate) as its value, object or subject. Dates encoded as values ("applies_from":
// "2026-08-02") count, so a date-heavy question can reach the timeline layout threshold.
func Dated(c model.Claim) bool {
	if c.AsOf != "" || norm.IsDate(c.Object) || norm.IsDate(c.Subject) {
		return true
	}
	v, ok := c.Value.(string)
	return ok && norm.IsDate(v)
}

// Candidates proposes one Interface per pair of regions that share at least one normalized label.
func Candidates(graphs []model.RegionGraph) []model.Interface {
	type info struct {
		labels map[string]string // normalized label -> object id
		rels   map[string]bool
	}
	byRegion := map[string]info{}
	var regions []string
	for _, g := range graphs {
		in := info{labels: map[string]string{}, rels: map[string]bool{}}
		for _, o := range g.Objects {
			in.labels[strings.TrimPrefix(o.ID, g.Region+":")] = o.ID
		}
		for _, m := range g.Morphisms {
			in.rels[m.RelationType] = true
		}
		byRegion[g.Region] = in
		regions = append(regions, g.Region)
	}
	sort.Strings(regions)
	out := []model.Interface{}
	for i := range regions {
		for j := i + 1; j < len(regions); j++ {
			a, b := byRegion[regions[i]], byRegion[regions[j]]
			var shared []string
			for n := range a.labels {
				if _, ok := b.labels[n]; ok {
					shared = append(shared, n)
				}
			}
			if len(shared) == 0 {
				continue
			}
			sort.Strings(shared)
			iface := model.Interface{Regions: []string{regions[i], regions[j]}, IdentifiedRelations: []string{}, DeclaredBy: "candidate:graph-build",
				Reason: fmt.Sprintf("%d normalized label(s) appear in both regions", len(shared))}
			for _, n := range shared {
				iface.IdentifiedObjects = append(iface.IdentifiedObjects, model.IdentifiedObject{Label: n, Objects: []string{a.labels[n], b.labels[n]}})
			}
			for r := range a.rels {
				if b.rels[r] {
					iface.IdentifiedRelations = append(iface.IdentifiedRelations, r)
				}
			}
			sort.Strings(iface.IdentifiedRelations)
			out = append(out, iface)
		}
	}
	return out
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func dedupe(s []string) []string {
	out := s[:0]
	for i, x := range s {
		if i == 0 || x != s[i-1] {
			out = append(out, x)
		}
	}
	return out
}
