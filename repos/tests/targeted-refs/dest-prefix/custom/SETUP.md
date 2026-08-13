# Scenario

**Feature**: `TargetRefPrefix=refs/other/targets` stores dest under that prefix using the same hash rule.

```
FetchTargetedRefs(master, TargetRefPrefix=refs/other/targets)
  -> dest refs/other/targets/{sha256("master")}
  -> no refs/gitops/targets/ dest
```

## Preconditions

- Same two-branch remote and injected `ReposRoot`.
- Prefix is a neutral custom value (not the default).

## Steps

1. Set `TargetRefPrefix` to `refs/other/targets`.
2. Fetch `master`.

## Context

- Hash identity is independent of prefix: suffix is still sha256 of input `master`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.Refs = []string{"master"}
	req.TargetRefPrefix = CustomTargetRefPrefix
	return nil
}
```
