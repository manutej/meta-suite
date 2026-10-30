// Package claims implements D2's typed writes (`meta claim add`) and the stage-boundary merge
// (`meta claims merge`). A claim is stamped from a verified, kernel-issued stage token, unit-normalized
// (D-043), checked against the cover's relation vocabulary and validated by G1 before it is appended to its
// agent's shard under a file lock.
package claims

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/manutej/meta-suite/next/kernel/internal/config"
	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/ids"
	"github.com/manutej/meta-suite/next/kernel/internal/kerr"
	"github.com/manutej/meta-suite/next/kernel/internal/lock"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
	"github.com/manutej/meta-suite/next/kernel/internal/norm"
	runs "github.com/manutej/meta-suite/next/kernel/internal/rundir"
	"github.com/manutej/meta-suite/next/kernel/internal/schema"
	"github.com/manutej/meta-suite/next/kernel/internal/stage"
)

// Rejection is a verdict on an input claim (exit 1), with the schema-style errors the agent can act on.
type Rejection struct {
	Valid  bool           `json:"valid"`
	Errors []schema.Error `json:"errors"`
}

// AddResult is printed on success.
type AddResult struct {
	ID         string            `json:"id"`
	Region     string            `json:"region"`
	Shard      string            `json:"shard"`
	Normalized *model.Normalized `json:"normalized,omitempty"`
	Notes      []string          `json:"notes,omitempty"`
}

// Vocab is the relation vocabulary from cover.json.
type Vocab map[string]struct {
	ValueType     string `json:"value_type"`
	UnitDimension string `json:"unit_dimension"`
}

// LoadVocab reads cover.json's relation_vocab.
func LoadVocab(r *runs.Run) (Vocab, error) {
	v, _, err := loadCover(r)
	return v, err
}

// loadCover reads cover.json's relation_vocab and its shared label list (normalized; nil when absent).
func loadCover(r *runs.Run) (Vocab, map[string]bool, error) {
	var cover struct {
		RelationVocab []struct {
			Name          string `json:"name"`
			ValueType     string `json:"value_type"`
			UnitDimension string `json:"unit_dimension"`
		} `json:"relation_vocab"`
		Labels []string `json:"labels"`
	}
	if err := r.ReadJSONFile(runs.CoverFile, &cover); err != nil {
		return nil, nil, err
	}
	var labels map[string]bool
	if len(cover.Labels) > 0 {
		labels = map[string]bool{}
		for _, l := range cover.Labels {
			labels[norm.Label(l)] = true
		}
	}
	v := Vocab{}
	for _, rel := range cover.RelationVocab {
		v[rel.Name] = struct {
			ValueType     string `json:"value_type"`
			UnitDimension string `json:"unit_dimension"`
		}{rel.ValueType, rel.UnitDimension}
	}
	return v, labels, nil
}

