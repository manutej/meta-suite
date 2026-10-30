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
// Time (D-087): claims are grouped by key AND as_of window (lens.pullback.as_of_window; norm.Window), so
// claims about different windows never verify each other, and a value that differs between windows is a
// change over time, not a disagreement. Undated claims form their own window. A finding's as_of is its
// window, and its id includes the window only when there is one (undated finding ids are unchanged).
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

// Conflict is a same-key disagreement on value, passed on as a candidate (G3 in a later slice).
type Conflict struct {
	ID     string         `json:"id"`
	Key    model.Key      `json:"key"`
	AsOf   string         `json:"as_of,omitempty"`
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
	Rule         string         `json:"rule"`
	FindingsFile string         `json:"findings_file"`
	ConflictFile string         `json:"conflicts_file"`
}

type group struct {
	value  any
	unit   string
	claims []model.Claim
}

// fibre is one grouping cell: a key inside one as_of window.
type fibre struct {
	key    model.Key
	window string
}

func (f fibre) String() string {
	s := keyString(f.key)
	if f.window != "" {
		s += "\x00as_of=" + f.window
	}
	return s
}

// Verify computes findings and conflicts from claims.
func Verify(claims []model.Claim, p Params) ([]model.Finding, Conflicts) {
	byFibre := map[fibre][]model.Claim{}
	for _, c := range claims {
		k := model.Key{Subject: norm.Label(c.Subject), Relation: c.Relation, Object: norm.Label(c.Object)}
		f := fibre{key: k, window: norm.Window(c.AsOf, p.AsOfWindow)}
		byFibre[f] = append(byFibre[f], c)
	}
	fibres := make([]fibre, 0, len(byFibre))
	for f := range byFibre {
		fibres = append(fibres, f)
	}
	sort.Slice(fibres, func(i, j int) bool { return fibres[i].String() < fibres[j].String() })

	var out []model.Finding
	conf := Conflicts{Conflicts: []Conflict{}}
	for _, fb := range fibres {
		cs := byFibre[fb]
		sort.Slice(cs, func(i, j int) bool { return cs[i].ID < cs[j].ID })
		groups := groupValues(cs, p.NumericRelTol)
		contested := len(groups) > 1
		var sides []ConflictSide
		for _, g := range groups {
			f := makeFinding(fb, g, p, contested)
			out = append(out, f)
			if contested {
				sides = append(sides, ConflictSide{FindingID: f.ID, Value: g.value, Unit: g.unit, AsOf: asOfs(g.claims),
					ClaimIDs: f.Evidence, IndependenceKeys: f.IndependenceKeys, Regions: f.Regions})
			}
		}
		if contested {
			reason := "sources give different values for the same statement"
			units := map[string]bool{}
			for _, g := range groups {
				units[g.unit] = true
			}
			if len(units) > 1 {
				reason = "sources give values in different units that the kernel cannot convert"
			}
			conf.Conflicts = append(conf.Conflicts, Conflict{ID: "C-" + fsx.SHA256([]byte(fb.String()))[:12], Key: fb.key, AsOf: fb.window,
				Reason: reason, Sides: sides})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, conf
}

func makeFinding(fb fibre, g group, p Params, contested bool) model.Finding {
	k := fb.key
	ev := make([]string, 0, len(g.claims))
	ikeys, regs := map[string]bool{}, map[string]bool{}
	for _, c := range g.claims {
		ev = append(ev, c.ID)
		ikeys[norm.IndependenceKey(c.Source.IndependenceKey, p.StripPrefixes)] = true
		regs[c.Region] = true
	}
	sort.Strings(ev)
	kl, rl := sortedKeys(ikeys), sortedKeys(regs)
	kind, confidence := "single_source", "single_source"
	if len(kl) >= p.MinIndependent {
		kind, confidence = "verified", "settled"
	}
	if contested {
		confidence = "contested"
	}
	f := model.Finding{
		ID:               "F-" + fsx.SHA256([]byte(fb.String() + "\x00" + valueString(g.value) + "\x00" + g.unit))[:12],
		Kind:             kind,
		Evidence:         ev,
		Regions:          rl,
		Confidence:       confidence,
		Key:              &model.Key{Subject: k.Subject, Relation: k.Relation, Object: k.Object},
		Value:            g.value,
		Unit:             g.unit,
		AsOf:             fb.window,
		IndependenceKeys: kl,
	}
	if len(rl) == 1 {
		f.Region = rl[0]
	}
	f.PlainText = plainText(g, kind, len(kl), contested, fb.window)
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

func asOfs(cs []model.Claim) []string {
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
