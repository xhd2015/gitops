# Scenario

**Feature**: direct side chains which meet known history are deduplicated.

```text
# first direct merge recovers a -> b
base -> a -> b
base -> local -> merge1(P1=local, P2=b)

# second direct merge recovers c then stops at already recovered a
a -> c
merge1 -> merge2(P1=merge1, P2=c)
```

## Steps

1. Construct two direct merges whose second-parent paths reconverge at `a`.

```go
import (
    "testing"

    "github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    base := commit(t, req.Dir, "base")
    a := commit(t, req.Dir, "a", base)
    b := commit(t, req.Dir, "b", a)
    local := commit(t, req.Dir, "local", base)
    mergeOne := commit(t, req.Dir, "merge1", local, b)
    c := commit(t, req.Dir, "c", a)
    head := commit(t, req.Dir, "merge2", mergeOne, c)
    req.Base, req.Head = base, head
    req.Expected = []string{head, mergeOne, local, c, b, a}
    return nil
}
```
