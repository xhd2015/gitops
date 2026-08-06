# Scenario

**Feature**: unstage a single file, verify it is gone from index but still on disk

```
# stage a.txt -> restore a.txt -> a.txt gone from index, a.txt still on disk
git add a.txt -> git restore --staged -- a.txt -> diff --cached is empty, a.txt exists
```

## Preconditions

- A git repository with an initial commit.
- `a.txt` has been created and staged.

## Steps

1. Create and stage `a.txt`.
2. Set `req.Paths = ["a.txt"]`.
3. Run `RestoreStaged`, then verify `a.txt` is not in the staged list.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if err := stageFile(req.Dir, "a.txt", "hello"); err != nil {
		return err
	}
	req.Paths = []string{"a.txt"}
	return nil
}
```
