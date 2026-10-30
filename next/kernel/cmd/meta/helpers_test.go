package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// res is one in-process invocation of the CLI.
type res struct {
	code           int
	stdout, stderr string
}

func meta(t *testing.T, stdin string, args ...string) res {
	t.Helper()
	var o, e bytes.Buffer
	// The first element may name a two-word command ("run init"); split it into real argv words.
	argv := append(strings.Fields(args[0]), args[1:]...)
	code := run(argv, strings.NewReader(stdin), &o, &e)
	return res{code, o.String(), e.String()}
}

// must runs a command that has to succeed.
func must(t *testing.T, stdin string, args ...string) res {
	t.Helper()
	r := meta(t, stdin, args...)
	if r.code != 0 {
		t.Fatalf("meta %v: exit %d\nstdout: %s\nstderr: %s", args, r.code, r.stdout, r.stderr)
	}
	return r
}

func (r res) json(t *testing.T) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(r.stdout), &m); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, r.stdout)
	}
	return m
}

// errKind asserts exit 2 with the error JSON shape on stderr and an empty stdout, and returns the kind.
func (r res) errKind(t *testing.T) string {
	t.Helper()
	if r.code != 2 || r.stdout != "" {
		t.Fatalf("want exit 2 with empty stdout, got %d stdout=%q stderr=%q", r.code, r.stdout, r.stderr)
	}
	var m struct {
		Error struct {
			Kind    string `json:"kind"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(r.stderr), &m); err != nil || m.Error.Kind == "" || m.Error.Message == "" {
		t.Fatalf("stderr is not an error JSON: %q", r.stderr)
	}
	return m.Error.Kind
}

const fixtures = "testdata/synthetic/e2e"

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(fixtures, name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// newRun creates a run in a temp dir with the fixture brief and (optionally) cover, and returns the run dir.
func newRun(t *testing.T, mode string, withCover bool) string {
	t.Helper()
	root := t.TempDir()
	brief := filepath.Join(root, "brief.in.json")
	os.WriteFile(brief, readFixture(t, "brief.json"), 0o644)
	r := must(t, "", "run init", "--brief", brief, "--mode", mode, "--root", root, "--seed", "7")
	dir := strings.TrimSpace(r.stdout)
	if withCover {
		os.WriteFile(filepath.Join(dir, "cover.json"), readFixture(t, "cover.json"), 0o644)
		must(t, "", "validate", filepath.Join(dir, "cover.json"))
	}
	return dir
}

func token(t *testing.T, dir, stage, region, agent string) string {
	t.Helper()
	args := []string{"stage begin", dir, stage, "--agent", agent, "--model", "test-model"}
	if region != "" {
		args = append(args, "--region", region)
	}
	return strings.TrimSpace(must(t, "", args...).stdout)
}

func readJSONL(t *testing.T, path string) []map[string]any {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var out []map[string]any
	for _, l := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		if l == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(l), &m); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		out = append(out, m)
	}
	return out
}

func readJSON(t *testing.T, path string, v any) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
}
