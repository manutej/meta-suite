// Package config loads the kernel configuration. The embedded config/defaults.toml is used unless the
// META_CONFIG environment variable names an override file, which then replaces it whole. Every key is
// required: a missing key is an error, never a silent zero (P6: thresholds live in config, not in code).
package config

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/BurntSushi/toml"
	defaults "github.com/manutej/meta-suite/next/kernel/config"
	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"golang.org/x/text/currency"
)

// Mode is one row of the modes table (spec/v2/03-architecture.md, Modes).
type Mode struct {
	ReviewPasses int    `toml:"review_passes" json:"review_passes"`
	MaxRegions   int    `toml:"max_regions" json:"max_regions"`
	K            int    `toml:"k" json:"k"`
	G3Layers     string `toml:"g3_layers" json:"g3_layers"`
	G3PairCap    int    `toml:"g3_pair_cap" json:"g3_pair_cap"`
	G4           bool   `toml:"g4" json:"g4"`
	R2           string `toml:"r2" json:"r2"`
	LoadedLayers string `toml:"loaded_layers" json:"loaded_layers"`
	P50Minutes   int    `toml:"p50_minutes" json:"p50_minutes"`
	P90Minutes   int    `toml:"p90_minutes" json:"p90_minutes"`
}

// UnitDim is one canonical-unit dimension.
type UnitDim struct {
	Canonical string             `toml:"canonical" json:"canonical"`
	Factors   map[string]float64 `toml:"factors" json:"factors"`
}

// UnitCodes is a code-only dimension (D-087, D-097): a unit matching Pattern (and, when List is "iso4217",
// listed in ISO 4217, from the CLDR data embedded in golang.org/x/text/currency) is its own canonical unit and
// is never converted.
type UnitCodes struct {
	Pattern string         `toml:"pattern" json:"pattern"`
	List    string         `toml:"list" json:"list"`
	re      *regexp.Regexp `toml:"-"`
}

// Match reports whether unit is a code of this dimension (case-sensitive: callers upper-case first).
func (u UnitCodes) Match(unit string) bool {
	if u.re == nil {
		u.re = regexp.MustCompile(u.Pattern)
	}
	if !u.re.MatchString(unit) {
		return false
	}
	if u.List == "iso4217" {
		_, err := currency.ParseISO(unit)
		return err == nil
	}
	return true
}

// Config is the whole configuration.
type Config struct {
	Schema      string            `toml:"schema" json:"schema"`
	ModeAliases map[string]string `toml:"mode_aliases" json:"mode_aliases"`
	Modes       map[string]Mode   `toml:"modes" json:"modes"`
	Loop        struct {
		D2ReentryMax int `toml:"d2_reentry_max" json:"d2_reentry_max"`
	} `toml:"loop" json:"loop"`
	Intake struct {
		ReExpandMax int `toml:"re_expand_max" json:"re_expand_max"`
	} `toml:"intake" json:"intake"`
	Roles  map[string]string `toml:"roles" json:"roles"`
	Claims struct {
		AllowedStages []string `toml:"allowed_stages" json:"allowed_stages"`
	} `toml:"claims" json:"claims"`
	Lens struct {
		Pullback struct {
			MinIndependent int     `toml:"min_independent" json:"min_independent"`
			NumericRelTol  float64 `toml:"numeric_rel_tol" json:"numeric_rel_tol"`
			AsOfWindow     string  `toml:"as_of_window" json:"as_of_window"`
		} `toml:"pullback" json:"pullback"`
	} `toml:"lens" json:"lens"`
	Graph struct {
		MergeThreshold float64 `toml:"merge_threshold" json:"merge_threshold"`
	} `toml:"graph" json:"graph"`
	Independence struct {
		StripPrefixes []string `toml:"strip_prefixes" json:"strip_prefixes"`
	} `toml:"independence" json:"independence"`
	Gates struct {
		G2 struct {
			InjectionPatterns []string `toml:"injection_patterns" json:"injection_patterns"`
		} `toml:"g2" json:"g2"`
		G6 struct {
			Terms   []string `toml:"terms" json:"terms"`
			Tokens  []string `toml:"tokens" json:"tokens"`
			Phrases []string `toml:"phrases" json:"phrases"`
		} `toml:"g6" json:"g6"`
	} `toml:"gates" json:"gates"`
	Ledger struct {
		Kinds []string `toml:"kinds" json:"kinds"`
	} `toml:"ledger" json:"ledger"`
	Artifact struct {
		Forms  map[string]string  `toml:"forms" json:"forms"`
		Layout map[string]float64 `toml:"layout" json:"layout"`
	} `toml:"artifact" json:"artifact"`
	Report struct {
		Gates       []string `toml:"gates" json:"gates"`
		Implemented []string `toml:"implemented" json:"implemented"`
	} `toml:"report" json:"report"`
	Check struct {
		AtomCount     int `toml:"atom_count" json:"atom_count"`
		SkillMaxWords int `toml:"skill_max_words" json:"skill_max_words"`
	} `toml:"check" json:"check"`
	Units     map[string]UnitDim   `toml:"units" json:"units"`
	UnitCodes map[string]UnitCodes `toml:"unit_codes" json:"unit_codes"`

	// Hash is the sha256 of the bytes loaded; Source says where they came from.
	Hash   string `toml:"-" json:"-"`
	Source string `toml:"-" json:"-"`
}

