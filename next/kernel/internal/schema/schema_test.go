package schema

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Fixtures in testdata/ are SYNTHETIC (hand-written to exercise the schemas); see testdata/README.txt.

func TestEveryCoreTypeEmbedded(t *testing.T) {
	want := []string{"ArtifactPlan", "Brief", "Claim", "ClaimInput", "Conflicts", "Cover", "Coverage", "DomainGraph", "EditScript",
		"Finding", "HowChecked", "Interface", "Judgment", "LedgerEvent", "Probes", "RegionGraph", "Retraction", "RunRecord",
		"SectionMap", "Spec"}
	if got := Types(); !reflect.DeepEqual(got, want) {
		t.Fatalf("types = %v", got)
	}
}

func TestFixtures(t *testing.T) {
	files, _ := filepath.Glob("testdata/*.json")
	if len(files) < 20 {
		t.Fatalf("only %d fixtures", len(files))
	}
	for _, f := range files {
		base := filepath.Base(f)
		parts := strings.SplitN(strings.TrimSuffix(base, ".json"), ".", 3) // <type>.<valid|invalid>.<case>
		if len(parts) != 3 {
			t.Fatalf("bad fixture name %s", base)
		}
		b, _ := os.ReadFile(f)
		errs, err := Validate(parts[0], b)
		if err != nil {
			t.Fatalf("%s: %v", base, err)
		}
		switch parts[1] {
		case "valid":
			if len(errs) != 0 {
				t.Errorf("%s: want valid, got %+v", base, errs)
			}
		case "invalid":
			if len(errs) == 0 {
				t.Errorf("%s: want invalid, got valid", base)
			}
		}
	}
}

func TestDeterministicErrorOrder(t *testing.T) {
	b, _ := os.ReadFile("testdata/claim.invalid.many_errors.json")
	first, _ := Validate("Claim", b)
	if len(first) < 3 {
		t.Fatalf("want several errors, got %+v", first)
	}
	for i := 0; i < 50; i++ {
		again, _ := Validate("claim", b)
		if !reflect.DeepEqual(first, again) {
			t.Fatalf("run %d differs:\n%+v\n%+v", i, first, again)
		}
	}
}

func TestUnknownType(t *testing.T) {
	if _, err := Validate("Nope", []byte("{}")); err == nil {
		t.Fatal("unknown type accepted")
	}
}

func TestArray(t *testing.T) {
	errs, n, err := ValidateArray("SectionMap", []byte(`[{"sections":[{"id":"a","tag":"main","finding_ids":[]}]},{"sections":[]}]`))
	if err != nil || n != 2 || len(errs) == 0 || *errs[0].Index != 1 {
		t.Fatalf("errs=%+v n=%d err=%v", errs, n, err)
	}
}
