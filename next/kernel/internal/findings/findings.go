// Package findings implements L2 as code (`meta findings verify`). Claims are grouped by key
// (normalized subject + relation + normalized object). Inside a key, claims that agree on the value form one
// value group:
//
//   - no value: every claim on the key agrees (the statement itself is the fact);
//   - numbers: same unit, and |a-b| <= lens.pullback.numeric_rel_tol × max(|a|,|b|), chained in sorted order;
//   - text: equal normalized labels; booleans: equal.
//
// A value group is `verified` when its claims carry at least lens.pullback.min_independent distinct
// (normalized) independence keys, and `single_source` otherwise. When one key has two or more value groups,
// the disagreement is written to conflicts.json as a candidate and every finding on that key is marked
// confidence "contested". The kernel never picks a winner or averages.
//
// Time (D-087, D-097): every claim has an as_of window (norm.Window: UTC-normalized, truncated to
// lens.pullback.as_of_window). Two windows are compatible when one contains the other; an undated claim is
// compatible with every window. Inside a value group, claims split into time components (connected under
// compatibility), and each component is one finding: claims about disjoint periods never verify each other.
// Across value groups on one key:
//
//   - two components with compatible windows are a disagreement: a "disagreement" record in conflicts.json,
//     and both findings get confidence "contested";
//   - two dated components with disjoint windows are a change over time: a "changed_over_time" record listing
//     each value with its period, so the page can show the change. It is never dropped.
//
// A finding's as_of is the hull of its claims' windows (norm.WindowHull); its id includes as_of only when
// there is one, so undated finding ids are unchanged. Every claim lands in exactly one finding.
package findings

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
	"github.com/manutej/meta-suite/next/kernel/internal/norm"
)

// Params are the config values L2 uses.
type Params struct {
	MinIndependent int
	NumericRelTol  float64
	StripPrefixes  []string
	AsOfWindow     string // day | month | year | none
}

// ConflictSide is one value group on a contested key.
type ConflictSide struct {
	FindingID        string   `json:"finding_id"`
	Value            any      `json:"value"`
	Unit             string   `json:"unit,omitempty"`
	AsOf             []string `json:"as_of"`
	ClaimIDs         []string `json:"claim_ids"`
	IndependenceKeys []string `json:"independence_keys"`
	Regions          []string `json:"regions"`
}

// Conflict kinds in conflicts.json.
const (
	KindDisagreement = "disagreement"      // different values for overlapping periods (or undated)
	KindChanged      = "changed_over_time" // different values for disjoint periods
)

// Conflict is a same-key difference in value, passed on as a candidate (G3 in a later slice): a
// disagreement, or a change over time.
type Conflict struct {
	ID     string         `json:"id"`
	Kind   string         `json:"kind"`
	Key    model.Key      `json:"key"`
	Reason string         `json:"reason"`
	Sides  []ConflictSide `json:"sides"`
}

// Conflicts is conflicts.json.
type Conflicts struct {
	Conflicts []Conflict `json:"conflicts"`
}

// Summary is printed by the command.
type Summary struct {
	Claims       int            `json:"claims"`
	Keys         int            `json:"keys"`
	Findings     int            `json:"findings"`
	ByKind       map[string]int `json:"by_kind"`
	Conflicts    int            `json:"conflicts"`
	Changes      int            `json:"changed_over_time"`
	Rule         string         `json:"rule"`
	FindingsFile string         `json:"findings_file"`
	ConflictFile string         `json:"conflicts_file"`
}

type group struct {
	value  any
	unit   string
	claims []model.Claim
}

// comp is one time slice of a value group: the claims that state the value for exactly one window (their
// primary evidence), plus the coarser or undated claims of the same value that contain that window
// (supporting). One finding.
type comp struct {
	group      int
	value      any
	unit       string
	claims     []model.Claim // primary: window == asOf
	supporting []model.Claim // same value, window strictly contains asOf (or undated)
	asOf       string
	finding    model.Finding
}

