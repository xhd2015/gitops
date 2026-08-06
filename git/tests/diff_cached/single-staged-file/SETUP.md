# Scenario

**Feature**: one staged file yields non-nil CachedDiff with Files populated

```
# stage file.txt -> DiffCached -> *CachedDiff{Files:[{NewPath:file.txt, Kind set}]}
write file.txt -> git add file.txt -> parse staged patch -> Files include file.txt
```

## Preconditions

- A git repository with an initial commit.

## Steps

1. Create and stage `file.txt` with content `hello`.
2. Run `git.DiffCached(req.Dir)`.

## Context

- Unified/raw body fields are not asserted yet (P3).
- Kind must be non-empty; exact kind vocabulary is free string for this phase
  (e.g. add/modify) as long as the staged path is reflected.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	return stageFile(req.Dir, "file.txt", "hello")
}
```
