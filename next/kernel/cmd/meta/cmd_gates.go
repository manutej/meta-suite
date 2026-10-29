package main

import (
	"os"
	"strings"

	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/gates"
	"github.com/manutej/meta-suite/next/kernel/internal/kerr"
	"github.com/manutej/meta-suite/next/kernel/internal/model"
	"github.com/manutej/meta-suite/next/kernel/internal/report"
	runs "github.com/manutej/meta-suite/next/kernel/internal/rundir"
	"github.com/manutej/meta-suite/next/kernel/internal/sheaf"
)

// cmdGateG2 runs the provenance gate. Exit 0 only on "pass"; "fail" and "unchecked" exit 1 with the report.
func cmdGateG2(e *env, args []string) int {
	f, err := parseFlags(args, nil, nil)
	if err != nil || len(f.pos) != 1 {
		return e.fail(usageErr("gate G2"))
	}
	r, err := openRun(f.pos[0])
	if err != nil {
		return e.fail(err)
	}
	rep, err := gates.G2(r, e.cfg.Lens.Pullback.MinIndependent, e.cfg.Independence.StripPrefixes)
	if err != nil {
		return e.failRun(r, "gate G2", err)
	}
	summary := "claims " + itoa(rep.Counts.Claims) + ", publishers " + itoa(rep.Counts.IndependenceKeys) + ", problems " + itoa(len(rep.Problems))
	if !rep.Complete {
		var parts []string
		for _, c := range rep.Checks {
			if c.Status == "unchecked" {
				parts = append(parts, c.Name)
			}
		}
		summary += "; not yet checkable: " + strings.Join(parts, ", ")
	}
	g := model.GateResult{Gate: "G2", Verdict: rep.Verdict, At: runs.Stamp(runs.Now()), Summary: summary, InputsHash: r.HashFiles(rep.Inputs), Files: rep.Inputs}
	if err := r.RecordGate(g); err != nil {
		return e.failRun(r, "gate G2", err)
	}
	if err := e.log(r, "gate G2", rep.Verdict == "pass", map[string]any{"verdict": rep.Verdict, "complete": rep.Complete,
		"problems": len(rep.Problems), "counts": rep.Counts, "inputs_hash": g.InputsHash}); err != nil {
		return e.fail(err)
	}
	if rep.Verdict != "pass" {
		return e.out(rep, kerr.ExitVerdict)
	}
	return e.out(rep, 0)
}

// cmdGateG6 scans user-facing files for banned vocabulary. --run reads the brief's vocab_allowlist and
// records the verdict in that run; --allow adds exemptions (each logged with reason "--allow").
func cmdGateG6(e *env, args []string) int {
	f, err := parseFlags(args, []string{"run", "allow"}, nil)
	if err != nil || len(f.pos) == 0 {
		return e.fail(usageErr("gate G6"))
	}
	var allow []gates.Allow
	var r *runs.Run
	if dir := f.vals["run"]; dir != "" {
		if r, err = openRun(dir); err != nil {
			return e.fail(err)
		}
		var brief struct {
			VocabAllowlist []gates.Allow `json:"vocab_allowlist"`
		}
		if r.Exists(runs.BriefFile) {
			if err := r.ReadJSONFile(runs.BriefFile, &brief); err != nil {
				return e.failRun(r, "gate G6", err)
			}
		}
		allow = append(allow, brief.VocabAllowlist...)
	} else if dir := runs.Find(f.pos[0]); dir != "" {
		r, _ = runs.Open(dir)
	}
	for _, t := range strings.Split(f.vals["allow"], ",") {
		if t = strings.TrimSpace(t); t != "" {
			allow = append(allow, gates.Allow{Term: t, Reason: "--allow"})
		}
	}
	for _, p := range f.pos {
		if st, err := os.Stat(p); err != nil || st.IsDir() {
			if err == nil {
				err = kerr.New("usage", "%s is a directory; pass files", p)
			}
			return e.failRun(r, "gate G6", kerr.Wrap("io", err))
		}
	}
	cfg := gates.G6Config{Terms: e.cfg.Gates.G6.Terms, Tokens: e.cfg.Gates.G6.Tokens, Phrases: e.cfg.Gates.G6.Phrases}
	rep, err := gates.G6(f.pos, os.ReadFile, cfg, allow)
	if err != nil {
		return e.failRun(r, "gate G6", kerr.Wrap("io", err))
	}
	if r != nil {
		files := make([]string, len(rep.Files))
		var buf []byte
		for i, p := range rep.Files {
			files[i] = r.Rel(p)
			b, _ := os.ReadFile(p)
			buf = append(buf, []byte(files[i]+"\x00"+fsx.SHA256(b)+"\x00")...)
		}
		g := model.GateResult{Gate: "G6", Verdict: rep.Verdict, At: runs.Stamp(runs.Now()), Files: files, InputsHash: fsx.SHA256(buf),
			Summary: itoa(len(rep.Hits)) + " hit(s) in " + itoa(len(files)) + " file(s), " + itoa(len(rep.Allowed)) + " allowed"}
		if err := r.RecordGate(g); err != nil {
			return e.failRun(r, "gate G6", err)
		}
		allowLog := make([]map[string]string, 0, len(rep.Allowed))
		for _, h := range rep.Allowed {
			allowLog = append(allowLog, map[string]string{"term": h.Term, "match": h.Match, "reason": h.Reason, "file": r.Rel(h.File)})
		}
		if err := e.log(r, "gate G6", rep.Verdict == "pass", map[string]any{"verdict": rep.Verdict, "files": files, "hits": len(rep.Hits),
			"exemptions": allowLog}); err != nil {
			return e.fail(err)
		}
	}
	if rep.Verdict != "pass" {
		return e.out(rep, kerr.ExitVerdict)
	}
	return e.out(rep, 0)
}

func cmdReport(e *env, args []string) int {
	f, err := parseFlags(args, nil, nil)
	if err != nil || len(f.pos) != 1 {
		return e.fail(usageErr("report"))
	}
	r, err := openRun(f.pos[0])
	if err != nil {
		return e.fail(err)
	}
	hc, err := report.Build(r, e.cfg)
	if err != nil {
		return e.failRun(r, "report", kerr.Wrap("io", err))
	}
	if err := fsx.WriteJSON(r.Path(runs.HowCheckedFile), hc); err != nil {
		return e.failRun(r, "report", kerr.Wrap("io", err))
	}
	if err := e.log(r, "report", true, map[string]any{"out": runs.HowCheckedFile, "gates_passed": hc.GatesPassed,
		"gates_failed": hc.GatesFailed, "gates_unchecked": hc.GatesUnchecked, "verified": hc.Counts.Verified}); err != nil {
		return e.fail(err)
	}
	return e.out(hc, 0)
}

// cmdSheafAnalyze is the standalone gluing check (D-048): the vendored sheaf engine's CLI, unchanged in
// flags, exit codes and report format (frontier-sheaf/report-v1). Its input errors are one "P_INPUT_01 …"
// line on stderr (exit 2), as upstream.
func cmdSheafAnalyze(e *env, args []string) int {
	return sheaf.Main(append([]string{"analyze"}, args...), e.stdout, e.stderr)
}
