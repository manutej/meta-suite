package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/manutej/meta-suite/next/kernel/internal/cover"
	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/kerr"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
	runs "github.com/manutej/meta-suite/next/kernel/internal/rundir"
	"github.com/manutej/meta-suite/next/kernel/internal/schema"
)

// inferType maps well-known run file names to their type when --type is omitted.
func inferType(path string) string {
	slash := filepath.ToSlash(path)
	base := filepath.Base(path)
	switch {
	case strings.Contains(slash, "graph/regions/"):
		return "RegionGraph"
	case base == "run.json":
		return "RunRecord"
	case base == "brief.json":
		return "Brief"
	case base == "spec.json":
		return "Spec"
	case base == "cover.json":
		return "Cover"
	case base == "claims.jsonl" || strings.Contains(slash, "claims/"):
		return "Claim"
	case base == "merged.json":
		return "DomainGraph"
	case base == "interfaces.json" || base == "interface_candidates.json":
		return "Interface"
	case base == "findings.jsonl":
		return "Finding"
	case base == "plan.json":
		return "ArtifactPlan"
	case base == "section_map.json":
		return "SectionMap"
	case base == "ledger.jsonl":
		return "LedgerEvent"
	case strings.HasPrefix(base, "edits."):
		return "EditScript"
	case base == "coverage.json":
		return "Coverage"
	case base == "probes.json":
		return "Probes"
	case base == "conflicts.json":
		return "Conflicts"
	case base == "how_checked.json":
		return "HowChecked"
	case base == "retractions.jsonl":
		return "Retraction"
	}
	return ""
}

// ValidateResult is printed by `meta validate`.
type ValidateResult struct {
	File    string         `json:"file"`
	Type    string         `json:"type"`
	Valid   bool           `json:"valid"`
	Records int            `json:"records"`
	Errors  []schema.Error `json:"errors"`
	Hash    string         `json:"hash,omitempty"`
	Cover   *cover.Report  `json:"cover,omitempty"`
	Xref    *Xref          `json:"xref,omitempty"`
}

