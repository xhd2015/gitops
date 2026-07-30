# Scenario

**Feature**: two linked worktrees on the same branch via worktree add --force

```
# two linked checkouts of feature (main stays on master)
WorktreesOnBranch(feature) -> len=2
```

## Preconditions

- Git supports `git worktree add --force <path> <existing-branch>`.

## Steps

1. Init main on `master`.
2. Add linked worktree with new branch `feature`.
3. Add second linked worktree on same `feature` with `--force`.
4. Query `WorktreesOnBranch(feature)`.

## Context

- Scenario 3 from P1 requirement; no refuse error — data only.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	main := initRepo(t)
	linked1 := addLinkedBranch(t, main, "feature")
	linked2 := addLinkedExistingBranch(t, main, "feature", true)
	req.Op = "worktrees_on_branch"
	req.Dir = main
	req.Branch = "feature"
	req.MainPath = main
	req.LinkedPath = linked1
	req.LinkedPath2 = linked2
	return nil
}
```
