# Scenario

**Feature**: untracked-only dirties porcelain-clean but not diff-clean

```
# write untracked.txt only
IsDiffClean=true, IsPorcelainClean=false
```

## Preconditions

- Clean repo from root bootstrap.

## Steps

1. Write untracked file `untracked.txt` (do not `git add`).

## Context

- Scenario 10 from P1 requirement.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	writeFile(t, req.Dir, "untracked.txt", "only untracked\n")
	return nil
}
```
