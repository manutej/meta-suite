package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	runs "github.com/manutej/meta-suite/next/kernel/internal/rundir"
	"github.com/manutej/meta-suite/next/kernel/internal/schema"
)

// Xref says what a run-aware validation checked ids against (D-087, K11).
type Xref struct {
	Run     string   `json:"run"`
	Against []string `json:"against"`
	Skipped []string `json:"skipped,omitempty"`
}

// runIDs is every id a reviewer, a section map or a coverage map may name in a run.
type runIDs struct {
	findings, sections, conflicts, claims, regions, nodes, leaves, page map[string]bool
	have                                                                map[string]bool // which source files were read
}

var (
	pageIDRe  = regexp.MustCompile(`\bid\s*=\s*["']([^"']+)["']`)
	citedIDRe = regexp.MustCompile(`\b(?:F|C)-[0-9a-f]{12}\b`)
)

func loadRunIDs(r *runs.Run) *runIDs {
	ids := &runIDs{findings: map[string]bool{}, sections: map[string]bool{}, conflicts: map[string]bool{}, claims: map[string]bool{},
		regions: map[string]bool{}, nodes: map[string]bool{}, leaves: map[string]bool{}, page: map[string]bool{}, have: map[string]bool{}}
	lineIDs := func(rel string, into map[string]bool) {
		lines, err := fsx.ReadLines(r.Path(rel))
		if err != nil {
			return
		}
		ids.have[rel] = true
		for _, l := range lines {
			var x struct {
				ID string `json:"id"`
			}
			if json.Unmarshal(l.Raw, &x) == nil && x.ID != "" {
				into[x.ID] = true
			}
		}
	}
	lineIDs(runs.FindingsFile, ids.findings)
	lineIDs(runs.ClaimsFile, ids.claims)
	var sm struct {
		Sections []struct {
			ID string `json:"id"`
		} `json:"sections"`
	}
	if r.ReadJSONFile(runs.SectionMapFile, &sm) == nil {
		ids.have[runs.SectionMapFile] = true
		for _, s := range sm.Sections {
			ids.sections[s.ID] = true
		}
	}
	var cf struct {
		Conflicts []struct {
			ID string `json:"id"`
		} `json:"conflicts"`
	}
	if r.ReadJSONFile(runs.ConflictsFile, &cf) == nil {
		ids.have[runs.ConflictsFile] = true
		for _, c := range cf.Conflicts {
			ids.conflicts[c.ID] = true
		}
	}
	type node struct {
		ID       string `json:"id"`
		Children []node `json:"children"`
	}
	var cover struct {
		Tree    node `json:"tree"`
		Regions []struct {
			ID string `json:"id"`
		} `json:"regions"`
	}
	if r.ReadJSONFile(runs.CoverFile, &cover) == nil {
		ids.have[runs.CoverFile] = true
		for _, reg := range cover.Regions {
			ids.regions[reg.ID] = true
		}
		var walk func(n node)
		walk = func(n node) {
			ids.nodes[n.ID] = true
			if len(n.Children) == 0 {
				ids.leaves[n.ID] = true
			}
			for _, c := range n.Children {
				walk(c)
			}
		}
		walk(cover.Tree)
	}
	if b, err := os.ReadFile(r.Path(runs.ArtifactPage)); err == nil {
		ids.have[runs.ArtifactPage] = true
		for _, m := range pageIDRe.FindAllSubmatch(b, -1) {
			ids.page[string(m[1])] = true
		}
	}
	return ids
}

func (x *runIDs) against(files ...string) (have, skipped []string) {
	for _, f := range files {
		if x.have[f] {
			have = append(have, f)
		} else {
			skipped = append(skipped, f+" (absent)")
		}
	}
	return have, skipped
}