// Rule is the agreement rule in plain words; `findings verify` and how_checked.json both state it (D-097).
func Rule(minIndependent int, window string) string {
	s := fmt.Sprintf("A point counts as verified only when at least %d independent publishers state it", minIndependent)
	if window == "" || window == "none" {
		return s + "."
	}
	return s + fmt.Sprintf(" for the same period: dates are compared by calendar %s, a statement dated by a longer period (a year) counts for every %s inside it, and an undated statement counts for any period.", window, window)
}

// Verify computes findings and conflicts from claims.
func Verify(claims []model.Claim, p Params) ([]model.Finding, Conflicts) {
	byKey := map[model.Key][]model.Claim{}
	for _, c := range claims {
		k := model.Key{Subject: norm.Label(c.Subject), Relation: c.Relation, Object: norm.Label(c.Object)}
		byKey[k] = append(byKey[k], c)
	}
	keys := make([]model.Key, 0, len(byKey))
	for k := range byKey {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keyString(keys[i]) < keyString(keys[j]) })

	var out []model.Finding
	conf := Conflicts{Conflicts: []Conflict{}}
	for _, k := range keys {
		cs := byKey[k]
		sort.Slice(cs, func(i, j int) bool { return cs[i].ID < cs[j].ID })
		var comps []*comp
		for gi, g := range groupValues(cs, p.NumericRelTol) {
			comps = append(comps, timeSlices(gi, g, p.AsOfWindow)...)
		}
		// Pairwise relations between components of different values.
		n := len(comps)
		disagree, change := newUF(n), newUF(n)
		inDisagree, inChange := make([]bool, n), make([]bool, n)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if comps[i].group == comps[j].group {
					continue
				}
				if overlaps(comps[i], comps[j]) {
					disagree.union(i, j)
					inDisagree[i], inDisagree[j] = true, true
				} else {
					change.union(i, j)
					inChange[i], inChange[j] = true, true
				}
			}
		}
		for i, c := range comps {
			c.finding = makeFinding(k, c, p, inDisagree[i])
			out = append(out, c.finding)
		}
		conf.Conflicts = append(conf.Conflicts, records(k, comps, disagree, inDisagree, KindDisagreement)...)
		conf.Conflicts = append(conf.Conflicts, records(k, comps, change, inChange, KindChanged)...)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, conf
}

// timeSlices splits a value group into one slice per distinct window (D-097). Every claim is primary
// evidence in the slice of its own window. A coarser or undated claim also supports each finer slice its window
// contains, without joining it: claims in two specific, disjoint windows never verify each other, and a
// slice's evidence plus supporting claims always mutually overlap (calendar windows nest).
func timeSlices(gi int, g group, window string) []*comp {
	byWindow := map[string]*comp{}
	var order []*comp
	ws := make([]string, len(g.claims))
	for i, c := range g.claims { // claims are sorted by id, so slices come out in a stable order
		ws[i] = norm.Window(c.AsOf, window)
		if byWindow[ws[i]] == nil {
			byWindow[ws[i]] = &comp{group: gi, unit: g.unit, asOf: ws[i]}
			order = append(order, byWindow[ws[i]])
		}
		byWindow[ws[i]].claims = append(byWindow[ws[i]].claims, c)
	}
	for _, sl := range order {
		for i, c := range g.claims {
			if ws[i] != sl.asOf && len(ws[i]) < len(sl.asOf) && norm.WindowsOverlap(ws[i], sl.asOf) {
				sl.supporting = append(sl.supporting, c)
			}
		}
		sl.value = g.value
		if _, isNum := g.value.(float64); isNum && len(order) > 1 {
			sl.value = representative(sl.claims)
		}
	}
	sort.SliceStable(order, func(i, j int) bool { return order[i].asOf < order[j].asOf })
	return order
}

func overlaps(a, b *comp) bool { return norm.WindowsOverlap(a.asOf, b.asOf) }

