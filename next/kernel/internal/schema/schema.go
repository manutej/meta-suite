// Package schema is gate G1: every typed object validates against its embedded JSON Schema (2020-12).
// Errors are sorted by line, instance location, keyword location, then message, so the same input always
// gives the same bytes (D-050).
package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"sync"

	"github.com/manutej/meta-suite/next/kernel/schemas"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Error is one validation error.
type Error struct {
	Line             int    `json:"line,omitempty"`
	Index            *int   `json:"index,omitempty"`
	InstanceLocation string `json:"instance_location"`
	KeywordLocation  string `json:"keyword_location"`
	Message          string `json:"message"`
}

var (
	once     sync.Once
	compiled map[string]*jsonschema.Schema
	names    map[string]string // lower -> canonical title
	loadErr  error
)

func load() {
	compiled = map[string]*jsonschema.Schema{}
	names = map[string]string{}
	c := jsonschema.NewCompiler()
	c.AssertFormat()
	files, err := fs.Glob(schemas.FS, "*.schema.json")
	if err != nil {
		loadErr = err
		return
	}
	type entry struct{ file, title string }
	var entries []entry
	for _, f := range files {
		raw, err := schemas.FS.ReadFile(f)
		if err != nil {
			loadErr = err
			return
		}
		doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
		if err != nil {
			loadErr = fmt.Errorf("schema %s: %w", f, err)
			return
		}
		var head struct {
			Title string `json:"title"`
		}
		json.Unmarshal(raw, &head)
		if err := c.AddResource(f, doc); err != nil {
			loadErr = err
			return
		}
		entries = append(entries, entry{f, head.Title})
	}
	for _, e := range entries {
		s, err := c.Compile(e.file)
		if err != nil {
			loadErr = fmt.Errorf("compile %s: %w", e.file, err)
			return
		}
		compiled[strings.ToLower(e.title)] = s
		names[strings.ToLower(e.title)] = e.title
	}
}

// Types returns the canonical names of every embedded type, sorted.
func Types() []string {
	once.Do(load)
	out := make([]string, 0, len(names))
	for _, n := range names {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// Canonical returns the canonical name of a type (case-insensitive), or an error listing the known types.
func Canonical(typ string) (string, error) {
	once.Do(load)
	if loadErr != nil {
		return "", loadErr
	}
	if n, ok := names[strings.ToLower(typ)]; ok {
		return n, nil
	}
	return "", fmt.Errorf("unknown type %q (known: %s)", typ, strings.Join(Types(), ", "))
}

// Validate checks one JSON document (bytes) against type typ.
func Validate(typ string, doc []byte) ([]Error, error) {
	n, err := Canonical(typ)
	if err != nil {
		return nil, err
	}
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(doc))
	if err != nil {
		return []Error{{InstanceLocation: "", KeywordLocation: "", Message: "not valid JSON: " + err.Error()}}, nil
	}
	return ValidateValue(n, inst)
}

// ValidateValue checks an already-decoded value (as produced by jsonschema.UnmarshalJSON or encoding/json
// with UseNumber) against type typ.
func ValidateValue(typ string, inst any) ([]Error, error) {
	n, err := Canonical(typ)
	if err != nil {
		return nil, err
	}
	s := compiled[strings.ToLower(n)]
	verr := s.Validate(inst)
	if verr == nil {
		return nil, nil
	}
	var ve *jsonschema.ValidationError
	if !errors.As(verr, &ve) {
		return nil, verr
	}
	var out []Error
	for _, u := range ve.BasicOutput().Errors {
		msg := ""
		if u.Error != nil {
			msg = u.Error.String()
		}
		if msg == "" || strings.HasPrefix(msg, "validation failed") {
			continue // container node; its children carry the concrete errors
		}
		out = append(out, Error{InstanceLocation: u.InstanceLocation, KeywordLocation: u.KeywordLocation, Message: msg})
	}
	if len(out) == 0 {
		out = append(out, Error{Message: verr.Error()})
	}
	Sort(out)
	return dedupe(out), nil
}

// Sort orders errors deterministically.
func Sort(errs []Error) {
	idx := func(e Error) int {
		if e.Index == nil {
			return -1
		}
		return *e.Index
	}
	sort.SliceStable(errs, func(i, j int) bool {
		a, b := errs[i], errs[j]
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if idx(a) != idx(b) {
			return idx(a) < idx(b)
		}
		if a.InstanceLocation != b.InstanceLocation {
			return a.InstanceLocation < b.InstanceLocation
		}
		if a.KeywordLocation != b.KeywordLocation {
			return a.KeywordLocation < b.KeywordLocation
		}
		return a.Message < b.Message
	})
}

func dedupe(errs []Error) []Error {
	out := errs[:0]
	for i, e := range errs {
		if i > 0 && e.Line == errs[i-1].Line && e.InstanceLocation == errs[i-1].InstanceLocation &&
			e.KeywordLocation == errs[i-1].KeywordLocation && e.Message == errs[i-1].Message {
			continue
		}
		out = append(out, e)
	}
	return out
}

// IsArray reports whether doc's top level is a JSON array.
func IsArray(doc []byte) bool {
	t := bytes.TrimLeft(doc, " \t\r\n")
	return len(t) > 0 && t[0] == '['
}

// ValidateArray validates each element of a JSON array document, tagging errors with the element index.
func ValidateArray(typ string, doc []byte) ([]Error, int, error) {
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(doc))
	if err != nil {
		return []Error{{Message: "not valid JSON: " + err.Error()}}, 0, nil
	}
	arr, ok := inst.([]any)
	if !ok {
		return []Error{{Message: "expected a JSON array"}}, 0, nil
	}
	var out []Error
	for i, el := range arr {
		errs, err := ValidateValue(typ, el)
		if err != nil {
			return nil, 0, err
		}
		for _, e := range errs {
			k := i
			e.Index = &k
			out = append(out, e)
		}
	}
	Sort(out)
	return out, len(arr), nil
}
