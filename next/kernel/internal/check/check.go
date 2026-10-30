// Package check is `meta check`, the repo-local gate (M1, D-052): it runs in seconds and is what an agent
// runs before every push. Checks: go vet and go test (which include the sheaf parity tests) on the kernel,
// every next/skills/*/SKILL.md within check.skill_max_words words, runs/ and the built binaries gitignored
// (and the launcher not), the spec atom count equal to check.atom_count and to the index's atom links, and
// every D-### citation in the spec resolving to a decision row.
package check

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Result is one check's outcome.
type Result struct {
	Name   string   `json:"name"`
	Status string   `json:"status"` // pass | fail | unchecked | delegated
	Detail string   `json:"detail"`
	Items  []string `json:"items,omitempty"`
}

// Report is the gate's output.
type Report struct {
	Verdict string   `json:"verdict"`
	Repo    string   `json:"repo"`
	Checks  []Result `json:"checks"`
}

// Options configure a check run.
type Options struct {
	Repo          string
	AtomCount     int
	SkillMaxWords int
	SkipGo        bool // the caller (make check) has just run go vet and go test itself
	GoTimeout     time.Duration
}

// FindRepo walks up from start to the directory holding spec/v2/04-atoms.md and next/kernel/go.mod.
func FindRepo(start string) string {
	abs, err := filepath.Abs(start)
	if err != nil {
		return ""
	}
	for {
		if exists(filepath.Join(abs, "spec", "v2", "04-atoms.md")) && exists(filepath.Join(abs, "next", "kernel", "go.mod")) {
			return abs
		}
		p := filepath.Dir(abs)
		if p == abs {
			return ""
		}
		abs = p
	}
}

// Run executes every check.
func Run(o Options) *Report {
	rep := &Report{Repo: o.Repo}
	kernel := filepath.Join(o.Repo, "next", "kernel")
	if o.SkipGo {
		reason := "run by the caller (make check) immediately before this command"
		rep.Checks = append(rep.Checks, Result{Name: "go_vet", Status: "delegated", Detail: reason}, Result{Name: "go_test", Status: "delegated", Detail: reason})
	} else {
		rep.Checks = append(rep.Checks, goCmd(kernel, o.GoTimeout, "go_vet", "vet", "./..."))
		rep.Checks = append(rep.Checks, goTest(kernel, o.GoTimeout))
	}
	rep.Checks = append(rep.Checks, skillWords(o.Repo, o.SkillMaxWords), gitignore(o.Repo), atoms(o.Repo, o.AtomCount), citations(o.Repo))
	rep.Verdict = "pass"
	for _, c := range rep.Checks {
		if c.Status == "fail" || c.Status == "unchecked" {
			rep.Verdict = "fail"
		}
	}
	return rep
}

func goCmd(dir string, timeout time.Duration, name string, args ...string) Result {
	gobin, err := exec.LookPath("go")
	if err != nil {
		return Result{Name: name, Status: "unchecked", Detail: "the go toolchain is not on PATH"}
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, gobin, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return Result{Name: name, Status: "fail", Detail: "go " + strings.Join(args, " ") + ": " + err.Error(), Items: tail(out, 40)}
	}
	return Result{Name: name, Status: "pass", Detail: "go " + strings.Join(args, " ")}
}

func goTest(dir string, timeout time.Duration) Result {
	r := goCmd(dir, timeout, "go_test", "test", "-count=1", "./...")
	if r.Status != "pass" {
		return r
	}
	// List the packages that have tests, and require the sheaf parity package among them.
	gobin, _ := exec.LookPath("go")
	cmd := exec.Command(gobin, "list", "-f", "{{if .TestGoFiles}}{{.ImportPath}}{{end}}", "./...")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return Result{Name: "go_test", Status: "fail", Detail: "go list: " + err.Error()}
	}
	var pkgs []string
	for _, l := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			pkgs = append(pkgs, l)
		}
	}
	sort.Strings(pkgs)
	hasSheaf := false
	for _, p := range pkgs {
		hasSheaf = hasSheaf || strings.HasSuffix(p, "/internal/sheaf")
	}
	if !hasSheaf {
		return Result{Name: "go_test", Status: "fail", Detail: "the sheaf parity tests (internal/sheaf) are missing", Items: pkgs}
	}
	return Result{Name: "go_test", Status: "pass", Detail: fmt.Sprintf("go test -count=1 ./... passed in %d packages, including the sheaf parity tests", len(pkgs)), Items: pkgs}
}

func skillWords(repo string, max int) Result {
	files, _ := filepath.Glob(filepath.Join(repo, "next", "skills", "*", "SKILL.md"))
	sort.Strings(files)
	if len(files) == 0 {
		return Result{Name: "skill_md_words", Status: "unchecked", Detail: "no next/skills/*/SKILL.md found"}
	}
	r := Result{Name: "skill_md_words", Status: "pass", Detail: fmt.Sprintf("every SKILL.md (whole file, frontmatter included) has at most %d words", max)}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return Result{Name: "skill_md_words", Status: "fail", Detail: err.Error()}
		}
		n := len(strings.Fields(string(b)))
		rel, _ := filepath.Rel(repo, f)
		r.Items = append(r.Items, fmt.Sprintf("%s: %d words", filepath.ToSlash(rel), n))
		if n > max {
			r.Status = "fail"
			r.Detail = fmt.Sprintf("a SKILL.md exceeds %d words", max)
		}
	}
	return r
}

