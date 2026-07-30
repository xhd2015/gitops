## Expected

- No error
- Exactly one entry
- Path is main, Branch is `master`, IsMain is true

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
		t.Fatalf("List len = %d, want 1; paths=%v", len(resp.Entries), entryPaths(resp.Entries))
	}
	e := resp.Entries[0]
	if !samePath(t, e.Path, req.MainPath) {
		t.Fatalf("Path = %q, want main %q", e.Path, req.MainPath)
	}
	if e.Branch != "master" {
		t.Fatalf("Branch = %q, want master", e.Branch)
	}
	if !e.IsMain {
		t.Fatal("IsMain = false, want true")
	}
}
```
