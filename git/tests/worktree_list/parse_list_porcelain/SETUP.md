# Scenario

**Feature**: pure ParseListPorcelain of worktree list porcelain text

```
# porcelain string (no git subprocess)
ParseListPorcelain(text) -> []Entry
```

## Preconditions

- No live git required for these leaves.

## Steps

1. Leaves set `req.Op=parse_list` and `req.Porcelain`.

## Context

- Branch refs strip `refs/heads/`; detached leaves Branch empty.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	t.Helper()
	req.Op = "parse_list"
	return nil
}
```
