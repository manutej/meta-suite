//go:build windows

package lock

import (
	"errors"
	"os"
	"time"
)

// On Windows the lock is an exclusively created marker file (O_EXCL), retried until the timeout. A marker
// older than the timeout is treated as abandoned by a crashed process and removed.
func acquire(name string, timeout time.Duration) (func(), error) {
	deadline := time.Now().Add(timeout)
	for {
		f, err := os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			f.Close()
			return func() { os.Remove(name) }, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		if st, serr := os.Stat(name); serr == nil && time.Since(st.ModTime()) > timeout {
			os.Remove(name)
			continue
		}
		if time.Now().After(deadline) {
			return nil, errors.New("timed out waiting for lock")
		}
		time.Sleep(5 * time.Millisecond)
	}
}
