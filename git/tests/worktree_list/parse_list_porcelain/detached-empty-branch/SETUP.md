# Scenario

**Feature**: detached porcelain entry has empty Branch

```
# worktree block with "detached" line (no branch refs/heads/)
ParseListPorcelain -> Branch == ""
```

## Preconditions

- None (pure string input).

## Steps

1. Set porcelain for a detached worktree block.

## Context

- Detached never matches named-branch filters.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "parse_list"
	req.Porcelain = "" +
		"worktree /path/to/detached\n" +
		"HEAD aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n" +
		"detached\n"
	return nil
}
```
