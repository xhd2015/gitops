# Scenario

**Feature**: shared temp-repo bootstrap for dual clean leaves

```
# root creates clean master repo; leaves may add dirty files
IsDiffClean / IsPorcelainClean on req.Dir
```

## Preconditions

- `git` available in PATH

## Steps

1. Create temp directory and initialize git repository on `master`.
2. Write `README.md`, commit with message `init`.
3. Set `req.Dir` to the repository path.
4. Leaf setups may add untracked or modified files.

## Context

- Module: `github.com/xhd2015/gitops`
- Package under test: `github.com/xhd2015/gitops/git/worktree`
- Existing `IsClean` is historically diff-based; P1 names dual APIs distinctly

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/xgo/support/cmd"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	t.Helper()
	dir, err := os.MkdirTemp("", "gitops-worktree-clean-*")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	if err := cmd.Dir(dir).Run("git", "init"); err != nil {
		return err
	}
	if err := cmd.Dir(dir).Run("git", "branch", "-M", "master"); err != nil {
		return err
	}
	if err := cmd.Dir(dir).Run("git", "config", "user.email", "test@test.com"); err != nil {
		return err
	}
	if err := cmd.Dir(dir).Run("git", "config", "user.name", "test"); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("init\n"), 0644); err != nil {
		return err
	}
	if err := cmd.Dir(dir).Run("git", "add", "-A"); err != nil {
		return err
	}
	if err := cmd.Dir(dir).Run("git", "commit", "-m", "init"); err != nil {
		return err
	}

	req.Dir = dir
	return nil
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
```
