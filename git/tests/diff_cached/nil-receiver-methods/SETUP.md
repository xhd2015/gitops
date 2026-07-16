# Scenario

**Feature**: nil *CachedDiff method receivers are safe zero values

```
# no git call — pure method contract on a nil pointer
nil *model.CachedDiff -> FileCount() == 0
nil *model.CachedDiff -> Unified() == ""
nil *model.CachedDiff -> UnifiedTruncated(24) == ""
```

## Preconditions

- Does not require a live git repository for assertions (root SETUP still runs).
- `req.Mode = "nil-methods"` so Run returns `Diff == nil` without calling git.

## Steps

1. Set `req.Mode` to `nil-methods`.
2. Run returns `(*Response){Diff: nil}` with nil error.
3. Assert calls FileCount / Unified / UnifiedTruncated on the nil receiver.

## Context

- Classic RED: methods do not exist yet on `*model.CachedDiff`.
- Contract from P3 API: nil receiver → FileCount 0, Unified `""`, UnifiedTruncated `""`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "nil-methods"
	return nil
}
```
