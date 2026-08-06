# Scenario

**Feature**: unstage multiple files at once

```
# stage a.txt b.txt -> restore a.txt b.txt -> both gone from index, both on disk
git add a.txt b.txt -> git restore --staged -- a.txt b.txt -> diff --cached empty, both exist
```

## Preconditions

- A git repository with an initial commit.
- `a.txt` and `b.txt` have been created and staged.

## Steps

1. Create and stage `a.txt` and `b.txt`.
2. Set `req.Paths = ["a.txt", "b.txt"]`.
3. Run `RestoreStaged`, then verify both are unstaged.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if err := stageFile(req.Dir, "a.txt", "aaa"); err != nil {
		return err
	}
	if err := stageFile(req.Dir, "b.txt", "bbb"); err != nil {
		return err
	}
	req.Paths = []string{"a.txt", "b.txt"}
	return nil
}
```
