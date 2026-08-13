package repos

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type repoLock struct {
	file *os.File
}

func acquireRepoLock(cacheDir string) (*repoLock, error) {
	lock, _, err := acquireRepoLockContextTimed(context.Background(), cacheDir)
	return lock, err
}

func acquireRepoLockContext(ctx context.Context, cacheDir string) (*repoLock, error) {
	lock, _, err := acquireRepoLockContextTimed(ctx, cacheDir)
	return lock, err
}

// acquireRepoLockContextTimed is like acquireRepoLockContext but also returns how
// long the caller blocked waiting for the exclusive flock (includes successful
// non-blocking acquire as a near-zero wait).
func acquireRepoLockContextTimed(ctx context.Context, cacheDir string) (*repoLock, time.Duration, error) {
	start := time.Now()
	if err := os.MkdirAll(filepath.Dir(cacheDir), 0o755); err != nil {
		return nil, time.Since(start), fmt.Errorf("create cache parent dir: %w", err)
	}
	lockPath := cacheDir + ".lock"
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, time.Since(start), fmt.Errorf("open lock file: %w", err)
	}
	for {
		err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return &repoLock{file: f}, time.Since(start), nil
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) && !errors.Is(err, syscall.EAGAIN) {
			_ = f.Close()
			return nil, time.Since(start), fmt.Errorf("acquire repo lock: %w", err)
		}
		timer := time.NewTimer(25 * time.Millisecond)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			_ = f.Close()
			return nil, time.Since(start), fmt.Errorf("acquire repo lock: %w", ctx.Err())
		case <-timer.C:
		}
	}
}

func releaseRepoLock(lock *repoLock) error {
	if lock == nil || lock.file == nil {
		return nil
	}
	err := syscall.Flock(int(lock.file.Fd()), syscall.LOCK_UN)
	closeErr := lock.file.Close()
	if err != nil {
		return fmt.Errorf("release repo lock: %w", err)
	}
	if closeErr != nil {
		return fmt.Errorf("close lock file: %w", closeErr)
	}
	return nil
}
