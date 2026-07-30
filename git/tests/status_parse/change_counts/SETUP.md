# Scenario

**Feature**: four-bucket ChangeCounts via ParseChangeCounts

```
porcelain text -> ParseChangeCounts -> ChangeCounts{Added,Changed,Renamed,Deleted}
```

## Preconditions

- None.

## Steps

1. Grouping sets default Op; leaves set porcelain.

## Context

- Inspect-compatible: `??`→Added (not a separate untracked field).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	t.Helper()
	req.Op = "change_counts"
	return nil
}
```
