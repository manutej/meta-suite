// Package report builds how_checked.json (`meta report`): a plain-language account of how a run was
// checked, built ONLY from recorded gate results (run.json gate_results) and counts computed from the run's
// files. It never re-judges anything. A gate with no recorded result is "unchecked" with the reason; a gate
// whose inputs changed after it ran is reported "stale" (and counted as unchecked).
package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/manutej/meta-suite/next/kernel/internal/claims"
	"github.com/manutej/meta-suite/next/kernel/internal/config"
	"github.com/manutej/meta-suite/next/kernel/internal/findings"
	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/graph"
	"github.com/manutej/meta-suite/next/kernel/internal/ledger"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
	"github.com/manutej/meta-suite/next/kernel/internal/norm"
	runs "github.com/manutej/meta-suite/next/kernel/internal/rundir"
)

// Gate is one gate's line in how_checked.json.
type Gate struct {
	Gate    string   `json:"gate"`
	Name    string   `json:"name"`
	Status  string   `json:"status"` // passed | failed | unchecked
	Reason  string   `json:"reason,omitempty"`
	At      string   `json:"at,omitempty"`
	Summary string   `json:"summary,omitempty"`
	Files   []string `json:"files,omitempty"`
}

// Counts are the run's numbers.
type Counts struct {
	Claims             int            `json:"claims"`
	Sources            int            `json:"sources"`
	IndependentSources int            `json:"independent_sources"`
	RegionsPlanned     int            `json:"regions_planned"`
	RegionsGathered    int            `json:"regions_gathered"`
	RegionsNotGathered []string       `json:"regions_not_gathered"`
	ClaimsByRegion     map[string]int `json:"claims_by_region"`
	Findings           int            `json:"findings"`
	Verified           int            `json:"verified"`
	SingleSource       int            `json:"single_source"`
	Contested          int            `json:"contested_findings"`
	Conflicts          int            `json:"conflicts"`
	SharedEntities     int            `json:"shared_entities"`
	RegionsCombined    *bool          `json:"regions_combined"`
	UnitConversions    int            `json:"unit_conversions"`
	DatedClaims        int            `json:"dated_claims"`
	RetractedClaims    int            `json:"retracted_claims"`
	Deviations         int            `json:"deviations"`
	DeviationsByStage  map[string]int `json:"deviations_by_stage"`
}

// Deviation is one recorded departure from an op (D-080), listed for reviewers to skim.
type Deviation struct {
	Seq   int    `json:"seq"`
	At    string `json:"at"`
	Stage string `json:"stage"`
	Note  string `json:"note"`
}

// Source is one entry of the sources list.
type Source struct {
	Locator     string   `json:"locator"`
	Publisher   string   `json:"publisher"`
	RetrievedAt string   `json:"retrieved_at"`
	ClaimIDs    []string `json:"claim_ids"`
}

// HowChecked is how_checked.json.
type HowChecked struct {
	RunID          string             `json:"run_id"`
	GeneratedAt    string             `json:"generated_at"`
	Mode           string             `json:"mode"`
	Kernel         string             `json:"kernel"`
	ConfigHash     string             `json:"config_hash"`
	SpecHash       string             `json:"spec_hash,omitempty"`
	SchemaVersion  string             `json:"schema_version"`
	Rule           string             `json:"rule"`
	MinIndependent int                `json:"min_independent"`
	Counts         Counts             `json:"counts"`
	Gates          []Gate             `json:"gates"`
	GatesPassed    int                `json:"gates_passed"`
	GatesFailed    int                `json:"gates_failed"`
	GatesUnchecked int                `json:"gates_unchecked"`
	Summary        []string           `json:"summary"`
	Limits         []string           `json:"limits"`
	Deviations     []Deviation        `json:"deviations"`
	Retractions    []model.Retraction `json:"retractions"`
	Sources        []Source           `json:"sources"`
}