// records emits one conflicts.json record per connected set of related components.
func records(k model.Key, comps []*comp, uf *uf, member []bool, kind string) []Conflict {
	sets := map[int][]*comp{}
	var roots []int
	for i, c := range comps {
		if !member[i] {
			continue
		}
		r := uf.find(i)
		if sets[r] == nil {
			roots = append(roots, r)
		}
		sets[r] = append(sets[r], c)
	}
	var out []Conflict
	for _, r := range roots {
		set := sets[r]
		if kind == KindChanged {
			sort.SliceStable(set, func(i, j int) bool { return set[i].asOf < set[j].asOf })
		}
		id := keyString(k)
		var asOfs []string
		units := map[string]bool{}
		var sides []ConflictSide
		for _, c := range set {
			if c.asOf != "" {
				asOfs = append(asOfs, c.asOf)
			}
			units[c.unit] = true
			f := c.finding
			sides = append(sides, ConflictSide{FindingID: f.ID, Value: c.value, Unit: c.unit, AsOf: asOfsOf(c.claims),
				ClaimIDs: f.Evidence, IndependenceKeys: f.IndependenceKeys, Regions: f.Regions})
		}
		sort.Strings(asOfs)
		for _, a := range asOfs { // undated-only disagreements keep their pre-D-087 ids
			id += "\x00as_of=" + a
		}
		reason := "sources give different values for the same statement"
		if kind == KindChanged {
			id += "\x00changed"
			reason = "sources give different values for different periods: a change over time, each value with its date"
		} else if len(units) > 1 {
			reason = "sources give values in different units that the kernel cannot convert"
		}
		out = append(out, Conflict{ID: "C-" + fsx.SHA256([]byte(id))[:12], Kind: kind, Key: k, Reason: reason, Sides: sides})
	}
	return out
}

type uf struct{ parent []int }

func newUF(n int) *uf {
	u := &uf{parent: make([]int, n)}
	for i := range u.parent {
		u.parent[i] = i
	}
	return u
}

func (u *uf) find(i int) int {
	for u.parent[i] != i {
		u.parent[i] = u.parent[u.parent[i]]
		i = u.parent[i]
	}
	return i
}

func (u *uf) union(a, b int) {
	ra, rb := u.find(a), u.find(b)
	if ra < rb {
		u.parent[rb] = ra
	} else if rb < ra {
		u.parent[ra] = rb
	}
}

func makeFinding(k model.Key, c *comp, p Params, contested bool) model.Finding {
	ev := make([]string, 0, len(c.claims))
	ikeys, regs := map[string]bool{}, map[string]bool{}
	for _, cl := range c.claims {
		ev = append(ev, cl.ID)
		ikeys[norm.IndependenceKey(cl.Source.IndependenceKey, p.StripPrefixes)] = true
		regs[cl.Region] = true
	}
	sort.Strings(ev)
	var sup []string
	for _, cl := range c.supporting { // supporting claims count as publishers, deduplicated by key
		sup = append(sup, cl.ID)
		ikeys[norm.IndependenceKey(cl.Source.IndependenceKey, p.StripPrefixes)] = true
	}
	sort.Strings(sup)
	kl, rl := sortedKeys(ikeys), sortedKeys(regs)
	kind, confidence := "single_source", "single_source"
	if len(kl) >= p.MinIndependent {
		kind, confidence = "verified", "settled"
	}
	if contested {
		confidence = "contested"
	}
	idBase := keyString(k)
	if c.asOf != "" {
		idBase += "\x00as_of=" + c.asOf
	}
	f := model.Finding{
		ID:               "F-" + fsx.SHA256([]byte(idBase + "\x00" + valueString(c.value) + "\x00" + c.unit))[:12],
		Kind:             kind,
		Evidence:         ev,
		Supporting:       sup,
		Regions:          rl,
		Confidence:       confidence,
		Key:              &model.Key{Subject: k.Subject, Relation: k.Relation, Object: k.Object},
		Value:            c.value,
		Unit:             c.unit,
		AsOf:             c.asOf,
		IndependenceKeys: kl,
	}
	if len(rl) == 1 {
		f.Region = rl[0]
	}
	f.PlainText = plainText(group{value: c.value, unit: c.unit, claims: c.claims}, kind, len(kl), contested, c.asOf)
	return f
}

