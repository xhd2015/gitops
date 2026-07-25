# Scenario

**Feature**: an ordinary first-parent path remains unchanged.

```text
# no merge commits exist between base and head
base -> a -> b
```

## Steps

1. Construct a linear base-to-head graph.

```go
import (
    "testing"

    "github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    base := commit(t, req.Dir, "base")
    a := commit(t, req.Dir, "a", base)
    b := commit(t, req.Dir, "b", a)
    req.Base, req.Head = base, b
    req.Expected = []string{b, a}
    return nil
}
```
