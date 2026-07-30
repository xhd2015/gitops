# Scenario

**Feature**: dead worktree path still appears in List and IsDead

```
# add linked then rm -rf path (not pruned)
List -> still includes dead path
IsDead(deadPath) -> true
```

## Preconditions

- Main repo with a linked worktree on `feature`.

## Steps

1. Init main; add linked on `feature`.
2. `os.RemoveAll(linked)` without `git worktree prune`.
3. Set `req.Op=list` and record dead path.
4. Assert also checks `IsDead` via a second call pattern in Assert using helper.

## Context

- Dead entries remain registered until prune.

```go
import (
	"os"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	main := initRepo(t)
	linked := addLinkedBranch(t, main, "feature")
	if err := os.RemoveAll(linked); err != nil {
		return err
	}
	req.Op = "list"
	req.Dir = main
	req.MainPath = main
	req.DeadPath = linked
	req.LinkedPath = linked
	return nil
}
```
