# Scenario

**Feature**: `FetchTargetedRefs` with an empty refs list returns an error.

```
FetchTargetedRefs(cloneURL, refs=[], opts) -> error
```

## Preconditions

- Root Setup default `Refs=["master"]` must be cleared.

## Steps

1. Set `Refs` to empty.
2. Run `FetchTargetedRefs`.

## Context

- Empty input is invalid; do not fetch `+refs/heads/*` as a fallback.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.Refs = nil
	return nil
}
```
