# Scenario

**Feature**: remote has `master` and `feature`; fetching `feature` only creates that dest target.

```
remote: master + feature
FetchTargetedRefs(["feature"])
  -> dest for feature exists
  -> dest for unused master not created
```

## Preconditions

- Root Setup seeded both `master` and `feature` on the `file://` remote.

## Steps

1. Request only `feature` (short branch name).
2. Assert dest hash of input `feature` exists and dest hash of `master` does not.

## Context

- Unused branch must not appear as a dest target (and must not appear as `refs/remotes/origin/*`).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.Refs = []string{"feature"}
	req.TargetRefPrefix = ""
	return nil
}
```
