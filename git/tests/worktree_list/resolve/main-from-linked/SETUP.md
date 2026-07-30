# Scenario

**Feature**: ResolveMainRepo from a linked worktree path

```
# linked worktree
ResolveMainRepo(linked) -> main repo path
```

## Preconditions

- Main + linked on `feature`.

## Steps

1. Init main; add linked branch `feature`.
2. Set `req.Op=resolve_main`, `req.Path=linked`.

## Context

- Path comparison uses canon/symlink resolution.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	main := initRepo(t)
	linked := addLinkedBranch(t, main, "feature")
	req.Op = "resolve_main"
	req.Path = linked
	req.MainPath = main
	req.LinkedPath = linked
	req.Dir = main
	return nil
}
```
