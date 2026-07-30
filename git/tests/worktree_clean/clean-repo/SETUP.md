# Scenario

**Feature**: clean repository is both diff-clean and porcelain-clean

```
# no mutations after init commit
IsDiffClean=true, IsPorcelainClean=true
```

## Preconditions

- Root bootstrap left the repo clean.

## Steps

1. No additional mutations.

## Context

- Scenario 11 from P1 requirement.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	t.Helper()
	t.Logf("clean repo at %s — no mutations", req.Dir)
	return nil
}
```
