package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	defaults "github.com/manutej/meta-suite/next/kernel/config"
)

func TestEmbeddedDefaultsMatchModesTable(t *testing.T) {
	c, err := Parse(defaults.TOML, "embedded")
	if err != nil {
		t.Fatal(err)
	}
	// spec/v2/03-architecture.md, Modes table: review_passes, max_regions, k, g3_pair_cap.
	want := map[string][4]int{"fast": {1, 3, 1, 50}, "deep_narrow": {2, 6, 3, 200}, "deep_broad": {3, 12, 5, 600}}
	for name, w := range want {
		m := c.Modes[name]
		if got := [4]int{m.ReviewPasses, m.MaxRegions, m.K, m.G3PairCap}; got != w {
			t.Errorf("modes.%s = %v, spec table says %v", name, got, w)
		}
	}
	if c.Modes["fast"].G4 || !c.Modes["deep_narrow"].G4 {
		t.Error("g4 must be off in fast and on in deep")
	}
	if m, _ := c.ResolveMode("deep"); m != "deep_narrow" {
		t.Errorf("deep resolves to %s", m)
	}
	if len(c.Hash) != 64 {
		t.Errorf("hash %q", c.Hash)
	}
}

func TestMissingKeyRejected(t *testing.T) {
	b := strings.Replace(string(defaults.TOML), "min_independent = 2", "", 1)
	if _, err := Parse([]byte(b), "test"); err == nil || !strings.Contains(err.Error(), "min_independent") {
		t.Fatalf("want missing-key error, got %v", err)
	}
}

func TestUnknownKeyRejected(t *testing.T) {
	b := string(defaults.TOML) + "\n[extra]\nfoo = 1\n"
	if _, err := Parse([]byte(b), "test"); err == nil || !strings.Contains(err.Error(), "extra.foo") {
		t.Fatalf("want unknown-key error, got %v", err)
	}
}

func TestEvenKRejected(t *testing.T) {
	b := strings.Replace(string(defaults.TOML), "k = 3", "k = 4", 1)
	if _, err := Parse([]byte(b), "test"); err == nil {
		t.Fatal("even k accepted")
	}
}

func TestOverrideFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "c.toml")
	b := strings.Replace(string(defaults.TOML), "min_independent = 2", "min_independent = 3", 1)
	os.WriteFile(p, []byte(b), 0o644)
	t.Setenv(EnvOverride, p)
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.Lens.Pullback.MinIndependent != 3 || c.Source != p {
		t.Fatalf("override not applied: %d %s", c.Lens.Pullback.MinIndependent, c.Source)
	}
}
