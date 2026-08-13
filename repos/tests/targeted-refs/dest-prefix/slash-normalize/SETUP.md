# Scenario

**Feature**: a trailing slash on `TargetRefPrefix` does not produce a double-slash dest ref.

```
FetchTargetedRefs(master, TargetRefPrefix=refs/gitops/targets/)
  -> dest refs/gitops/targets/{sha256("master")}
  -> not refs/gitops/targets//{hash}
```

## Preconditions

- Same two-branch remote and injected `ReposRoot`.

## Steps

1. Set `TargetRefPrefix` to `refs/gitops/targets/` (trailing slash).
2. Fetch `master`.

## Context

- Prefix with/without trailing slash is one dest form.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.Refs = []string{"master"}
	req.TargetRefPrefix = DefaultTargetRefPrefix + "/"
	return nil
}
```
