# Scenario

**Feature**: ListIgnoredDirs returns an empty slice (no error) when dir is not a git repo

## Preconditions
- `dir` is a plain directory, NOT a git repo (no `.git`).

## Steps
1. `dir` (created by the shared Setup) is left as a plain non-repo dir.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// req.Dir is a plain temp dir: deliberately do NOT run git init so it is not a
	// git repo. ListIgnoredDirs must return an empty slice with no error.
	if _, err := os.Stat(filepath.Join(req.Dir, ".git")); err == nil {
		t.Fatalf("precondition failed: %s is unexpectedly a git repo", req.Dir)
	}
	return nil
}
```
