package main

import (
	"errors"
	"os"
	"sort"

	"github.com/manutej/meta-suite/next/kernel/internal/claims"
	"github.com/manutej/meta-suite/next/kernel/internal/findings"
	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/graph"
	"github.com/manutej/meta-suite/next/kernel/internal/iface"
	"github.com/manutej/meta-suite/next/kernel/internal/kerr"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
	runs "github.com/manutej/meta-suite/next/kernel/internal/rundir"
	"github.com/manutej/meta-suite/next/kernel/internal/schema"
)

// writeTyped validates v against typ (G1) and writes it; a kernel-written file that would not validate is a
// kernel fault, never written.
func writeTyped(r *runs.Run, rel, typ string, v any) error {
	b, err := fsx.Marshal(v)
	if err != nil {
		return kerr.Wrap("internal", err)
	}
	var errs []schema.Error
	if schema.IsArray(b) {
		errs, _, err = schema.ValidateArray(typ, b)
	} else {
		errs, err = schema.Validate(typ, b)
	}
	if err != nil {
		return kerr.Wrap("internal", err)
	}
	if len(errs) > 0 {
		return kerr.New("kernel_fault", "%s would not validate as %s", rel, typ).With(map[string]any{"errors": errs})
	}
	return kerr.Wrap("io", fsx.WriteAtomic(r.Path(rel), b))
}

// GraphBuildResult is printed by `meta graph build`.
type GraphBuildResult struct {
	Regions             []string                      `json:"regions"`
	RegionsWithoutClaim []string                      `json:"regions_without_claims"`
	Files               []string                      `json:"files"`
	Objects             int                           `json:"objects"`
	Morphisms           int                           `json:"morphisms"`
	Candidates          int                           `json:"interface_candidates"`
	ShapeMetrics        map[string]model.ShapeMetrics `json:"shape_metrics"`
	MergeRule           string                        `json:"merge_rule"`
}

func cmdGraphBuild(e *env, args []string) int {
	f, err := parseFlags(args, nil, nil)
	if err != nil || len(f.pos) != 1 {
		return e.fail(usageErr("graph build"))
	}
	r, err := openRun(f.pos[0])
	if err != nil {
		return e.fail(err)
	}
	if e.cfg.Graph.MergeThreshold != 1.0 {
		return e.failRun(r, "graph build", kerr.New("config", "graph.merge_threshold = %v: this build merges entities by exact normalized label only (1.0); relation-profile merging arrives with D3's full version", e.cfg.Graph.MergeThreshold))
	}
	all, err := claims.Load(r)
	if err != nil {
		return e.failRun(r, "graph build", err)
	}
	var cover struct {
		Regions []struct {
			ID string `json:"id"`
		} `json:"regions"`
	}
	if err := r.ReadJSONFile(runs.CoverFile, &cover); err != nil {
		return e.failRun(r, "graph build", err)
	}
	byRegion := map[string][]model.Claim{}
	for _, c := range all {
		byRegion[c.Region] = append(byRegion[c.Region], c)
	}
	res := GraphBuildResult{Regions: []string{}, RegionsWithoutClaim: []string{}, Files: []string{}, ShapeMetrics: map[string]model.ShapeMetrics{},
		MergeRule: "within a region, labels that are equal after normalization are one entity; nothing is merged across regions"}
	known := map[string]bool{}
	for _, reg := range cover.Regions {
		known[reg.ID] = true
		if len(byRegion[reg.ID]) == 0 {
			res.RegionsWithoutClaim = append(res.RegionsWithoutClaim, reg.ID)
		}
	}
	for reg := range byRegion {
		if !known[reg] {
			return e.failRun(r, "graph build", kerr.New("precondition", "claims name region %q, which is not in %s", reg, runs.CoverFile))
		}
		res.Regions = append(res.Regions, reg)
	}
	sort.Strings(res.Regions)
	sort.Strings(res.RegionsWithoutClaim)
	// A rebuild replaces every region file, so a region that lost its claims does not linger.
	for _, old := range r.Glob(runs.RegionsDir + "/*.json") {
		if err := os.Remove(r.Path(old)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return e.failRun(r, "graph build", kerr.Wrap("io", err))
		}
	}
	var graphs []model.RegionGraph
	for _, reg := range res.Regions {
		g := graph.BuildRegion(reg, byRegion[reg])
		rel := runs.RegionsDir + "/" + reg + ".json"
		if err := writeTyped(r, rel, "RegionGraph", g); err != nil {
			return e.failRun(r, "graph build", err)
		}
		graphs = append(graphs, g)
		res.Files = append(res.Files, rel)
		res.Objects += len(g.Objects)
		res.Morphisms += len(g.Morphisms)
		res.ShapeMetrics[reg] = g.ShapeMetrics
	}
	cands := graph.Candidates(graphs)
	if err := writeTyped(r, runs.CandidatesFile, "Interface", cands); err != nil {
		return e.failRun(r, "graph build", err)
	}
	res.Files = append(res.Files, runs.CandidatesFile)
	res.Candidates = len(cands)
	if err := e.log(r, "graph build", true, map[string]any{"regions": res.Regions, "objects": res.Objects, "morphisms": res.Morphisms,
		"interface_candidates": res.Candidates, "regions_without_claims": res.RegionsWithoutClaim}); err != nil {
		return e.fail(err)
	}
	return e.out(res, 0)
}

