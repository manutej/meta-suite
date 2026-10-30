// Package ledger is S1's append-only event log, runs/<run_id>/ledger.jsonl. Each line is a LedgerEvent
// (validated by G1 before it is written) carrying a sequence number and the sha256 of the previous line, so
// any edit, deletion or reordering is detectable. Every append happens under the ledger's file lock and first
// re-verifies the whole chain, so nothing is ever appended to a tampered ledger.
package ledger

import (
	"bytes"
	"errors"

	"io"
	"os"

	"github.com/manutej/meta-suite/next/kernel/internal/fsx"
	"github.com/manutej/meta-suite/next/kernel/internal/kerr"
	"github.com/manutej/meta-suite/next/kernel/internal/lock"
	runs "github.com/manutej/meta-suite/next/kernel/internal/rundir"
	"github.com/manutej/meta-suite/next/kernel/internal/schema"
)

// Event is one ledger line.
type Event struct {
	Seq     int            `json:"seq"`
	At      string         `json:"at"`
	Run     string         `json:"run"`
	Command string         `json:"command"`
	OK      bool           `json:"ok"`
	Data    map[string]any `json:"data"`
	Prev    string         `json:"prev"`
}

type head struct {
	Seq  int    `json:"seq"`
	Hash string `json:"hash"`
}

// Append writes one event to the run's ledger.
func Append(r *runs.Run, command string, ok bool, data map[string]any) error {
	if data == nil {
		data = map[string]any{}
	}
	path := r.Path(runs.LedgerFile)
	return lock.With(path, func() error {
		h, err := scanHead(path)
		if err != nil {
			return err
		}
		prev := h.Hash
		if h.Seq == 0 {
			prev = "genesis"
		}
		ev := Event{Seq: h.Seq + 1, At: runs.Stamp(runs.Now()), Run: r.Record.RunID, Command: command, OK: ok, Data: data, Prev: prev}
		line, err := fsx.Compact(ev)
		if err != nil {
			return kerr.Wrap("internal", err)
		}
		errs, err := schema.Validate("LedgerEvent", line)
		if err != nil {
			return kerr.Wrap("internal", err)
		}
		if len(errs) > 0 {
			return kerr.New("kernel_fault", "ledger event would not validate").With(map[string]any{"errors": errs})
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return kerr.Wrap("io", err)
		}
		if _, err := f.Write(append(line, '\n')); err != nil {
			f.Close()
			return kerr.Wrap("io", err)
		}
		return kerr.Wrap("io", f.Close())
	})
}

func scanHead(path string) (head, error) {
	lines, err := fsx.ReadLines(path)
	if errors.Is(err, os.ErrNotExist) {
		return head{}, nil
	}
	if err != nil {
		return head{}, kerr.Wrap("io", err)
	}
	if len(lines) == 0 {
		return head{}, nil
	}
	if _, err := Verify(path); err != nil {
		return head{}, err
	}
	last := lines[len(lines)-1]
	var ev Event
	if err := jsonUnmarshal(last.Raw, &ev); err != nil {
		return head{}, kerr.Wrap("ledger_corrupt", err)
	}
	return head{Seq: ev.Seq, Hash: fsx.SHA256(last.Raw)}, nil
}

// Verify checks the hash chain and sequence of a ledger file and returns the number of events.
func Verify(path string) (int, error) {
	lines, err := fsx.ReadLines(path)
	if err != nil {
		return 0, kerr.Wrap("io", err)
	}
	prev := "genesis"
	for i, l := range lines {
		var ev Event
		if err := jsonUnmarshal(l.Raw, &ev); err != nil {
			return i, kerr.New("ledger_corrupt", "line %d: %v", l.N, err)
		}
		if ev.Seq != i+1 {
			return i, kerr.New("ledger_corrupt", "line %d: seq %d, expected %d", l.N, ev.Seq, i+1)
		}
		if ev.Prev != prev {
			return i, kerr.New("ledger_corrupt", "line %d: hash chain broken (prev %s, expected %s)", l.N, ev.Prev, prev)
		}
		prev = fsx.SHA256(l.Raw)
	}
	return len(lines), nil
}

// Show writes the ledger to w after verifying its chain.
func Show(r *runs.Run, w io.Writer) error {
	path := r.Path(runs.LedgerFile)
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return kerr.New("precondition", "%s has no ledger", r.Dir)
		}
		return kerr.Wrap("io", err)
	}
	if _, err := Verify(path); err != nil {
		return err
	}
	if len(b) > 0 && !bytes.HasSuffix(b, []byte("\n")) {
		b = append(b, '\n')
	}
	_, err = w.Write(b)
	return err
}

// Events reads all events.
func Events(r *runs.Run) ([]Event, error) {
	lines, err := fsx.ReadLines(r.Path(runs.LedgerFile))
	if err != nil {
		return nil, kerr.Wrap("io", err)
	}
	out := make([]Event, 0, len(lines))
	for _, l := range lines {
		var ev Event
		if err := jsonUnmarshal(l.Raw, &ev); err != nil {
			return nil, kerr.New("ledger_corrupt", "line %d: %v", l.N, err)
		}
		out = append(out, ev)
	}
	return out, nil
}
