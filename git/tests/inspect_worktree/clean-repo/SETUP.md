# Scenario

**Feature**: InspectWorktree on a clean repository after initial commit

```
# root Setup already committed "init" on master — no further mutations
InspectWorktree(req.Dir) -> IsRepo=true, IsClean=true
```

## Preconditions

- Git repo at `req.Dir` with one commit on `master`, message `init`.

## Steps

1. No additional setup; repo remains clean after root bootstrap.

## Context

- Validates branch, 7-char commit hash, and subject message.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Logf("clean repo at %s after root bootstrap — no mutations", req.Dir)
	return nil
}
```