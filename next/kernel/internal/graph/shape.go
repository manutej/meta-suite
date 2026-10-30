package graph

import (
	"math"
	"sort"

	"github.com/manutej/meta-suite/next/kernel/internal/model"
)

// Shape computes the four shape metrics A1's layout rules read (thresholds in config artifact.layout).
// Definitions (documented here because A1's rules depend on them):
//
//   - hub_concentration: the share of morphisms that touch the single most-connected object, in [0,1].
//     1 means every relation runs through one hub; 0 for an empty graph.
//   - longest_chain: the number of morphisms on the longest directed path through the graph after collapsing
//     each cycle to one node (self-loops ignored), so it is well-defined on any graph.
//   - two_sidedness: (deg(h2)/deg(h1)) × (share of morphisms touching h1 or h2), where h1 and h2 are the two
//     most-connected objects, in [0,1]. 1 means two equally weighted hubs account for every relation (a
//     natural split view); 0 when one hub dominates or there are fewer than two objects.
//   - temporal_density: the share of morphisms backed by at least one dated claim (Dated: an as_of, or an ISO
//     date as value, object or subject), in [0,1].
//
// Values are rounded to 6 decimals so output bytes are stable across platforms.
func Shape(objects []model.Object, morphs []model.Morphism, datedMorphisms []string) model.ShapeMetrics {
	var s model.ShapeMetrics
	if len(morphs) == 0 {
		return s
	}
	deg := map[string]int{}
	for _, m := range morphs {
		deg[m.Src]++
		if m.Dst != m.Src {
			deg[m.Dst]++
		}
	}
	type od struct {
		id string
		d  int
	}
	ranked := make([]od, 0, len(deg))
	for id, d := range deg {
		ranked = append(ranked, od{id, d})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].d != ranked[j].d {
			return ranked[i].d > ranked[j].d
		}
		return ranked[i].id < ranked[j].id
	})
	e := float64(len(morphs))
	h1 := ranked[0].id
	s.HubConcentration = round6(float64(ranked[0].d) / e)
	if len(ranked) > 1 {
		h2 := ranked[1].id
		touching := 0
		for _, m := range morphs {
			if m.Src == h1 || m.Dst == h1 || m.Src == h2 || m.Dst == h2 {
				touching++
			}
		}
		s.TwoSidedness = round6(float64(ranked[1].d) / float64(ranked[0].d) * float64(touching) / e)
	}
	s.LongestChain = longestChain(objects, morphs)
	s.TemporalDensity = round6(float64(len(datedMorphisms)) / e)
	return s
}

func round6(x float64) float64 { return math.Round(x*1e6) / 1e6 }

// longestChain: Tarjan SCC condensation, then the longest path (in edges) in the condensation DAG.
func longestChain(objects []model.Object, morphs []model.Morphism) int {
	nodes := map[string]bool{}
	for _, o := range objects {
		nodes[o.ID] = true
	}
	adj := map[string][]string{}
	for _, m := range morphs {
		nodes[m.Src], nodes[m.Dst] = true, true
		if m.Src != m.Dst {
			adj[m.Src] = append(adj[m.Src], m.Dst)
		}
	}
	ids := make([]string, 0, len(nodes))
	for n := range nodes {
		ids = append(ids, n)
	}
	sort.Strings(ids)
	for k := range adj {
		sort.Strings(adj[k])
	}
	index, low, comp := map[string]int{}, map[string]int{}, map[string]int{}
	onStack := map[string]bool{}
	var stack []string
	next, ncomp := 0, 0
	var strong func(v string)
	strong = func(v string) {
		index[v], low[v] = next, next
		next++
		stack = append(stack, v)
		onStack[v] = true
		for _, w := range adj[v] {
			if _, seen := index[w]; !seen {
				strong(w)
				low[v] = min(low[v], low[w])
			} else if onStack[w] {
				low[v] = min(low[v], index[w])
			}
		}
		if low[v] == index[v] {
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				comp[w] = ncomp
				if w == v {
					break
				}
			}
			ncomp++
		}
	}
	for _, v := range ids {
		if _, seen := index[v]; !seen {
			strong(v)
		}
	}
	// Tarjan numbers components in reverse topological order: edges go from higher to lower component ids.
	cadj := make([]map[int]bool, ncomp)
	for i := range cadj {
		cadj[i] = map[int]bool{}
	}
	for v, ws := range adj {
		for _, w := range ws {
			if comp[v] != comp[w] {
				cadj[comp[v]][comp[w]] = true
			}
		}
	}
	best := make([]int, ncomp)
	longest := 0
	for c := 0; c < ncomp; c++ { // successors have lower ids, so they are final already
		for d := range cadj[c] {
			if best[d]+1 > best[c] {
				best[c] = best[d] + 1
			}
		}
		longest = max(longest, best[c])
	}
	return longest
}
