package ids

import (
	"testing"
	"time"
)

func TestULIDFormatAndMonotonic(t *testing.T) {
	prev := ""
	for i := 0; i < 2000; i++ {
		id := New()
		if !Valid(id) {
			t.Fatalf("invalid ULID %q", id)
		}
		if id <= prev {
			t.Fatalf("not increasing: %q after %q", id, prev)
		}
		prev = id
	}
}

func TestULIDTimePrefix(t *testing.T) {
	// 2^48-1 ms is the maximum time; its encoding starts with "7ZZZZZZZZZ".
	if got := encode(1<<48-1, [10]byte{}); got[:10] != "7ZZZZZZZZZ" || got[10:] != "0000000000000000" {
		t.Fatalf("encode max = %s", got)
	}
	old := Now
	defer func() { Now = old }()
	Now = func() time.Time { return time.UnixMilli(0) }
	lastMs = 0
	if id := New(); id[:10] != "0000000000" {
		t.Fatalf("epoch prefix = %s", id[:10])
	}
}
