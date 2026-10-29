// Package fsx holds the kernel's file helpers: deterministic JSON encoding, atomic writes and JSONL I/O.
package fsx

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Marshal encodes v as indented JSON with a trailing newline. encoding/json sorts map keys and keeps struct
// field order, so equal values give identical bytes (D-050). HTML escaping is off so text stays readable.
func Marshal(v any) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Compact encodes v as one line of JSON without a trailing newline.
func Compact(v any) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(b.Bytes(), "\n"), nil
}

// WriteAtomic writes data to path through a temp file and rename, creating parent directories.
func WriteAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Chmod(name, 0o644); err != nil {
		os.Remove(name)
		return err
	}
	return os.Rename(name, path)
}

// WriteJSON marshals v deterministically and writes it atomically.
func WriteJSON(path string, v any) error {
	b, err := Marshal(v)
	if err != nil {
		return err
	}
	return WriteAtomic(path, b)
}

// WriteJSONL writes one compact JSON value per line, atomically.
func WriteJSONL[T any](path string, items []T) error {
	var b bytes.Buffer
	for _, it := range items {
		line, err := Compact(it)
		if err != nil {
			return err
		}
		b.Write(line)
		b.WriteByte('\n')
	}
	return WriteAtomic(path, b.Bytes())
}

// ReadJSON decodes the JSON file at path into v.
func ReadJSON(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

// Line is one non-empty line of a JSONL file with its 1-based line number.
type Line struct {
	N   int
	Raw []byte
}

// ReadLines returns the non-empty lines of a JSONL file.
func ReadLines(path string) ([]Line, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []Line
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	n := 0
	for sc.Scan() {
		n++
		raw := bytes.TrimSpace(sc.Bytes())
		if len(raw) == 0 {
			continue
		}
		out = append(out, Line{N: n, Raw: append([]byte(nil), raw...)})
	}
	return out, sc.Err()
}

// SHA256 returns the lowercase hex sha256 of b.
func SHA256(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

// Exists reports whether path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