// MergeOutcome is printed by `meta interface merge`.
type MergeOutcome struct {
	Merged       bool                `json:"merged"`
	Plain        bool                `json:"plain"`
	Reason       string              `json:"reason,omitempty"`
	Regions      []string            `json:"regions"`
	Interfaces   int                 `json:"interfaces"`
	Problems     []iface.Problem     `json:"problems,omitempty"`
	Collisions   []iface.Collision   `json:"collisions,omitempty"`
	Objects      int                 `json:"objects,omitempty"`
	Morphisms    int                 `json:"morphisms,omitempty"`
	Shared       int                 `json:"identifications,omitempty"`
	ShapeMetrics *model.ShapeMetrics `json:"shape_metrics,omitempty"`
	Out          string              `json:"out,omitempty"`
}

func cmdInterfaceMerge(e *env, args []string) int {
	f, err := parseFlags(args, nil, []string{"plain"})
	if err != nil || len(f.pos) != 1 {
		return e.fail(usageErr("interface merge"))
	}
	r, err := openRun(f.pos[0])
	if err != nil {
		return e.fail(err)
	}
	plain := f.set["plain"]
	var graphs []model.RegionGraph
	for _, rel := range r.Glob(runs.RegionsDir + "/*.json") {
		var g model.RegionGraph
		if err := r.ReadJSONFile(rel, &g); err != nil {
			return e.failRun(r, "interface merge", err)
		}
		graphs = append(graphs, g)
	}
	if len(graphs) == 0 {
		return e.failRun(r, "interface merge", kerr.New("precondition", "no region graphs in %s (run `meta graph build`)", runs.RegionsDir))
	}
	out := MergeOutcome{Plain: plain}
	for _, g := range graphs {
		out.Regions = append(out.Regions, g.Region)
	}
	sort.Strings(out.Regions)
	var interfaces []model.Interface
	if plain {
		if err := r.ReadJSONFile(runs.CandidatesFile, &interfaces); err != nil {
			return e.failRun(r, "interface merge", err)
		}
		for i := range interfaces {
			interfaces[i].DeclaredBy = "plain"
			if interfaces[i].Reason == "" {
				interfaces[i].Reason = "plain merge declares every candidate"
			}
		}
		if interfaces == nil {
			interfaces = []model.Interface{}
		}
		if err := writeTyped(r, runs.InterfacesFile, "Interface", interfaces); err != nil {
			return e.failRun(r, "interface merge", err)
		}
	} else {
		if !r.Exists(runs.InterfacesFile) {
			if len(graphs) > 1 {
				out.Reason = "no declared interface: write " + runs.InterfacesFile + " (or use --plain to declare every candidate)"
				return e.refuseMerge(r, out)
			}
			interfaces = []model.Interface{}
		} else {
			b, err := os.ReadFile(r.Path(runs.InterfacesFile))
			if err != nil {
				return e.failRun(r, "interface merge", kerr.Wrap("io", err))
			}
			errs, _, err := schema.ValidateArray("Interface", b)
			if err != nil {
				return e.failRun(r, "interface merge", kerr.Wrap("internal", err))
			}
			if len(errs) > 0 {
				out.Reason = runs.InterfacesFile + " does not validate"
				for _, er := range errs {
					idx := -1
					if er.Index != nil {
						idx = *er.Index
					}
					out.Problems = append(out.Problems, iface.Problem{Interface: idx, Message: er.InstanceLocation + ": " + er.Message})
				}
				return e.refuseMerge(r, out)
			}
			if err := r.ReadJSONFile(runs.InterfacesFile, &interfaces); err != nil {
				return e.failRun(r, "interface merge", err)
			}
			if len(interfaces) == 0 && len(graphs) > 1 {
				out.Reason = "no declared interface: " + runs.InterfacesFile + " is empty"
				return e.refuseMerge(r, out)
			}
		}
	}
	out.Interfaces = len(interfaces)
	var distinct []iface.Distinct
	if r.Exists(runs.DistinctFile) {
		b, err := os.ReadFile(r.Path(runs.DistinctFile))
		if err != nil {
			return e.failRun(r, "interface merge", kerr.Wrap("io", err))
		}
		if distinct, err = iface.ParseDistinct(b); err != nil {
			return e.failRun(r, "interface merge", kerr.New("invalid_json", "%s: %v", runs.DistinctFile, err))
		}
	}
	dated := map[string]bool{}
	if all, err := claims.Load(r); err == nil {
		for _, c := range all {
			if graph.Dated(c) {
				dated[c.ID] = true
			}
		}
	} else {
		return e.failRun(r, "interface merge", err)
	}
	res := iface.Merge(graphs, interfaces, distinct, plain, dated)
	if res.Graph == nil {
		out.Problems, out.Collisions = res.Problems, res.Collisions
		out.Reason = "declared interfaces cannot be used"
		if len(res.Collisions) > 0 {
			out.Reason = "the identifications would join objects recorded as distinct"
		}
		return e.refuseMerge(r, out)
	}
	if err := writeTyped(r, runs.MergedFile, "DomainGraph", res.Graph); err != nil {
		return e.failRun(r, "interface merge", err)
	}
	out.Merged, out.Out = true, runs.MergedFile
	out.Objects, out.Morphisms, out.Shared = len(res.Graph.Objects), len(res.Graph.Morphisms), res.Graph.Merge.Identifications
	sm := res.Graph.ShapeMetrics
	out.ShapeMetrics = &sm
	if err := e.log(r, "interface merge", true, map[string]any{"plain": plain, "interfaces": out.Interfaces, "objects": out.Objects,
		"morphisms": out.Morphisms, "identifications": out.Shared, "shape_metrics": sm}); err != nil {
		return e.fail(err)
	}
	return e.out(out, 0)
}

