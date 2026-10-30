// Command meta is the meta-suite v2 kernel: schemas, gates, ledger and bookkeeping for the meta-suite-next
// skills (spec/v2/03-architecture.md). The model does judgment; the kernel does bookkeeping and verdicts.
//
// Output contract (every command):
//   - exit 0: success; the result is on stdout (JSON unless the usage line says otherwise).
//   - exit 1: a verdict, not a crash: invalid object, failed or unchecked gate, refused merge. The full JSON
//     report is on stdout.
//   - exit 2: the command could not run (usage, missing input, I/O, broken precondition). stderr carries
//     {"error":{"kind":…,"message":…}} and stdout is empty.
//
// Every command that names a run directory appends an event to that run's ledger (except `ledger show`,
// which only reads). Output is deterministic and sorted (D-050).
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/manutej/meta-suite/next/kernel/internal/config"
	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/kerr"
	"github.com/manutej/meta-suite/next/kernel/internal/ledger"
	runs "github.com/manutej/meta-suite/next/kernel/internal/rundir"
)

// Set at link time: -ldflags "-X main.version=… -X main.buildHash=…".
var (
	version   = "dev"
	buildHash = "unset"
)

// Usage lines, one per command. `meta help` prints them; the report quotes them.
var usage = []string{
	"meta version",
	"meta help [<command>]          (or: meta <command> --help)",
	"meta config show [--mode fast|deep|deep_narrow|deep_broad]",
	"meta run init --brief <brief.json> [--mode fast|deep] [--surface <name>] [--seed <n>] [--root <dir>] [--json]       (prints the run dir)",
	"meta run scope <run_dir> narrow|broad",
	"meta validate <file> [--type <Type>] [--mode <mode>]",
	"meta stage begin <run_dir> <stage> [--region <r>] [--agent <a>] [--model <m>] [--prompt-file <f>] [--json]   (prints the token)",
	"meta claim add <run_dir> --token <token>            (claim JSON on stdin)",
	"meta claims merge <run_dir>",
	"meta graph build <run_dir>",
	"meta interface merge <run_dir> [--plain]",
	"meta findings verify <run_dir>",
	"meta gate G2 <run_dir>",
	"meta gate G6 <file>... [--run <run_dir>] [--allow <term>[,<term>...]]",
	"meta ledger show <run_dir>",
	"meta ledger append <run_dir> --kind <kind> [--stage <stage>] [--note <text>]   (optional JSON object on stdin)",
	"meta report <run_dir>",
	"meta sheaf analyze <in.json> [--out <report.json>] [--k K] [--null-reps R] [--seed S] [--top N] [--boundary <b.json>]",
	"meta check [--repo <repo_root>] [--skip-go]",
}

func main() { os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)) }

// env carries the I/O of one invocation.
type env struct {
	stdin          io.Reader
	stdout, stderr io.Writer
	cfg            *config.Config
}

type command func(e *env, args []string) int

