# Scenario

**Feature**: ParseListPorcelain maps two worktree blocks

```
# main + feature porcelain blocks
ParseListPorcelain -> 2 entries; Branch without refs/heads/
```

## Preconditions

- None (pure string input).

## Steps

1. Set porcelain with two worktree blocks (main + feature).

## Context

- HEAD and Path fields populated from porcelain lines.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "parse_list"
	req.Porcelain = "" +
		"worktree /path/to/repo\n" +
		"HEAD 1234567890abcdef1234567890abcdef12345678\n" +
		"branch refs/heads/master\n" +
		"\n" +
		"worktree /path/to/repo-feature\n" +
		"HEAD abcdef1234567890abcdef1234567890abcdef12\n" +
		"branch refs/heads/feature\n"
	return nil
}
```
