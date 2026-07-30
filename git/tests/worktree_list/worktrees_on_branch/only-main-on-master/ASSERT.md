## Expected

- `WorktreesOnBranch(master)` length 1, path is main
- `WorktreesOnBranch(does-not-exist)` length 0
- No error

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	onMaster := resp.ByBranch["master"]
	if len(onMaster) != 1 {
		t.Fatalf("WorktreesOnBranch(master) len = %d, want 1; paths=%v", len(onMaster), entryPaths(onMaster))
	}
	if !samePath(t, onMaster[0].Path, req.MainPath) {
		t.Fatalf("master entry path = %q, want %q", onMaster[0].Path, req.MainPath)
	}
	if onMaster[0].Branch != "master" {
		t.Fatalf("Branch = %q, want master", onMaster[0].Branch)
	}
	other := resp.ByBranch["does-not-exist"]
	if len(other) != 0 {
		t.Fatalf("WorktreesOnBranch(does-not-exist) len = %d, want 0; paths=%v", len(other), entryPaths(other))
	}
}
```
