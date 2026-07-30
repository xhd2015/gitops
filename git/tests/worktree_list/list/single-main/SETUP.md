# Scenario

**Feature**: List on a single main checkout

```
# only main on master
List(main) -> [Entry{Path=main, Branch=master, IsMain=true}]
```

## Preconditions

- Clean repo with one commit on `master`.

## Steps

1. Init temp repo on `master`.
2. Set `req.Op=list`, `req.Dir=main`.

## Context

- No linked worktrees registered.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	main := initRepo(t)
	req.Op = "list"
	req.Dir = main
	req.MainPath = main
	return nil
}
```
