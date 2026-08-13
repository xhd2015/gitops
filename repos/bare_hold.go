package repos

import (
	"path/filepath"
	"sync"
)

// Process-local bare cache holder refcount. Keyed by filepath.Clean(cacheDir).
// Used to refuse corrupt-bare RemoveAll while live AcquireCommitWorktree holds exist.
var (
	bareHoldMu    sync.Mutex
	bareHoldCount = make(map[string]int)
)

func bareHoldKey(cacheDir string) string {
	return filepath.Clean(cacheDir)
}

// AcquireBareHold increments the in-process holder count for cacheDir.
func AcquireBareHold(cacheDir string) {
	key := bareHoldKey(cacheDir)
	bareHoldMu.Lock()
	bareHoldCount[key]++
	bareHoldMu.Unlock()
}

// ReleaseBareHold decrements the in-process holder count for cacheDir.
// Count is clamped at 0 (never negative).
func ReleaseBareHold(cacheDir string) {
	key := bareHoldKey(cacheDir)
	bareHoldMu.Lock()
	if n := bareHoldCount[key]; n > 1 {
		bareHoldCount[key] = n - 1
	} else {
		delete(bareHoldCount, key)
	}
	bareHoldMu.Unlock()
}

// BareHoldCount returns the current in-process holder count for cacheDir.
func BareHoldCount(cacheDir string) int {
	key := bareHoldKey(cacheDir)
	bareHoldMu.Lock()
	n := bareHoldCount[key]
	bareHoldMu.Unlock()
	return n
}
