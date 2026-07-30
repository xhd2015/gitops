## Expected

- `List` still includes the removed linked path
- `worktree.IsDead(deadPath)` is true
- Main remains present

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
	if !containsPath(t, resp.Entries, req.MainPath) {
		t.Fatalf("List missing main; paths=%v", entryPaths(resp.Entries))
	}
	if !containsPath(t, resp.Entries, req.DeadPath) {
		t.Fatalf("List missing dead linked %q; paths=%v", req.DeadPath, entryPaths(resp.Entries))
	}
	if !worktree.IsDead(req.DeadPath) {
		t.Fatalf("IsDead(%q) = false, want true", req.DeadPath)
	}
	if worktree.IsDead(req.MainPath) {
		t.Fatalf("IsDead(main) = true, want false")
	}
}
```
