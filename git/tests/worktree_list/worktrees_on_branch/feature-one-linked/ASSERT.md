## Expected

- Length 1
- Path is the linked worktree (not main)
- Branch is `feature`
- IsMain is false

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Entries) != 1 {
		t.Fatalf("WorktreesOnBranch(feature) len = %d, want 1; paths=%v", len(resp.Entries), entryPaths(resp.Entries))
	}
	e := resp.Entries[0]
	if !samePath(t, e.Path, req.LinkedPath) {
		t.Fatalf("path = %q, want linked %q", e.Path, req.LinkedPath)
	}
	if samePath(t, e.Path, req.MainPath) {
		t.Fatal("expected linked path, got main")
	}
	if e.Branch != "feature" {
		t.Fatalf("Branch = %q, want feature", e.Branch)
	}
	if e.IsMain {
		t.Fatal("IsMain = true, want false for linked")
	}
}
```
