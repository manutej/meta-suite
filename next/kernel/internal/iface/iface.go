// Package iface implements L3's `meta interface merge`: the only merge in the pipeline, and only along
// declared interfaces. Identifications from every interface are closed transitively with union-find; the
// merge is refused when the closure would join two objects recorded as distinct (graph/distinct.json, J1
// "distinct" verdicts), and the refusal names the chain of identifications that caused the collision (D-041).
//
// --plain (the v2−lens path, D-058/D-068) declares every interface candidate as-is. A plain merge with zero
// candidates succeeds: the merged graph is then the disjoint union of the region graphs. A non-plain merge of
// two or more regions with no declared interface is refused.
package iface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/manutej/meta-suite/next/kernel/internal/graph"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
)

// Distinct is a recorded "these are different things" constraint.
type Distinct struct {
	Objects     []string `json:"objects"`
	Reason      string   `json:"reason,omitempty"`
	JudgmentIDs []string `json:"judgment_ids,omitempty"`
}

// ParseDistinct decodes graph/distinct.json strictly.
func ParseDistinct(b []byte) ([]Distinct, error) {
	var out []Distinct
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&out); err != nil {
		return nil, err
	}
	for i, d := range out {
		if len(d.Objects) != 2 || d.Objects[0] == d.Objects[1] {
			return nil, fmt.Errorf("entry %d: objects must name two different object ids", i)
		}
	}
	return out, nil
}

// Link is one identification edge, with the interface that declared it.
type Link struct {
	A, B      string
	Interface int
	Label     string
}

// Collision explains a refused, incoherent merge.
type Collision struct {
	Distinct []string    `json:"distinct"`
	Reason   string      `json:"reason,omitempty"`
	Chain    []ChainStep `json:"chain"`
}

// ChainStep is one identification on the path between two objects judged distinct.
type ChainStep struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Interface int    `json:"interface"`
	Label     string `json:"label"`
}

// Problem is a reason an interface cannot be used.
type Problem struct {
	Interface int    `json:"interface"`
	Message   string `json:"message"`
}

// Result of a merge attempt.
type Result struct {
	Graph      *model.DomainGraph
	Problems   []Problem
	Collisions []Collision
}

