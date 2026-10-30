package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"runtime"
	"strconv"
	"strings"

	"github.com/manutej/meta-suite/next/kernel/internal/claims"
	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/kerr"
	"github.com/manutej/meta-suite/next/kernel/internal/ledger"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
	runs "github.com/manutej/meta-suite/next/kernel/internal/rundir"
	"github.com/manutej/meta-suite/next/kernel/internal/schema"
	"github.com/manutej/meta-suite/next/kernel/internal/stage"
	"github.com/manutej/meta-suite/next/kernel/schemas"
)

func cmdVersion(e *env, args []string) int {
	if len(args) > 0 {
		return e.fail(usageErr("version"))
	}
	fmt.Fprintf(e.stdout, "meta %s %s/%s build=%s go=%s schema=%s\n", version, runtime.GOOS, runtime.GOARCH, buildHash, runtime.Version(), schemas.Version)
	return 0
}

func cmdHelp(e *env, args []string) int {
	for _, u := range usage {
		fmt.Fprintln(e.stdout, u)
	}
	fmt.Fprint(e.stdout, "\n`meta <command> --help` prints the command's guidance: purpose, arguments, outputs, error kinds, gotchas.\n"+exitContract)
	return 0
}

func cmdConfigShow(e *env, args []string) int {
	f, err := parseFlags(args, []string{"mode"}, nil)
	if err != nil || len(f.pos) > 0 {
		return e.fail(usageErr("config show"))
	}
	out := map[string]any{"config_hash": e.cfg.Hash, "source": e.cfg.Source, "config": e.cfg}
	if m, ok := f.vals["mode"]; ok {
		name, err := e.cfg.ResolveMode(m)
		if err != nil {
			return e.fail(kerr.New("usage", "%v", err))
		}
		out["mode"] = name
		out["mode_values"] = e.cfg.Modes[name]
	}
	return e.out(out, 0)
}

// cmdRunInit creates runs/<ULID>/ under --root (default: the working directory), copies the brief to
// brief.json (after G1), and writes the RunRecord skeleton.
func cmdRunInit(e *env, args []string) int {
	f, err := parseFlags(args, []string{"brief", "mode", "surface", "seed", "root"}, []string{"json"})
	if err != nil || len(f.pos) > 0 || f.vals["brief"] == "" {
		return e.fail(usageErr("run init"))
	}
	mode := "deep"
	if m, ok := f.vals["mode"]; ok {
		mode = m
	}
	if mode != "fast" && mode != "deep" {
		return e.fail(kerr.New("usage", "--mode must be fast or deep (deep resolves to modes.deep_narrow; `meta run scope` switches to broad)"))
	}
	modeRow, err := e.cfg.ResolveMode(mode)
	if err != nil {
		return e.fail(kerr.New("config", "%v", err))
	}
	brief, err := os.ReadFile(f.vals["brief"])
	if err != nil {
		return e.fail(kerr.New("io", "--brief: %v", err))
	}
	errs, err := schema.Validate("Brief", brief)
	if err != nil {
		return e.fail(kerr.Wrap("internal", err))
	}
	if len(errs) > 0 {
		return e.out(map[string]any{"created": false, "file": f.vals["brief"], "type": "Brief", "valid": false, "errors": errs}, kerr.ExitVerdict)
	}
	var seed int64
	if s, ok := f.vals["seed"]; ok {
		if seed, err = strconv.ParseInt(s, 10, 64); err != nil || seed < 0 {
			return e.fail(kerr.New("usage", "--seed must be a non-negative integer"))
		}
	} else {
		var b [8]byte
		rand.Read(b[:])
		seed = int64(binary.BigEndian.Uint64(b[:]) >> 1)
	}
	surface := f.vals["surface"]
	if surface == "" {
		surface = os.Getenv("META_SURFACE")
	}
	if surface == "" {
		surface = "unknown"
	}
	root := f.vals["root"]
	if root == "" {
		if root, err = os.Getwd(); err != nil {
			return e.fail(kerr.Wrap("io", err))
		}
	}
	rec := model.RunRecord{
		CreatedAt: runs.Stamp(runs.Now()), BriefHash: fsx.SHA256(brief), Mode: modeRow, Surface: surface,
		Kernel: "present", KernelVersion: version, ModelMap: map[string]string{}, PromptHashes: map[string]string{},
		ConfigHash: e.cfg.Hash, ConfigSource: e.cfg.Source, SchemaVersion: schemas.Version, GateResults: []model.GateResult{},
		Usage: model.Usage{Source: "estimated", Calls: 0}, Seed: seed,
	}
	r, err := runs.Create(root, rec)
	if err != nil {
		return e.fail(err)
	}
	if err := fsx.WriteAtomic(r.Path(runs.BriefFile), brief); err != nil {
		return e.failRun(r, "run init", kerr.Wrap("io", err))
	}
	if err := r.RecordGate(model.GateResult{Gate: "G1", Verdict: "pass", At: runs.Stamp(runs.Now()), Summary: "Brief", Files: []string{runs.BriefFile}}); err != nil {
		return e.failRun(r, "run init", err)
	}
	if err := e.log(r, "run init", true, map[string]any{"mode": modeRow, "brief_hash": rec.BriefHash, "config_hash": e.cfg.Hash,
		"config_source": e.cfg.Source, "schema_version": schemas.Version, "kernel_version": version, "surface": surface, "seed": seed}); err != nil {
		return e.fail(err)
	}
	if f.set["json"] {
		return e.out(map[string]any{"run_id": r.Record.RunID, "run_dir": r.Dir, "mode": modeRow, "config_hash": e.cfg.Hash}, 0)
	}
	fmt.Fprintln(e.stdout, r.Dir)
	return 0
}