// plainText is a neutral sentence built only from the claim's own words and the counts.
func plainText(g group, kind string, nKeys int, contested bool, window string) string {
	c := g.claims[0]
	s := fmt.Sprintf("%s %s %s", strings.TrimSpace(c.Subject), strings.ReplaceAll(c.Relation, "_", " "), strings.TrimSpace(c.Object))
	if g.value != nil {
		s += ": " + valueString(g.value)
		if g.unit != "" {
			s += " " + g.unit
		}
	}
	if window != "" {
		s += ", as of " + window
	}
	src := "1 source"
	if nKeys != 1 {
		src = fmt.Sprintf("%d independent sources", nKeys)
	}
	switch {
	case contested:
		s += fmt.Sprintf(" (%s; other sources disagree)", src)
	case kind == "verified":
		s += fmt.Sprintf(" (%s agree)", src)
	default:
		s += " (" + src + " only)"
	}
	return s
}

func groupValues(cs []model.Claim, tol float64) []group {
	var none []model.Claim
	nums := map[string][]model.Claim{} // by unit
	texts := map[string][]model.Claim{}
	textRaw := map[string]string{}
	bools := map[bool][]model.Claim{}
	for _, c := range cs {
		switch v := c.Value.(type) {
		case nil:
			none = append(none, c)
		case float64:
			nums[c.Unit] = append(nums[c.Unit], c)
		case bool:
			bools[v] = append(bools[v], c)
		case string:
			n := norm.Label(v)
			texts[n] = append(texts[n], c)
			if _, ok := textRaw[n]; !ok || v < textRaw[n] {
				textRaw[n] = v
			}
		default:
			texts[fmt.Sprint(v)] = append(texts[fmt.Sprint(v)], c)
		}
	}
	var gs []group
	if len(none) > 0 {
		gs = append(gs, group{claims: none})
	}
	for _, b := range []bool{false, true} {
		if len(bools[b]) > 0 {
			gs = append(gs, group{value: b, claims: bools[b]})
		}
	}
	for _, n := range sortedKeysClaims(texts) {
		gs = append(gs, group{value: textRaw[n], claims: texts[n]})
	}
	for _, unit := range sortedKeysClaims(nums) {
		list := nums[unit]
		sort.SliceStable(list, func(i, j int) bool { return list[i].Value.(float64) < list[j].Value.(float64) })
		start := 0
		for i := 1; i <= len(list); i++ {
			if i < len(list) {
				a, b := list[i-1].Value.(float64), list[i].Value.(float64)
				if math.Abs(a-b) <= tol*math.Max(math.Abs(a), math.Abs(b)) {
					continue
				}
			}
			chunk := append([]model.Claim(nil), list[start:i]...)
			gs = append(gs, group{value: representative(chunk), unit: unit, claims: chunk})
			start = i
		}
	}
	for i := range gs {
		sort.Slice(gs[i].claims, func(a, b int) bool { return gs[i].claims[a].ID < gs[i].claims[b].ID })
	}
	return gs
}

// representative is the median of the group's values (lower median for even counts): a value a source
// actually reported, never an average.
func representative(cs []model.Claim) float64 {
	vals := make([]float64, len(cs))
	for i, c := range cs {
		vals[i] = c.Value.(float64)
	}
	sort.Float64s(vals)
	return vals[(len(vals)-1)/2]
}

func asOfsOf(cs []model.Claim) []string {
	m := map[string]bool{}
	for _, c := range cs {
		if c.AsOf != "" {
			m[c.AsOf] = true
		}
	}
	return sortedKeys(m)
}

func keyString(k model.Key) string { return k.Subject + "\x00" + k.Relation + "\x00" + k.Object }

func valueString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case float64:
		return strconv.FormatFloat(x, 'g', -1, 64)
	case bool:
		return strconv.FormatBool(x)
	default:
		return fmt.Sprint(x)
	}
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func sortedKeysClaims(m map[string][]model.Claim) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
