# Scenario

**Feature**: IsSubmodule reports whether a path is a tracked submodule gitlink of
dir, distinguishing real submodules from nested separate repos.

## Preconditions

- `git` is available in PATH
- `dir` is a git repository
- Each leaf builds the `ext/` layout (submodule vs nested separate repo) and sets
  `req.Dir` (the root repo) and `req.Path = "ext"`.

## Steps

1. Create a temporary directory for the root repo and git-init it with an initial
   commit.
2. Leaf `Setup` configures the `ext/` directory (real submodule vs nested separate
   repo).
3. `req.Path` is set to `"ext"`.

## Context

- Go module: `github.com/xhd2015/gitops`
- Package under test: `git`
- `IsSubmodule(dir, path)` returns `(bool, error)` — true iff path is a tracked
  submodule gitlink of dir.

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
)

// Setup creates an isolated git repo at req.Dir with an initial commit, and sets
// req.Path = "ext". Leaf Setup functions wire ext/ as a submodule or a nested
// separate repo.
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir, err := os.MkdirTemp("", "gittest")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	exec.Command("git", "-C", dir, "init").Run()
	exec.Command("git", "-C", dir, "branch", "-M", "master").Run()
	exec.Command("git", "-C", dir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", dir, "config", "user.name", "Test User").Run()
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("x"), 0644)
	exec.Command("git", "-C", dir, "add", ".").Run()
	exec.Command("git", "-C", dir, "commit", "-m", "init").Run()
	req.Dir = dir
	req.Path = "ext"
	return nil
}
```