// Add validates and appends one claim read from input. It returns (result, nil, nil) on success,
// (nil, rejection, nil) for an invalid claim, or an error when the command cannot run.
func Add(r *runs.Run, cfg *config.Config, token string, input []byte) (*AddResult, *Rejection, error) {
	t, err := stage.Decode(strings.TrimSpace(token))
	if err != nil {
		return nil, nil, err
	}
	if err := stage.Verify(r, t); err != nil {
		return nil, nil, err
	}
	allowed := false
	for _, s := range cfg.Claims.AllowedStages {
		allowed = allowed || s == t.Stage
	}
	if !allowed {
		return nil, nil, kerr.New("bad_token", "stage %s may not write claims (claims.allowed_stages = %v)", t.Stage, cfg.Claims.AllowedStages)
	}
	// G1 on the input shape (no id, no stage_token).
	errs, err := schema.Validate("ClaimInput", input)
	if err != nil {
		return nil, nil, kerr.Wrap("internal", err)
	}
	if len(errs) > 0 {
		return nil, &Rejection{Errors: errs}, nil
	}
	var c model.Claim
	dec := json.NewDecoder(bytes.NewReader(input))
	if err := dec.Decode(&c); err != nil {
		return nil, &Rejection{Errors: []schema.Error{{Message: "not a claim: " + err.Error()}}}, nil
	}
	var notes []string
	if c.Region != "" && c.Region != t.Region {
		notes = append(notes, fmt.Sprintf("region %q replaced by the stage token's region %q", c.Region, t.Region))
	}
	c.Region = t.Region
	vocab, labels, err := loadCover(r)
	if err != nil {
		return nil, nil, err
	}
	if rej := checkVocab(&c, vocab); rej != nil {
		return nil, rej, nil
	}
	known, rej := Normalize(&c, cfg, vocab)
	if rej != nil {
		return nil, rej, nil
	}
	if !known {
		notes = append(notes, fmt.Sprintf("unit %q is not in the canonical-unit table; kept as given", c.Unit))
	}
	if labels != nil { // D-088: the cover's shared label list is how independent regions agree on keys
		for _, f := range [][2]string{{"subject", c.Subject}, {"object", c.Object}} {
			if !labels[norm.Label(f[1])] {
				notes = append(notes, fmt.Sprintf("%s %q is not in cover.json labels; use a listed label when one names the same thing", f[0], f[1]))
			}
		}
	}
	c.ID = ids.New()
	c.StageToken = t
	line, err := fsx.Compact(c)
	if err != nil {
		return nil, nil, kerr.Wrap("internal", err)
	}
	if errs, err := schema.Validate("Claim", line); err != nil || len(errs) > 0 {
		if err != nil {
			return nil, nil, kerr.Wrap("internal", err)
		}
		return nil, &Rejection{Errors: errs}, nil
	}
	shard := runs.ClaimsDir + "/" + t.Agent + ".jsonl"
	path := r.Path(shard)
	err = lock.With(path, func() error {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		if _, err := f.Write(append(line, '\n')); err != nil {
			f.Close()
			return err
		}
		return f.Close()
	})
	if err != nil {
		return nil, nil, kerr.Wrap("io", err)
	}
	sort.Strings(notes)
	return &AddResult{ID: c.ID, Region: c.Region, Shard: shard, Normalized: c.Normalized, Notes: notes}, nil, nil
}

func checkVocab(c *model.Claim, v Vocab) *Rejection {
	rel, ok := v[c.Relation]
	if !ok {
		names := make([]string, 0, len(v))
		for k := range v {
			names = append(names, k)
		}
		sort.Strings(names)
		return &Rejection{Errors: []schema.Error{{InstanceLocation: "/relation", KeywordLocation: "cover/relation_vocab",
			Message: fmt.Sprintf("relation %q is not in the cover's relation vocabulary %v; use the closest one", c.Relation, names)}}}
	}
	if rel.ValueType != c.ValueType {
		return &Rejection{Errors: []schema.Error{{InstanceLocation: "/value_type", KeywordLocation: "cover/relation_vocab/" + c.Relation,
			Message: fmt.Sprintf("relation %q has value_type %q in the vocabulary, claim says %q", c.Relation, rel.ValueType, c.ValueType)}}}
	}
	if c.ValueType == "relative" {
		if _, ok := c.Value.(float64); c.Value != nil && !ok {
			return &Rejection{Errors: []schema.Error{{InstanceLocation: "/value", KeywordLocation: "cover/relation_vocab/" + c.Relation,
				Message: "a relative value must be a number"}}}
		}
	}
	return nil
}

