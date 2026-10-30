// Package gates holds the deterministic gates of this build: G2 (provenance) and G6 (vocabulary).
package gates

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"

	"github.com/manutej/meta-suite/next/kernel/internal/claims"
	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/ledger"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
	"github.com/manutej/meta-suite/next/kernel/internal/norm"
	runs "github.com/manutej/meta-suite/next/kernel/internal/rundir"
	"github.com/manutej/meta-suite/next/kernel/internal/schema"
)

// Problem is one G2 violation.
type Problem struct {
	Check  string `json:"check"`
	File   string `json:"file"`
	ID     string `json:"id,omitempty"`
	Detail string `json:"detail"`
}

// Check is the status of one G2 check.
type Check struct {
	Name     string `json:"name"`
	Status   string `json:"status"` // pass | fail | unchecked
	Checked  int    `json:"checked"`
	Problems int    `json:"problems"`
	Reason   string `json:"reason,omitempty"`
}

// Warning is a non-failing G2 observation (D-081: injection_suspect).
type Warning struct {
	Kind    string `json:"kind"`
	File    string `json:"file"`
	ID      string `json:"id"`
	Field   string `json:"field"`
	Pattern string `json:"pattern"`
	Match   string `json:"match"`
}

// G2Report is G2's output.
type G2Report struct {
	Gate     string    `json:"gate"`
	Verdict  string    `json:"verdict"`
	Complete bool      `json:"complete"`
	Checks   []Check   `json:"checks"`
	Problems []Problem `json:"problems"`
	Warnings []Warning `json:"warnings"`
	Counts   G2Counts  `json:"counts"`
	Inputs   []string  `json:"inputs"`
}

// G2Counts summarises what was checked.
type G2Counts struct {
	Claims           int `json:"claims"`
	Sources          int `json:"sources"`
	IndependenceKeys int `json:"independence_keys"`
	RegionGraphs     int `json:"region_graphs"`
	Objects          int `json:"objects"`
	Morphisms        int `json:"morphisms"`
	Findings         int `json:"findings"`
}

// G2Inputs are the run-relative files G2 reads (for the inputs hash).
func G2Inputs(r *runs.Run) []string {
	in := []string{runs.ClaimsFile, runs.MergedFile, runs.FindingsFile, runs.RetractFile}
	in = append(in, r.Glob(runs.ClaimsDir+"/*.jsonl")...)
	in = append(in, r.Glob(runs.RegionsDir+"/*.json")...)
	in = append(in, r.Glob(runs.TokensDir+"/*.json")...)
	sort.Strings(in)
	return in
}

