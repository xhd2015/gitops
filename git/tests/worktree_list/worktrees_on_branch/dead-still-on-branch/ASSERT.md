## Expected

- Length ≥ 1
- Dead linked path is present
- Branch is `feature` on that entry

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/gitops/git/worktree"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Entries) < 1 {
		t.Fatal("WorktreesOnBranch(feature) empty, want dead linked present")
	}
	if !containsPath(t, resp.Entries, req.DeadPath) {
		t.Fatalf("missing dead path %q; paths=%v", req.DeadPath, entryPaths(resp.Entries))
	}
	e, ok := findByPath(t, resp.Entries, req.DeadPath)
	if !ok {
		t.Fatal("dead entry not found")
	}
	if e.Branch != "feature" {
		t.Fatalf("dead entry Branch = %q, want feature", e.Branch)
	}
	if !worktree.IsDead(req.DeadPath) {
		t.Fatalf("IsDead(%q) = false, want true", req.DeadPath)
	}
}
```