var gateNames = map[string]string{
	"G1": "Every record matches its declared format",
	"G2": "Every statement traces to a named source, and independent sources really differ",
	"G3": "Overlapping parts of the research agree with each other",
	"G4": "Answers built from parts match the answer to the whole question",
	"G5": "The plan still matches what was asked",
	"G6": "No internal jargon reaches the reader",
	"G7": "The page renders, reads and works on phone and desktop",
}

// Build computes how_checked.json for a run. dated_claims counts claims with dated content (graph.Dated: an
// as_of or an ISO date as value or object, D-087); retracted claims are listed, and counted nowhere else.
func Build(r *runs.Run, cfg *config.Config) (*HowChecked, error) {
	strip := cfg.Independence.StripPrefixes
	hc := &HowChecked{
		RunID: r.Record.RunID, GeneratedAt: runs.Stamp(runs.Now()), Mode: r.Record.Mode, Kernel: r.Record.Kernel,
		ConfigHash: r.Record.ConfigHash, SpecHash: r.Record.SpecHash, SchemaVersion: r.Record.SchemaVersion, MinIndependent: cfg.Lens.Pullback.MinIndependent,
		Rule:  fmt.Sprintf("A point counts as verified only when at least %d independent publishers state it.", cfg.Lens.Pullback.MinIndependent),
		Gates: []Gate{}, Summary: []string{}, Limits: []string{}, Sources: []Source{}, Retractions: []model.Retraction{},
	}
	c := &hc.Counts
	c.ClaimsByRegion = map[string]int{}
	c.RegionsNotGathered = []string{}

	var all []model.Claim
	if r.Exists(runs.ClaimsFile) {
		var err error
		if all, err = claims.Load(r); err != nil {
			return nil, err
		}
	}
	srcs := map[string]*Source{}
	keys := map[string]bool{}
	for _, cl := range all {
		c.ClaimsByRegion[cl.Region]++
		keys[norm.IndependenceKey(cl.Source.IndependenceKey, strip)] = true
		s := srcs[cl.Source.Locator]
		if s == nil {
			s = &Source{Locator: cl.Source.Locator, Publisher: norm.IndependenceKey(cl.Source.IndependenceKey, strip), RetrievedAt: cl.Source.RetrievedAt}
			srcs[cl.Source.Locator] = s
		}
		if cl.Source.RetrievedAt < s.RetrievedAt {
			s.RetrievedAt = cl.Source.RetrievedAt
		}
		s.ClaimIDs = append(s.ClaimIDs, cl.ID)
		if cl.Normalized != nil {
			c.UnitConversions++
		}
		if graph.Dated(cl) {
			c.DatedClaims++
		}
	}
	c.Claims, c.Sources, c.IndependentSources = len(all), len(srcs), len(keys)
	for _, s := range srcs {
		sort.Strings(s.ClaimIDs)
		hc.Sources = append(hc.Sources, *s)
	}
	sort.Slice(hc.Sources, func(i, j int) bool {
		if hc.Sources[i].Publisher != hc.Sources[j].Publisher {
			return hc.Sources[i].Publisher < hc.Sources[j].Publisher
		}
		return hc.Sources[i].Locator < hc.Sources[j].Locator
	})

	var cover struct {
		Regions []struct {
			ID string `json:"id"`
		} `json:"regions"`
	}
	if r.Exists(runs.CoverFile) {
		if err := r.ReadJSONFile(runs.CoverFile, &cover); err != nil {
			return nil, err
		}
	}
	c.RegionsPlanned = len(cover.Regions)
	for _, reg := range cover.Regions {
		if c.ClaimsByRegion[reg.ID] > 0 {
			c.RegionsGathered++
		} else {
			c.RegionsNotGathered = append(c.RegionsNotGathered, reg.ID)
		}
	}
	sort.Strings(c.RegionsNotGathered)

	if r.Exists(runs.FindingsFile) {
		lines, err := fsx.ReadLines(r.Path(runs.FindingsFile))
		if err != nil {
			return nil, err
		}
		for _, l := range lines {
			var f model.Finding
			if err := json.Unmarshal(l.Raw, &f); err != nil {
				return nil, fmt.Errorf("%s line %d: %w", runs.FindingsFile, l.N, err)
			}
			c.Findings++
			switch f.Kind {
			case "verified":
				c.Verified++
			case "single_source":
				c.SingleSource++
			}
			if f.Confidence == "contested" {
				c.Contested++
			}
		}
	}
	if r.Exists(runs.ConflictsFile) {
		var cf findings.Conflicts
		if err := r.ReadJSONFile(runs.ConflictsFile, &cf); err != nil {
			return nil, err
		}
		c.Conflicts = len(cf.Conflicts)
	}
	if r.Exists(runs.MergedFile) {
		var dg model.DomainGraph
		if err := r.ReadJSONFile(runs.MergedFile, &dg); err != nil {
			return nil, err
		}
		for _, o := range dg.Objects {
			if len(o.Regions) > 1 {
				c.SharedEntities++
			}
		}
		combined := true
		c.RegionsCombined = &combined
	}

	evs, err := ledger.Events(r)
	if err != nil {
		return nil, err
	}
	c.DeviationsByStage = map[string]int{}
	hc.Deviations = []Deviation{}
	for _, ev := range evs {
		if ev.Command != "ledger append" || ev.Data["kind"] != "deviation" {
			continue
		}
		st, _ := ev.Data["stage"].(string)
		note, _ := ev.Data["note"].(string)
		hc.Deviations = append(hc.Deviations, Deviation{Seq: ev.Seq, At: ev.At, Stage: st, Note: note})
		c.DeviationsByStage[st]++
	}
	c.Deviations = len(hc.Deviations)

	tomb, err := claims.Retracted(r)
	if err != nil {
		return nil, err
	}
	for _, t := range tomb {
		hc.Retractions = append(hc.Retractions, t)
	}
	sort.Slice(hc.Retractions, func(i, j int) bool { return hc.Retractions[i].ID < hc.Retractions[j].ID })
	c.RetractedClaims = len(hc.Retractions)

	hc.Gates = gates(r, cfg)
	for _, g := range hc.Gates {
		switch g.Status {
		case "passed":
			hc.GatesPassed++
		case "failed":
			hc.GatesFailed++
		default:
			hc.GatesUnchecked++
		}
	}
	hc.Summary = summary(hc)
	hc.Limits = limits(hc)
	return hc, nil
}

