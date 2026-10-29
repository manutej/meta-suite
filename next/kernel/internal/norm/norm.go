// Package norm holds the kernel's normalizers: entity labels, claim keys and independence keys. They are
// deliberately simple and documented, because every merge and every "independent" count rests on them.
package norm

import (
	"net/url"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Label normalizes an entity label for exact-label matching: Unicode NFKC, lower case, whitespace runs
// collapsed to one space, and leading/trailing punctuation and symbols removed ("The  Foo." -> "the foo").
func Label(s string) string {
	s = norm.NFKC.String(s)
	s = strings.ToLower(s)
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimFunc(s, func(r rune) bool { return unicode.IsPunct(r) || unicode.IsSymbol(r) || unicode.IsSpace(r) })
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
