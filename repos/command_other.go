//go:build !linux && !darwin
// +build !linux,!darwin

package repos

import (
	"context"
	"fmt"
	"os/exec"
)

func runCommandTree(ctx context.Context, cmd *exec.Cmd) error {
	// Fail closed rather than silently promising process-tree cancellation on an
	// unsupported host. The Git execution path is supported on Linux and macOS.
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("Git process-tree execution is unsupported on this operating system")
}
