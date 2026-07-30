# Scenario

**Feature**: backup-style Counts via ParsePorcelain

```
porcelain text -> ParsePorcelain -> Counts
```

## Preconditions

- None.

## Steps

1. Grouping sets default Op to counts; leaves set porcelain.

## Context

- Backup taxonomy keeps Untracked distinct from Added.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	t.Helper()
	req.Op = "counts"
	return nil
}
```