// Normalize converts a numeric or relative claim's value to its dimension's canonical unit. The dimension is
// the vocabulary's unit_dimension when it names a configured dimension, otherwise the dimension whose table
// holds the claim's unit. A unit that belongs to a different dimension than the vocabulary declares is
// rejected. It returns known=false for a unit in no table (kept as given, Normalized nil); a canonical unit
// such as "%" is known and left unchanged (D-087).
//
// Code-only dimensions (config unit_codes, e.g. currency: ISO 4217 minus config's exclude list, D-097) are
// never converted. Codes are upper-cased on any relation ("eur" -> "EUR"). A relation declared in a code
// dimension needs a unit; a code on a relation of a convertible dimension is rejected.
func Normalize(c *model.Claim, cfg *config.Config, v Vocab) (bool, *Rejection) {
	if c.ValueType != "numeric" && c.ValueType != "relative" {
		return true, nil
	}
	declared := v[c.Relation].UnitDimension
	if uc, ok := cfg.UnitCodes[declared]; ok {
		up := strings.ToUpper(strings.TrimSpace(c.Unit))
		switch {
		case up == "":
			return false, &Rejection{Errors: []schema.Error{{InstanceLocation: "/unit", KeywordLocation: "config/unit_codes/" + declared,
				Message: fmt.Sprintf("relation %q is declared as %s: unit is required, a code such as EUR", c.Relation, declared)}}}
		case uc.Match(up):
			c.Unit = up
			return true, nil
		}
		if d := tableDim(cfg, c.Unit); d != "" {
			return true, &Rejection{Errors: []schema.Error{{InstanceLocation: "/unit", KeywordLocation: "config/units",
				Message: fmt.Sprintf("unit %q measures %s, but relation %q is declared as %s", c.Unit, d, c.Relation, declared)}}}
		}
		return false, &Rejection{Errors: []schema.Error{{InstanceLocation: "/unit", KeywordLocation: "config/unit_codes/" + declared,
			Message: fmt.Sprintf("relation %q is declared as %s: unit %q is not a code of it (pattern %s, list %q); use a code such as EUR, amount in base units", c.Relation, declared, c.Unit, uc.Pattern, uc.List)}}}
	}
	if c.Unit == "" {
		return true, nil
	}
	val, isNum := c.Value.(float64)
	if !isNum {
		return true, nil
	}
	found, coded := tableDim(cfg, c.Unit), false
	if found == "" {
		up := strings.ToUpper(strings.TrimSpace(c.Unit))
		for _, d := range sortedDims(cfg) {
			if uc, ok := cfg.UnitCodes[d]; ok && uc.Match(up) {
				found, coded = d, true
				c.Unit = up // "eur" groups with "EUR" on any relation
				break
			}
		}
	}
	if found == "" {
		return false, nil
	}
	if _, inUnits := cfg.Units[declared]; inUnits && declared != found {
		return true, &Rejection{Errors: []schema.Error{{InstanceLocation: "/unit", KeywordLocation: "config/units",
			Message: fmt.Sprintf("unit %q measures %s, but relation %q is declared as %s", c.Unit, found, c.Relation, declared)}}}
	}
	if coded {
		return true, nil
	}
	dim := cfg.Units[found]
	if c.Unit == dim.Canonical {
		return true, nil
	}
	c.Normalized = &model.Normalized{Dimension: found, FromUnit: c.Unit, FromValue: val}
	c.Value = val * dim.Factors[c.Unit]
	c.Unit = dim.Canonical
	return true, nil
}

// tableDim returns the convertible dimension whose factor table holds unit (case-sensitive), or "".
func tableDim(cfg *config.Config, unit string) string {
	dims := make([]string, 0, len(cfg.Units))
	for d := range cfg.Units {
		dims = append(dims, d)
	}
	sort.Strings(dims)
	for _, d := range dims {
		if _, ok := cfg.Units[d].Factors[unit]; ok {
			return d
		}
	}
	return ""
}

func sortedDims(cfg *config.Config) []string {
	dims := make([]string, 0, len(cfg.UnitCodes))
	for d := range cfg.UnitCodes {
		dims = append(dims, d)
	}
	sort.Strings(dims)
	return dims
}

// LineError locates a bad line in a shard.
type LineError struct {
	File   string         `json:"file"`
	Line   int            `json:"line"`
	Errors []schema.Error `json:"errors"`
}

// MergeResult is printed by `meta claims merge`.
type MergeResult struct {
	Claims           int            `json:"claims"`
	Retracted        int            `json:"retracted"`
	Shards           []string       `json:"shards"`
	ByRegion         map[string]int `json:"by_region"`
	IndependenceKeys int            `json:"independence_keys"`
	Out              string         `json:"out"`
}

