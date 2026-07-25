# Scenario

**Feature**: recovery does not recurse through an inner merge's P2.

```text
# hidden is only P2 of the recovered inner merge
base -> side -> inner(P1=side, P2=hidden)
base -> local -> outer(P1=local, P2=inner)
```

## Steps

1. Construct an outer merge whose recovered P2 path contains an inner merge.

```go
import (
    "testing"

    "github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    base := commit(t, req.Dir, "base")
    side := commit(t, req.Dir, "side", base)
    hidden := commit(t, req.Dir, "hidden", base)
    inner := commit(t, req.Dir, "inner", side, hidden)
    local := commit(t, req.Dir, "local", base)
    outer := commit(t, req.Dir, "outer", local, inner)
    req.Base, req.Head = base, outer
    req.Expected = []string{outer, local, inner, side}
    return nil
}
```
