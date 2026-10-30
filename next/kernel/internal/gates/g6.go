package gates

import (
	"bytes"

	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// G6Config is the vocabulary from config gates.g6.
type G6Config struct {
	Terms, Tokens, Phrases []string
}

// Allow is one exemption: a term from the brief's vocab_allowlist (or --allow), with its reason.
type Allow struct {
	Term   string `json:"term"`
	Reason string `json:"reason"`
}

// Hit is one banned-vocabulary match.
type Hit struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Col     int    `json:"col"`
	Match   string `json:"match"`
	Term    string `json:"term"`
	Kind    string `json:"kind"` // term | token | phrase
	Context string `json:"context"`
	Reason  string `json:"allowed_because,omitempty"`
}

// G6Report is G6's output.
type G6Report struct {
	Gate    string   `json:"gate"`
	Verdict string   `json:"verdict"`
	Files   []string `json:"files"`
	Hits    []Hit    `json:"hits"`
	Allowed []Hit    `json:"allowed"`
	Allow   []Allow  `json:"allowlist"`
}

type matcher struct {
	re   *regexp.Regexp
	term string
	kind string
}

func compile(c G6Config) []matcher {
	var ms []matcher
	for _, t := range c.Terms {
		ms = append(ms, matcher{regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(t) + `[\p{L}\p{N}_-]*`), t, "term"})
	}
	for _, t := range c.Tokens {
		ms = append(ms, matcher{regexp.MustCompile(`(?:^|[^\p{L}\p{N}_^])(` + regexp.QuoteMeta(t) + `)(?:$|[^\p{L}\p{N}_])`), t, "token"})
	}
	for _, p := range c.Phrases {
		words := strings.Fields(p)
		for i := range words {
			words[i] = regexp.QuoteMeta(words[i])
		}
		ms = append(ms, matcher{regexp.MustCompile(`(?i)\b` + strings.Join(words, `\s+`) + `\b`), p, "phrase"})
	}
	return ms
}

// G6 scans files for banned vocabulary. HTML files are scanned on what a reader sees: element text and the
// alt, title, aria-label, placeholder and content attributes, plus JSON data blocks; markup, CSS and script
// code are masked (offsets are preserved so line:col point into the original file). Other files are scanned
// whole. A hit is exempt when its config term or its matched word is on the allowlist.
func G6(files []string, read func(string) ([]byte, error), c G6Config, allow []Allow) (*G6Report, error) {
	rep := &G6Report{Gate: "G6", Files: []string{}, Hits: []Hit{}, Allowed: []Hit{}, Allow: allow}
	if rep.Allow == nil {
		rep.Allow = []Allow{}
	}
	sort.Slice(rep.Allow, func(i, j int) bool { return rep.Allow[i].Term < rep.Allow[j].Term })
	allowed := map[string]string{}
	for _, a := range allow {
		allowed[strings.ToLower(strings.TrimSpace(a.Term))] = a.Reason
	}
	ms := compile(c)
	sorted := append([]string(nil), files...)
	sort.Strings(sorted)
	for _, f := range sorted {
		b, err := read(f)
		if err != nil {
			return nil, err
		}
		rep.Files = append(rep.Files, f)
		text := MaskInternal(b)
		ext := strings.ToLower(filepath.Ext(f))
		if ext == ".html" || ext == ".htm" || ext == ".svg" {
			text = MaskHTML(text)
		}
		lineStarts := []int{0}
		for i, ch := range b {
			if ch == '\n' {
				lineStarts = append(lineStarts, i+1)
			}
		}
		for _, m := range ms {
			for _, loc := range m.re.FindAllSubmatchIndex(text, -1) {
				s, e := loc[0], loc[1]
				if m.kind == "token" {
					s, e = loc[2], loc[3]
				}
				line := sort.Search(len(lineStarts), func(i int) bool { return lineStarts[i] > s })
				col := s - lineStarts[line-1] + 1
				word := string(b[s:e])
				h := Hit{File: f, Line: line, Col: col, Match: word, Term: m.term, Kind: m.kind, Context: context(text, s, e)}
				if r, ok := allowed[strings.ToLower(m.term)]; ok {
					h.Reason = r
				} else if r, ok := allowed[strings.ToLower(word)]; ok {
					h.Reason = r
				}
				if h.Reason != "" {
					rep.Allowed = append(rep.Allowed, h)
				} else {
					rep.Hits = append(rep.Hits, h)
				}
			}
		}
	}
	for _, list := range [][]Hit{rep.Hits, rep.Allowed} {
		sort.Slice(list, func(i, j int) bool {
			a, b := list[i], list[j]
			if a.File != b.File {
				return a.File < b.File
			}
			if a.Line != b.Line {
				return a.Line < b.Line
			}
			if a.Col != b.Col {
				return a.Col < b.Col
			}
			return a.Term < b.Term
		})
	}
	rep.Hits = dedupeHits(rep.Hits)
	rep.Allowed = dedupeHits(rep.Allowed)
	rep.Verdict = "pass"
	if len(rep.Hits) > 0 {
		rep.Verdict = "fail"
	}
	return rep, nil
}

