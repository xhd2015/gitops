package git_isolated

import (
	"os/exec"
	"strings"
	"testing"
)

// MustRun runs git in dir and fatals on error.
func MustRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	if _, err := CombinedOutput(dir, args...); err != nil {
		t.Fatal(err)
	}
}

// MustOutput runs git in dir and returns trimmed stdout, fatalling on error.
func MustOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := Output(dir, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// WorktreeList returns raw git worktree list stdout.
func WorktreeList(t *testing.T, dir string) string {
	t.Helper()
	cmd := Command(dir, "worktree", "list")
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			t.Fatalf("git worktree list in %s: %v\n%s", dir, err, ee.Stderr)
		}
		t.Fatalf("git worktree list in %s: %v", dir, err)
	}
	return string(out)
}

// MustOutputError runs git expecting failure and returns trimmed combined output.
func MustOutputError(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := Command(dir, args...).CombinedOutput()
	if err == nil {
		t.Fatalf("git %v in %s: expected failure", args, dir)
	}
	return strings.TrimSpace(string(out))
}