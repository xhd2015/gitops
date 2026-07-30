# Scenario

**Feature**: WorktreesOnBranch filters registered entries by branch name

```
# fixture with 0..N worktrees on a named branch
WorktreesOnBranch(repo, branch) -> matching Entry slice (no policy error)
```

## Preconditions

- Root helpers for init / linked worktrees.

## Steps

1. Grouping only; leaves configure branches and fixtures.

## Context

- Detached never matches; dead still matches until prune; len>1 is allowed data.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	t.Helper()
	// Grouping node: leaves set Branch/Branches and fixtures.
	return nil
}
```