// EnvOverride names the environment variable that points at an override config file.
const EnvOverride = "META_CONFIG"

// Load reads the effective config: META_CONFIG if set, else the embedded defaults.
func Load() (*Config, error) {
	if p := os.Getenv(EnvOverride); p != "" {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("config: %s=%s: %w", EnvOverride, p, err)
		}
		return Parse(b, p)
	}
	return Parse(defaults.TOML, "embedded:config/defaults.toml")
}

// Parse decodes and checks config bytes.
func Parse(b []byte, source string) (*Config, error) {
	var c Config
	md, err := toml.NewDecoder(bytes.NewReader(b)).Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("config %s: %w", source, err)
	}
	if und := md.Undecoded(); len(und) > 0 {
		keys := make([]string, len(und))
		for i, k := range und {
			keys[i] = k.String()
		}
		sort.Strings(keys)
		return nil, fmt.Errorf("config %s: unknown keys %v", source, keys)
	}
	required := []string{
		"schema", "mode_aliases", "loop.d2_reentry_max", "intake.re_expand_max", "roles", "claims.allowed_stages",
		"lens.pullback.min_independent", "lens.pullback.numeric_rel_tol", "lens.pullback.as_of_window", "graph.merge_threshold",
		"independence.strip_prefixes", "gates.g2.injection_patterns", "ledger.kinds", "gates.g6.terms", "gates.g6.tokens", "gates.g6.phrases", "artifact.forms",
		"artifact.layout", "report.gates", "report.implemented", "check.atom_count", "check.skill_max_words", "units", "unit_codes",
	}
	for _, m := range []string{"fast", "deep_narrow", "deep_broad"} {
		for _, k := range []string{"review_passes", "max_regions", "k", "g3_layers", "g3_pair_cap", "g4", "r2",
			"loaded_layers", "p50_minutes", "p90_minutes"} {
			required = append(required, "modes."+m+"."+k)
		}
	}
	for _, k := range required {
		if !md.IsDefined(splitKey(k)...) {
			return nil, fmt.Errorf("config %s: missing required key %s", source, k)
		}
	}
	for name, m := range c.Modes {
		if m.K < 1 || m.K%2 == 0 {
			return nil, fmt.Errorf("config %s: modes.%s.k must be odd and >= 1", source, name)
		}
		if m.MaxRegions < 1 || m.ReviewPasses < 1 {
			return nil, fmt.Errorf("config %s: modes.%s caps must be >= 1", source, name)
		}
	}
	for alias, target := range c.ModeAliases {
		if _, ok := c.Modes[target]; !ok {
			return nil, fmt.Errorf("config %s: mode_aliases.%s -> unknown mode %q", source, alias, target)
		}
	}
	for i, p := range c.Gates.G2.InjectionPatterns {
		if _, err := regexp.Compile("(?i)" + p); err != nil {
			return nil, fmt.Errorf("config %s: gates.g2.injection_patterns[%d]: %v", source, i, err)
		}
	}
	switch c.Lens.Pullback.AsOfWindow {
	case "day", "month", "year", "none":
	default:
		return nil, fmt.Errorf("config %s: lens.pullback.as_of_window must be day, month, year or none", source)
	}
	for dim, uc := range c.UnitCodes {
		if _, clash := c.Units[dim]; clash {
			return nil, fmt.Errorf("config %s: unit dimension %q is in both units and unit_codes", source, dim)
		}
		re, err := regexp.Compile(uc.Pattern)
		if err != nil || uc.Pattern == "" {
			return nil, fmt.Errorf("config %s: unit_codes.%s.pattern must be a non-empty RE2 pattern", source, dim)
		}
		if uc.List != "" && uc.List != "iso4217" {
			return nil, fmt.Errorf("config %s: unit_codes.%s.list must be \"iso4217\" or empty", source, dim)
		}
		uc.re = re
		c.UnitCodes[dim] = uc
	}
	if c.Lens.Pullback.MinIndependent < 1 {
		return nil, fmt.Errorf("config %s: lens.pullback.min_independent must be >= 1", source)
	}
	c.Hash = fsx.SHA256(b)
	c.Source = source
	return &c, nil
}

func splitKey(k string) []string {
	var out []string
	start := 0
	for i := 0; i < len(k); i++ {
		if k[i] == '.' {
			out = append(out, k[start:i])
			start = i + 1
		}
	}
	return append(out, k[start:])
}

// ResolveMode maps a user-facing mode (fast, deep) or a table row name to a modes.<mode> row name.
func (c *Config) ResolveMode(m string) (string, error) {
	if t, ok := c.ModeAliases[m]; ok {
		return t, nil
	}
	if _, ok := c.Modes[m]; ok {
		return m, nil
	}
	names := make([]string, 0, len(c.ModeAliases)+len(c.Modes))
	for k := range c.ModeAliases {
		names = append(names, k)
	}
	for k := range c.Modes {
		names = append(names, k)
	}
	sort.Strings(names)
	return "", fmt.Errorf("unknown mode %q (one of %v)", m, names)
}

// Layout returns a required artifact.layout threshold.
func (c *Config) LayoutThreshold(key string) (float64, error) {
	v, ok := c.Artifact.Layout[key]
	if !ok {
		return 0, fmt.Errorf("config: missing artifact.layout.%s", key)
	}
	return v, nil
}
