# Scenario

**Feature**: observed fetch argv includes `--filter=blob:none` and does not include `+refs/heads/*`.

```
FetchTargetedRefs(master) + Progress
  -> argv contains --filter=blob:none
  -> argv does not contain +refs/heads/*
```

## Preconditions

- Cold targeted cache (first fetch) so Progress records the fetch command.
- Root `Run` injects a Progress buffer.

## Steps

1. Fetch `master` with default prefix.
2. Inspect Progress for conservative fetch argv.

## Context

- Reference argv: `git fetch --no-tags --filter=blob:none --depth=1 origin +<source>:<dest>`.
- A refspec `+refs/heads/master:<dest>` is allowed; the glob `+refs/heads/*` is not.

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
