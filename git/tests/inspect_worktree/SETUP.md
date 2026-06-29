# Scenario

**Feature**: shared git repo bootstrap for InspectWorktree leaves

```
# temp repo with initial commit on master
Setup -> git init -> branch -M master -> commit "init" -> req.Dir
```

## Preconditions

- `git` available in PATH

## Steps

1. Create temp directory and initialize git repository on `master`.
2. Write `README.md`, commit with message `init`.
3. Set `req.Dir` to the repository path.

## Context

- Module: `github.com/xhd2015/gitops`
- Package under test: `git`
- Leaf setups may replace `req.Dir` (e.g. not-a-repo) or mutate the repo.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/xgo/support/cmd"
)

func Setup(t *testing.T, req *Request) error {
	dir, err := os.MkdirTemp("", "gitops-inspect-*")
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
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("test"), 0644); err != nil {
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

func gitAddA(t *testing.T, dir string) {
	t.Helper()
	if err := cmd.Dir(dir).Run("git", "add", "-A"); err != nil {
		t.Fatal(err)
	}
}

func gitCommit(t *testing.T, dir, msg string) {
	t.Helper()
	if err := cmd.Dir(dir).Run("git", "commit", "-m", msg); err != nil {
		t.Fatal(err)
	}
}
```