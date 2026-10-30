// Package norm holds the kernel's normalizers: entity labels, claim keys and independence keys. They are
// deliberately simple and documented, because every merge and every "independent" count rests on them.
package norm

import (
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Label normalizes an entity label for exact-label matching: Unicode NFKC, lower case, whitespace runs
// collapsed to one space, and leading/trailing punctuation and symbols removed ("The  Foo." -> "the foo").
// Brackets and paired quotes are removed only when that keeps the label balanced (D-087): a pair wrapping the
// whole label goes ("(entity b)" -> "entity b"), an unmatched one at either end goes, and a matched one stays
// ("transparency obligations (Article 50)." -> "transparency obligations (article 50)").
func Label(s string) string {
	s = norm.NFKC.String(s)
	s = strings.ToLower(s)
	s = strings.Join(strings.Fields(s), " ")
	for {
		r := []rune(strings.TrimFunc(s, func(r rune) bool { return trimmable(r) }))
		if len(r) == 0 {
			return ""
		}
		first, last := r[0], r[len(r)-1]
		switch {
		case len(r) >= 2 && closer[first] == last && matchAt(r, 0) == len(r)-1:
			r = r[1 : len(r)-1] // one pair wraps the whole label
		case isCloser(last) && matchAt(r, len(r)-1) < 0:
			r = r[:len(r)-1] // unmatched closer at the end
		case closer[first] != 0 && matchAt(r, 0) < 0, isCloser(first) && closer[first] == 0:
			r = r[1:] // unmatched opener, or any closer, at the start
		case closer[last] != 0 && !isCloser(last):
			r = r[:len(r)-1] // an opener at the end opens nothing
		default:
			return string(r)
		}
		s = string(r)
	}
}

// closer maps each opening bracket or quote to its closing partner. The straight double quote pairs with
// itself; the apostrophe is ordinary punctuation (it is also a possessive).
var closer = map[rune]rune{'(': ')', '[': ']', '{': '}', '<': '>', '“': '”', '‘': '’', '«': '»', '"': '"'}

func isCloser(r rune) bool {
	for _, c := range closer {
		if c == r {
			return true
		}
	}
	return false
}

// trimmable is outer punctuation removed unconditionally: everything except brackets and paired quotes.
func trimmable(r rune) bool {
	if closer[r] != 0 || isCloser(r) {
		return false
	}
	return unicode.IsPunct(r) || unicode.IsSymbol(r) || unicode.IsSpace(r)
}

// matchAt returns the index of the partner of the bracket at i, or -1. Straight double quotes pair in order.
func matchAt(r []rune, i int) int {
	c := r[i]
	if c == '"' {
		var idx []int
		for j, x := range r {
			if x == '"' {
				idx = append(idx, j)
			}
		}
		for k, j := range idx {
			if j == i {
				if k%2 == 0 && k+1 < len(idx) {
					return idx[k+1]
				}
				if k%2 == 1 {
					return idx[k-1]
				}
			}
		}
		return -1
	}
	if cl, ok := closer[c]; ok { // opener: scan forward
		depth := 0
		for j := i; j < len(r); j++ {
			switch r[j] {
			case c:
				depth++
			case cl:
				depth--
				if depth == 0 {
					return j
				}
			}
		}
		return -1
	}
	var op rune // closer: scan backward
	for o, cl := range closer {
		if cl == c {
			op = o
		}
	}
	depth := 0
	for j := i; j >= 0; j-- {
		switch r[j] {
		case c:
			depth++
		case op:
			depth--
			if depth == 0 {
				return j
			}
		}
	}
	return -1
}

var isoDate = regexp.MustCompile(`^\d{4}-\d{2}(-\d{2}([T ]\d{2}:\d{2}(:\d{2}(\.\d+)?)?(Z|[+-]\d{2}:?\d{2})?)?)?$`)

// IsDate reports whether s is an ISO 8601 calendar date with at least a month ("2026-08", "2026-08-02",
// "2026-08-02T10:00Z"). A bare year is not a date here: it is too easily a count or an article number.
func IsDate(s string) bool { return isoDate.MatchString(strings.TrimSpace(s)) }

// Window truncates an as_of value to the configured agreement window (lens.pullback.as_of_window): "year"
// keeps YYYY, "month" YYYY-MM, "day" YYYY-MM-DD; "none" or an empty as_of gives "". An as_of coarser than the
// window keeps its own precision ("2026" stays "2026" under "month").
func Window(asOf, window string) string {
	asOf = strings.TrimSpace(asOf)
	n := map[string]int{"year": 4, "month": 7, "day": 10}[window]
	if asOf == "" || n == 0 {
		return ""
	}
	if len(asOf) < n {
		return asOf
	}
	return asOf[:n]
}

// IndependenceKey normalizes a publisher key: a URL becomes its host; then lower case, trimmed, trailing dots
// removed and each configured prefix (e.g. "www.") stripped once.
func IndependenceKey(s string, stripPrefixes []string) string {
	s = strings.TrimSpace(s)
	if strings.Contains(s, "://") {
		if u, err := url.Parse(s); err == nil && u.Host != "" {
			s = u.Hostname()
		}
	}
	s = strings.TrimRight(strings.ToLower(s), ".")
	for _, p := range stripPrefixes {
		if strings.HasPrefix(s, p) && len(s) > len(p) {
			s = s[len(p):]
			break
		}
	}
	return s
}

// Year returns the leading 4-digit year of an as_of value, or "" when absent or malformed.
func Year(asOf string) string {
	if len(asOf) < 4 {
		return ""
	}
	for _, c := range asOf[:4] {
		if c < '0' || c > '9' {
			return ""
		}
	}
	return asOf[:4]
}
