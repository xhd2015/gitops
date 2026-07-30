# Scenario

**Feature**: shared pure-parse setup for status package leaves

```
# leaf sets Op + Porcelain text
ParsePorcelain / ParseChangeCounts -> Counts / ChangeCounts
```

## Preconditions

- Pure functions; no git subprocess required.

## Steps

1. Root setup is a no-op; leaves set `req.Op` and `req.Porcelain`.

## Context

- Module: `github.com/xhd2015/gitops`
- Package under test: `github.com/xhd2015/gitops/git/status`
- Neutral names only — no `WrkCounts` / `FormatWrk`

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	t.Helper()
	// Pure-parse tree: leaves set Op and Porcelain.
	return nil
}
```
