# Scenario

**Feature**: show toplevel returns the repository root for nested directories

```
# caller asks the helper for the worktree root containing a nested directory
caller -> ShowToplevel(<repo>/go-pkgs)

# helper returns Git's raw command-style toplevel output
caller <- ShowToplevel: <repo>\n
```

## Preconditions

- `git` is available in `PATH`.

## Steps

1. Create a temporary Git repository.
2. Create a nested `go-pkgs` directory inside it.
3. Set the request directory to the nested directory.
4. Record the expected raw toplevel output as the canonical repository path plus a trailing newline.

## Context

- The temporary repository shape mirrors a monorepo with a nested Go module directory.
- The function under test must preserve its raw output contract, including the trailing newline.

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir, err := os.MkdirTemp("", "gitops-show-toplevel-*")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}

	nested := filepath.Join(dir, "go-pkgs")
	if err := os.MkdirAll(nested, 0755); err != nil {
		return err
	}

	canonRepo, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return err
	}
	canonNested, err := filepath.EvalSymlinks(nested)
	if err != nil {
		return err
	}

	req.RepoDir = canonRepo
	req.NestedDir = canonNested
	req.ExpectedTop = canonRepo + "\n"
	req.GitDir = filepath.Join(canonRepo, ".git")
	return nil
}
```