// G2 runs the provenance gate on a run directory (spec/v2/atoms/G2-gate-provenance.md):
//   - claims: every claim validates (G1), has a locator and an independence_key, carries a kernel-issued,
//     unedited stage token whose region is the claim's region; claims.jsonl equals the merged shards;
//   - graphs: every morphism references >= 1 claim that exists, in the same region (unsourced morphisms);
//     every object takes part in >= 1 morphism (orphans); every claim appears in its region's graph;
//   - findings: every evidence id exists; a verified finding's evidence has >= min_independent distinct
//     normalized independence keys, recomputed from the claims (fake independence); a tension has >= 1
//     claim on each side (unsourced sides).
//
// A part whose file does not exist yet is "unchecked" (with the reason). Verdict (D-087): "fail" when any
// check fails; else "unchecked" when no claims exist yet; else "incomplete" when any part is unchecked (only
// the parts that exist passed; never a clean pass); else "pass". complete says whether every part was checked.
//
// Retracted claims (retractions.jsonl, D-087) must still be valid and stamped in their shards, but must not
// appear in claims.jsonl, any graph or any finding's evidence; each retraction must name a shard claim.
//
// D-081: claim text fields matching any injection pattern (config gates.g2.injection_patterns) are listed
// as injection_suspect warnings. Warnings never change the verdict; they are for a reviewer to look at.
func G2(r *runs.Run, minIndependent int, strip []string, injection []string) (*G2Report, error) {
	rep := &G2Report{Gate: "G2", Problems: []Problem{}, Warnings: []Warning{}, Inputs: []string{}}
	pats := make([]*regexp.Regexp, 0, len(injection))
	for _, p := range injection {
		re, err := regexp.Compile("(?i)" + p)
		if err != nil {
			return nil, fmt.Errorf("gates.g2.injection_patterns: %v", err)
		}
		pats = append(pats, re)
	}
	for _, in := range G2Inputs(r) {
		if r.Exists(in) {
			rep.Inputs = append(rep.Inputs, in)
		}
	}
	add := func(p Problem) { rep.Problems = append(rep.Problems, p) }

	// --- claims ---
	shardClaims, shards, bad, err := claims.ReadShards(r)
	if err != nil {
		return nil, err
	}
	byID := map[string]model.Claim{}
	tomb, err := claims.Retracted(r)
	if err != nil {
		return nil, err
	}
	gone := func(id string) bool { _, ok := tomb[id]; return ok }
	claimCheck := Check{Name: "claims_sourced_and_stamped"}
	if len(shards) == 0 && !r.Exists(runs.ClaimsFile) {
		claimCheck.Status, claimCheck.Reason = "unchecked", "no claims gathered yet"
	} else {
		for _, b := range bad {
			for _, e := range b.Errors {
				add(Problem{Check: claimCheck.Name, File: b.File, Detail: fmt.Sprintf("line %d %s: %s", b.Line, e.InstanceLocation, e.Message)})
			}
		}
		inShard := map[string]bool{}
		for _, c := range shardClaims {
			inShard[c.ID] = true
			if c.Source.Locator == "" || c.Source.IndependenceKey == "" {
				add(Problem{Check: claimCheck.Name, File: runs.ClaimsDir, ID: c.ID, Detail: "missing locator or independence_key"})
			}
			if !gone(c.ID) {
				byID[c.ID] = c
			}
		}
		// D-097: a tombstone counts only with the ledger event that created it: an ok `claim retract` event
		// with the same id, reason and time, not a later "already retracted" call. A hand-written line in
		// retractions.jsonl has none.
		recorded := map[string]bool{}
		if len(tomb) > 0 {
			evs, err := ledger.Events(r)
			if err != nil {
				return nil, err
			}
			for _, ev := range evs {
				id, _ := ev.Data["id"].(string)
				reason, _ := ev.Data["reason"].(string)
				at, _ := ev.Data["at"].(string)
				repeat, _ := ev.Data["already_retracted"].(bool)
				if ev.Command == "claim retract" && ev.OK && id != "" && !repeat {
					recorded[id+"\x00"+reason+"\x00"+at] = true
				}
			}
		}
		for _, id := range sortedIDs(tomb) {
			if !inShard[id] {
				add(Problem{Check: claimCheck.Name, File: runs.RetractFile, ID: id, Detail: "retracts a claim that is in no shard"})
			}
			if !recorded[id+"\x00"+tomb[id].Reason+"\x00"+tomb[id].At] {
				add(Problem{Check: claimCheck.Name, File: runs.RetractFile, ID: id, Detail: "retraction has no matching ledger event: write retractions only with `meta claim retract`"})
			}
		}
		claimCheck.Checked = len(shardClaims) + len(bad)
		// claims.jsonl must be exactly the merged shards
		if !r.Exists(runs.ClaimsFile) {
			add(Problem{Check: claimCheck.Name, File: runs.ClaimsFile, Detail: "claims are not merged: run `meta claims merge`"})
		} else {
			merged, err := fsx.ReadLines(r.Path(runs.ClaimsFile))
			if err != nil {
				return nil, err
			}
			ids := map[string]bool{}
			for _, l := range merged {
				var c struct {
					ID string `json:"id"`
				}
				json.Unmarshal(l.Raw, &c)
				if gone(c.ID) {
					add(Problem{Check: claimCheck.Name, File: runs.ClaimsFile, ID: c.ID, Detail: "retracted claim is still merged: run `meta claims merge`"})
					continue
				}
				ids[c.ID] = true
				if errs, _ := schema.Validate("Claim", l.Raw); len(errs) > 0 {
					add(Problem{Check: claimCheck.Name, File: runs.ClaimsFile, ID: c.ID, Detail: fmt.Sprintf("line %d does not validate", l.N)})
				} else if sc, ok := byID[c.ID]; !ok {
					add(Problem{Check: claimCheck.Name, File: runs.ClaimsFile, ID: c.ID, Detail: "claim is not in any shard (written outside `meta claim add`?)"})
				} else {
					var mc model.Claim
					json.Unmarshal(l.Raw, &mc)
					if !reflect.DeepEqual(mc, sc) {
						add(Problem{Check: claimCheck.Name, File: runs.ClaimsFile, ID: c.ID, Detail: "differs from its shard: a file was edited after `meta claims merge`"})
					}
				}
			}
			if len(ids) != len(byID) {
				add(Problem{Check: claimCheck.Name, File: runs.ClaimsFile, Detail: fmt.Sprintf("%d claims merged, %d in shards: run `meta claims merge`", len(ids), len(byID))})
			}
		}
	}
	rep.Warnings = append(rep.Warnings, injectionScan(shardClaims, injection, pats)...)
	sources, keys := map[string]bool{}, map[string]bool{}
	for _, c := range byID {
		sources[c.Source.Locator] = true
		keys[norm.IndependenceKey(c.Source.IndependenceKey, strip)] = true
	}
	rep.Counts.Claims, rep.Counts.Sources, rep.Counts.IndependenceKeys = len(byID), len(sources), len(keys)

	// --- graphs ---
	graphCheck := Check{Name: "graph_no_orphans_no_unsourced"}
	regionFiles := r.Glob(runs.RegionsDir + "/*.json")
	if len(regionFiles) == 0 {
		graphCheck.Status, graphCheck.Reason = "unchecked", "graph not built yet"
	} else {
		inGraph := map[string]bool{}
		for _, rel := range regionFiles {
			var g model.RegionGraph
			if err := r.ReadJSONFile(rel, &g); err != nil {
				return nil, err
			}
			rep.Counts.RegionGraphs++
			rep.Counts.Objects += len(g.Objects)
			rep.Counts.Morphisms += len(g.Morphisms)
			used := map[string]bool{}
			for _, m := range g.Morphisms {
				graphCheck.Checked++
				used[m.Src], used[m.Dst] = true, true
				if len(m.ClaimIDs) == 0 {
					add(Problem{Check: graphCheck.Name, File: rel, ID: m.ID, Detail: "unsourced morphism (no claim ids)"})
				}
				for _, cid := range m.ClaimIDs {
					c, ok := byID[cid]
					switch {
					case gone(cid):
						add(Problem{Check: graphCheck.Name, File: rel, ID: m.ID, Detail: "references retracted claim " + cid + ": rebuild with `meta graph build`"})
					case !ok:
						add(Problem{Check: graphCheck.Name, File: rel, ID: m.ID, Detail: "references unknown claim " + cid})
					case c.Region != g.Region:
						add(Problem{Check: graphCheck.Name, File: rel, ID: m.ID, Detail: fmt.Sprintf("claim %s belongs to region %s", cid, c.Region)})
					default:
						inGraph[cid] = true
					}
				}
			}
			for _, o := range g.Objects {
				graphCheck.Checked++
				if !used[o.ID] {
					add(Problem{Check: graphCheck.Name, File: rel, ID: o.ID, Detail: "orphan object (in no morphism)"})
				}
			}
		}
		var missing []string
		for id := range byID {
			if !inGraph[id] {
				missing = append(missing, id)
			}
		}
		sort.Strings(missing)
		for _, id := range missing {
			add(Problem{Check: graphCheck.Name, File: runs.RegionsDir, ID: id, Detail: "claim is in no region graph: rebuild with `meta graph build`"})
		}
		if r.Exists(runs.MergedFile) {
			var dg model.DomainGraph
			if err := r.ReadJSONFile(runs.MergedFile, &dg); err != nil {
				return nil, err
			}
			used := map[string]bool{}
			for _, m := range dg.Morphisms {
				graphCheck.Checked++
				used[m.Src], used[m.Dst] = true, true
				if len(m.ClaimIDs) == 0 {
					add(Problem{Check: graphCheck.Name, File: runs.MergedFile, ID: m.ID, Detail: "unsourced morphism (no claim ids)"})
				}
				for _, cid := range m.ClaimIDs {
					if gone(cid) {
						add(Problem{Check: graphCheck.Name, File: runs.MergedFile, ID: m.ID, Detail: "references retracted claim " + cid + ": re-run `meta interface merge`"})
					} else if _, ok := byID[cid]; !ok {
						add(Problem{Check: graphCheck.Name, File: runs.MergedFile, ID: m.ID, Detail: "references unknown claim " + cid})
					}
				}
			}
			for _, o := range dg.Objects {
				graphCheck.Checked++
				if !used[o.ID] {
					add(Problem{Check: graphCheck.Name, File: runs.MergedFile, ID: o.ID, Detail: "orphan object (in no morphism)"})
				}
			}
		}
	}

	// --- findings ---
	findCheck := Check{Name: "findings_independent_and_sourced"}
	if !r.Exists(runs.FindingsFile) {
		findCheck.Status, findCheck.Reason = "unchecked", "findings not computed yet"
	} else {
		lines, err := fsx.ReadLines(r.Path(runs.FindingsFile))
		if err != nil {
			return nil, err
		}
		inFindings := map[string]int{} // claim id -> findings citing it
		for _, l := range lines {
			findCheck.Checked++
			rep.Counts.Findings++
			var f model.Finding
			if errs, _ := schema.Validate("Finding", l.Raw); len(errs) > 0 {
				add(Problem{Check: findCheck.Name, File: runs.FindingsFile, Detail: fmt.Sprintf("line %d does not validate", l.N)})
				continue
			}
			json.Unmarshal(l.Raw, &f)
			ks := map[string]bool{}
			for _, cid := range f.Evidence {
				c, ok := byID[cid]
				inFindings[cid]++
				if gone(cid) {
					add(Problem{Check: findCheck.Name, File: runs.FindingsFile, ID: f.ID, Detail: "evidence references retracted claim " + cid + ": re-run `meta findings verify`"})
					continue
				}
				if !ok {
					add(Problem{Check: findCheck.Name, File: runs.FindingsFile, ID: f.ID, Detail: "evidence references unknown claim " + cid})
					continue
				}
				ks[norm.IndependenceKey(c.Source.IndependenceKey, strip)] = true
			}
			for _, cid := range f.Supporting { // D-097: supporting claims count, but are not primary evidence
				c, ok := byID[cid]
				switch {
				case gone(cid):
					add(Problem{Check: findCheck.Name, File: runs.FindingsFile, ID: f.ID, Detail: "supporting claim " + cid + " is retracted: re-run `meta findings verify`"})
				case !ok:
					add(Problem{Check: findCheck.Name, File: runs.FindingsFile, ID: f.ID, Detail: "supporting claim " + cid + " is unknown"})
				default:
					ks[norm.IndependenceKey(c.Source.IndependenceKey, strip)] = true
				}
			}
			if f.Kind == "verified" && len(ks) < minIndependent {
				add(Problem{Check: findCheck.Name, File: runs.FindingsFile, ID: f.ID,
					Detail: fmt.Sprintf("verified with %d distinct independence key(s); needs %d (fake independence)", len(ks), minIndependent)})
			}
			if f.Kind == "tension" {
				if len(f.Sides) != 2 {
					add(Problem{Check: findCheck.Name, File: runs.FindingsFile, ID: f.ID, Detail: "tension needs two sides"})
				}
				for _, s := range f.Sides {
					if len(s.ClaimIDs) == 0 {
						add(Problem{Check: findCheck.Name, File: runs.FindingsFile, ID: f.ID, Detail: "unsourced side " + s.Label})
					}
				}
			}
		}
		// D-097: every live claim is primary evidence in exactly one finding; otherwise the findings predate
		// the claims (a merge after `findings verify`).
		ids := make([]string, 0, len(byID))
		for id := range byID {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for _, id := range ids {
			switch n := inFindings[id]; {
			case n == 0:
				add(Problem{Check: findCheck.Name, File: runs.FindingsFile, ID: id, Detail: "claim is in no finding: findings are stale: run `meta findings verify`"})
			case n > 1:
				add(Problem{Check: findCheck.Name, File: runs.FindingsFile, ID: id, Detail: fmt.Sprintf("claim is evidence in %d findings; each claim belongs to exactly one: run `meta findings verify`", n)})
			}
		}
	}

	sort.SliceStable(rep.Problems, func(i, j int) bool {
		a, b := rep.Problems[i], rep.Problems[j]
		if a.Check != b.Check {
			return a.Check < b.Check
		}
		if a.File != b.File {
			return a.File < b.File
		}
		if a.ID != b.ID {
			return a.ID < b.ID
		}
		return a.Detail < b.Detail
	})
	rep.Checks = []Check{claimCheck, graphCheck, findCheck}
	rep.Verdict, rep.Complete = "pass", true
	for i := range rep.Checks {
		c := &rep.Checks[i]
		for _, p := range rep.Problems {
			if p.Check == c.Name {
				c.Problems++
			}
		}
		if c.Status == "" {
			c.Status = "pass"
			if c.Problems > 0 {
				c.Status = "fail"
			}
		}
		if c.Status == "fail" {
			rep.Verdict = "fail"
		}
		if c.Status == "unchecked" {
			rep.Complete = false
		}
	}
	if rep.Verdict != "fail" {
		switch {
		case rep.Checks[0].Status == "unchecked":
			rep.Verdict = "unchecked"
		case !rep.Complete:
			rep.Verdict = "incomplete"
		}
	}
	return rep, nil
}

func sortedIDs(m map[string]model.Retraction) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// injectionScan checks each claim's text fields against the patterns, in claim-id then field order.
func injectionScan(cs []model.Claim, raw []string, pats []*regexp.Regexp) []Warning {
	var out []Warning
	for _, c := range cs {
		fields := [][2]string{{"subject", c.Subject}, {"relation", c.Relation}, {"object", c.Object}}
		if v, ok := c.Value.(string); ok {
			fields = append(fields, [2]string{"value", v})
		}
		fields = append(fields, [2]string{"source.locator", c.Source.Locator}, [2]string{"source.independence_key", c.Source.IndependenceKey})
		for _, f := range fields {
			for i, re := range pats {
				if m := re.FindString(f[1]); m != "" {
					out = append(out, Warning{Kind: "injection_suspect", File: runs.ClaimsDir + "/" + c.StageToken.Agent + ".jsonl", ID: c.ID, Field: f[0], Pattern: raw[i], Match: m})
				}
			}
		}
	}
	return out
}
