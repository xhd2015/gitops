# Scenario

**Feature**: representative porcelain maps to four ChangeCounts buckets

```
?? + M + R + D lines -> Added=1, Changed=1, Renamed=1, Deleted=1
```

## Preconditions

- None.

## Steps

1. Set one porcelain line per bucket (`??` counts as Added).

## Context

- Scenario 8 from P1 requirement; inspect foundation (scenario 12).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "change_counts"
	req.Porcelain = " M modified.txt\n?? untracked.txt\nR  old.txt -> new.txt\n D deleted.txt"
	return nil
}
```
