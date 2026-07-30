# Scenario

**Feature**: List includes main + linked; ListLinked excludes main

```
# main on master + linked on feature
List -> 2 entries (main + linked)
ListLinked -> 1 entry (linked only, IsMain=false)
```

## Preconditions

- Main repo with initial commit.

## Steps

1. Init main on `master`.
2. `git worktree add -b feature <linked> HEAD`.
3. Set `req.Op=list_and_linked`.

## Context

- Covers scenario: List includes main; ListLinked excludes main.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	main := initRepo(t)
	linked := addLinkedBranch(t, main, "feature")
	req.Op = "list_and_linked"
	req.Dir = main
	req.MainPath = main
	req.LinkedPath = linked
	return nil
}
```
