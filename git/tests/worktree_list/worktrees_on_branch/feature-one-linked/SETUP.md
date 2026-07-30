# Scenario

**Feature**: main on master + one linked on feature

```
# main master, linked feature
WorktreesOnBranch(feature) -> len=1, path=linked only
```

## Preconditions

- Main on `master` with one linked worktree on `feature`.

## Steps

1. Init main; add linked branch `feature`.
2. Query `WorktreesOnBranch(feature)`.

## Context

- Scenario 2 from P1 requirement.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	main := initRepo(t)
	linked := addLinkedBranch(t, main, "feature")
	req.Op = "worktrees_on_branch"
	req.Dir = main
	req.Branch = "feature"
	req.MainPath = main
	req.LinkedPath = linked
	return nil
}
```
