//go:build linux || darwin
// +build linux darwin

package repos

import (
	"context"
	"os/exec"
	"syscall"
	"time"
)

// runCommandTree owns a new process group. CommandContext alone kills only the
// leader and can leave git-remote-https holding inherited output descriptors.
// Wait always completes before ownership (and thus any repository lock) returns.
func runCommandTree(ctx context.Context, cmd *exec.Cmd) error {
	scope := currentProcessScope()
	if scope != nil {
		joined, cancel := scope.join(ctx)
		defer cancel()
		ctx = joined
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
	if scope != nil {
		if err := scope.start(cmd); err != nil {
			return err
		}
		defer scope.finish(cmd.Process.Pid)
	} else if err := cmd.Start(); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		timer := time.NewTimer(2 * time.Second)
		defer timer.Stop()
		select {
		case <-done:
			// The leader may have exited before its descendants. Kill the remaining
			// group as well before telling the caller that its resource is released.
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		case <-timer.C:
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
		}
		return ctx.Err()
	}
}
