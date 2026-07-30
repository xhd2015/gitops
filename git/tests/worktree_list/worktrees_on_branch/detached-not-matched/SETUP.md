# Scenario

**Feature**: detached linked worktree is not included for a named branch

```
# linked --detach HEAD
WorktreesOnBranch(master) -> only main (detached excluded)
```

## Preconditions

- Main on `master`; linked worktree added with `--detach`.

## Steps

1. Init main on `master`.
2. Add detached linked worktree at HEAD.
3. Query `WorktreesOnBranch(master)` and also multi-query to ensure detached not under any name used.

## Context

- Scenario 6 from P1 requirement; detached Branch is empty.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	main := initRepo(t)
	detached := addLinkedDetach(t, main)
	req.Op = "worktrees_on_branch_multi"
	req.Dir = main
	req.MainPath = main
	req.LinkedPath = detached
	req.Branches = []string{"master", "feature"}
	return nil
}
```
