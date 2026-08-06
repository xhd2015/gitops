# Scenario

**Feature**: DiffCached on a non-git directory returns an ordinary error

```
# plain temp dir without git init
Setup -> req.Dir = empty directory (no .git) -> DiffCached -> (nil, err)
# err is ordinary (git failure), not required to be *DiffCachedParseError
```

## Preconditions

- Directory exists but is not inside a git work tree.

## Steps

1. Replace `req.Dir` with a fresh temp directory that was never `git init`'d.
2. Run `git.DiffCached(req.Dir)`.

## Context

- Root setup initialized a repo; this leaf overrides `req.Dir`.

```go
import (
	"os"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir, err := os.MkdirTemp("", "not-a-git-repo-diff-cached")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	req.Dir = dir
	return nil
}
```
