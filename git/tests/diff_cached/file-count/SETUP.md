# Scenario

**Feature**: FileCount reports the number of staged file patches

```
# two staged files -> DiffCached -> *CachedDiff.FileCount() == 2
stage a.txt + b.txt -> DiffCached(dir) -> d.FileCount() == 2
```

## Preconditions

- A git repository with an initial commit.
- Two distinct new files will be staged.

## Steps

1. Stage `a.txt` with content `alpha`.
2. Stage `b.txt` with content `beta`.
3. Run `git.DiffCached(req.Dir)`.
4. Assert calls `resp.Diff.FileCount()`.

## Context

- Classic RED: `(*model.CachedDiff).FileCount` does not exist yet.
- Count is the number of `Files` entries (file patches), not hunks or lines.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if err := stageFile(req.Dir, "a.txt", "alpha\n"); err != nil {
		return err
	}
	if err := stageFile(req.Dir, "b.txt", "beta\n"); err != nil {
		return err
	}
	return nil
}
```
