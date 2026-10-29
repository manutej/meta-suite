package norm

import "testing"

func TestLabel(t *testing.T) {
	cases := map[string]string{
		"  The   Foo. ": "the foo",
		"ＡＢＣ":           "abc", // NFKC folds full-width letters
		"\"Entity-A\"":  "entity-a",
		"entity a":      "entity a",
		"Entity\tA\n":   "entity a",
		"(entity b)":    "entity b",
		"":              "",
	}
	for in, want := range cases {
		if got := Label(in); got != want {
			t.Errorf("Label(%q) = %q, want %q", in, got, want)
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
