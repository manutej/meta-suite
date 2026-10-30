package main

import (
	"bytes"
	"io"
	"path/filepath"
	"strconv"

	"github.com/manutej/meta-suite/next/kernel/internal/norm"
)

// maxStdin bounds what a command reads from stdin (a claim or an event is small).
const maxStdin = 4 << 20

func readAll(r io.Reader) ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	b, err := io.ReadAll(io.LimitReader(r, maxStdin+1))
	if err != nil {
		return nil, err
	}
	if len(b) > maxStdin {
		return nil, io.ErrShortBuffer
	}
	return b, nil
}

func trimSpace(b []byte) []byte { return bytes.TrimSpace(b) }

func normIndep(s string, strip []string) string { return norm.IndependenceKey(s, strip) }

func itoa(i int) string { return strconv.Itoa(i) }

func mustAbs(p string) string {
	a, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return a
}
