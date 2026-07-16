# Scenario

**Feature**: P1 CachedDiff model types and DiffCachedParseError contract

```
# type-only surface — no DiffCached(dir) call
caller -> construct model.CachedDiff / FilePatch / Hunk
caller -> construct git.DiffCachedParseError{Dir,Raw,Err}
DiffCachedParseError -> Error() / Unwrap() / errors.As
```

## Preconditions

- Go module `github.com/xhd2015/gitops` is the workspace under test.
- Types under test are expected in `model` (`CachedDiff`, `FilePatch`, `Hunk`)
  and `git` (`DiffCachedParseError`).
- No git repository or `git` CLI is required for these pure type tests.

## Steps

1. Root setup marks the harness as type-level (no filesystem repo).
2. Each leaf sets `req.Mode` and scenario-specific fields.
3. Root `Run` constructs the types and returns field/error observations.

## Context

- Package: `github.com/xhd2015/gitops/git` and `.../model`
- Out of scope (P2+): `DiffCached` signature change, Unified*, commit_msg
- FileCount helpers are not required yet

```go
import (
	"fmt"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	// Type-level tree: no git repo; leaves fill Mode + inputs.
	t.Helper()
	if req == nil {
		return fmt.Errorf("nil request")
	}
	return nil
}
```