func (e *env) refuseMerge(r *runs.Run, out MergeOutcome) int {
	if err := e.log(r, "interface merge", false, map[string]any{"refused": true, "reason": out.Reason, "plain": out.Plain,
		"problems": len(out.Problems), "collisions": len(out.Collisions)}); err != nil {
		return e.fail(err)
	}
	return e.out(out, kerr.ExitVerdict)
}

func cmdFindingsVerify(e *env, args []string) int {
	f, err := parseFlags(args, nil, nil)
	if err != nil || len(f.pos) != 1 {
		return e.fail(usageErr("findings verify"))
	}
	r, err := openRun(f.pos[0])
	if err != nil {
		return e.fail(err)
	}
	all, err := claims.Load(r)
	if err != nil {
		return e.failRun(r, "findings verify", err)
	}
	p := findings.Params{MinIndependent: e.cfg.Lens.Pullback.MinIndependent, NumericRelTol: e.cfg.Lens.Pullback.NumericRelTol,
		StripPrefixes: e.cfg.Independence.StripPrefixes, AsOfWindow: e.cfg.Lens.Pullback.AsOfWindow}
	fs, conf := findings.Verify(all, p)
	for _, fd := range fs {
		b, _ := fsx.Compact(fd)
		if errs, err := schema.Validate("Finding", b); err != nil || len(errs) > 0 {
			return e.failRun(r, "findings verify", kerr.New("kernel_fault", "finding %s would not validate", fd.ID).With(map[string]any{"errors": errs}))
		}
	}
	if fs == nil {
		fs = []model.Finding{}
	}
	if err := fsx.WriteJSONL(r.Path(runs.FindingsFile), fs); err != nil {
		return e.failRun(r, "findings verify", kerr.Wrap("io", err))
	}
	if err := writeTyped(r, runs.ConflictsFile, "Conflicts", conf); err != nil {
		return e.failRun(r, "findings verify", err)
	}
	s := findings.Summary{Claims: len(all), Findings: len(fs), ByKind: map[string]int{}, Conflicts: len(conf.Conflicts),
		FindingsFile: runs.FindingsFile, ConflictFile: runs.ConflictsFile}
	keys := map[model.Key]bool{}
	for _, fd := range fs {
		s.ByKind[fd.Kind]++
		keys[*fd.Key] = true
	}
	s.Keys = len(keys)
	s.Rule = ruleText(p.MinIndependent, p.AsOfWindow)
	if err := e.log(r, "findings verify", true, map[string]any{"findings": s.Findings, "by_kind": s.ByKind, "conflicts": s.Conflicts,
		"min_independent": p.MinIndependent, "numeric_rel_tol": p.NumericRelTol, "as_of_window": p.AsOfWindow}); err != nil {
		return e.fail(err)
	}
	return e.out(s, 0)
}

func ruleText(n int, window string) string {
	s := "verified = at least " + itoa(n) + " distinct independence keys agree on subject, relation, object and value"
	if window != "none" && window != "" {
		s += " within one as_of " + window
	}
	return s + "; otherwise single_source"
}