// MergeRejection is the verdict when a shard holds invalid or forged lines.
type MergeRejection struct {
	Merged bool        `json:"merged"`
	Bad    []LineError `json:"bad_lines"`
}

// ReadShards reads and checks every claim in claims/*.jsonl: G1, a kernel-issued token, region = token region.
func ReadShards(r *runs.Run) ([]model.Claim, []string, []LineError, error) {
	shards := r.Glob(runs.ClaimsDir + "/*.jsonl")
	issued, err := stage.Issued(r)
	if err != nil {
		return nil, nil, nil, err
	}
	var all []model.Claim
	var bad []LineError
	seen := map[string]string{}
	for _, sh := range shards {
		lines, err := fsx.ReadLines(r.Path(sh))
		if err != nil {
			return nil, nil, nil, kerr.Wrap("io", err)
		}
		for _, l := range lines {
			errs, err := schema.Validate("Claim", l.Raw)
			if err != nil {
				return nil, nil, nil, kerr.Wrap("internal", err)
			}
			if len(errs) > 0 {
				bad = append(bad, LineError{File: sh, Line: l.N, Errors: errs})
				continue
			}
			var c model.Claim
			json.Unmarshal(l.Raw, &c)
			if tok, ok := issued[c.StageToken.TID]; !ok || tok != c.StageToken {
				bad = append(bad, LineError{File: sh, Line: l.N, Errors: []schema.Error{{InstanceLocation: "/stage_token",
					Message: "stage token was not issued by the kernel for this run, or was edited"}}})
				continue
			}
			if c.Region != c.StageToken.Region {
				bad = append(bad, LineError{File: sh, Line: l.N, Errors: []schema.Error{{InstanceLocation: "/region",
					Message: "region differs from the stage token's region"}}})
				continue
			}
			if filepath.Base(sh) != c.StageToken.Agent+".jsonl" {
				bad = append(bad, LineError{File: sh, Line: l.N, Errors: []schema.Error{{InstanceLocation: "/stage_token/agent",
					Message: "claim is in another agent's shard"}}})
				continue
			}
			if prev, dup := seen[c.ID]; dup {
				bad = append(bad, LineError{File: sh, Line: l.N, Errors: []schema.Error{{InstanceLocation: "/id",
					Message: "duplicate claim id (also in " + prev + ")"}}})
				continue
			}
			seen[c.ID] = sh
			all = append(all, c)
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })
	return all, shards, bad, nil
}

// Merge writes claims.jsonl from the shards, sorted by claim id. It refuses (returns a rejection) if any
// shard line is invalid or forged.
func Merge(r *runs.Run, cfg *config.Config, independenceKey func(string) string) (*MergeResult, *MergeRejection, error) {
	all, shards, bad, err := ReadShards(r)
	if err != nil {
		return nil, nil, err
	}
	if len(shards) == 0 {
		return nil, nil, kerr.New("precondition", "no claim shards in %s/ (gather writes them with `meta claim add`)", runs.ClaimsDir)
	}
	if len(bad) > 0 {
		return nil, &MergeRejection{Bad: bad}, nil
	}
	tomb, err := Retracted(r)
	if err != nil {
		return nil, nil, err
	}
	live := all[:0:0]
	for _, c := range all {
		if _, gone := tomb[c.ID]; !gone {
			live = append(live, c)
		}
	}
	retracted := len(all) - len(live)
	all = live
	if err := fsx.WriteJSONL(r.Path(runs.ClaimsFile), all); err != nil {
		return nil, nil, kerr.Wrap("io", err)
	}
	by := map[string]int{}
	keys := map[string]bool{}
	for _, c := range all {
		by[c.Region]++
		keys[independenceKey(c.Source.IndependenceKey)] = true
	}
	return &MergeResult{Claims: len(all), Retracted: retracted, Shards: shards, ByRegion: by, IndependenceKeys: len(keys), Out: runs.ClaimsFile}, nil, nil
}

