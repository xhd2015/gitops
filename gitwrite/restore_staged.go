package gitwrite

import (
	"fmt"
	"os/exec"
	"strings"
)

// RestoreStaged removes paths from the git index in dir without touching the
// working copy.
//
// With at least one commit it runs `git restore --staged`, which resets each
// index entry to HEAD so modified tracked files stay tracked. On an unborn
// HEAD (git init, no commits) there is nothing to restore to; every staged
// path is new, so `git rm --cached` is exactly "unstage".
func RestoreStaged(dir string, paths ...string) error {
	if len(paths) == 0 {
		return nil
	}
	hasHead, err := repoHasHead(dir)
	if err != nil {
		return err
	}
	if !hasHead {
		args := append([]string{"-C", dir, "rm", "--cached", "--quiet", "--"}, paths...)
		cmd := exec.Command("git", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git rm --cached failed: %s: %w", strings.TrimSpace(string(output)), err)
		}
		return nil
	}
	args := append([]string{"-C", dir, "restore", "--staged", "--"}, paths...)
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git restore --staged failed: %s: %w", strings.TrimSpace(string(output)), err)
	}
	return nil
}

// repoHasHead reports whether HEAD in dir resolves to a commit. An unborn
// branch (zero commits) reports false: rev-parse -q exits 1 with no output.
// Other git failures (not a repo, git missing) propagate as errors.
func repoHasHead(dir string) (bool, error) {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--verify", "-q", "HEAD")
	output, err := cmd.CombinedOutput()
	if err == nil {
		return true, nil
	}
	if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 1 && strings.TrimSpace(string(output)) == "" {
		return false, nil
	}
	msg := strings.TrimSpace(string(output))
	if msg == "" {
		msg = err.Error()
	}
	return false, fmt.Errorf("git rev-parse --verify HEAD: %s: %w", msg, err)
}