// gitignore asks git which probe paths are ignored; without git it cannot tell (unchecked).
func gitignore(repo string) Result {
	name := "gitignore"
	gitbin, err := exec.LookPath("git")
	if err != nil {
		return Result{Name: name, Status: "unchecked", Detail: "git is not on PATH"}
	}
	mustIgnore := []string{"runs/01PROBE/run.json", "next/runs/01PROBE/run.json", "next/kernel/runs/01PROBE/run.json",
		"next/bin/meta-linux-amd64", "next/bin/meta-darwin-arm64", "next/bin/meta-windows-amd64.exe"}
	mustKeep := []string{"next/bin/meta", "next/kernel/internal/rundir/rundir.go", "next/kernel/config/defaults.toml"}
	check := func(p string) (bool, error) {
		cmd := exec.Command(gitbin, "check-ignore", "-q", "--no-index", p)
		cmd.Dir = repo
		err := cmd.Run()
		if err == nil {
			return true, nil
		}
		if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}
	r := Result{Name: name, Status: "pass", Detail: "runs/ and built binaries are ignored; the launcher and sources are not"}
	for _, p := range mustIgnore {
		ig, err := check(p)
		if err != nil {
			return Result{Name: name, Status: "unchecked", Detail: "git check-ignore: " + err.Error()}
		}
		if !ig {
			r.Status, r.Detail = "fail", "a path that must be ignored is not"
			r.Items = append(r.Items, "not ignored: "+p)
		}
	}
	for _, p := range mustKeep {
		ig, err := check(p)
		if err != nil {
			return Result{Name: name, Status: "unchecked", Detail: "git check-ignore: " + err.Error()}
		}
		if ig {
			r.Status, r.Detail = "fail", "a path that must be committed is ignored"
			r.Items = append(r.Items, "ignored: "+p)
		}
	}
	return r
}

var atomLink = regexp.MustCompile(`\]\(atoms/([A-Z0-9]+-[a-z0-9-]+\.md)\)`)

func atoms(repo string, want int) Result {
	name := "spec_atom_count"
	files, _ := filepath.Glob(filepath.Join(repo, "spec", "v2", "atoms", "*.md"))
	idx, err := os.ReadFile(filepath.Join(repo, "spec", "v2", "04-atoms.md"))
	if err != nil {
		return Result{Name: name, Status: "unchecked", Detail: err.Error()}
	}
	onDisk := map[string]bool{}
	for _, f := range files {
		onDisk[filepath.Base(f)] = true
	}
	linked := map[string]bool{}
	for _, m := range atomLink.FindAllSubmatch(idx, -1) {
		linked[string(m[1])] = true
	}
	r := Result{Name: name, Status: "pass", Detail: fmt.Sprintf("%d atom files, %d linked from 04-atoms.md, expected %d", len(onDisk), len(linked), want)}
	if len(onDisk) != want || len(linked) != want {
		r.Status = "fail"
	}
	for f := range onDisk {
		if !linked[f] {
			r.Status = "fail"
			r.Items = append(r.Items, "not in the index: "+f)
		}
	}
	for f := range linked {
		if !onDisk[f] {
			r.Status = "fail"
			r.Items = append(r.Items, "index links a missing file: "+f)
		}
	}
	sort.Strings(r.Items)
	return r
}

var (
	decisionRow = regexp.MustCompile(`(?m)^\| (D-\d{3}) \|`)
	decisionRef = regexp.MustCompile(`\bD-\d{3}\b`)
)

func citations(repo string) Result {
	name := "decision_citations"
	dec, err := os.ReadFile(filepath.Join(repo, "spec", "v2", "09-decisions.md"))
	if err != nil {
		return Result{Name: name, Status: "unchecked", Detail: err.Error()}
	}
	rows := map[string]bool{}
	for _, m := range decisionRow.FindAllSubmatch(dec, -1) {
		rows[string(m[1])] = true
	}
	var files []string
	for _, pat := range []string{"spec/v2/*.md", "spec/v2/atoms/*.md"} {
		m, _ := filepath.Glob(filepath.Join(repo, filepath.FromSlash(pat)))
		files = append(files, m...)
	}
	sort.Strings(files)
	missing := map[string]bool{}
	cited := map[string]bool{}
	for _, f := range files {
		b, _ := os.ReadFile(f)
		for _, m := range decisionRef.FindAll(b, -1) {
			cited[string(m)] = true
			if !rows[string(m)] {
				rel, _ := filepath.Rel(repo, f)
				missing[filepath.ToSlash(rel)+": "+string(m)] = true
			}
		}
	}
	r := Result{Name: name, Status: "pass", Detail: fmt.Sprintf("%d distinct D-### citations, all resolve to one of %d decision rows", len(cited), len(rows))}
	if len(missing) > 0 {
		r.Status, r.Detail = "fail", "citations without a decision row"
		for k := range missing {
			r.Items = append(r.Items, k)
		}
		sort.Strings(r.Items)
	}
	return r
}

func tail(b []byte, n int) []string {
	lines := strings.Split(strings.TrimRight(string(bytes.TrimSpace(b)), "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return lines
}

func exists(p string) bool { _, err := os.Stat(p); return err == nil }