// Load reads claims.jsonl (G1 on every line), skipping retracted claims even when claims.jsonl predates the
// retraction (G2 still flags that stale merge).
func Load(r *runs.Run) ([]model.Claim, error) {
	tomb, err := Retracted(r)
	if err != nil {
		return nil, err
	}
	lines, err := fsx.ReadLines(r.Path(runs.ClaimsFile))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, kerr.New("precondition", "%s is missing (run `meta claims merge` first)", runs.ClaimsFile)
		}
		return nil, kerr.Wrap("io", err)
	}
	out := make([]model.Claim, 0, len(lines))
	for _, l := range lines {
		errs, err := schema.Validate("Claim", l.Raw)
		if err != nil {
			return nil, kerr.Wrap("internal", err)
		}
		if len(errs) > 0 {
			return nil, kerr.New("invalid_claims", "%s line %d does not validate", runs.ClaimsFile, l.N).With(map[string]any{"errors": errs})
		}
		var c model.Claim
		json.Unmarshal(l.Raw, &c)
		if _, gone := tomb[c.ID]; gone {
			continue
		}
		out = append(out, c)
	}
	return out, nil
}

// Retracted reads retractions.jsonl (absent: none), keyed by claim id. Every line is validated (G1).
func Retracted(r *runs.Run) (map[string]model.Retraction, error) {
	out := map[string]model.Retraction{}
	lines, err := fsx.ReadLines(r.Path(runs.RetractFile))
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, kerr.Wrap("io", err)
	}
	for _, l := range lines {
		if errs, err := schema.Validate("Retraction", l.Raw); err != nil || len(errs) > 0 {
			if err != nil {
				return nil, kerr.Wrap("internal", err)
			}
			return nil, kerr.New("invalid_json", "%s line %d does not validate (it is written only by `meta claim retract`)", runs.RetractFile, l.N).With(map[string]any{"errors": errs})
		}
		var t model.Retraction
		json.Unmarshal(l.Raw, &t)
		out[t.ID] = t
	}
	return out, nil
}

// Retract appends a tombstone for claim id (D-087). The claim must be in a shard. Retracting it again with
// the same reason is a no-op that returns the existing tombstone and already=true (D-097); a different reason
// is refused. Nothing is deleted: shards stay append-only, and merge, graph, findings, G2 and report skip the
// claim.
func Retract(r *runs.Run, id, reason, stage string) (t *model.Retraction, already bool, err error) {
	if !ids.Valid(id) {
		return nil, false, kerr.New("usage", "--id %q is not a claim id (ULID)", id)
	}
	if strings.TrimSpace(reason) == "" {
		return nil, false, kerr.New("usage", "--reason is required: say in one sentence why the claim is withdrawn")
	}
	all, _, _, err := ReadShards(r)
	if err != nil {
		return nil, false, err
	}
	found := false
	for _, c := range all {
		found = found || c.ID == id
	}
	if !found {
		return nil, false, kerr.New("precondition", "claim %s is not in any shard of this run", id)
	}
	t = &model.Retraction{ID: id, Reason: strings.TrimSpace(reason), At: runs.Stamp(runs.Now()), Stage: stage}
	line, err := fsx.Compact(t)
	if err != nil {
		return nil, false, kerr.Wrap("internal", err)
	}
	if errs, err := schema.Validate("Retraction", line); err != nil || len(errs) > 0 {
		return nil, false, kerr.New("kernel_fault", "retraction would not validate").With(map[string]any{"errors": errs})
	}
	path := r.Path(runs.RetractFile)
	err = lock.With(path, func() error {
		tomb, err := Retracted(r)
		if err != nil {
			return err
		}
		if prev, dup := tomb[id]; dup {
			if prev.Reason == t.Reason {
				t, already = &prev, true
				return nil
			}
			return kerr.New("exists", "claim %s is already retracted, with another reason: %q", id, prev.Reason)
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return kerr.Wrap("io", err)
		}
		if _, err := f.Write(append(line, '\n')); err != nil {
			f.Close()
			return kerr.Wrap("io", err)
		}
		return kerr.Wrap("io", f.Close())
	})
	if err != nil {
		return nil, false, kerr.Wrap("io", err)
	}
	return t, already, nil
}
