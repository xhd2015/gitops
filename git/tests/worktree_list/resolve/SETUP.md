# Scenario

**Feature**: resolve helpers for main repo identity

```
# linked path -> ResolveMainRepo -> main checkout path
```

## Preconditions

- Root helpers for init / linked worktrees.

## Steps

1. Grouping only.

## Context

- ResolveMainRepo works from linked `.git` file.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	t.Helper()
	// Grouping node: leaves set resolve Op and paths.
	return nil
}
```