func commands() map[string]command {
	return map[string]command{
		"version":         cmdVersion,
		"help":            cmdHelp,
		"config show":     cmdConfigShow,
		"run init":        cmdRunInit,
		"run scope":       cmdRunScope,
		"validate":        cmdValidate,
		"stage begin":     cmdStageBegin,
		"claim add":       cmdClaimAdd,
		"claims merge":    cmdClaimsMerge,
		"graph build":     cmdGraphBuild,
		"interface merge": cmdInterfaceMerge,
		"findings verify": cmdFindingsVerify,
		"gate G2":         cmdGateG2,
		"gate G6":         cmdGateG6,
		"ledger show":     cmdLedgerShow,
		"ledger append":   cmdLedgerAppend,
		"report":          cmdReport,
		"sheaf analyze":   cmdSheafAnalyze,
		"check":           cmdCheck,
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	e := &env{stdin: stdin, stdout: stdout, stderr: stderr}
	if len(args) == 0 {
		return e.fail(kerr.New("usage", "no command; one of: %s", strings.Join(usage, " | ")))
	}
	cmds := commands()
	name, rest := args[0], args[1:]
	if len(args) > 1 {
		if _, ok := cmds[args[0]+" "+args[1]]; ok {
			name, rest = args[0]+" "+args[1], args[2:]
		}
	}
	c, ok := cmds[name]
	if !ok {
		var groups []string
		for k := range cmds {
			if strings.HasPrefix(k, args[0]+" ") {
				groups = append(groups, k)
			}
		}
		sort.Strings(groups)
		if len(groups) > 0 {
			return e.fail(kerr.New("usage", "unknown subcommand for %q; one of: %s", args[0], strings.Join(groups, ", ")))
		}
		return e.fail(kerr.New("usage", "unknown command %q (try `meta help`)", args[0]))
	}
	if wantsHelp(rest) || (name == "help" && len(rest) > 0) {
		target := name
		if name == "help" && !wantsHelp(rest) {
			target = strings.Join(rest, " ")
		}
		h, ok := helpFor(target)
		if !ok {
			return e.fail(kerr.New("usage", "no help for %q (try `meta help`)", target))
		}
		fmt.Fprint(e.stdout, h)
		return 0
	}
	if name != "version" && name != "help" {
		cfg, err := config.Load()
		if err != nil {
			return e.fail(kerr.New("config", "%v", err))
		}
		e.cfg = cfg
	}
	return c(e, rest)
}

// fail prints an error JSON on stderr and returns exit 2.
func (e *env) fail(err error) int {
	var ke *kerr.Error
	if !errors.As(err, &ke) {
		ke = kerr.New("internal", "%v", err)
	}
	b, _ := fsx.Compact(map[string]any{"error": ke})
	fmt.Fprintln(e.stderr, string(b))
	return kerr.ExitError
}

// out prints v as indented JSON on stdout and returns code.
func (e *env) out(v any, code int) int {
	b, err := fsx.Marshal(v)
	if err != nil {
		return e.fail(kerr.Wrap("internal", err))
	}
	e.stdout.Write(b)
	return code
}

// log appends a ledger event; a ledger failure is a command failure.
func (e *env) log(r *runs.Run, command string, ok bool, data map[string]any) error {
	if err := ledger.Append(r, command, ok, data); err != nil {
		return kerr.Wrap("ledger", err)
	}
	return nil
}

// failRun records a failed command in the run's ledger (best effort) and prints the error.
func (e *env) failRun(r *runs.Run, command string, err error) int {
	if r != nil {
		var ke *kerr.Error
		kind := "internal"
		if errors.As(err, &ke) {
			kind = ke.Kind
		}
		ledger.Append(r, command, false, map[string]any{"error_kind": kind, "error": err.Error()})
	}
	return e.fail(err)
}

// flags is a tiny parser: positionals plus --name value / --name=value / boolean --name.
type flags struct {
	pos  []string
	vals map[string]string
	set  map[string]bool
}

func parseFlags(args []string, valued []string, boolean []string) (*flags, error) {
	f := &flags{vals: map[string]string{}, set: map[string]bool{}}
	isVal, isBool := map[string]bool{}, map[string]bool{}
	for _, v := range valued {
		isVal[v] = true
	}
	for _, b := range boolean {
		isBool[b] = true
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			f.pos = append(f.pos, args[i+1:]...)
			break
		}
		if !strings.HasPrefix(a, "--") || a == "-" {
			f.pos = append(f.pos, a)
			continue
		}
		name, val, hasEq := strings.Cut(strings.TrimPrefix(a, "--"), "=")
		switch {
		case isBool[name]:
			if hasEq {
				return nil, kerr.New("usage", "--%s takes no value", name)
			}
			f.set[name] = true
		case isVal[name]:
			if !hasEq {
				if i+1 >= len(args) {
					return nil, kerr.New("usage", "--%s needs a value", name)
				}
				i++
				val = args[i]
			}
			if _, dup := f.vals[name]; dup {
				return nil, kerr.New("usage", "--%s given twice", name)
			}
			f.vals[name] = val
			f.set[name] = true
		default:
			return nil, kerr.New("usage", "unknown flag --%s", name)
		}
	}
	return f, nil
}

func usageErr(cmd string) error {
	for _, u := range usage {
		if strings.HasPrefix(u, "meta "+cmd+" ") || u == "meta "+cmd {
			return kerr.New("usage", "usage: %s", u)
		}
	}
	return kerr.New("usage", "usage: meta %s", cmd)
}

// openRun opens a run directory or returns a kernel error.
func openRun(dir string) (*runs.Run, error) {
	if dir == "" {
		return nil, kerr.New("usage", "missing <run_dir>")
	}
	return runs.Open(dir)
}
