// Package cover implements the kernel's checks on D1's Cover (D-042): the covering condition (every universe
// element is in some region, and regions use only universe elements), the computed overlaps and triple
// overlaps (never declared by the model), connectivity of the overlap graph, tree sanity, and the region cap
// modes.<mode>.max_regions.
package cover

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Region is one region of a cover.
type Region struct {
	ID       string   `json:"id"`
	Elements []string `json:"elements"`
}

// Node is a tree node.
type Node struct {
	ID       string `json:"id"`
	Children []Node `json:"children"`
}

// Cover is the subset of the Cover type the checks need.
type Cover struct {
	Universe       []string  `json:"universe"`
	Tree           Node      `json:"tree"`
	Regions        []Region  `json:"regions"`
	Overlaps       []Overlap `json:"overlaps"`
	TripleOverlaps []Overlap `json:"triple_overlaps"`
}

// Overlap is a computed pairwise or triple overlap.
type Overlap struct {
	Regions  []string `json:"regions"`
	Elements []string `json:"elements"`
}

// Issue is one failed cover check.
type Issue struct {
	Check   string `json:"check"`
	Message string `json:"message"`
}

// Report is the cover check's output.
type Report struct {
	Overlaps       []Overlap  `json:"overlaps"`
	TripleOverlaps []Overlap  `json:"triple_overlaps"`
	Components     [][]string `json:"components"`
	MaxRegions     int        `json:"max_regions,omitempty"`
	Issues         []Issue    `json:"issues"`
}

// Check runs every cover check. maxRegions <= 0 skips the region cap (no mode known).
func Check(doc []byte, maxRegions int) (*Report, error) {
	var c Cover
	if err := json.Unmarshal(doc, &c); err != nil {
		return nil, err
	}
	rep := &Report{MaxRegions: maxRegions, Issues: []Issue{}}
	add := func(check, f string, a ...any) { rep.Issues = append(rep.Issues, Issue{check, fmt.Sprintf(f, a...)}) }

	universe := map[string]bool{}
	for _, u := range c.Universe {
		universe[u] = true
	}
	ids := map[string]bool{}
	sets := map[string]map[string]bool{}
	regionIDs := make([]string, 0, len(c.Regions))
	for _, r := range c.Regions {
		if ids[r.ID] {
			add("unique_region_ids", "region id %q is used twice", r.ID)
			continue
		}
		ids[r.ID] = true
		regionIDs = append(regionIDs, r.ID)
		sets[r.ID] = map[string]bool{}
		for _, e := range r.Elements {
			if !universe[e] {
				add("regions_within_universe", "region %q uses element %q, which is not in the universe", r.ID, e)
			}
			sets[r.ID][e] = true
		}
	}
	sort.Strings(regionIDs)
	if maxRegions > 0 && len(regionIDs) > maxRegions {
		add("max_regions", "%d regions, the mode allows at most %d", len(regionIDs), maxRegions)
	}
	var uncovered []string
	for _, u := range c.Universe {
		covered := false
		for _, id := range regionIDs {
			covered = covered || sets[id][u]
		}
		if !covered {
			uncovered = append(uncovered, u)
		}
	}
	sort.Strings(uncovered)
	for _, u := range uncovered {
		add("covering", "universe element %q is in no region", u)
	}

	inter := func(regs ...string) []string {
		var out []string
		for e := range sets[regs[0]] {
			all := true
			for _, r := range regs[1:] {
				all = all && sets[r][e]
			}
			if all {
				out = append(out, e)
			}
		}
		sort.Strings(out)
		return out
	}
	rep.Overlaps = []Overlap{}
	rep.TripleOverlaps = []Overlap{}
	parent := map[string]string{}
	for _, id := range regionIDs {
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
	for i := range regionIDs {
		for j := i + 1; j < len(regionIDs); j++ {
			a, b := regionIDs[i], regionIDs[j]
			if e := inter(a, b); len(e) > 0 {
				rep.Overlaps = append(rep.Overlaps, Overlap{Regions: []string{a, b}, Elements: e})
				ra, rb := find(a), find(b)
				if ra != rb {
					parent[ra] = rb
				}
			}
			for k := j + 1; k < len(regionIDs); k++ {
				if e := inter(a, b, regionIDs[k]); len(e) > 0 {
					rep.TripleOverlaps = append(rep.TripleOverlaps, Overlap{Regions: []string{a, b, regionIDs[k]}, Elements: e})
				}
			}
		}
	}
	comps := map[string][]string{}
	for _, id := range regionIDs {
		comps[find(id)] = append(comps[find(id)], id)
	}
	for _, m := range comps {
		sort.Strings(m)
		rep.Components = append(rep.Components, m)
	}
	sort.Slice(rep.Components, func(i, j int) bool { return rep.Components[i][0] < rep.Components[j][0] })
	if len(rep.Components) > 1 {
		parts := make([]string, len(rep.Components))
		for i, m := range rep.Components {
			parts[i] = "{" + strings.Join(m, ",") + "}"
		}
		add("overlap_connected", "regions do not connect through shared elements: %d groups %s; add a bridging element", len(rep.Components), strings.Join(parts, " "))
	}

	// Declared overlaps, if any, must equal the computed ones (they are the kernel's to compute).
	if c.Overlaps != nil && !sameOverlaps(c.Overlaps, rep.Overlaps) {
		add("overlaps_computed", "overlaps are computed by the kernel; the declared overlaps differ (omit them)")
	}
	if c.TripleOverlaps != nil && !sameOverlaps(c.TripleOverlaps, rep.TripleOverlaps) {
		add("overlaps_computed", "triple_overlaps are computed by the kernel; the declared ones differ (omit them)")
	}

	// Tree: unique node ids; every leaf id is in the universe (leaves are the sub-questions D1 puts there).
	seen := map[string]bool{}
	var walk func(n Node)
	walk = func(n Node) {
		if seen[n.ID] {
			add("tree_unique_ids", "tree node id %q is used twice", n.ID)
		}
		seen[n.ID] = true
		if len(n.Children) == 0 && n.ID != c.Tree.ID && !universe[n.ID] {
			add("tree_leaves_in_universe", "tree leaf %q is not in the universe", n.ID)
		}
		for _, ch := range n.Children {
			walk(ch)
		}
	}
	walk(c.Tree)
	sort.SliceStable(rep.Issues, func(i, j int) bool {
		if rep.Issues[i].Check != rep.Issues[j].Check {
			return rep.Issues[i].Check < rep.Issues[j].Check
		}
		return rep.Issues[i].Message < rep.Issues[j].Message
	})
	return rep, nil
}

func sameOverlaps(a, b []Overlap) bool {
	key := func(os []Overlap) []string {
		out := make([]string, len(os))
		for i, o := range os {
			r := append([]string(nil), o.Regions...)
			e := append([]string(nil), o.Elements...)
			sort.Strings(r)
			sort.Strings(e)
			out[i] = strings.Join(r, ",") + "|" + strings.Join(e, ",")
		}
		sort.Strings(out)
		return out
	}
	ka, kb := key(a), key(b)
	if len(ka) != len(kb) {
		return false
	}
	for i := range ka {
		if ka[i] != kb[i] {
			return false
		}
	}
	return true
}
