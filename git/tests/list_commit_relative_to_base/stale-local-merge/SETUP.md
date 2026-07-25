# Scenario

**Bug**: a stale local merge places the remotely advanced branch line in P2.

```text
# remote advances from base through the covered anchor
base -> 9f -> 61 -> 2b

# stale local line is P1; its merge of remote is the new head
base -> 19b -> 698
              \-> 2b  (P2 of 698)
```

## Steps

1. Construct the stale-local merge shape from the incident.

```go
import (
    "testing"

    "github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    base := commit(t, req.Dir, "base")
    anchor := commit(t, req.Dir, "9f", base)
    sixtyOne := commit(t, req.Dir, "61", anchor)
    remote := commit(t, req.Dir, "2b", sixtyOne)
    local := commit(t, req.Dir, "19b", base)
    head := commit(t, req.Dir, "698", local, remote)
    req.Base, req.Head = base, head
    req.Expected = []string{head, local, remote, sixtyOne, anchor}
    return nil
}
```
