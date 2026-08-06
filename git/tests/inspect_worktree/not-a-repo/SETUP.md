# Scenario

**Feature**: InspectWorktree on a non-git directory

```
# plain temp dir without git init
Setup -> req.Dir = empty directory (no .git)
```

## Preconditions

- Directory exists but is not inside a git work tree.

## Steps

1. Replace `req.Dir` with a fresh temp directory that was never `git init`'d.

## Context

- Root setup initialized a repo; this leaf overrides `req.Dir`.

```go
import (
	"os"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir, err := os.MkdirTemp("", "gitops-inspect-nogit-*")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	req.Dir = dir
	return nil
}
```