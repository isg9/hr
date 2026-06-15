// Package errlog appends timestamped error entries to a log file
// (typically <vault>/.hr/err.txt). Designed to be safe to call with a
// nil receiver — callers don't need to guard.
package errlog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Log struct {
	path string
	mu   sync.Mutex
}

// New returns a Log that appends entries to path. The file and its
// parent directory are created lazily on first Write.
func New(path string) *Log {
	return &Log{path: path}
}

// Write appends one entry tagged with context. A nil log is a no-op.
// Errors writing to the log itself are deliberately swallowed to avoid
// recursive failure paths during sync.
func (l *Log) Write(context string, err error) {
	if l == nil || err == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if mkErr := os.MkdirAll(filepath.Dir(l.path), 0o755); mkErr != nil {
		return
	}
	f, openErr := os.OpenFile(
		l.path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if openErr != nil {
		return
	}
	defer f.Close()
	ts := time.Now().UTC().Format(time.RFC3339)
	fmt.Fprintf(f, "%s\t%s\t%s\n", ts, context, err.Error())
}
