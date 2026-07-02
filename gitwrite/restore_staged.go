package gitwrite

import (
	"fmt"
	"os/exec"
	"strings"
)

func RestoreStaged(dir string, paths ...string) error {
	args := append([]string{"-C", dir, "restore", "--staged", "--"}, paths...)
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git restore --staged failed: %s: %w", strings.TrimSpace(string(output)), err)
	}
	return nil
}
