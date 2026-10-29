package cover

import (
	"strings"
	"testing"
)

// SYNTHETIC covers (element ids are placeholders), exercising each D-042 check.

const tree = `"tree":{"id":"root","question":"q","answer_type":"judgment","children":[{"id":"q1","question":"a","answer_type":"fact","children":[]},{"id":"q2","question":"b","answer_type":"fact","children":[]},{"id":"q3","question":"c","answer_type":"fact","children":[]}]}`

func doc(regions string) []byte {
	return []byte(`{"universe":["q1","q2","q3","kind:x"],` + tree + `,"regions":` + regions + `,"relation_vocab":[{"name":"r","value_type":"qualitative"}]}`)
}

func TestConnectedCoverComputesOverlaps(t *testing.T) {
	rep, err := Check(doc(`[{"id":"a","elements":["q1","kind:x"]},{"id":"b","elements":["q2","kind:x"]},{"id":"c","elements":["q3","kind:x","q2"]}]`), 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Issues) != 0 {
		t.Fatalf("issues: %+v", rep.Issues)
	}
	if len(rep.Overlaps) != 3 || len(rep.TripleOverlaps) != 1 {
		t.Fatalf("overlaps %+v triples %+v", rep.Overlaps, rep.TripleOverlaps)
	}
	if got := strings.Join(rep.Overlaps[1].Elements, ","); got != "kind:x" { // a&c
		t.Fatalf("a&c = %s", got)
	}
	if got := strings.Join(rep.Overlaps[2].Elements, ","); got != "kind:x,q2" { // b&c
		t.Fatalf("b&c = %s", got)
	}
}

func TestUncoveredAndDisconnected(t *testing.T) {
	rep, _ := Check(doc(`[{"id":"a","elements":["q1"]},{"id":"b","elements":["q2","kind:x"]}]`), 6)
	checks := map[string]bool{}
	for _, i := range rep.Issues {
		checks[i.Check] = true
	}
	if !checks["covering"] || !checks["overlap_connected"] {
		t.Fatalf("issues: %+v", rep.Issues)
	}
}

func TestMaxRegionsAndForeignElement(t *testing.T) {
	rep, _ := Check(doc(`[{"id":"a","elements":["q1","q2","q3","kind:x","zz"]},{"id":"b","elements":["q1"]}]`), 1)
	checks := map[string]bool{}
	for _, i := range rep.Issues {
		checks[i.Check] = true
	}
	if !checks["max_regions"] || !checks["regions_within_universe"] {
		t.Fatalf("issues: %+v", rep.Issues)
	}
}

func TestDeclaredOverlapsMustMatch(t *testing.T) {
	d := strings.Replace(string(doc(`[{"id":"a","elements":["q1","q2","q3","kind:x"]},{"id":"b","elements":["q1"]}]`)), `"relation_vocab"`, `"overlaps":[{"regions":["a","b"],"elements":["q2"]}],"relation_vocab"`, 1)
	rep, _ := Check([]byte(d), 0)
	if len(rep.Issues) != 1 || rep.Issues[0].Check != "overlaps_computed" {
		t.Fatalf("issues: %+v", rep.Issues)
	}
}
