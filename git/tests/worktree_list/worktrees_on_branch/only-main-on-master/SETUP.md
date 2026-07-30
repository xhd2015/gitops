# Scenario

**Feature**: single main on master — filter hits master only

```
# only main on master
WorktreesOnBranch(master) -> len=1
WorktreesOnBranch(does-not-exist) -> len=0
```

## Preconditions

- Single main checkout on `master`.

## Steps

1. Init repo on `master`.
2. Query both `master` and a nonexistent branch via multi op.

## Context

- Scenario 1 from P1 requirement.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	main := initRepo(t)
	req.Op = "worktrees_on_branch_multi"
	req.Dir = main
	req.MainPath = main
	req.Branches = []string{"master", "does-not-exist"}
	return nil
}
```