// cmdRunScope switches a deep run between deep_narrow and deep_broad (I2 decides scope).
func cmdRunScope(e *env, args []string) int {
	f, err := parseFlags(args, nil, nil)
	if err != nil || len(f.pos) != 2 {
		return e.fail(usageErr("run scope"))
	}
	r, err := openRun(f.pos[0])
	if err != nil {
		return e.fail(err)
	}
	target := map[string]string{"narrow": "deep_narrow", "broad": "deep_broad"}[f.pos[1]]
	if target == "" {
		return e.failRun(r, "run scope", kerr.New("usage", "scope must be narrow or broad"))
	}
	if r.Record.Mode == "fast" {
		return e.failRun(r, "run scope", kerr.New("precondition", "a fast run has no scope; fast mode never scales up silently"))
	}
	from := r.Record.Mode
	if err := r.Update(func(rec *model.RunRecord) error { rec.Mode = target; return nil }); err != nil {
		return e.failRun(r, "run scope", err)
	}
	if err := e.log(r, "run scope", true, map[string]any{"from": from, "to": target}); err != nil {
		return e.fail(err)
	}
	return e.out(map[string]any{"mode": target, "mode_values": e.cfg.Modes[target]}, 0)
}

func cmdStageBegin(e *env, args []string) int {
	f, err := parseFlags(args, []string{"region", "agent", "model", "prompt-file"}, []string{"json"})
	if err != nil || len(f.pos) != 2 {
		return e.fail(usageErr("stage begin"))
	}
	r, err := openRun(f.pos[0])
	if err != nil {
		return e.fail(err)
	}
	agent, mdl := f.vals["agent"], f.vals["model"]
	if agent == "" {
		agent = "main"
	}
	if mdl == "" {
		mdl = "unrecorded"
	}
	t, tok, err := stage.Begin(r, stage.Options{Stage: f.pos[1], Region: f.vals["region"], Agent: agent, Model: mdl, PromptFile: f.vals["prompt-file"]})
	if err != nil {
		return e.failRun(r, "stage begin", err)
	}
	data := map[string]any{"token": t, "prompt_file": f.vals["prompt-file"]}
	if err := e.log(r, "stage begin", true, data); err != nil {
		return e.fail(err)
	}
	if f.set["json"] {
		return e.out(map[string]any{"token": tok, "fields": t}, 0)
	}
	fmt.Fprintln(e.stdout, tok)
	return 0
}

func cmdClaimAdd(e *env, args []string) int {
	f, err := parseFlags(args, []string{"token"}, nil)
	if err != nil || len(f.pos) != 1 || f.vals["token"] == "" {
		return e.fail(usageErr("claim add"))
	}
	r, err := openRun(f.pos[0])
	if err != nil {
		return e.fail(err)
	}
	input, err := readAll(e.stdin)
	if err != nil {
		return e.failRun(r, "claim add", kerr.Wrap("io", err))
	}
	res, rej, err := claims.Add(r, e.cfg, f.vals["token"], input)
	if err != nil {
		return e.failRun(r, "claim add", err)
	}
	if rej != nil {
		if err := e.log(r, "claim add", false, map[string]any{"rejected": true, "errors": rej.Errors, "input_hash": fsx.SHA256(input)}); err != nil {
			return e.fail(err)
		}
		return e.out(rej, kerr.ExitVerdict)
	}
	if err := e.log(r, "claim add", true, map[string]any{"id": res.ID, "region": res.Region, "shard": res.Shard, "notes": res.Notes}); err != nil {
		return e.fail(err)
	}
	return e.out(res, 0)
}

// cmdClaimRetract appends a tombstone for one claim (D-087). Shards stay append-only; every later command
// skips the claim. Run `meta claims merge` and the steps after it again.
func cmdClaimRetract(e *env, args []string) int {
	f, err := parseFlags(args, []string{"id", "reason", "stage"}, nil)
	if err != nil || len(f.pos) != 1 || f.vals["id"] == "" || strings.TrimSpace(f.vals["reason"]) == "" {
		return e.fail(usageErr("claim retract"))
	}
	r, err := openRun(f.pos[0])
	if err != nil {
		return e.fail(err)
	}
	if s := f.vals["stage"]; s != "" && !contains(stage.Stages, s) {
		return e.failRun(r, "claim retract", kerr.New("usage", "unknown --stage %q (one of %v)", s, stage.Stages))
	}
	t, already, err := claims.Retract(r, f.vals["id"], f.vals["reason"], f.vals["stage"])
	if err != nil {
		return e.failRun(r, "claim retract", err)
	}
	data := map[string]any{"id": t.ID, "reason": t.Reason, "stage": t.Stage}
	if already {
		data["already_retracted"] = true
	}
	if err := e.log(r, "claim retract", true, data); err != nil {
		return e.fail(err)
	}
	tomb, err := claims.Retracted(r)
	if err != nil {
		return e.fail(err)
	}
	out := map[string]any{"retracted": t, "retractions": len(tomb),
		"next": "run `meta claims merge`, then graph build, interface merge, findings verify, gate G2 and report again"}
	if already {
		out["note"] = "already retracted with this reason; nothing was added"
	}
	return e.out(out, 0)
}

