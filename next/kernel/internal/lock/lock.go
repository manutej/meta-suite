// Package lock provides an exclusive advisory lock on a file path, used for every append (ledger, claim
// shards) and every read-modify-write of run.json. The lock file is <path>.lock.
package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Timeout bounds how long a caller waits for a lock before failing.
var Timeout = 30 * time.Second

// With runs fn while holding the exclusive lock for path.
func With(path string, fn func() error) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	release, err := acquire(path+".lock", Timeout)
	if err != nil {
		return fmt.Errorf("lock %s: %w", path, err)
	}
	defer release()
	return fn()
}
