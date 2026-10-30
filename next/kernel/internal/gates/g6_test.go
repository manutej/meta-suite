package gates

import (
	"strconv"
	"strings"
	"testing"

	defaults "github.com/manutej/meta-suite/next/kernel/config"
	"github.com/manutej/meta-suite/next/kernel/internal/config"
)

// SYNTHETIC text fixtures: sentences written to contain (or avoid) each class of banned vocabulary.

func cfg(t *testing.T) G6Config {
	c, err := config.Parse(defaults.TOML, "embedded")
	if err != nil {
		t.Fatal(err)
	}
	return G6Config{Terms: c.Gates.G6.Terms, Tokens: c.Gates.G6.Tokens, Phrases: c.Gates.G6.Phrases}
}

func scan(t *testing.T, files map[string]string, allow []Allow) *G6Report {
	names := make([]string, 0, len(files))
	for k := range files {
		names = append(names, k)
	}
	rep, err := G6(names, func(f string) ([]byte, error) { return []byte(files[f]), nil }, cfg(t), allow)
	if err != nil {
		t.Fatal(err)
	}
	return rep
}

func TestStemsTokensPhrases(t *testing.T) {
	rep := scan(t, map[string]string{"a.md": "The sheaves agree.\nA functorial view.\nH0 is 1 but H01 is a code.\nThe square fails  to\ncommute here.\nA morphism."}, nil)
	got := []string{}
	for _, h := range rep.Hits {
		got = append(got, h.Term+"@"+strings.Join([]string{itoa(h.Line), itoa(h.Col)}, ":"))
	}
	want := "sheav@1:5 functor@2:3 H0@3:1 fails to commute@4:12 morphism@6:3"
	if strings.Join(got, " ") != want {
		t.Fatalf("hits %q\nwant %q", strings.Join(got, " "), want)
	}
	if rep.Verdict != "fail" {
		t.Fatal("verdict")
	}
}

func TestCleanTextPasses(t *testing.T) {
	rep := scan(t, map[string]string{"a.md": "Independent sources agree on the figure. The stalking horse bid closed."}, nil)
	// "stalking" starts with the stem "stalk": a known false positive the allowlist exists for.
	if len(rep.Hits) != 1 || rep.Hits[0].Match != "stalking" {
		t.Fatalf("%+v", rep.Hits)
	}
	rep = scan(t, map[string]string{"a.md": "The stalking horse bid closed."}, []Allow{{Term: "stalking", Reason: "finance term in the brief"}})
	if rep.Verdict != "pass" || len(rep.Allowed) != 1 || rep.Allowed[0].Reason == "" {
		t.Fatalf("%+v", rep)
	}
}

func TestAllowByConfigTerm(t *testing.T) {
	rep := scan(t, map[string]string{"a.md": "Harmonic minor and harmonics."}, []Allow{{Term: "harmonic", Reason: "music brief"}})
	if rep.Verdict != "pass" || len(rep.Allowed) != 2 {
		t.Fatalf("%+v", rep)
	}
}

func TestHTMLMasksMarkupButScansVisibleText(t *testing.T) {
	html := `<!doctype html><html><head><title>Morphism notes</title>
<style>.sheaf{color:red}</style><script>var pullback=1;</script>
<script type="application/json">{"note":"a cocone"}</script></head>
<body class="functor"><!-- yoneda --><p data-x="operad">Plain text.</p><img alt="a stalk diagram" src="x.png"></body></html>`
	rep := scan(t, map[string]string{"p.html": html}, nil)
	terms := []string{}
	for _, h := range rep.Hits {
		terms = append(terms, h.Term)
	}
	if strings.Join(terms, ",") != "morphism,cocone,stalk" {
		t.Fatalf("terms %v", terms)
	}
	if rep.Hits[0].Line != 1 || rep.Hits[0].Col != 35 {
		t.Fatalf("location %+v", rep.Hits[0])
	}
}

func TestDeterministic(t *testing.T) {
	files := map[string]string{"b.md": "sheaf", "a.md": "operad and sheaf"}
	a := scan(t, files, nil)
	for i := 0; i < 20; i++ {
		b := scan(t, files, nil)
		if len(a.Hits) != len(b.Hits) || a.Hits[0] != b.Hits[0] || a.Files[0] != "a.md" {
			t.Fatal("order differs")
		}
	}
}

func itoa(i int) string { return strconv.Itoa(i) }

// D-097: <internal> blocks are skipped only in skill Markdown and, in HTML, only inside a comment; a bare
// <internal> element in HTML renders, and other Markdown is a reader's file.
func TestInternalBlocksSkipped(t *testing.T) {
	rep := scan(t, map[string]string{
		"skills/x/ops/D1.md": "Plain words.\n<!-- <internal>A pullback over the\nsheaf of claims.</internal> -->\nA morphism here.",
		"skills/x/SKILL.md":  "<internal>operad</internal> fine",
		"run/assumptions.md": "<internal>cochain</internal>",
		"page.html":          "<p>visible</p><internal>colimit</internal><p>A functor.</p>",
		"commented.html":     "<p>ok</p><!-- <internal>adjunction</internal> --><p>fine</p>",
	}, nil)
	got := []string{}
	for _, h := range rep.Hits {
		got = append(got, h.File+":"+h.Term+"@"+itoa(h.Line))
	}
	if strings.Join(got, " ") != "page.html:colimit@1 page.html:functor@1 run/assumptions.md:cochain@1 skills/x/ops/D1.md:morphism@4" {
		t.Fatalf("hits %v", got)
	}
}
