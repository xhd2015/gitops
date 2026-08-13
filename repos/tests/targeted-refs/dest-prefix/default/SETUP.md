# Scenario

**Feature**: empty `TargetRefPrefix` stores the dest under `refs/gitops/targets/<64 hex>`, not a custom prefix.

```
FetchTargetedRefs(master, TargetRefPrefix="")
  -> dest refs/gitops/targets/{sha256("master")}
  -> no refs/other/targets/ dest
```

## Preconditions

- Root Setup seeded `master` + `feature` and injected `ReposRoot`.
- `TargetRefPrefix` left empty (default).

## Steps

1. Fetch `master` with empty prefix.
2. Assert dest is default-prefix + hash of input `master`.

## Context

- Default dest prefix is `refs/gitops/targets`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.Refs = []string{"master"}
	req.TargetRefPrefix = ""
	return nil
}
```
