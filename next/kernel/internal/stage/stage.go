// Package stage issues stage tokens (O1). A token is base64url(JSON{tid, run, stage, region, agent, model,
// prompt_hash}). The kernel keeps a copy of every token it issued in tokens/<tid>.json, so a token that was
// not issued by the kernel, or was edited, is refused (forged provenance) and G2 can re-check every claim.
package stage

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"sort"

	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/ids"
	"github.com/manutej/meta-suite/next/kernel/internal/kerr"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
	runs "github.com/manutej/meta-suite/next/kernel/internal/rundir"
)

// Stages are the stage names a token can carry (spec/v2/03-architecture.md, Pipeline).
var Stages = []string{"I1", "I2", "I3", "D1", "D2", "D3", "L1", "L2", "L3", "L4", "G3", "G4", "A1", "A2", "R1", "R2", "S1", "S2"}

var (
	agentRe  = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)
	regionRe = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)
)

// Options are the inputs of `meta stage begin`.
type Options struct {
	Stage, Region, Agent, Model, PromptFile string
}

// Begin issues a token, stores it, records prompt hash and model in run.json. It returns the token and the
// encoded string.
func Begin(r *runs.Run, o Options) (model.StageToken, string, error) {
	var zero model.StageToken
	if !contains(Stages, o.Stage) {
		return zero, "", kerr.New("usage", "unknown stage %q (one of %v)", o.Stage, Stages)
	}
	if !agentRe.MatchString(o.Agent) {
		return zero, "", kerr.New("usage", "--agent %q must match %s", o.Agent, agentRe)
	}
	if o.Model == "" {
		return zero, "", kerr.New("usage", "--model is required")
	}
	if o.Region != "" && !regionRe.MatchString(o.Region) {
		return zero, "", kerr.New("usage", "--region %q must match %s", o.Region, regionRe)
	}
	if o.Stage == "D2" {
		if o.Region == "" {
			return zero, "", kerr.New("usage", "stage D2 needs --region (claims are written per region)")
		}
		var cover struct {
			Regions []struct {
				ID string `json:"id"`
			} `json:"regions"`
		}
		if err := r.ReadJSONFile(runs.CoverFile, &cover); err != nil {
			return zero, "", kerr.New("precondition", "stage D2 needs %s (write and validate the cover first): %v", runs.CoverFile, err)
		}
		known := make([]string, 0, len(cover.Regions))
		for _, reg := range cover.Regions {
			known = append(known, reg.ID)
		}
		if !contains(known, o.Region) {
			sort.Strings(known)
			return zero, "", kerr.New("unknown_region", "region %q is not in %s (regions: %v)", o.Region, runs.CoverFile, known)
		}
	}
	hash := "none"
	if o.PromptFile != "" {
		b, err := os.ReadFile(o.PromptFile)
		if err != nil {
			return zero, "", kerr.New("io", "--prompt-file: %v", err)
		}
		hash = fsx.SHA256(b)
	}
	t := model.StageToken{TID: ids.New(), Run: r.Record.RunID, Stage: o.Stage, Region: o.Region, Agent: o.Agent, Model: o.Model, PromptHash: hash}
	if err := fsx.WriteJSON(r.Path(runs.TokensDir+"/"+t.TID+".json"), t); err != nil {
		return zero, "", kerr.Wrap("io", err)
	}
	err := r.Update(func(rec *model.RunRecord) error {
		if rec.ModelMap == nil {
			rec.ModelMap = map[string]string{}
		}
		if rec.PromptHashes == nil {
			rec.PromptHashes = map[string]string{}
		}
		rec.ModelMap[o.Agent] = o.Model
		if hash != "none" {
			rec.PromptHashes[promptKey(o)] = hash
		}
		return nil
	})
	if err != nil {
		return zero, "", err
	}
	return t, Encode(t), nil
}

func promptKey(o Options) string {
	k := o.Stage
	if o.Region != "" {
		k += ":" + o.Region
	}
	return k + "@" + o.Agent
}

// Encode serialises a token.
func Encode(t model.StageToken) string {
	b, _ := json.Marshal(t)
	return base64.RawURLEncoding.EncodeToString(b)
}

// Decode parses a token string (accepting padded or unpadded, standard or URL alphabets).
func Decode(s string) (model.StageToken, error) {
	var t model.StageToken
	var b []byte
	var err error
	for _, enc := range []*base64.Encoding{base64.RawURLEncoding, base64.URLEncoding, base64.RawStdEncoding, base64.StdEncoding} {
		if b, err = enc.DecodeString(s); err == nil {
			break
		}
	}
	if err != nil {
		return t, kerr.New("bad_token", "token is not base64: %v", err)
	}
	dec := json.NewDecoder(bytesReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&t); err != nil {
		return t, kerr.New("bad_token", "token is not a stage token: %v", err)
	}
	return t, nil
}

// Verify checks that a token was issued by the kernel for this run, unchanged.
func Verify(r *runs.Run, t model.StageToken) error {
	if t.Run != r.Record.RunID {
		return kerr.New("bad_token", "token belongs to run %s, not %s", t.Run, r.Record.RunID)
	}
	if !ids.Valid(t.TID) {
		return kerr.New("bad_token", "token id %q is not a ULID", t.TID)
	}
	var issued model.StageToken
	if err := fsx.ReadJSON(r.Path(runs.TokensDir+"/"+t.TID+".json"), &issued); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return kerr.New("bad_token", "token %s was not issued by the kernel for this run", t.TID)
		}
		return kerr.Wrap("io", err)
	}
	if issued != t {
		return kerr.New("bad_token", "token %s does not match the one the kernel issued (edited?)", t.TID)
	}
	return nil
}

// Issued returns every token the kernel issued for the run, keyed by tid.
func Issued(r *runs.Run) (map[string]model.StageToken, error) {
	out := map[string]model.StageToken{}
	for _, rel := range r.Glob(runs.TokensDir + "/*.json") {
		var t model.StageToken
		if err := fsx.ReadJSON(r.Path(rel), &t); err != nil {
			return nil, kerr.Wrap("io", err)
		}
		out[t.TID] = t
	}
	return out, nil
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
