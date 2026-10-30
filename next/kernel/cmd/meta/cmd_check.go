package main

import (
	"os"
	"time"

	"github.com/manutej/meta-suite/next/kernel/internal/check"
	"github.com/manutej/meta-suite/next/kernel/internal/kerr"
)

func cmdCheck(e *env, args []string) int {
	f, err := parseFlags(args, []string{"repo"}, []string{"skip-go"})
	if err != nil || len(f.pos) > 0 {
		return e.fail(usageErr("check"))
	}
	repo := f.vals["repo"]
	if repo == "" {
		wd, _ := os.Getwd()
		repo = check.FindRepo(wd)
		if repo == "" {
			return e.fail(kerr.New("precondition", "not inside the meta-suite repo (no spec/v2/04-atoms.md and next/kernel above %s); pass --repo", wd))
		}
	} else if check.FindRepo(repo) != mustAbs(repo) {
		return e.fail(kerr.New("precondition", "%s is not the meta-suite repo root", repo))
	}
	rep := check.Run(check.Options{Repo: mustAbs(repo), AtomCount: e.cfg.Check.AtomCount, SkillMaxWords: e.cfg.Check.SkillMaxWords,
		SkipGo: f.set["skip-go"], GoTimeout: 10 * time.Minute})
	if rep.Verdict != "pass" {
		return e.out(rep, kerr.ExitVerdict)
	}
	return e.out(rep, 0)
}
