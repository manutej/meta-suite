package norm

import (
	"strings"
	"testing"
)

func TestLabel(t *testing.T) {
	cases := map[string]string{
		"  The   Foo. ": "the foo",
		"ＡＢＣ":           "abc", // NFKC folds full-width letters
		"\"Entity-A\"":  "entity-a",
		"entity a":      "entity a",
		"Entity\tA\n":   "entity a",
		"(entity b)":    "entity b",
		"":              "",
		// D-087: brackets stay balanced.
		"transparency obligations (Article 50)":  "transparency obligations (article 50)",
		"transparency obligations (Article 50).": "transparency obligations (article 50)",
		"Commission guidelines (Article 6)":      "commission guidelines (article 6)",
		"(Article 50) duties":                    "(article 50) duties",
		"rule (a)(b)":                            "rule (a)(b)",
		"((entity c))":                           "entity c",
		"(a) and (b)":                            "(a) and (b)",
		"article 50)":                            "article 50",
		"(article 50":                            "article 50",
		"[draft] code of practice":               "[draft] code of practice",
		"\"so-called\" agents":                   "\"so-called\" agents",
		"“quoted”":                               "quoted",
		"members' states'":                       "members' states",
		"...":                                    "",
		"()":                                     "",
		")x(":                                    "x",
	}
	for in, want := range cases {
		if got := Label(in); got != want {
			t.Errorf("Label(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestLabelBalanced: trimming never unbalances a label: every input below is balanced inside its outer junk,
// and so is its normalized form (D-087). Unbalanced brackets inside a label are the source's own text and stay.
func TestLabelBalanced(t *testing.T) {
	for _, in := range []string{"x (y)", "(x) y)", "((x)", "x (y (z))", "x [y]}", ")x (y)(", "a (b) c (d).", "\"x\" y\"", "<x> (y)!", "x (y)…"} {
		out := Label(in)
		for o, c := range closer {
			if o == '"' {
				if strings.Count(out, `"`)%2 != 0 {
					t.Errorf("Label(%q) = %q: odd quotes", in, out)
				}
				continue
			}
			depth := 0
			for _, r := range out {
				if r == o {
					depth++
				} else if r == c {
					depth--
				}
				if depth < 0 {
					break
				}
			}
			if depth != 0 {
				t.Errorf("Label(%q) = %q: unbalanced %c%c", in, out, o, c)
			}
		}
	}
}

func TestIsDateAndWindow(t *testing.T) {
	for in, want := range map[string]bool{"2026-08-02": true, "2026-08": true, "2026-08-02T10:00Z": true, "2026": false, "Article 50": false, "2026-8-2": false, "": false} {
		if got := IsDate(in); got != want {
			t.Errorf("IsDate(%q) = %v", in, got)
		}
	}
	cases := [][3]string{{"2026-08-02", "month", "2026-08"}, {"2026-08-02", "day", "2026-08-02"}, {"2026-08-02", "year", "2026"},
		{"2026-08-02", "none", ""}, {"", "month", ""}, {"2026", "month", "2026"},
		// D-097: timestamps are normalized to UTC before windowing.
		{"2026-01-31T23:00-05:00", "month", "2026-02"}, {"2026-02-01T04:00Z", "month", "2026-02"},
		{"2026-01-31T23:00-05:00", "day", "2026-02-01"}, {"2026-12-31T22:00:00-03:00", "year", "2027"},
		{"2026-05-31 10:00", "day", "2026-05-31"}}
	for _, c := range cases {
		if got := Window(c[0], c[1]); got != c[2] {
			t.Errorf("Window(%q,%q) = %q, want %q", c[0], c[1], got, c[2])
		}
	}
}

func TestWindowsOverlapAndHull(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want bool
	}{{"2026", "2026-03", true}, {"2026-03", "2026", true}, {"2026-05", "2026-06", false}, {"", "2026-06", true},
		{"2026-05", "2026-05", true}, {"2025", "2026-01", false}, {"2026-05-31", "2026-06-01", false}} {
		if got := WindowsOverlap(c.a, c.b); got != c.want {
			t.Errorf("WindowsOverlap(%q,%q) = %v", c.a, c.b, got)
		}
	}
	for _, c := range []struct {
		in   []string
		want string
	}{{nil, ""}, {[]string{"", ""}, ""}, {[]string{"2026-03"}, "2026-03"}, {[]string{"2026", "2026-03"}, "2026"},
		{[]string{"2026-03", "2026-05", ""}, "2026"}, {[]string{"2026-03-01", "2026-03-09"}, "2026-03"}, {[]string{"2025-06", "2026-04"}, "2025-06/2026-04"}} {
		if got := WindowHull(c.in); got != c.want {
			t.Errorf("WindowHull(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestIndependenceKey(t *testing.T) {
	strip := []string{"www."}
	cases := map[string]string{
		"www.Example.org":                 "example.org",
		"https://www.example.org/a/b?c=1": "example.org",
		"example.org.":                    "example.org",
		" pub-a.example ":                 "pub-a.example",
		"www.":                            "www", // trailing dot trimmed, nothing left to strip
	}
	for in, want := range cases {
		if got := IndependenceKey(in, strip); got != want {
			t.Errorf("IndependenceKey(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestYear(t *testing.T) {
	for in, want := range map[string]string{"2024-06-01": "2024", "2024": "2024", "24": "", "n/a 2024": ""} {
		if got := Year(in); got != want {
			t.Errorf("Year(%q) = %q", in, got)
		}
	}
}