func cmdClaimsMerge(e *env, args []string) int {
	f, err := parseFlags(args, nil, nil)
	if err != nil || len(f.pos) != 1 {
		return e.fail(usageErr("claims merge"))
	}
	r, err := openRun(f.pos[0])
	if err != nil {
		return e.fail(err)
	}
	res, rej, err := claims.Merge(r, e.cfg, e.indepKey)
	if err != nil {
		return e.failRun(r, "claims merge", err)
	}
	if rej != nil {
		if err := e.log(r, "claims merge", false, map[string]any{"refused": true, "bad_lines": len(rej.Bad)}); err != nil {
			return e.fail(err)
		}
		return e.out(rej, kerr.ExitVerdict)
	}
	if err := e.log(r, "claims merge", true, map[string]any{"claims": res.Claims, "retracted": res.Retracted, "shards": res.Shards, "by_region": res.ByRegion,
		"claims_hash": r.HashFiles([]string{runs.ClaimsFile})}); err != nil {
		return e.fail(err)
	}
	return e.out(res, 0)
}

func cmdLedgerShow(e *env, args []string) int {
	f, err := parseFlags(args, nil, nil)
	if err != nil || len(f.pos) != 1 {
		return e.fail(usageErr("ledger show"))
	}
	r, err := openRun(f.pos[0])
	if err != nil {
		return e.fail(err)
	}
	if err := ledger.Show(r, e.stdout); err != nil {
		return e.fail(err)
	}
	return 0
}

// cmdLedgerAppend records an operator event (S1): --kind is one of config ledger.kinds; deviation (D-080)
// records a place the model departed from an op, with --stage and --note. Extra fields come from a JSON object
// on stdin, read ONLY with --stdin (D-087): an open, empty, non-terminal stdin (a harness pipe) used to block
// the command forever.
func cmdLedgerAppend(e *env, args []string) int {
	f, err := parseFlags(args, []string{"kind", "stage", "note", "event"}, []string{"stdin"})
	if err != nil || len(f.pos) != 1 {
		return e.fail(usageErr("ledger append"))
	}
	kind := f.vals["kind"]
	if kind == "" && f.vals["event"] != "" {
		kind = "note" // --event is the older spelling: a named note
	}
	if kind == "" {
		return e.fail(usageErr("ledger append"))
	}
	known := false
	for _, k := range e.cfg.Ledger.Kinds {
		known = known || k == kind
	}
	if !known {
		return e.fail(kerr.New("usage", "--kind %q is not one of %v (config ledger.kinds)", kind, e.cfg.Ledger.Kinds))
	}
	r, err := openRun(f.pos[0])
	if err != nil {
		return e.fail(err)
	}
	if kind == "deviation" && (f.vals["stage"] == "" || strings.TrimSpace(f.vals["note"]) == "") {
		return e.failRun(r, "ledger append", kerr.New("usage", "--kind deviation needs --stage <stage> and --note \"what departed from the op, and why\""))
	}
	if s := f.vals["stage"]; s != "" && !contains(stage.Stages, s) {
		return e.failRun(r, "ledger append", kerr.New("usage", "unknown --stage %q (one of %v)", s, stage.Stages))
	}
	data := map[string]any{}
	if f.set["stdin"] {
		in, err := readAll(e.stdin)
		if err != nil {
			return e.failRun(r, "ledger append", kerr.Wrap("io", err))
		}
		if len(trimSpace(in)) == 0 {
			return e.failRun(r, "ledger append", kerr.New("usage", "--stdin given but stdin is empty; pipe one JSON object or drop --stdin"))
		}
		if err := json.Unmarshal(in, &data); err != nil {
			return e.failRun(r, "ledger append", kerr.New("invalid_json", "stdin must be a JSON object: %v", err))
		}
	}
	data["kind"] = kind
	for _, k := range []string{"stage", "note", "event"} {
		if v := f.vals[k]; v != "" {
			data[k] = v
		}
	}
	if err := e.log(r, "ledger append", true, data); err != nil {
		return e.fail(err)
	}
	evs, err := ledger.Events(r)
	if err != nil {
		return e.fail(err)
	}
	return e.out(map[string]any{"seq": len(evs), "kind": kind}, 0)
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func (e *env) indepKey(s string) string {
	return normIndep(s, e.cfg.Independence.StripPrefixes)
}
