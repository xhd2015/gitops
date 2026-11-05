//go:build linux || darwin
// +build linux darwin

package repos

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ProcessScope is installed only in a single-task runner process. It extends
// parent cancellation to every Git helper, including legacy Background callers.
// Installing it in a multi-task HTTP worker is intentionally not supported.
type ProcessScope struct {
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.Mutex
	closed    bool
	uncertain bool
	readOnly  bool
	forbidGit bool
	groups    map[int]struct{}
	wg        sync.WaitGroup
}

var processScopeState struct {
	sync.RWMutex
	scope *ProcessScope
}

func InstallProcessScope(ctx context.Context) (*ProcessScope, error) {
	processScopeState.Lock()
	defer processScopeState.Unlock()
	if processScopeState.scope != nil {
		return nil, fmt.Errorf("Git process scope already installed")
	}
	owned, cancel := context.WithCancel(ctx)
	scope := &ProcessScope{ctx: owned, cancel: cancel, groups: map[int]struct{}{}}
	processScopeState.scope = scope
	return scope, nil
}
func (s *ProcessScope) RestrictToReads(noGit bool) {
	s.mu.Lock()
	s.readOnly = true
	s.forbidGit = noGit
	s.mu.Unlock()
}
func readOnlyCommand(args []string) bool {
	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg == "-c" || arg == "-C" || arg == "--git-dir" || arg == "--work-tree" {
			i++
			continue
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		switch arg {
		case "version", "rev-parse", "cat-file", "show", "diff", "merge-base", "log", "for-each-ref", "rev-list", "ls-tree", "status", "check-attr", "check-ignore", "describe", "ls-files", "grep":
			return true
		default:
			return false
		}
	}
	return false
}

func currentProcessScope() *ProcessScope {
	processScopeState.RLock()
	defer processScopeState.RUnlock()
	return processScopeState.scope
}
func (s *ProcessScope) join(ctx context.Context) (context.Context, context.CancelFunc) {
	joined, cancel := context.WithCancel(ctx)
	go func() {
		select {
		case <-s.ctx.Done():
			cancel()
		case <-joined.Done():
		}
	}()
	return joined, cancel
}
func (s *ProcessScope) start(cmd *exec.Cmd) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed || s.ctx.Err() != nil {
		return context.Canceled
	}
	if filepath.Base(cmd.Path) == "git" && (s.forbidGit || s.readOnly && !readOnlyCommand(cmd.Args)) {
		return fmt.Errorf("Git preparation is forbidden during scheduled computation")
	}
	// Start/register and sealing are serialized. No PID can be born in the gap
	// between enumerating owned groups and declaring the scope drained.
	if err := cmd.Start(); err != nil {
		return err
	}
	s.groups[cmd.Process.Pid] = struct{}{}
	s.wg.Add(1)
	return nil
}
func (s *ProcessScope) finish(pid int) {
	// A successfully exited leader may still have descendants. They must not
	// survive a command's ownership interval or its worker reservation.
	_ = syscall.Kill(-pid, syscall.SIGKILL)
	deadline := time.NewTimer(2 * time.Second)
	defer deadline.Stop()
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	gone := false
	for !gone {
		if syscall.Kill(-pid, 0) == syscall.ESRCH {
			gone = true
			break
		}
		select {
		case <-deadline.C:
			goto finished
		case <-tick.C:
		}
	}
finished:
	s.mu.Lock()
	delete(s.groups, pid)
	if !gone {
		s.uncertain = true
	}
	s.mu.Unlock()
	s.wg.Done()
}

// CloseAndWait seals the scope, cancels even legacy helpers and waits for their
// actual command Wait paths. A timeout is not stop proof; callers fail closed.
func (s *ProcessScope) CloseAndWait(ctx context.Context) bool {
	s.mu.Lock()
	s.closed = true
	s.cancel()
	s.mu.Unlock()
	done := make(chan struct{})
	go func() { s.wg.Wait(); close(done) }()
	select {
	case <-done:
		s.mu.Lock()
		defer s.mu.Unlock()
		return !s.uncertain
	case <-ctx.Done():
		return false
	}
}