// Merge merges region graphs along interfaces. claimsAsOf maps claim id -> has as_of (for temporal density).
func Merge(regions []model.RegionGraph, interfaces []model.Interface, distinct []Distinct, plain bool, dated map[string]bool) Result {
	sort.Slice(regions, func(i, j int) bool { return regions[i].Region < regions[j].Region })
	objRegion := map[string]string{}
	objects := map[string]model.Object{}
	var regionIDs []string
	for _, g := range regions {
		regionIDs = append(regionIDs, g.Region)
		for _, o := range g.Objects {
			objRegion[o.ID] = g.Region
			objects[o.ID] = o
		}
	}
	var res Result
	var links []Link
	for i, in := range interfaces {
		if len(in.Regions) != 2 || in.Regions[0] == in.Regions[1] {
			res.Problems = append(res.Problems, Problem{i, "an interface joins exactly two different regions"})
			continue
		}
		for _, r := range in.Regions {
			if !contains(regionIDs, r) {
				res.Problems = append(res.Problems, Problem{i, fmt.Sprintf("region %q has no graph (run `meta graph build`)", r)})
			}
		}
		for _, io := range in.IdentifiedObjects {
			if len(io.Objects) != 2 {
				res.Problems = append(res.Problems, Problem{i, fmt.Sprintf("identification %q must name two objects", io.Label)})
				continue
			}
			ok := true
			for k, o := range io.Objects {
				if objRegion[o] == "" {
					res.Problems = append(res.Problems, Problem{i, fmt.Sprintf("object %q is not in any region graph", o)})
					ok = false
				} else if objRegion[o] != in.Regions[k] {
					res.Problems = append(res.Problems, Problem{i, fmt.Sprintf("object %q belongs to region %q, not %q", o, objRegion[o], in.Regions[k])})
					ok = false
				}
			}
			if ok {
				links = append(links, Link{io.Objects[0], io.Objects[1], i, io.Label})
			}
		}
	}
	for _, d := range distinct {
		for _, o := range d.Objects {
			if objRegion[o] == "" {
				res.Problems = append(res.Problems, Problem{-1, fmt.Sprintf("distinct.json names unknown object %q", o)})
			}
		}
	}
	sortProblems(res.Problems)
	if len(res.Problems) > 0 {
		return res
	}

	// Union-find over object ids; the representative is the smallest id in the class.
	parent := map[string]string{}
	for id := range objects {
		parent[id] = id
	}
	var find func(string) string
	find = func(x string) string {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	union := func(a, b string) {
		ra, rb := find(a), find(b)
		if ra == rb {
			return
		}
		if ra < rb {
			parent[rb] = ra
		} else {
			parent[ra] = rb
		}
	}
	for _, l := range links {
		union(l.A, l.B)
	}
	for _, d := range distinct {
		if find(d.Objects[0]) == find(d.Objects[1]) {
			res.Collisions = append(res.Collisions, Collision{Distinct: append([]string(nil), d.Objects...), Reason: d.Reason, Chain: chain(links, d.Objects[0], d.Objects[1])})
		}
	}
	if len(res.Collisions) > 0 {
		sort.Slice(res.Collisions, func(i, j int) bool {
			return strings.Join(res.Collisions[i].Distinct, "|") < strings.Join(res.Collisions[j].Distinct, "|")
		})
		return res
	}

	classes := map[string][]string{}
	for id := range objects {
		classes[find(id)] = append(classes[find(id)], id)
	}
	dg := &model.DomainGraph{Regions: regionIDs, Objects: []model.Object{}, Morphisms: []model.Morphism{}}
	identifications := 0
	for rep, members := range classes {
		sort.Strings(members)
		labels, regs := map[string]bool{}, map[string]bool{}
		for _, m := range members {
			for _, l := range objects[m].Labels {
				labels[l] = true
			}
			regs[objRegion[m]] = true
		}
		o := model.Object{ID: rep, Kind: "entity", Labels: keys(labels), Regions: keys(regs)}
		if len(members) > 1 {
			o.Members = members
			identifications += len(members) - 1
		}
		dg.Objects = append(dg.Objects, o)
	}
	sort.Slice(dg.Objects, func(i, j int) bool { return dg.Objects[i].ID < dg.Objects[j].ID })
	type mk struct{ src, rel, dst string }
	grouped := map[mk]map[string]bool{}
	for _, g := range regions {
		for _, m := range g.Morphisms {
			k := mk{find(m.Src), m.RelationType, find(m.Dst)}
			if grouped[k] == nil {
				grouped[k] = map[string]bool{}
			}
			for _, c := range m.ClaimIDs {
				grouped[k][c] = true
			}
		}
	}
	var datedIDs []string
	for k, cs := range grouped {
		m := model.Morphism{ID: graph.MorphismID("merged", k.src, k.rel, k.dst), Src: k.src, Dst: k.dst, RelationType: k.rel, ClaimIDs: keys(cs)}
		dg.Morphisms = append(dg.Morphisms, m)
		for _, c := range m.ClaimIDs {
			if dated[c] {
				datedIDs = append(datedIDs, m.ID)
				break
			}
		}
	}
	sort.Slice(dg.Morphisms, func(i, j int) bool { return dg.Morphisms[i].ID < dg.Morphisms[j].ID })
	dg.ShapeMetrics = graph.Shape(dg.Objects, dg.Morphisms, datedIDs)
	dg.Merge = model.MergeInfo{Plain: plain, Interfaces: len(interfaces), Identifications: identifications}
	res.Graph = dg
	return res
}

// chain finds the shortest identification path from a to b (BFS, deterministic neighbour order).
func chain(links []Link, a, b string) []ChainStep {
	adj := map[string][]ChainStep{}
	for _, l := range links {
		adj[l.A] = append(adj[l.A], ChainStep{l.A, l.B, l.Interface, l.Label})
		adj[l.B] = append(adj[l.B], ChainStep{l.B, l.A, l.Interface, l.Label})
	}
	for k := range adj {
		s := adj[k]
		sort.Slice(s, func(i, j int) bool {
			if s[i].To != s[j].To {
				return s[i].To < s[j].To
			}
			return s[i].Interface < s[j].Interface
		})
	}
	prev := map[string]ChainStep{}
	seen := map[string]bool{a: true}
	queue := []string{a}
	for len(queue) > 0 && !seen[b] {
		x := queue[0]
		queue = queue[1:]
		for _, st := range adj[x] {
			if !seen[st.To] {
				seen[st.To] = true
				prev[st.To] = st
				queue = append(queue, st.To)
			}
		}
	}
	var out []ChainStep
	for x := b; x != a; {
		st, ok := prev[x]
		if !ok {
			return nil
		}
		out = append([]ChainStep{st}, out...)
		x = st.From
	}
	return out
}

func sortProblems(p []Problem) {
	sort.Slice(p, func(i, j int) bool {
		if p[i].Interface != p[j].Interface {
			return p[i].Interface < p[j].Interface
		}
		return p[i].Message < p[j].Message
	})
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func contains(l []string, s string) bool {
	for _, x := range l {
		if x == s {
			return true
		}
	}
	return false
}
