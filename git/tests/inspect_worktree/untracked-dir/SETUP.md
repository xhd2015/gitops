# Scenario

**Feature**: untracked directory expands added count to files inside

```
# one untracked file + one untracked dir with 3 files -> Added=4, Uncommitted=2
```

## Preconditions

- Clean repo from root bootstrap.

## Steps

1. Create `solo.txt` (untracked file).
2. Create `nested/pkg/{a.go,b.go,c.go}` (untracked directory tree).

## Context

- Porcelain shows 2 lines; added count reflects 4 files total.

```go
import (
	"os"
	"testing"

	"github.com/xhd2015/doctest/session"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := req.Dir

	writeFile(t, dir, "solo.txt", "one")
	if err := os.MkdirAll(filepath.Join(dir, "nested", "pkg"), 0755); err != nil {
		return err
	}
	writeFile(t, dir, "nested/pkg/a.go", "a")
	writeFile(t, dir, "nested/pkg/b.go", "b")
	writeFile(t, dir, "nested/pkg/c.go", "c")
	return nil
}
```