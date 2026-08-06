# Scenario

**Feature**: ListIgnoredDirs returns the slash-relative ignored directory paths of a
git repo, or an empty slice when dir is not a git repo.

## Preconditions

- `git` is available in PATH
- `dir` is a directory (a git repo for the ignored-dir leaf; a plain dir for the
  not-a-git-repo leaf)

## Steps

1. Create a temporary directory.
2. Leaf `Setup` configures the repo (git init + .gitignore + fixtures) or leaves it
   a plain non-repo dir, then sets `req.Dir`.

## Context

- Go module: `github.com/xhd2015/gitops`
- Package under test: `git`
- `ListIgnoredDirs(dir)` returns `([]string, error)` — slash-relative ignored dir
  paths, no trailing slash, no `./`.

```go
import (
	"os"

	"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
)

// Setup creates an isolated temp dir and assigns it to req.Dir. Leaf Setup
// functions configure the repo (or leave it a plain non-repo dir).
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir, err := os.MkdirTemp("", "gittest")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	req.Dir = dir
	return nil
}

var _ = exec.Command
var _ = filepath.Join
```
