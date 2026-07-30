# Scenario

**Feature**: dead linked worktree still returned by WorktreesOnBranch

```
# linked on feature then rm -rf (not pruned)
WorktreesOnBranch(feature) -> includes dead path
```

## Preconditions

- Linked worktree on `feature` removed from disk without prune.

## Steps

1. Init main; add linked on `feature`.
2. Remove linked directory.
3. Query `WorktreesOnBranch(feature)`.

## Context

- Scenario 5 from P1 requirement.

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
	req.Op = "worktrees_on_branch"
	req.Dir = main
	req.Branch = "feature"
	req.MainPath = main
	req.DeadPath = linked
	req.LinkedPath = linked
	return nil
}
```
