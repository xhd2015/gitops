# Scenario

**Feature**: List / ListLinked inventory of registered worktrees

```
# leaf builds main (+ optional linked)
List(main) / ListLinked(main) -> Entry slices
```

## Preconditions

- Temp main repo helpers from root SETUP.

## Steps

1. Grouping only; leaves set `req.Op` and fixture paths.

## Context

- `List` includes main; `ListLinked` excludes main.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	t.Helper()
	// Grouping node: leaves set Op and build fixtures.
	return nil
}
```