func gates(r *runs.Run, cfg *config.Config) []Gate {
	implemented := map[string]bool{}
	for _, g := range cfg.Report.Implemented {
		implemented[g] = true
	}
	latest := map[string]model.GateResult{}
	g1 := map[string]model.GateResult{} // per file
	for _, res := range r.Record.GateResults {
		if res.Gate == "G1" {
			for _, f := range res.Files {
				g1[f] = res
			}
			continue
		}
		if prev, ok := latest[res.Gate]; !ok || res.At >= prev.At {
			latest[res.Gate] = res
		}
	}
	var out []Gate
	for _, id := range cfg.Report.Gates {
		g := Gate{Gate: id, Name: gateNames[id]}
		switch {
		case !implemented[id]:
			g.Status, g.Reason = "unchecked", "not part of this build"
		case id == "G1":
			if len(g1) == 0 {
				g.Status, g.Reason = "unchecked", "no format check was recorded for this run"
				break
			}
			files := make([]string, 0, len(g1))
			failed := []string{}
			for f, res := range g1 {
				files = append(files, f)
				if res.Verdict != "pass" {
					failed = append(failed, f)
				}
			}
			sort.Strings(files)
			sort.Strings(failed)
			g.Files = files
			g.Status = "passed"
			g.Summary = fmt.Sprintf("%d file(s) checked, and every claim was checked as it was recorded", len(files))
			if len(failed) > 0 {
				g.Status = "failed"
				g.Summary = fmt.Sprintf("%d of %d file(s) failed: %s", len(failed), len(files), strings.Join(failed, ", "))
			}
		default:
			res, ok := latest[id]
			if !ok {
				g.Status, g.Reason = "unchecked", "this check was not run"
				break
			}
			g.At, g.Summary = res.At, res.Summary
			switch res.Verdict {
			case "pass":
				g.Status = "passed"
			case "fail":
				g.Status = "failed"
			case "incomplete":
				g.Status, g.Reason = "unchecked", "only part of this check could run, because some of the run's files did not exist yet"
			default:
				g.Status, g.Reason = "unchecked", "the check could not run on this run's files"
			}
			if id == "G2" && res.InputsHash != "" && res.InputsHash != r.HashFiles(res.Files) {
				g.Status, g.Reason = "unchecked", "the files changed after this check ran (stale); run it again"
			}
		}
		out = append(out, g)
	}
	return out
}

