# Scenario

**Feature**: IsSubmodule returns false for a nested separate repo (own .git, untracked)

## Preconditions
- `ext/` carries its OWN `.git` from a separate `git init`, and is NOT tracked by
  `dir` (no gitlink in dir's index).

## Steps
1. In `dir`, create `ext/` and run a separate `git init ext` so `ext/` has its own
   `.git`.
2. Do NOT add `ext` to `dir`'s index.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir
	ext := filepath.Join(dir, "ext")
	os.MkdirAll(ext, 0755)
	os.WriteFile(filepath.Join(ext, "go.mod"), []byte("module example.com/ext\n\ngo 1.22\n"), 0644)
	exec.Command("git", "-C", ext, "init").Run()
	exec.Command("git", "-C", ext, "branch", "-M", "master").Run()
	exec.Command("git", "-C", ext, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", ext, "config", "user.name", "Test User").Run()
	exec.Command("git", "-C", ext, "add", ".").Run()
	exec.Command("git", "-C", ext, "commit", "-m", "init").Run()
	// Deliberately do NOT add ext to dir's index.
	return nil
}
```
