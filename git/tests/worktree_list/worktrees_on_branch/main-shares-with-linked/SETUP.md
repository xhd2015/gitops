# Scenario

**Feature**: main shares branch with linked via checkout --ignore-other-worktrees

```
# linked on feature; main checks out feature ignoring other worktrees
WorktreesOnBranch(feature) -> len=2 (main + linked)
```

## Preconditions

- Git supports `git checkout --ignore-other-worktrees`.

## Steps

1. Init main on `master`.
2. Add linked on `feature`.
3. From main: `git checkout --ignore-other-worktrees feature`.
4. Query `WorktreesOnBranch(feature)`.

## Context

- Scenario 4 from P1 requirement.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	main := initRepo(t)
	linked := addLinkedBranch(t, main, "feature")
	gitRun(t, main, "checkout", "--ignore-other-worktrees", "feature")
	req.Op = "worktrees_on_branch"
	req.Dir = main
	req.Branch = "feature"
	req.MainPath = main
	req.LinkedPath = linked
	return nil
}
```