// cmdValidate is G1. A .jsonl file (or a JSON array for a non-array type) is validated record by record.
// Spec: the kernel computes the content hash (sha256 of the canonical JSON without "hash"); a supplied hash
// must match; for a run's spec.json the hash is recorded in run.json (spec_hash) and the ledger (D-087).
// Cover: also runs the D1/D-042 checks and prints the computed overlaps. EditScript, SectionMap, Coverage and
// Probes: inside a run (or with --run), the ids they name are checked against the run's files (D-087). When
// the file lies in a run directory, the verdict is recorded in run.json and the ledger.
func cmdValidate(e *env, args []string) int {
	f, err := parseFlags(args, []string{"type", "mode", "run"}, nil)
	if err != nil || len(f.pos) != 1 {
		return e.fail(usageErr("validate"))
	}
	file := f.pos[0]
	typ := f.vals["type"]
	if typ == "" {
		if typ = inferType(file); typ == "" {
			return e.fail(kerr.New("usage", "cannot infer the type of %s; pass --type (one of %s)", file, strings.Join(schema.Types(), ", ")))
		}
	}
	canon, err := schema.Canonical(typ)
	if err != nil {
		return e.fail(kerr.New("usage", "%v", err))
	}
	doc, err := os.ReadFile(file)
	if err != nil {
		return e.fail(kerr.New("io", "%v", err))
	}
	runDir := runs.Find(file)
	var r *runs.Run
	if runDir != "" {
		if r, err = runs.Open(runDir); err != nil {
			r = nil // a run.json we cannot read: validate anyway, record nothing
		}
	}
	xr := r // the run whose files ids are checked against
	if d := f.vals["run"]; d != "" {
		if xr, err = openRun(d); err != nil {
			return e.fail(err)
		}
	}
	res := ValidateResult{File: file, Type: canon, Errors: []schema.Error{}}
	switch {
	case strings.HasSuffix(strings.ToLower(file), ".jsonl"):
		lines, err := fsx.ReadLines(file)
		if err != nil {
			return e.fail(kerr.Wrap("io", err))
		}
		for _, l := range lines {
			errs, err := schema.Validate(canon, l.Raw)
			if err != nil {
				return e.fail(kerr.Wrap("internal", err))
			}
			for _, er := range errs {
				er.Line = l.N
				res.Errors = append(res.Errors, er)
			}
		}
		res.Records = len(lines)
	case schema.IsArray(doc):
		errs, n, err := schema.ValidateArray(canon, doc)
		if err != nil {
			return e.fail(kerr.Wrap("internal", err))
		}
		res.Errors, res.Records = append(res.Errors, errs...), n
	default:
		errs, err := schema.Validate(canon, doc)
		if err != nil {
			return e.fail(kerr.Wrap("internal", err))
		}
		res.Errors, res.Records = append(res.Errors, errs...), 1
		if len(errs) == 0 && canon == "Spec" {
			h, supplied, err := specHash(doc)
			if err != nil {
				return e.fail(kerr.Wrap("internal", err))
			}
			res.Hash = h
			if supplied != "" && supplied != h {
				res.Errors = append(res.Errors, schema.Error{InstanceLocation: "/hash", KeywordLocation: "kernel/spec_hash",
					Message: "hash does not match the content (" + h + "); omit it and the kernel computes it"})
			}
		}
		if len(errs) == 0 && canon == "Cover" {
			maxRegions := 0
			mode := f.vals["mode"]
			if mode == "" && r != nil {
				mode = r.Record.Mode
			}
			if mode != "" {
				row, err := e.cfg.ResolveMode(mode)
				if err != nil {
					return e.fail(kerr.New("usage", "%v", err))
				}
				maxRegions = e.cfg.Modes[row].MaxRegions
			}
			rep, err := cover.Check(doc, maxRegions)
			if err != nil {
				return e.fail(kerr.Wrap("internal", err))
			}
			res.Cover = rep
			for _, is := range rep.Issues {
				res.Errors = append(res.Errors, schema.Error{InstanceLocation: "/regions", KeywordLocation: "kernel/cover/" + is.Check, Message: is.Message})
			}
		}
	}
	if len(res.Errors) == 0 && xr != nil {
		xerrs, info := crossCheck(xr, canon, doc)
		res.Errors, res.Xref = append(res.Errors, xerrs...), info
	}
	schema.Sort(res.Errors)
	res.Valid = len(res.Errors) == 0
	if r != nil {
		rel := r.Rel(file)
		verdict := "pass"
		if !res.Valid {
			verdict = "fail"
		}
		if err := r.RecordGate(model.GateResult{Gate: "G1", Verdict: verdict, At: runs.Stamp(runs.Now()), Summary: canon, Files: []string{rel},
			InputsHash: fsx.SHA256(doc)}); err != nil {
			return e.fail(err)
		}
		data := map[string]any{"file": rel, "type": canon, "valid": res.Valid, "records": res.Records, "errors": len(res.Errors), "sha256": fsx.SHA256(doc)}
		if res.Hash != "" {
			data["spec_hash"] = res.Hash
			if res.Valid && rel == runs.SpecFile {
				if err := r.Update(func(rec *model.RunRecord) error { rec.SpecHash = res.Hash; return nil }); err != nil {
					return e.fail(err)
				}
				data["spec_hash_recorded"] = runs.RunFile
			}
		}
		if err := e.log(r, "validate", res.Valid, data); err != nil {
			return e.fail(err)
		}
	}
	if !res.Valid {
		return e.out(res, kerr.ExitVerdict)
	}
	return e.out(res, 0)
}

// specHash returns sha256 of the canonical JSON of a spec without its "hash" field, and the supplied hash.
func specHash(doc []byte) (string, string, error) {
	dec := json.NewDecoder(bytes.NewReader(doc))
	dec.UseNumber()
	var m map[string]any
	if err := dec.Decode(&m); err != nil {
		return "", "", err
	}
	supplied, _ := m["hash"].(string)
	delete(m, "hash")
	b, err := fsx.Compact(m)
	if err != nil {
		return "", "", err
	}
	return fsx.SHA256(b), supplied, nil
}
