# Scenario

**Feature**: mixed porcelain maps to backup Counts

```
2x " M" + 1x "??" -> Modified=2, Untracked=1
```

## Preconditions

- None.

## Steps

1. Set porcelain with two modified and one untracked line.

## Context

- Scenario 9 (optional backup-style leaf) from P1 requirement.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "counts"
	req.Porcelain = " M a.txt\n M b.txt\n?? c.txt"
	return nil
}
```