// crossCheck validates the ids inside a typed run file against the run's other files. It runs only for types
// that name ids (EditScript, SectionMap, Coverage, Probes) and returns nil info for the rest.
func crossCheck(r *runs.Run, typ string, doc []byte) ([]schema.Error, *Xref) {
	ids := loadRunIDs(r)
	var errs []schema.Error
	bad := func(loc, msg string) {
		errs = append(errs, schema.Error{InstanceLocation: loc, KeywordLocation: "kernel/xref", Message: msg})
	}
	info := &Xref{Run: r.Dir}
	switch typ {
	case "EditScript":
		info.Against, info.Skipped = ids.against(runs.FindingsFile, runs.SectionMapFile, runs.ConflictsFile, runs.ClaimsFile, runs.CoverFile, runs.ArtifactPage)
		var es struct {
			Edits []struct {
				TargetStage string `json:"target_stage"`
				TargetID    string `json:"target_id"`
				Rationale   string `json:"rationale"`
			} `json:"edits"`
		}
		json.Unmarshal(doc, &es)
		for i, e := range es.Edits {
			kinds := ids.kinds(e.TargetID)
			if len(kinds) == 0 {
				bad(fmt.Sprintf("/edits/%d/target_id", i), fmt.Sprintf("%q is not a finding, section, conflict, claim, region, sub-question or page element id in this run", e.TargetID))
			} else if allowed := stagesFor(kinds); !allowed[e.TargetStage] {
				bad(fmt.Sprintf("/edits/%d/target_stage", i), fmt.Sprintf("target_stage %s does not fit a %s id; one of %s (page wording goes to A2)",
					e.TargetStage, strings.Join(kinds, " or "), strings.Join(sortedMapKeys(allowed), ", ")))
			}
			for _, c := range uniq(citedIDRe.FindAllString(e.Rationale, -1)) {
				if !ids.findings[c] && !ids.conflicts[c] {
					bad(fmt.Sprintf("/edits/%d/rationale", i), fmt.Sprintf("cites %s, which is not a finding or conflict id in this run", c))
				}
			}
		}
	case "SectionMap":
		info.Against, info.Skipped = ids.against(runs.FindingsFile)
		if !ids.have[runs.FindingsFile] {
			return nil, info
		}
		var sm struct {
			Sections []struct {
				FindingIDs []string `json:"finding_ids"`
			} `json:"sections"`
		}
		json.Unmarshal(doc, &sm)
		for i, s := range sm.Sections {
			for j, f := range s.FindingIDs {
				if !ids.findings[f] {
					bad(fmt.Sprintf("/sections/%d/finding_ids/%d", i, j), fmt.Sprintf("%s is not in findings.jsonl", f))
				}
			}
		}
	case "Coverage":
		info.Against, info.Skipped = ids.against(runs.CoverFile, runs.FindingsFile)
		var cov map[string][]string
		json.Unmarshal(doc, &cov)
		for _, leaf := range sortedMapKeys(cov) {
			if ids.have[runs.CoverFile] && !ids.leaves[leaf] {
				bad("/"+leaf, fmt.Sprintf("%q is not a leaf sub-question of cover.json's tree", leaf))
			}
			if ids.have[runs.FindingsFile] {
				for j, f := range cov[leaf] {
					if !ids.findings[f] {
						bad(fmt.Sprintf("/%s/%d", leaf, j), fmt.Sprintf("%s is not in findings.jsonl", f))
					}
				}
			}
		}
	case "Probes":
		info.Against, info.Skipped = ids.against(runs.CoverFile)
		if !ids.have[runs.CoverFile] {
			return nil, info
		}
		var pr map[string][]string
		json.Unmarshal(doc, &pr)
		for _, reg := range sortedMapKeys(pr) {
			if !ids.regions[reg] {
				bad("/"+reg, fmt.Sprintf("%q is not a region id in cover.json", reg))
			}
		}
	default:
		return nil, nil
	}
	return errs, info
}

// targetStages maps the kind of id an edit targets to the stages that can change it (D-097, index.md's stage
// table): page elements and sections are page work (A1, A2); a finding can be reworded (A2), re-gathered (D2),
// relabelled (D3), split (L3) or retracted from (L2); a claim is re-gathered, relabelled or retracted; a
// region or sub-question belongs to the split (D1) or to gathering; a sub-question also to the spec (I2).
var targetStages = map[string][]string{
	"page element": {"A1", "A2"},
	"section":      {"A1", "A2"},
	"finding":      {"A2", "A1", "D2", "D3", "L3", "L2"},
	"conflict":     {"A2", "D2", "L3", "L2"},
	"claim":        {"D2", "D3", "L2"},
	"region":       {"D1", "D2"},
	"sub-question": {"I2", "D1", "D2"},
}

// kinds lists what an id is in this run (an id can be several: a page element "f-<finding>" is both).
func (x *runIDs) kinds(id string) []string {
	var k []string
	add := func(ok bool, kind string) {
		if ok {
			k = append(k, kind)
		}
	}
	add(x.page[id], "page element")
	add(x.sections[id] || (len(id) > 2 && id[:2] == "s-" && x.sections[id[2:]]), "section")
	add(x.findings[id] || (len(id) > 2 && id[:2] == "f-" && x.findings[id[2:]]), "finding")
	add(x.conflicts[id], "conflict")
	add(x.claims[id], "claim")
	add(x.regions[id], "region")
	add(x.nodes[id], "sub-question")
	return k
}

func stagesFor(kinds []string) map[string]bool {
	out := map[string]bool{}
	for _, k := range kinds {
		for _, s := range targetStages[k] {
			out[s] = true
		}
	}
	return out
}

func uniq(s []string) []string {
	m := map[string]bool{}
	var out []string
	for _, x := range s {
		if !m[x] {
			m[x] = true
			out = append(out, x)
		}
	}
	sort.Strings(out)
	return out
}

func sortedMapKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