func plural(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, one)
	}
	return fmt.Sprintf("%d %s", n, many)
}

func summary(hc *HowChecked) []string {
	c := hc.Counts
	s := []string{
		fmt.Sprintf("%s from %s, published by %s.", plural(c.Claims, "statement was gathered", "statements were gathered"),
			plural(c.Sources, "source", "sources"), plural(c.IndependentSources, "independent publisher", "independent publishers")),
		fmt.Sprintf("The question was split into %s; %d of them returned evidence.", plural(c.RegionsPlanned, "part", "parts"), c.RegionsGathered),
		fmt.Sprintf("%s verified (%s), %s on a single source.", plural(c.Verified, "point is", "points are"),
			strings.TrimSuffix(strings.ToLower(hc.Rule[:1])+hc.Rule[1:], "."), plural(c.SingleSource, "rests", "rest")),
	}
	if c.Conflicts > 0 {
		s = append(s, fmt.Sprintf("%s where sources disagree; both sides are shown and no winner is picked.", plural(c.Conflicts, "point is contested", "points are contested")))
	} else if c.Findings > 0 {
		s = append(s, "No sources disagreed on the same point.")
	}
	s = append(s, fmt.Sprintf("Checks: %d passed, %d failed, %d not run.", hc.GatesPassed, hc.GatesFailed, hc.GatesUnchecked))
	if c.RetractedClaims > 0 {
		s = append(s, fmt.Sprintf("%s withdrawn after gathering, each with its reason, and not counted.",
			plural(c.RetractedClaims, "statement was", "statements were")))
	}
	if c.Deviations > 0 {
		s = append(s, fmt.Sprintf("The process departed from its plan in %s, each logged with a reason.", plural(c.Deviations, "place", "places")))
	}
	return s
}

func limits(hc *HowChecked) []string {
	c := hc.Counts
	var l []string
	if len(c.RegionsNotGathered) > 0 {
		l = append(l, "No evidence was gathered for: "+strings.Join(c.RegionsNotGathered, ", ")+".")
	}
	if c.RegionsCombined == nil {
		l = append(l, "The parts of the research were not combined into one picture.")
	} else if c.SharedEntities == 0 && c.RegionsGathered > 1 {
		l = append(l, "The parts were researched separately; no shared subjects were found between them.")
	}
	if c.Findings == 0 {
		l = append(l, "No points were assessed for agreement between sources.")
	} else if c.Verified == 0 {
		l = append(l, "Nothing reached the verified bar; treat every point as provisional.")
	}
	for _, g := range hc.Gates {
		if g.Status == "unchecked" && g.Reason != "not part of this build" {
			l = append(l, fmt.Sprintf("Not checked: %s (%s).", strings.ToLower(g.Name[:1])+g.Name[1:], g.Reason))
		}
	}
	l = append(l, "Disagreements are listed, not resolved: weighing them is left to the reader.")
	return l
}
