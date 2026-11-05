package repos

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestRepoLockBackpressureDoesNotReleaseHolder(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "repo")
	held, err := acquireRepoLock(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer releaseRepoLock(held)
	begin := time.Now()
	_, _, err = acquireRepoLockContextTimed(WithRepoLockWait(context.Background(), 10*time.Millisecond), dir)
	if !errors.Is(err, ErrRepoBusy) {
		t.Fatalf("error=%v", err)
	}
	if time.Since(begin) > time.Second {
		t.Fatal("short lock wait was unbounded")
	}
	// A second attempt must still see the original owner's lock.
	_, _, err = acquireRepoLockContextTimed(WithRepoLockWait(context.Background(), 0), dir)
	if !errors.Is(err, ErrRepoBusy) {
		t.Fatalf("writer lost its lock: %v", err)
	}
	releaseRepoLock(held)
	lock, _, err := acquireRepoLockContextTimed(WithRepoLockWait(context.Background(), 0), dir)
	if err != nil {
		t.Fatal(err)
	}
	releaseRepoLock(lock)
}
func TestCancelledStaticPreparationDoesNotCreateCache(t *testing.T) {
	root := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := EnsureStaticBareCache(ctx, "file:///unused", BareCacheOptions{ReposRoot: root})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v", err)
	}
	entries, _ := os.ReadDir(root)
	if len(entries) != 0 {
		t.Fatal("cancelled preparation wrote a cache")
	}
}
func TestGitCommandCancellationTerminatesDescendants(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pidFile := filepath.Join(t.TempDir(), "child.pid")
	cmd := exec.Command("sh", "-c", `sleep 60 & child=$!; echo "$child" > "$1"; wait`, "probe", pidFile)
	done := make(chan error, 1)
	go func() { done <- runCommandTree(ctx, cmd) }()
	deadline := time.After(3 * time.Second)
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()
	var child int
wait:
	for {
		select {
		case <-deadline:
			t.Fatal("child did not start")
		case <-tick.C:
			b, e := os.ReadFile(pidFile)
			if e == nil {
				child, _ = strconv.Atoi(strings.TrimSpace(string(b)))
				if child > 0 {
					break wait
				}
			}
		}
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error=%v", err)
		}
	case <-time.After(4 * time.Second):
		t.Fatal("cancellation left command blocked")
	}
	// A briefly unreaped zombie may still have a PID on Linux; it cannot keep
	// doing work. On macOS the process disappears once the shell reaps it.
	until := time.Now().Add(2 * time.Second)
	for time.Now().Before(until) {
		if syscall.Kill(child, 0) != nil {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	if b, e := os.ReadFile(filepath.Join("/proc", strconv.Itoa(child), "stat")); e == nil && strings.Contains(string(b), ") Z ") {
		return
	}
	t.Fatalf("descendant %d remains alive after cancellation", child)
}
