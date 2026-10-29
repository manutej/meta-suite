// Package ids issues ULIDs (26 chars, Crockford base32: 48-bit millisecond time + 80 random bits).
// Within one process, IDs are strictly increasing (monotonic in the same millisecond).
package ids

import (
	"crypto/rand"
	"sync"
	"time"
)

const alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

var (
	mu       sync.Mutex
	lastMs   uint64
	lastRand [10]byte
)

// Now is the clock; tests may replace it.
var Now = time.Now

// New returns a fresh ULID.
func New() string {
	mu.Lock()
	defer mu.Unlock()
	ms := uint64(Now().UnixMilli())
	var r [10]byte
	if ms <= lastMs {
		ms = lastMs
		r = lastRand
		for i := 9; i >= 0; i-- { // increment the 80-bit random part
			r[i]++
			if r[i] != 0 {
				break
			}
		}
	} else if _, err := rand.Read(r[:]); err != nil {
		panic("ids: crypto/rand failed: " + err.Error())
	}
	lastMs, lastRand = ms, r
	return encode(ms, r)
}

func encode(ms uint64, r [10]byte) string {
	var b [16]byte
	b[0], b[1], b[2], b[3], b[4], b[5] = byte(ms>>40), byte(ms>>32), byte(ms>>24), byte(ms>>16), byte(ms>>8), byte(ms)
	copy(b[6:], r[:])
	out := make([]byte, 26)
	// 128 bits -> 26 base32 digits (the first digit carries 3 bits).
	var acc uint64
	bits := uint(0)
	pos := 25
	for i := 15; i >= 0; i-- {
		acc |= uint64(b[i]) << bits
		bits += 8
		for bits >= 5 && pos >= 0 {
			out[pos] = alphabet[acc&31]
			acc >>= 5
			bits -= 5
			pos--
		}
	}
	for pos >= 0 {
		out[pos] = alphabet[acc&31]
		acc >>= 5
		pos--
	}
	return string(out)
}

// Valid reports whether s is a well-formed ULID.
func Valid(s string) bool {
	if len(s) != 26 || s[0] > '7' {
		return false
	}
	for i := 0; i < len(s); i++ {
		ok := false
		for j := 0; j < len(alphabet); j++ {
			if s[i] == alphabet[j] {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}
