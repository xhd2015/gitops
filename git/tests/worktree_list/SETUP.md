# Scenario

**Feature**: shared helpers for gitops worktree inventory leaves

```
# leaf Setup builds a temp main repo (+ optional linked worktrees)
leaf Setup -> git init/commit/worktree add -> req.Dir, paths
# Run calls List / ListLinked / WorktreesOnBranch / parse / resolve
```

## Preconditions

- `git` available in PATH
- Classic TDD: production APIs under `git/worktree` are expected missing/incomplete

## Steps

1. Root setup is a no-op; each leaf creates its own temp repository.
2. Helpers provide init, commit, worktree add, path canon, and path membership.

## Context

- Module: `github.com/xhd2015/gitops`
- Package under test: `github.com/xhd2015/gitops/git/worktree`
- No wrk product strings or refuse errors in scope

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/gitops/git/worktree"
	"github.com/xhd2015/xgo/support/cmd"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	t.Helper()
	// Leaves own fixture creation; root only exposes helpers.
	return nil
}

func initRepo(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "gitops-worktree-list-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	if err := cmd.Dir(dir).Run("git", "init"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Dir(dir).Run("git", "branch", "-M", "master"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Dir(dir).Run("git", "config", "user.email", "test@test.com"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Dir(dir).Run("git", "config", "user.name", "test"); err != nil {
		t.Fatal(err)
	}
	writeFile(t, dir, "README.md", "init\n")
	if err := cmd.Dir(dir).Run("git", "add", "-A"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Dir(dir).Run("git", "commit", "-m", "init"); err != nil {
		t.Fatal(err)
	}
	return canonPath(t, dir)
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	if err := cmd.Dir(dir).Run("git", args...); err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
}

func addLinkedBranch(t *testing.T, mainDir, branch string) string {
	t.Helper()
	linked, err := os.MkdirTemp("", "gitops-wt-linked-*")
	if err != nil {
		t.Fatal(err)
	}
	// worktree add requires non-existent path; MkdirTemp created it — remove first.
	if err := os.RemoveAll(linked); err != nil {
		t.Fatal(err)
	}
	gitRun(t, mainDir, "worktree", "add", "-b", branch, linked, "HEAD")
	return canonPath(t, linked)
}

func addLinkedExistingBranch(t *testing.T, mainDir, branch string, force bool) string {
	t.Helper()
	linked, err := os.MkdirTemp("", "gitops-wt-linked-*")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(linked); err != nil {
		t.Fatal(err)
	}
	args := []string{"worktree", "add"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, linked, branch)
	gitRun(t, mainDir, args...)
	return canonPath(t, linked)
}

func addLinkedDetach(t *testing.T, mainDir string) string {
	t.Helper()
	linked, err := os.MkdirTemp("", "gitops-wt-detach-*")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(linked); err != nil {
		t.Fatal(err)
	}
	gitRun(t, mainDir, "worktree", "add", "--detach", linked, "HEAD")
	return canonPath(t, linked)
}

func canonPath(t *testing.T, path string) string {
	t.Helper()
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		return resolved
	}
	// Dead or not-yet-created: resolve parent + base.
	dir, base := filepath.Dir(abs), filepath.Base(abs)
	if resolvedDir, err := filepath.EvalSymlinks(dir); err == nil {
		return filepath.Join(resolvedDir, base)
	}
	return filepath.Clean(abs)
}

func samePath(t *testing.T, a, b string) bool {
	t.Helper()
	return canonPath(t, a) == canonPath(t, b)
}

func entryPaths(entries []worktree.Entry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Path
	}
	return out
}

func containsPath(t *testing.T, entries []worktree.Entry, want string) bool {
	t.Helper()
	for _, e := range entries {
		if samePath(t, e.Path, want) {
			return true
		}
	}
	return false
}

func findByPath(t *testing.T, entries []worktree.Entry, want string) (worktree.Entry, bool) {
	t.Helper()
	for _, e := range entries {
		if samePath(t, e.Path, want) {
			return e, true
		}
	}
	return worktree.Entry{}, false
}
```