// dedupeHits keeps one hit per location (the first term in sort order), so "sheaf" and "sheav" style
// overlapping stems report once.
func dedupeHits(h []Hit) []Hit {
	out := h[:0]
	for i, x := range h {
		if i > 0 && x.File == h[i-1].File && x.Line == h[i-1].Line && x.Col == h[i-1].Col {
			continue
		}
		out = append(out, x)
	}
	return out
}

func context(text []byte, s, e int) string {
	a, b := s-30, e+30
	if a < 0 {
		a = 0
	}
	if b > len(text) {
		b = len(text)
	}
	return strings.Join(strings.Fields(string(text[a:b])), " ")
}

var (
	blockRe   = regexp.MustCompile(`(?is)<(script|style)\b([^>]*)>(.*?)</(script|style)\s*>`)
	commentRe = regexp.MustCompile(`(?s)<!--.*?-->`)
	tagRe     = regexp.MustCompile(`(?s)<[^>]*>`)
	attrRe    = regexp.MustCompile(`(?is)\b(alt|title|aria-label|placeholder|content)\s*=\s*("[^"]*"|'[^']*')`)
	jsonType  = regexp.MustCompile(`(?i)type\s*=\s*["']?application/(ld\+)?json`)
)

var internalRe = regexp.MustCompile(`(?is)<internal\b[^>]*>.*?</internal\s*>`)

// MaskInternal blanks every <internal>…</internal> block (D-097): text marked internal is for the model and
// maintainers, never shown to a reader, so G6 skips it in any file. Offsets and newlines are kept.
func MaskInternal(b []byte) []byte {
	locs := internalRe.FindAllIndex(b, -1)
	if len(locs) == 0 {
		return b
	}
	out := append([]byte(nil), b...)
	for _, l := range locs {
		for i := l[0]; i < l[1]; i++ {
			if out[i] != '\n' {
				out[i] = ' '
			}
		}
	}
	return out
}

// MaskHTML blanks everything in an HTML document a reader does not see, keeping byte offsets and newlines.
func MaskHTML(b []byte) []byte {
	out := append([]byte(nil), b...)
	blank := func(s, e int) {
		for i := s; i < e; i++ {
			if out[i] != '\n' {
				out[i] = ' '
			}
		}
	}
	keep := map[[2]int]bool{}
	for _, m := range blockRe.FindAllSubmatchIndex(b, -1) {
		isData := bytes.EqualFold(b[m[2]:m[3]], []byte("script")) && jsonType.Match(b[m[4]:m[5]])
		if isData {
			keep[[2]int{m[6], m[7]}] = true
			blank(m[0], m[6])
			blank(m[7], m[1])
		} else {
			blank(m[0], m[1])
		}
	}
	for _, m := range commentRe.FindAllIndex(out, -1) {
		blank(m[0], m[1])
	}
	for _, m := range tagRe.FindAllIndex(out, -1) {
		inData := false
		for k := range keep {
			if m[0] >= k[0] && m[1] <= k[1] {
				inData = true
			}
		}
		if inData {
			continue
		}
		vals := attrRe.FindAllSubmatchIndex(out[m[0]:m[1]], -1)
		cursor := m[0]
		for _, v := range vals {
			vs, ve := m[0]+v[4]+1, m[0]+v[5]-1 // inside the quotes
			blank(cursor, vs)
			cursor = ve
		}
		blank(cursor, m[1])
	}
	return out
}
