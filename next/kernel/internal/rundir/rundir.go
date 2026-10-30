// Package rundir manages a run directory: runs/<run_id>/ with run.json (the RunRecord), the ledger, claim
// shards, graphs, findings and gate results. run.json is only changed under its file lock.
package rundir

import (
	"encoding/json"
	"errors"

	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/ids"
	"github.com/manutej/meta-suite/next/kernel/internal/kerr"
	"github.com/manutej/meta-suite/next/kernel/internal/lock"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
	"github.com/manutej/meta-suite/next/kernel/internal/schema"
)

// File names inside a run directory.
const (
	RunFile        = "run.json"
	BriefFile      = "brief.json"
	SpecFile       = "spec.json"
	CoverFile      = "cover.json"
	LedgerFile     = "ledger.jsonl"
	ClaimsDir      = "claims"
	ClaimsFile     = "claims.jsonl"
	TokensDir      = "tokens"
	RegionsDir     = "graph/regions"
	CandidatesFile = "graph/interface_candidates.json"
	InterfacesFile = "graph/interfaces.json"
	DistinctFile   = "graph/distinct.json"
	MergedFile     = "graph/merged.json"
	FindingsFile   = "findings.jsonl"
	ConflictsFile  = "conflicts.json"
	HowCheckedFile = "how_checked.json"
	PlanFile       = "plan.json"
	SectionMapFile = "section_map.json"
	CoverageFile   = "coverage.json"
	ProbesFile     = "probes.json"
	RetractFile    = "retractions.jsonl"
	ArtifactPage   = "artifact/index.html"
)

// WorkDirs are created by `meta run init`: every folder an op writes into, so a shell redirect such as
// `> "$RUN/gates/G2.json"` never fails before the kernel runs (D-087).
var WorkDirs = []string{ClaimsDir, TokensDir, "gates", "prompts", "review", "artifact", "artifact/shots", RegionsDir}

// Now is the clock; tests may replace it.
var Now = func() time.Time { return time.Now().UTC() }

// Stamp formats t as RFC 3339 UTC with milliseconds.
func Stamp(t time.Time) string { return t.UTC().Format("2006-01-02T15:04:05.000Z") }

// Run is an opened run directory.
type Run struct {
	Dir    string
	Record model.RunRecord
}

// Path joins a run-relative path.
func (r *Run) Path(rel string) string { return filepath.Join(r.Dir, filepath.FromSlash(rel)) }

// Open loads run.json from dir.
func Open(dir string) (*Run, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, kerr.Wrap("io", err)
	}
	var rec model.RunRecord
	if err := fsx.ReadJSON(filepath.Join(abs, RunFile), &rec); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, kerr.New("not_a_run", "%s has no %s (create one with `meta run init`)", dir, RunFile)
		}
		return nil, kerr.Wrap("io", err)
	}
	if !ids.Valid(rec.RunID) {
		return nil, kerr.New("not_a_run", "%s: run_id %q is not a ULID", filepath.Join(dir, RunFile), rec.RunID)
	}
	return &Run{Dir: abs, Record: rec}, nil
}

// Find walks up from path to the nearest directory holding run.json. It returns "" when there is none.
func Find(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	if st, err := os.Stat(abs); err == nil && !st.IsDir() {
		abs = filepath.Dir(abs)
	}
	for {
		if fsx.Exists(filepath.Join(abs, RunFile)) {
			return abs
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return ""
		}
		abs = parent
	}
}

// Create makes runs/<ULID>/ under root and writes the RunRecord skeleton (validated first).
func Create(root string, rec model.RunRecord) (*Run, error) {
	rec.RunID = ids.New()
	dir := filepath.Join(root, "runs", rec.RunID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, kerr.Wrap("io", err)
	}
	for _, d := range WorkDirs {
		if err := os.MkdirAll(filepath.Join(dir, filepath.FromSlash(d)), 0o755); err != nil {
			return nil, kerr.Wrap("io", err)
		}
	}
	r := &Run{Dir: dir, Record: rec}
	if err := r.write(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Run) write() error {
	b, err := fsx.Marshal(r.Record)
	if err != nil {
		return kerr.Wrap("internal", err)
	}
	errs, err := schema.Validate("RunRecord", b)
	if err != nil {
		return kerr.Wrap("internal", err)
	}
	if len(errs) > 0 {
		return kerr.New("kernel_fault", "run.json would not validate").With(map[string]any{"errors": errs})
	}
	return kerr.Wrap("io", fsx.WriteAtomic(r.Path(RunFile), b))
}

// Update re-reads run.json under its lock, applies fn and writes it back (validated).
func (r *Run) Update(fn func(*model.RunRecord) error) error {
	return lock.With(r.Path(RunFile), func() error {
		var rec model.RunRecord
		if err := fsx.ReadJSON(r.Path(RunFile), &rec); err != nil {
			return kerr.Wrap("io", err)
		}
		if err := fn(&rec); err != nil {
			return err
		}
		r.Record = rec
		return r.write()
	})
}

// RecordGate appends a gate result to run.json.
func (r *Run) RecordGate(g model.GateResult) error {
	return r.Update(func(rec *model.RunRecord) error {
		rec.GateResults = append(rec.GateResults, g)
		return nil
	})
}

// HashFiles returns the sha256 over the run-relative paths and contents of the given files (missing files
// hash as absent), so a gate can tell later whether its inputs changed.
func (r *Run) HashFiles(rels []string) string {
	sorted := append([]string(nil), rels...)
	sort.Strings(sorted)
	var buf []byte
	for _, rel := range sorted {
		b, err := os.ReadFile(r.Path(rel))
		if err != nil {
			buf = append(buf, []byte(rel+"\x00absent\x00")...)
			continue
		}
		buf = append(buf, []byte(rel+"\x00"+fsx.SHA256(b)+"\x00")...)
	}
	return fsx.SHA256(buf)
}

// Glob returns run-relative paths matching pattern (run-relative), sorted.
func (r *Run) Glob(pattern string) []string {
	m, _ := filepath.Glob(r.Path(pattern))
	out := make([]string, 0, len(m))
	for _, p := range m {
		rel, err := filepath.Rel(r.Dir, p)
		if err == nil {
			out = append(out, filepath.ToSlash(rel))
		}
	}
	sort.Strings(out)
	return out
}

// ReadJSONFile decodes a run-relative JSON file.
func (r *Run) ReadJSONFile(rel string, v any) error {
	if err := fsx.ReadJSON(r.Path(rel), v); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return kerr.New("precondition", "%s is missing in %s", rel, r.Dir)
		}
		var se *json.SyntaxError
		if errors.As(err, &se) {
			return kerr.New("invalid_json", "%v", err)
		}
		return kerr.Wrap("io", err)
	}
	return nil
}

// Exists reports whether a run-relative path exists.
func (r *Run) Exists(rel string) bool { return fsx.Exists(r.Path(rel)) }

// Rel returns path relative to the run dir when it is inside it, else path unchanged.
func (r *Run) Rel(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	rel, err := filepath.Rel(r.Dir, abs)
	if err != nil || len(rel) >= 2 && rel[:2] == ".." {
		return path
	}
	return filepath.ToSlash(rel)
}

// MustNotExist returns an error if a run-relative path exists (used to refuse overwriting).
func (r *Run) MustNotExist(rel string) error {
	if r.Exists(rel) {
		return kerr.New("exists", "%s already exists", rel)
	}
	return nil
}
